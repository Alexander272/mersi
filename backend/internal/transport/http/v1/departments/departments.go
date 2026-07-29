package departments

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DepartmentHandlers struct {
	service services.Department
}

func NewDepartmentHandlers(service services.Department) *DepartmentHandlers {
	return &DepartmentHandlers{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Department, middleware *middleware.Middleware) {
	handlers := NewDepartmentHandlers(service)

	departments := api.Group("/departments", middleware.CheckPermissions(constants.Department, constants.Read))
	{
		departments.GET("", handlers.GetAll)
		departments.GET("/:id", handlers.GetById)
		departments.GET("/sso", handlers.GetBySSOId)

		write := departments.Group("", middleware.CheckPermissions(constants.Department, constants.Write))
		{
			write.POST("", handlers.Create)
			write.PUT("/:id", handlers.Update)
			write.DELETE("/:id", handlers.Delete)
		}
	}
}

func (h *DepartmentHandlers) GetAll(c *gin.Context) {
	realm := c.GetHeader("realm")
	err := uuid.Validate(realm)
	if err != nil {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	dto := &models.GetDepartmentsDTO{RealmId: realm}

	departments, err := h.service.GetAll(c, dto)
	if err != nil {
		response.SendError(c, err, dto)
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: departments})
}

func (h *DepartmentHandlers) GetById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendError(c, models.ErrInvalidInput)
		return
	}
	dto := &models.GetDepartmentByIdDTO{Id: id}

	department, err := h.service.GetById(c, dto)
	if err != nil {
		response.SendError(c, err, dto)
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: department})
}

func (h *DepartmentHandlers) GetBySSOId(c *gin.Context) {
	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	user, ok := u.(models.User)
	if !ok {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}

	departments, err := h.service.GetBySSOId(c, user.ID)
	if err != nil {
		response.SendError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: departments})
}

func (h *DepartmentHandlers) Create(c *gin.Context) {
	dto := &models.DepartmentDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}
	realm := c.GetHeader("realm")
	err := uuid.Validate(realm)
	if err != nil {
		response.SendError(c, models.ErrInvalidInput)
		return
	}
	dto.RealmId = realm

	id, err := h.service.Create(c, dto)
	if err != nil {
		response.SendError(c, err, dto)
		return
	}
	logger.Info("Подразделение создано",
		logger.StringAttr("id", dto.Id),
		logger.StringAttr("name", dto.Name),
		logger.AnyAttr("data", dto),
	)

	c.JSON(http.StatusCreated, response.IdResponse{Id: id, Message: "Новое подразделение создано"})
}

func (h *DepartmentHandlers) Update(c *gin.Context) {
	dto := &models.DepartmentDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}
	id := c.Param("id")
	if id == "" {
		response.SendError(c, models.ErrInvalidInput)
		return
	}
	dto.Id = id

	if err := h.service.Update(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	logger.Info("Подразделение обновлено",
		logger.StringAttr("id", dto.Id),
		logger.StringAttr("name", dto.Name),
		logger.AnyAttr("data", dto),
	)

	c.JSON(http.StatusOK, response.IdResponse{Message: "Подразделение обновлено"})
}

func (h *DepartmentHandlers) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	if err := h.service.Delete(c, id); err != nil {
		response.SendError(c, err, id)
		return
	}
	logger.Info("Подразделение удалено", logger.StringAttr("id", id))

	c.JSON(http.StatusNoContent, response.StatusResponse{})
}
