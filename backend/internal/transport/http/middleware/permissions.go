package middleware

import (
	"github.com/gin-gonic/gin"
)

type Permission struct {
	Section string
	Method  string
}

func (m *Middleware) CheckPermissions(menuItem, method string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// realm := c.GetHeader("realm")
		// err := uuid.Validate(realm)
		// if err != nil {
		// 	response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "invalid id param")
		// 	return
		// }

		// identity, err := c.Cookie(constants.IdentityCookie)
		// if err != nil || identity == "" {
		// 	response.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "сессия не найдена")
		// 	return
		// }
		// id := &models.Identity{}
		// err = id.Parse(identity)
		// if err != nil {
		// 	response.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "сессия не найдена")
		// 	return
		// }

		// role := ""
		// for _, item := range id.Roles {
		// 	if item.RealmId == realm {
		// 		role = item.Name
		// 	}
		// }

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
		// u, exists := c.Get(constants.CtxUser)
		// if !exists {
		// 	response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "сессия не найдена")
		// 	return
		// }

		// user := u.(models.User)
		// access := false
		// for _, item := range perm {
		// 	a, err := m.services.Permission.Enforce(user.Role, item.Section, item.Method)
		// 	if err != nil {
		// 		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		// 		return
		// 	}
		// 	logger.Debug("permissions", logger.StringAttr("section", item.Section), logger.StringAttr("method", item.Method), logger.BoolAttr("access", a))

		// 	if a {
		// 		access = true
		// 		break
		// 	}
		// }

		// if !access {
		// 	response.NewErrorResponse(c, http.StatusForbidden, "access denied", "нет доступа к данному разделу")
		// 	return
		// }

		c.Next()
	}
}
