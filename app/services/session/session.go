package session

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mises-id/sns/app/models"
	"github.com/mises-id/sns/config/env"
	"github.com/mises-id/sns/lib/codes"
	"github.com/mises-id/sns/lib/mises"
)

var (
	secret      = env.Envs.JWTSecret
	misesClient mises.Client
)

func init() {
	misesClient = mises.New()
}

func SignIn(ctx context.Context, misesid, auth_code string) (string, error) {
	err := misesClient.Auth(misesid, auth_code)
	if err != nil {
		return "", err
	}
	user, err := models.FindOrCreateUserByMisesid(ctx, misesid)
	if err != nil {
		return "", err
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":      user.UID,
		"misesid":  user.Misesid,
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	})
	return at.SignedString([]byte(secret))
}

func Auth(ctx context.Context, authToken string) (*models.User, error) {
	claim, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		if err.Error() == "Token is expired" {
			return nil, codes.ErrTokenExpired
		}
		return nil, err
	}
	mapClaims := claim.Claims.(jwt.MapClaims)
	return &models.User{
		UID:      uint64(mapClaims["uid"].(float64)),
		Misesid:  mapClaims["misesid"].(string),
		Username: mapClaims["username"].(string),
	}, nil
}
