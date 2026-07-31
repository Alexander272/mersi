package services

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/pkg/auth"
	"github.com/golang-jwt/jwt/v5"
)

type SessionService struct {
	keycloak *auth.KeycloakClient
	user     User
}

func NewSessionService(keycloak *auth.KeycloakClient, user User) *SessionService {
	return &SessionService{
		keycloak: keycloak,
		user:     user,
	}
}

type Session interface {
	SignIn(ctx context.Context, u *models.SignIn) (*models.User, error)
	SignOut(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, req *models.RefreshDTO) (*models.User, error)
	DecodeAccessToken(ctx context.Context, token string) (*models.User, error)
}

func (s *SessionService) SignIn(ctx context.Context, u *models.SignIn) (*models.User, error) {
	res, err := s.keycloak.Client.Login(ctx, s.keycloak.ClientId, s.keycloak.ClientSecret, s.keycloak.Realm, u.Username, u.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to login to keycloak. error: %w", err)
	}

	decodedUser, err := s.DecodeAccessToken(ctx, res.AccessToken)
	if err != nil {
		return nil, err
	}

	user, err := s.user.GetRoles(ctx, &models.GetUserInfoDTO{UserID: decodedUser.ID, Realm: u.Realm})
	if err != nil {
		return nil, err
	}

	user.Name = decodedUser.Name
	user.AccessToken = res.AccessToken
	user.RefreshToken = res.RefreshToken

	return user, nil
}

func (s *SessionService) SignOut(ctx context.Context, refreshToken string) error {
	err := s.keycloak.Client.Logout(ctx, s.keycloak.ClientId, s.keycloak.ClientSecret, s.keycloak.Realm, refreshToken)
	if err != nil {
		return fmt.Errorf("failed to logout to keycloak. error: %w", err)
	}
	return nil
}

func (s *SessionService) Refresh(ctx context.Context, req *models.RefreshDTO) (*models.User, error) {
	res, err := s.keycloak.Client.RefreshToken(ctx, req.Token, s.keycloak.ClientId, s.keycloak.ClientSecret, s.keycloak.Realm)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token in keycloak. error: %w", err)
	}

	decodedUser, err := s.DecodeAccessToken(ctx, res.AccessToken)
	if err != nil {
		return nil, err
	}

	user, err := s.user.GetRoles(ctx, &models.GetUserInfoDTO{UserID: decodedUser.ID, Realm: req.Realm})
	if err != nil {
		return nil, err
	}

	user.Name = decodedUser.Name
	user.AccessToken = res.AccessToken
	user.RefreshToken = res.RefreshToken

	return user, nil
}

func (s *SessionService) DecodeAccessToken(ctx context.Context, token string) (*models.User, error) {
	_, claims, err := s.keycloak.Client.DecodeAccessToken(ctx, token, s.keycloak.Realm)
	if err != nil {
		return nil, fmt.Errorf("failed to decode access token. error: %w", err)
	}

	return claimsToUser(*claims, os.Getenv("SERVICE_ID"))
}

func claimsToUser(claims jwt.MapClaims, serviceName string) (*models.User, error) {
	user := &models.User{}
	var role, username, userId string
	c := claims

	if ra, ok := c["realm_access"]; ok {
		if realmAccess, ok := ra.(map[string]interface{}); ok {
			if r, ok := realmAccess["roles"]; ok {
				if roles, ok := r.([]interface{}); ok {
					for _, rr := range roles {
						if roleStr, ok := rr.(string); ok && strings.Contains(roleStr, serviceName) {
							role = strings.Replace(roleStr, serviceName+"_", "", 1)
							break
						}
					}
				}
			}
		}
	}

	if u, ok := c["preferred_username"]; ok {
		if s, ok := u.(string); ok {
			username = s
		}
	}

	if uId, ok := c["sub"]; ok {
		if s, ok := uId.(string); ok {
			userId = s
		}
	}

	user.ID = userId
	user.Role = role
	user.Name = username

	return user, nil
}
