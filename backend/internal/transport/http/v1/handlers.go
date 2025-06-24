package v1

import (
	"github.com/Alexander272/mersi/backend/internal/config"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/auth"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/columns"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/forms"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/realm"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/sections"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/si"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services   *services.Services
	conf       *config.Config
	middleware *middleware.Middleware
}

type Deps struct {
	Services   *services.Services
	Conf       *config.Config
	Middleware *middleware.Middleware
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		services:   deps.Services,
		conf:       deps.Conf,
		middleware: deps.Middleware,
	}
}

func (h *Handler) Init(group *gin.RouterGroup) {
	v1 := group.Group("/v1")

	auth.Register(v1, auth.Deps{Service: h.services.Session, Middleware: h.middleware, Auth: h.conf.Auth})

	secure := v1.Group("", h.middleware.VerifyToken)

	realm.Register(secure, h.services.Realm, h.conf.Auth, h.middleware)
	// accesses.Register(secure, h.services.Accesses, h.middleware)
	sections.Register(secure, h.services.Section, h.middleware)
	columns.Register(secure, h.services.Columns, h.middleware)
	forms.Register(secure, h.services, h.middleware)
	si.Register(secure, h.services, h.middleware)
	// instruments.Register(secure, h.services.Instrument, h.middleware)
	// documents.Register(secure, h.services.Document, h.middleware)
}
