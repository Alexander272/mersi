package utils

import (
	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/gin-gonic/gin"
)

func GetActor(c *gin.Context) *models.Actor {
	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.SendError(c, models.ErrSessionEmpty)
		return nil
	}
	user, ok := u.(models.User)
	if !ok {
		response.SendError(c, models.ErrInvalidUserType)
		return nil
	}

	actor := &models.Actor{
		ID:   user.ID,
		Name: user.Name,
	}
	return actor
}

func GetUser(c *gin.Context) *models.User {
	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.SendError(c, models.ErrSessionEmpty)
		return nil
	}
	user, ok := u.(models.User)
	if !ok {
		response.SendError(c, models.ErrInvalidUserType)
		return nil
	}
	return &user
}
