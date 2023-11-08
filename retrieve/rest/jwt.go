package rest

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtClaims struct {
	Username string `json:"userame"`
	jwt.StandardClaims
}

const (
	secret        = "[f^(F1OhYzEL58Gi8!J5S]YnU:BfWsM4#WQ/Vsf9r3+*g!Qh60LeM#pw5O+5~~4"
	sessionLength = 12 * time.Hour
)

func CreateJwt(username string) (string, error) {
	claims := JwtClaims{
		username,
		jwt.StandardClaims{
			//Id::"main_user_id",
			ExpiresAt: time.Now().Add(sessionLength).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err := rawToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}
