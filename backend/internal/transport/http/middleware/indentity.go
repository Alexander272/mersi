package middleware

import (
	"strings"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/gin-gonic/gin"
)

func (m *Middleware) VerifyToken(c *gin.Context) {
	token := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", 1)

	user, err := m.services.Session.DecodeAccessToken(c, token)
	if err != nil {
		domain := m.auth.Domain
		if !strings.Contains(c.Request.Host, domain) {
			domain = c.Request.Host
		}

		c.SetCookie(constants.AuthCookie, "", -1, "/", domain, m.auth.Secure, true)
		response.SendError(c, models.ErrSessionEmpty)
		return
	}

	permission, err := c.Cookie(constants.IdentityCookie)
	if err != nil {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}

	c.Set(constants.CtxUser, *user)
	c.Set(constants.IdentityCookie, permission)
	c.Next()
}
