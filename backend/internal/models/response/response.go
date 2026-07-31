package response

import (
	"errors"
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/pkg/error_bot"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

type DataResponse struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total,omitempty"`
}

type IdResponse struct {
	Id      string `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

func SendError(c *gin.Context, err error, request ...any) {
	var status int
	var code string
	var message string

	var httpErr models.HTTPError
	if errors.As(err, &httpErr) {
		status = httpErr.Status()
		code = httpErr.Code()
		message = httpErr.Message()
	} else {
		status = http.StatusInternalServerError
		code = "U001"
		message = "Внутренняя ошибка сервера"
	}

	loggerValues := []any{
		logger.StringAttr("url", c.Request.URL.Path),
		logger.StringAttr("method", c.Request.Method),
		logger.StringAttr("ip", c.ClientIP()),
		logger.StringAttr("error", err.Error()),
		logger.StringAttr("code", code),
	}

	if status >= 500 {
		logger.Error("request_failed", loggerValues...)
		error_bot.Send(c, err.Error(), extractRequest(request))
	} else {
		logger.Info("request_failed", loggerValues...)
	}

	c.AbortWithStatusJSON(status, ErrorResponse{Message: message, Code: code})
}

func extractRequest(req []any) any {
	if len(req) > 0 {
		return req[0]
	}
	return nil
}
