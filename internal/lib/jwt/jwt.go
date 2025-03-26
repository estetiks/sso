package jwt_sso

import (
	"time"

	"github.com/estetiks/sso/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, app models.App, tokenTTL time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["tokenTTL"] = time.Now().Add(tokenTTL).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(app.Secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
