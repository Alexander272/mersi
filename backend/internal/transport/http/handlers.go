package http

import (
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/Alexander272/mersi/backend/internal/config"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	httpV1 "github.com/Alexander272/mersi/backend/internal/transport/http/v1"
	"github.com/Alexander272/mersi/backend/pkg/auth"
	"github.com/Alexander272/mersi/backend/pkg/error_bot"
	"github.com/Alexander272/mersi/backend/pkg/limiter"
	"github.com/Alexander272/mersi/backend/web"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	keycloak *auth.KeycloakClient
	services *services.Services
}

func NewHandler(services *services.Services, keycloak *auth.KeycloakClient) *Handler {
	return &Handler{
		services: services,
		keycloak: keycloak,
	}
}

func (h *Handler) Init(conf *config.Config) *gin.Engine {
	router := gin.Default()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	router.Use(
		gin.CustomRecovery(h.ErrorHandler),
		securityHeaders(),
	)

	// Init router
	router.GET("/api/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router, conf)
	h.initStatic(router)

	return router
}

func (h *Handler) ErrorHandler(c *gin.Context, origErr any) {
	err := fmt.Errorf("unexpected error: %v", origErr)
	error_bot.Send(c, err.Error(), gin.H{"PANIC": true, "Stack trace": string(debug.Stack())})
	debug.PrintStack()
	response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла непредвиденная ошибка: "+err.Error())
}

func securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "no-referrer-when-downgrade")
		c.Header("Content-Security-Policy",
			"default-src 'self' http: https: data: blob: 'unsafe-inline'")
		c.Header("Strict-Transport-Security",
			"max-age=31536000; includeSubDomains")
		c.Next()
	}
}

func (h *Handler) initAPI(router *gin.Engine, conf *config.Config) {
	middleware := middleware.NewMiddleware(h.services, conf.Auth, h.keycloak)
	handlerV1 := httpV1.NewHandler(httpV1.Deps{Services: h.services, Conf: conf, Middleware: middleware})

	api := router.Group("/api")
	api.Use(limiter.Limit(conf.Limiter.RPS, conf.Limiter.Burst, conf.Limiter.TTL))
	handlerV1.Init(api)
}

var appStartTime = time.Now()

const (
	frontendRoot = "frontend"
	indexFile    = "index.html"
	assetsPrefix = "assets/"
)

func (h *Handler) initStatic(router *gin.Engine) {
	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.Status(http.StatusNotFound)
			return
		}

		filePath := strings.TrimPrefix(c.Request.URL.Path, "/")
		if filePath == "" {
			filePath = indexFile
		}
		filePath = path.Clean(filePath)

		var f fs.File
		var err error
		openPath := frontendRoot + "/" + filePath
		encoding := negotiateEncoding(c.Request.Header.Get("Accept-Encoding"))

		if encoding == "br" {
			f, err = web.Frontend.Open(openPath + ".br")
			if err == nil {
				c.Header("Content-Encoding", "br")
			}
		}
		if f == nil && encoding == "gzip" {
			f, err = web.Frontend.Open(openPath + ".gz")
			if err == nil {
				c.Header("Content-Encoding", "gzip")
			}
		}
		if f == nil {
			f, err = web.Frontend.Open(openPath)
			if err != nil {
				f, err = web.Frontend.Open(frontendRoot + "/" + indexFile)
				if err != nil {
					c.Status(http.StatusNotFound)
					return
				}
				filePath = indexFile
			}
		}
		defer f.Close()

		c.Header("Vary", "Accept-Encoding")

		if strings.HasPrefix(filePath, assetsPrefix) {
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			c.Header("Cache-Control", "no-cache")
		}

		if ctype := mime.TypeByExtension(path.Ext(filePath)); ctype != "" {
			c.Header("Content-Type", ctype)
		}

		if rs, ok := f.(io.ReadSeeker); ok {
			http.ServeContent(c.Writer, c.Request, path.Base(filePath), appStartTime, rs)
		} else {
			io.Copy(c.Writer, f)
		}
	})
}

// negotiateEncoding parses Accept-Encoding and returns the best compression
// we can offer: "br", "gzip", or "" (no compression).
func negotiateEncoding(header string) string {
	if header == "" {
		return ""
	}

	type pref struct {
		name    string
		quality float64
	}
	var prefs []pref

	for _, field := range strings.Split(header, ",") {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		q := 1.0
		name := field

		if idx := strings.Index(field, ";"); idx != -1 {
			name = strings.TrimSpace(field[:idx])
			// find q-value
			if qi := strings.Index(field[idx:], "q="); qi != -1 {
				if parsed, err := strconv.ParseFloat(strings.TrimSpace(field[idx+qi+2:]), 64); err == nil {
					q = parsed
				}
			}
		}

		// identity;q=0 disables content-encoding
		if name == "identity" && q == 0 {
			return ""
		}
		if q > 0 && (name == "br" || name == "gzip" || name == "*") {
			prefs = append(prefs, pref{name, q})
		}
	}

	bestQ := 0.0
	best := ""
	for _, p := range prefs {
		if p.quality > bestQ {
			if p.name == "*" {
				best = "gzip"
			} else {
				best = p.name
			}
			bestQ = p.quality
		} else if p.quality == bestQ && p.name == "br" && best != "br" {
			best = "br"
		}
	}
	return best
}
