package middleware

import (
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Permission struct {
	Section string
	Method  string
}

func (m *Middleware) CheckPermissions(menuItem, method string) gin.HandlerFunc {
	return func(c *gin.Context) {
		realm := c.GetHeader("realm")
		u, exists := c.Get(constants.CtxUser)
		if !exists {
			response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "сессия не найдена")
			return
		}
		user := u.(models.User)

		access, err := m.services.Permission.Enforce(user.ID, realm, menuItem, method)
		if err != nil {
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
			return
		}
		logger.Debug("permissions",
			logger.StringAttr("user", user.ID),
			logger.StringAttr("realm", realm),
			logger.StringAttr("menu", menuItem),
			logger.StringAttr("method", method),
			logger.BoolAttr("access", access),
		)

		if !access {
			response.NewErrorResponse(c, http.StatusForbidden, "access denied", "нет доступа к данному разделу")
			return
		}

		//* old version

		// u, exists := c.Get(constants.CtxUser)
		// if !exists {
		// 	response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "сессия не найдена")
		// 	return
		// }
		// role := u.(models.User).Role

		// access, err := m.services.Permission.Enforce(role, menuItem, method)
		// if err != nil {
		// 	response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		// 	return
		// }
		// logger.Debug("permissions", logger.StringAttr("menu", menuItem), logger.StringAttr("method", method), logger.BoolAttr("access", access))

		// if !access {
		// 	response.NewErrorResponse(c, http.StatusForbidden, "access denied", "нет доступа к данному разделу")
		// 	return
		// }

		c.Next()
	}
}

func (m *Middleware) CheckPermissionsArray(perm []*Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		realm := c.GetHeader("realm")
		u, exists := c.Get(constants.CtxUser)
		if !exists {
			response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "сессия не найдена")
			return
		}

		user := u.(models.User)
		access := false
		for _, item := range perm {
			a, err := m.services.Permission.Enforce(user.ID, realm, item.Section, item.Method)
			if err != nil {
				response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
				return
			}
			// a, err := m.services.Permission.Enforce(user.Role, item.Section, item.Method)
			// if err != nil {
			// 	response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
			// 	return
			// }
			logger.Debug("permissions", logger.StringAttr("section", item.Section), logger.StringAttr("method", item.Method), logger.BoolAttr("access", a))

			if a {
				access = true
				break
			}
		}

		if !access {
			response.NewErrorResponse(c, http.StatusForbidden, "access denied", "нет доступа к данному разделу")
			return
		}

		c.Next()
	}
}
