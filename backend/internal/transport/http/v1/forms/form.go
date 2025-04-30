package forms

import (
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/forms/create"
	"github.com/gin-gonic/gin"
)

func Register(api *gin.RouterGroup, services *services.Services, middleware *middleware.Middleware) {
	forms := api.Group("forms")
	create.Register(forms, services.CreateForm, middleware)
}
