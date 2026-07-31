package middleware

import (
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
			response.SendError(c, models.ErrSessionEmpty)
			return
		}
		user, ok := u.(models.User)
		if !ok {
			response.SendError(c, models.ErrInvalidUserType)
			return
		}

		access, err := m.services.Permission.Enforce(user.ID, realm, menuItem, method)
		if err != nil {
			response.SendError(c, err)
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
			response.SendError(c, models.ErrForbidden)
			return
		}

		c.Next()
	}
}

func (m *Middleware) CheckPermissionsArray(perm []*Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		realm := c.GetHeader("realm")
		u, exists := c.Get(constants.CtxUser)
		if !exists {
			response.SendError(c, models.ErrSessionEmpty)
			return
		}

		user, ok := u.(models.User)
		if !ok {
			response.SendError(c, models.ErrInvalidUserType)
			return
		}
		access := false
		for _, item := range perm {
			a, err := m.services.Permission.Enforce(user.ID, realm, item.Section, item.Method)
			if err != nil {
				response.SendError(c, err)
				return
			}
			logger.Debug("permissions",
				logger.StringAttr("user", user.ID),
				logger.StringAttr("realm", realm),
				logger.StringAttr("section", item.Section),
				logger.StringAttr("method", item.Method),
				logger.BoolAttr("access", a),
			)

			if a {
				access = true
				break
			}
		}

		if !access {
			response.SendError(c, models.ErrForbidden)
			return
		}

		c.Next()
	}
}

func (m *Middleware) DepartmentAccess(c *gin.Context) {
	u, exists := c.Get(constants.CtxUser)
	if !exists {
		c.Set(constants.CtxDepartmentAccess, []string{})
		c.Set(constants.CtxHasWritePermission, false)
		c.Next()
		return
	}
	user, ok := u.(models.User)
	if !ok {
		c.Set(constants.CtxDepartmentAccess, []string{})
		c.Set(constants.CtxHasWritePermission, false)
		c.Next()
		return
	}

	realm := c.GetHeader("realm")

	access, err := m.services.Permission.Enforce(user.ID, realm, constants.SI, constants.Write)
	if err != nil {
		logger.Error("failed to check write permission", logger.ErrAttr(err))
	}
	if err == nil && access {
		c.Set(constants.CtxHasWritePermission, true)
	} else {
		c.Set(constants.CtxHasWritePermission, false)
	}

	deptAccess, err := m.services.DepartmentAccess.GetByUserId(c, &models.GetDepartmentAccessDTO{UserId: user.ID, RealmId: realm})
	if err != nil {
		c.Set(constants.CtxDepartmentAccess, []string{})
		c.Next()
		return
	}

	ids := make([]string, len(deptAccess))
	for i, da := range deptAccess {
		ids[i] = da.DepartmentId
	}
	c.Set(constants.CtxDepartmentAccess, ids)
	c.Next()
}
