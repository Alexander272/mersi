package si

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/si/documents"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/si/instruments"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/si/verifications"
	"github.com/Alexander272/mersi/backend/pkg/error_bot"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service services.SI
}

func NewHandler(service services.SI) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, services *services.Services, middleware *middleware.Middleware) {
	handler := NewHandler(services.SI)

	si := api.Group("si", middleware.CheckPermissions(constants.SI, constants.Read))
	{
		si.GET("", handler.get)
		si.GET("/:id", handler.getById)

		write := si.Group("", middleware.CheckPermissions(constants.SI, constants.Write))
		{
			write.POST("", handler.create)
			write.PUT("/:id", handler.update)
			write.PUT("/position", handler.changePosition)
			write.DELETE("/:id", handler.delete)
		}
	}

	instruments.Register(si, services.Instrument, middleware)
	documents.Register(si, services.Document, middleware)
	verifications.Register(si, services.Verification, middleware)
}

func (h *Handler) get(c *gin.Context) {
	section := c.Query("section")
	err := uuid.Validate(section)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Сессия не найдена")
		return
	}

	params := &models.GetSiDTO{
		SectionId: section,
		Page:      &models.Page{},
		Sort:      []*models.Sort{},
		Filters:   []*models.Filter{},
	}

	page := c.Query("page")
	size := c.Query("size")

	all := c.Query("all")

	sortLine := c.Query("sort_by")
	filters := c.QueryMap("filters")

	limit, err := strconv.Atoi(size)
	if err != nil {
		params.Page.Limit = 15
	} else {
		params.Page.Limit = limit
	}

	p, err := strconv.Atoi(page)
	if err != nil {
		params.Page.Offset = 0
	} else {
		params.Page.Offset = (p - 1) * params.Page.Limit
	}

	if sortLine != "" {
		sort := strings.Split(sortLine, ",")
		for _, v := range sort {
			field, found := strings.CutPrefix(v, "-")
			t := "ASC"
			if found {
				t = "DESC"
			}

			params.Sort = append(params.Sort, &models.Sort{
				Field: field,
				Type:  t,
			})
		}
	}

	// можно сделать массив с именами полей, а потом передавать для каждого поля значение фильтра, например
	// filter[0]=nextVerificationDate&nextVerificationDate[lte]=somevalue&nextVerificationDate[qte]=somevalue&filter[1]=name&name[eq]=somevalue
	// qte - >=; lte - <=
	// нужен еще тип как-то передать
	// как вариант можно передать filter[nextVerificationDate]=date, filter[name]=string
	// только надо проверить как это все будет читаться на сервере и записываться на клиенте

	// можно сделать следующие варианты compareType (это избавит от необходимости знать тип поля)
	// number or date: eq, qte, lte
	// string: like, con, start, end
	// list: in

	for k, v := range filters {
		valueMap := c.QueryMap(k)
		values := []*models.FilterValue{}
		for key, value := range valueMap {
			if k == "place" {
				statusFilter := &models.Filter{Field: "status", FieldType: "list", Values: []*models.FilterValue{{CompareType: "in"}}}
				tmp := []string{}
				if strings.Contains(value, "_reserve") {
					tmp = append(tmp, "reserve")
				}
				if strings.Contains(value, "_moved") {
					tmp = append(tmp, "moved")
				}
				statusFilter.Values[0].Value = strings.Join(tmp, ",")
				params.Filters = append(params.Filters, statusFilter)

				value = strings.Replace(value, "_reserve", "", -1)
				value = strings.Replace(value, "_moved", "", -1)
				value = strings.Trim(value, ",")
				k = "department"

				if value != "" {
					// params.Filters = append(params.Filters, &models.SIFilter{Field: "last_place", Values: []*models.SIFilterValue{}})
					values = append(values, &models.FilterValue{CompareType: key, Value: value})
				}
			}

			values = append(values, &models.FilterValue{
				CompareType: key,
				Value:       value,
			})
		}
		if values[0].Value == "" {
			continue
		}

		f := &models.Filter{
			Field:     k,
			FieldType: v,
			Values:    values,
		}

		params.Filters = append(params.Filters, f)
	}

	if all != "true" {
		params.Filters = append(params.Filters, &models.Filter{
			Field:     "status",
			FieldType: "",
			Values:    []*models.FilterValue{{CompareType: "nlike", Value: "reserve"}},
		}, &models.Filter{ //TODO задавать id подразделения не очень хорошая идея
			Field:     "department",
			FieldType: "list",
			Values:    []*models.FilterValue{{CompareType: "nin", Value: "cc718041-f3da-4490-b647-380297bd3344"}},
		})

	}

	status := c.Query("status")
	if status == "" {
		params.Status = models.InstrumentStatusWork
	} else {
		params.Status = models.InstrumentStatus(status)
	}

	data, err := h.service.Get(c, params)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), params)
		return
	}
	total := 0
	if len(data) > 0 {
		total = data[0].Total
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: total})
}

func (h *Handler) getById(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Id не валиден")
		return
	}
	req := &models.GetSiByIdDTO{Id: id}

	data, err := h.service.GetById(c, req)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) create(c *gin.Context) {
	dto := &models.SiDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "Сессия не найдена")
		return
	}
	user := u.(models.User)
	dto.Instrument.UserId = user.ID

	if err := h.service.Create(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}

	logger.Info("СИ сохранено",
		logger.StringAttr("user_id", user.ID),
		logger.AnyAttr("instrument-dto", dto.Instrument),
		logger.AnyAttr("verification-dto", dto.Verification),
		// logger.AnyAttr("location-dto", dto.Location),
	)
	c.JSON(http.StatusCreated, response.IdResponse{Message: "Данные о си успешно сохранены"})
}

func (h *Handler) update(c *gin.Context) {
	dto := &models.SiDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "Сессия не найдена")
		return
	}
	user := u.(models.User)

	if err := h.service.Update(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}

	logger.Info("СИ обновлено",
		logger.StringAttr("user_id", user.ID),
		logger.AnyAttr("instrument-dto", dto.Instrument),
		logger.AnyAttr("verification-dto", dto.Verification),
	)
	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные о си обновлены"})
}

func (h *Handler) changePosition(c *gin.Context) {
	dto := &models.ChangePositionDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	if err := h.service.ChangePosition(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Номер позиции изменен"})
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Id не валиден")
		return
	}

	response.NewErrorResponse(c, http.StatusNotImplemented, "not implemented", "not implemented")
}
