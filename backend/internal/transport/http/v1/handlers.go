package v1

import (
	"github.com/Alexander272/mersi/backend/internal/config"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/accesses"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/auth"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/columns"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/context_menu"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/forms"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/preservation"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/realm"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/repair"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/sections"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/si"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/tools_menu"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/transfer_to_dep"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/transfer_to_save"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/write_off"
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
	accesses.Register(secure, h.services.Accesses, h.middleware)
	sections.Register(secure, h.services.Section, h.middleware)
	columns.Register(secure, h.services.Columns, h.middleware)
	forms.Register(secure, h.services, h.middleware)
	si.Register(secure, h.services, h.middleware)
	context_menu.Register(secure, h.services.ContextMenu, h.middleware)
	tools_menu.Register(secure, h.services.ToolsMenu, h.middleware)
	repair.Register(secure, h.services.Repair, h.middleware)
	preservation.Register(secure, h.services.Preservation, h.middleware)
	transfer_to_save.Register(secure, h.services.TransferToSave, h.middleware)
	transfer_to_dep.Register(secure, h.services.TransferToDepartment, h.middleware)
	write_off.Register(secure, h.services.WriteOff, h.middleware)
}
