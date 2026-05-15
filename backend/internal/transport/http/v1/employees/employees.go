package employees

import (
	"errors"
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/pkg/error_bot"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EmployeeHandlers struct {
	service services.Employee
}

func NewEmployeeHandlers(service services.Employee) *EmployeeHandlers {
	return &EmployeeHandlers{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Employee, middleware *middleware.Middleware) {
	handlers := NewEmployeeHandlers(service)

	employees := api.Group("/employees", middleware.CheckPermissions(constants.Employee, constants.Read))
	{
		employees.GET("", handlers.GetAll)
		employees.GET("/unique", handlers.GetUnique)
		employees.GET("/:id", handlers.GetById)

		write := employees.Group("", middleware.CheckPermissions(constants.Employee, constants.Write))
		{
			write.POST("", handlers.Create)
			write.PUT("/:id", handlers.Update)
			write.DELETE("/:id", handlers.Delete)
		}
	}
}

func (h *EmployeeHandlers) GetAll(c *gin.Context) {
	filter := make(map[string]string, 0)

	realm := c.GetHeader("realm")
	err := uuid.Validate(realm)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "invalid id param")
		return
	}
	filter["realm_id"] = realm

	departmentId := c.Query("departmentId")
	if departmentId != "" {
		filter["department_id"] = departmentId
	}

	employees, err := h.service.GetAll(c, &models.GetEmployeesDTO{Filters: filter})
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), filter)
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: employees})
}

func (h *EmployeeHandlers) GetUnique(c *gin.Context) {
	realm := c.Query("realm")
	if err := uuid.Validate(realm); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "invalid realm param")
		return
	}
	dto := &models.GetUniqueEmployeeDTO{Realm: realm}

	employees, err := h.service.GetUnique(c, dto)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: employees, Total: len(employees)})
}

// func (h *EmployeeHandlers) GetByDepartment(c *gin.Context) {
// 	departmentId := c.Param("departmentId")
// 	if departmentId == "" {
// 		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Id подразделения не задан")
// 		return
// 	}

// 	users, err := h.service.GetByDepartment(c, departmentId)
// 	if err != nil {
// 		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
// 		error_bot.Send(c, err.Error(), departmentId)
// 		return
// 	}

// 	c.JSON(http.StatusOK, response.DataResponse{Data: users})
// }

func (h *EmployeeHandlers) GetById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Id пользователя не задан")
		return
	}

	data, err := h.service.GetById(c, id)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			response.NewErrorResponse(c, http.StatusNotFound, err.Error(), "Сотрудник не найден")
			return
		}
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), id)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *EmployeeHandlers) Create(c *gin.Context) {
	dto := &models.WriteEmployeeDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	if err := h.service.Create(c, dto); err != nil {
		if errors.Is(err, models.ErrAlreadyExists) {
			response.NewErrorResponse(c, http.StatusConflict, err.Error(), "Сотрудник с таким именем уже существует в данном подразделении")
			return
		}
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}

	logger.Info("Работник создан",
		logger.StringAttr("id", dto.Id),
		logger.StringAttr("name", dto.Name),
		logger.AnyAttr("data", dto),
	)

	c.JSON(http.StatusCreated, response.IdResponse{Message: "Пользователь создан"})
}

func (h *EmployeeHandlers) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Id пользователя не задан")
		return
	}

	dto := &models.WriteEmployeeDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}
	dto.Id = id

	if err := h.service.Update(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}

	logger.Info("Работник обновлен",
		logger.StringAttr("id", dto.Id),
		logger.StringAttr("name", dto.Name),
		logger.AnyAttr("data", dto),
	)

	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные пользователя обновлены"})
}

func (h *EmployeeHandlers) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Id пользователя не задан")
		return
	}

	if err := h.service.Delete(c, id); err != nil {
		if errors.Is(err, models.ErrHasInstrument) {
			response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Нельзя удалить работника у которого числятся инструменты")
			return
		}

		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), id)
		return
	}
	logger.Info("Работник удален", logger.StringAttr("id", id))

	c.JSON(http.StatusNoContent, response.IdResponse{})
}
