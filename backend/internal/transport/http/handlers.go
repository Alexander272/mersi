package http

import (
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"runtime/debug"
	"strings"
	"time"

	"github.com/Alexander272/mersi/backend/internal/config"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	httpV1 "github.com/Alexander272/mersi/backend/internal/transport/http/v1"
	"github.com/Alexander272/mersi/backend/pkg/accept_encoding"
	"github.com/Alexander272/mersi/backend/pkg/auth"
	"github.com/Alexander272/mersi/backend/pkg/limiter"
	"github.com/Alexander272/mersi/backend/pkg/logger"
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
	router := gin.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	router.Use(
		gin.LoggerWithConfig(gin.LoggerConfig{
			Skip: func(c *gin.Context) bool {
				path := c.Request.URL.Path
				if strings.HasPrefix(path, "/api") {
					return false
				}
				return c.Writer.Status() < http.StatusBadRequest
			},
		}),
		gin.CustomRecovery(h.ErrorHandler),
		securityHeaders(),
	)

	if err := router.SetTrustedProxies(conf.Http.TrustedProxies); err != nil {
		logger.Warn("invalid trusted proxies config", logger.ErrAttr(err))
	}

	h.initAPI(router, conf)
	h.initStatic(router, conf)

	return router
}

func (h *Handler) ErrorHandler(c *gin.Context, origErr any) {
	err := fmt.Errorf("unexpected error: %v", origErr)

	rawStack := string(debug.Stack())                        // 1. Получаем стек в виде байтов
	cleanStack := strings.ReplaceAll(rawStack, "\t", "    ") // 2. Заменяем все табуляции на 4 пробела для красоты
	stackLines := strings.Split(cleanStack, "\n")            // 3. Превращаем в срез строк, разделяя по символу \n

	response.SendError(c, err, gin.H{"PANIC": true, "Stack trace": stackLines})
	debug.PrintStack()
}

func securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy",
			"default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; "+
				"img-src 'self' data: blob:; font-src 'self'; connect-src 'self' ws: wss:; "+
				"frame-ancestors 'none'; base-uri 'self'; form-action 'self';")
		c.Header("Permissions-Policy",
			"camera=(), microphone=(), geolocation=(), gyroscope=(), "+
				"accelerometer=(), magnetometer=(), usb=(), payment=(), "+
				"display-capture=(), document-domain=()")
		c.Next()
	}
}

func (h *Handler) initAPI(router *gin.Engine, conf *config.Config) {
	middleware := middleware.NewMiddleware(h.services, conf.Auth, h.keycloak)
	handler := httpV1.NewHandler(httpV1.Deps{Services: h.services, Conf: conf, Middleware: middleware})

	api := router.Group("/api")
	api.Use(limiter.Limit(conf.ApiLimiter.RPS, conf.ApiLimiter.Burst, conf.ApiLimiter.TTL))
	handler.Init(api)

	api.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}

var appStartTime = time.Now()

const (
	frontendRoot = "frontend"
	indexFile    = "index.html"
	assetsPrefix = "assets/"
)

var allowedStaticExts = map[string]bool{
	".html": true, ".js": true, ".css": true, ".png": true, ".jpg": true,
	".jpeg": true, ".svg": true, ".gif": true, ".ico": true, ".webp": true,
	".woff": true, ".woff2": true, ".ttf": true, ".eot": true, ".map": true,
}

func (h *Handler) initStatic(router *gin.Engine, conf *config.Config) {
	router.NoRoute(limiter.Limit(conf.StaticLimiter.RPS, conf.StaticLimiter.Burst, conf.StaticLimiter.TTL), func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.Status(http.StatusNotFound)
			return
		}

		filePath := strings.TrimPrefix(c.Request.URL.Path, "/")
		if filePath == "" {
			filePath = indexFile
		}
		filePath = path.Clean(filePath)

		// 🔒 Блокируем скрытые файлы/директории (начинаются с точки)
		if strings.HasPrefix(filePath, ".") || strings.Contains(filePath, "/.") {
			c.Status(http.StatusNotFound)
			return
		}

		if ext := path.Ext(filePath); filePath != indexFile && ext != "" && !allowedStaticExts[ext] {
			c.Status(http.StatusNotFound)
			return
		}

		var f fs.File
		var err error
		openPath := frontendRoot + "/" + filePath
		encoding := accept_encoding.Negotiate(c.Request.Header.Get("Accept-Encoding"))

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
