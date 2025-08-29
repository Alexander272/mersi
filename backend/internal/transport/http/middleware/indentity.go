package middleware

import (
	"net/http"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (m *Middleware) VerifyToken(c *gin.Context) {
	token := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", 1)

	// TODO надо попробовать забирать из keycloak ключи и проверять токен здесь
	result, err := m.keycloak.Client.RetrospectToken(c, token, m.keycloak.ClientId, m.keycloak.ClientSecret, m.keycloak.Realm)
	if err != nil {
		domain := m.auth.Domain
		if !strings.Contains(c.Request.Host, domain) {
			domain = c.Request.Host
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie(constants.AuthCookie, "", -1, "/", domain, m.auth.Secure, true)
		response.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "сессия не найдена")
		return
	}

	if !*result.Active {
		response.NewErrorResponse(c, http.StatusUnauthorized, "token is not active", "время сессии истекло, повторите вход")
		return
	}

	user, err := m.services.Session.DecodeAccessToken(c, token)
	if err != nil {
		response.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "токен доступа не валиден")
		return
	}

	realm := c.GetHeader("realm")
	err = uuid.Validate(realm)
	if err == nil {
		identity, err := c.Cookie(constants.IdentityCookie)
		if err != nil || identity == "" {
			response.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "сессия не найдена")
			return
		}
		id := &models.Identity{}
		err = id.Parse(identity)
		if err != nil {
			response.NewErrorResponse(c, http.StatusUnauthorized, err.Error(), "сессия не найдена")
			return
		}

		role := ""
		for _, item := range id.Roles {
			if item.RealmId == realm {
				role = item.Name
			}
		}
		user.Role = role
	}

	c.Set(constants.CtxUser, *user)
	c.Next()
}
