package rexJwts

import "github.com/golang-jwt/jwt/v4"

func GetJwtToken(secretKey string, iat, seconds int64, sub string, payload string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["sub"] = sub
	claims["iat"] = iat
	claims["payload"] = payload
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
