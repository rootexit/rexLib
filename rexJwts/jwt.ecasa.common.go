package rexJwts

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rootexit/rexLib/rexCrypto"
	"time"
)

func JwtECDSACommonCreateToken(claims *jwt.StandardClaims, privateKey string) (string, int64, error) {
	key, err := rexCrypto.ParseECDSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", -1, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", -1, err
	}
	return tokenString, claims.ExpiresAt, nil
}

func JwtECDSACommonParseAndVerifyToken(tokenString, certPem string) (*jwt.StandardClaims, error) {
	key, err := rexCrypto.ParseECDSAPublicKeyFromCert(certPem)
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		// 验证过期时间
		if time.Now().Unix() > claims.ExpiresAt {
			fmt.Println("Token has expired")
			return nil, ErrorTokenHasExpired
		}
		// 验证开始时间
		if time.Now().Unix() < claims.NotBefore {
			fmt.Println("Token not active yet")
			return nil, ErrorTokenNotActiveYet
		}
		return claims, nil
	} else {
		return nil, err
	}
}

func JwtECDSACommonParse(tokenString, certPem string) (*jwt.Token, error) {
	key, err := rexCrypto.ParseECDSAPublicKeyFromCert(certPem)
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func JwtECDSACommonVerify(token *jwt.Token, Audience string) (*jwt.StandardClaims, error) {
	var err error
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		// 验证受众
		if claims.Audience != Audience {
			fmt.Println("Invalid audience")
			return nil, ErrorTokenInvalidAudience
		}

		// 验证过期时间
		if time.Now().Unix() > claims.ExpiresAt {
			fmt.Println("Token has expired")
			return nil, ErrorTokenHasExpired
		}

		// 验证开始时间
		if time.Now().Unix() < claims.NotBefore {
			fmt.Println("Token not active yet")
			return nil, ErrorTokenNotActiveYet
		}

		// 验证tokenId和Subject, 无法验证, 返回标准结构体
		return claims, nil
	} else {
		return nil, err
	}
}
