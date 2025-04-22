package sections

import (
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/pkg/error_bot"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service services.Section
}

func NewHandler(service services.Section) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Section, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	sections := api.Group("/sections", middleware.CheckPermissions(constants.Sections, constants.Read))
	{
		sections.GET("", handler.Get)

		write := sections.Group("", middleware.CheckPermissions(constants.Sections, constants.Write))
		{
			write.POST("", handler.create)
			write.PUT("/:id", handler.update)
			write.DELETE("/:id", handler.delete)
		}
	}
}

func (h *Handler) Get(c *gin.Context) {
	realm := c.Query("realm")
	if uuid.Validate(realm) != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Отправлены некорректные данные")
		return
	}
	dto := &models.GetSectionsDTO{RealmID: realm}

	data, err := h.service.Get(c, dto)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) create(c *gin.Context) {}

func (h *Handler) update(c *gin.Context) {}

func (h *Handler) delete(c *gin.Context) {}
