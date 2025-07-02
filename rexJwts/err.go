package rexJwts

import "errors"

var (
	ErrorTokenInvalid              = errors.New("invalid token")
	ErrorTokenInvalidSignature     = errors.New("invalid signature")
	ErrorJwtClaimsInvalid          = errors.New("invalid jwt claims")
	ErrorTokenInvalidAudience      = errors.New("invalid audience")
	ErrorTokenHasExpired           = errors.New("token has expired")
	ErrorTokenInvalidIssuer        = errors.New("invalid issuer")
	ErrorTokenNotActiveYet         = errors.New("token not active yet")
	ErrorInvalidPublicKeyPEMFormat = errors.New("invalid public key PEM format")
	ErrorPublicKeyNotECDSA         = errors.New("public key not ECDSA type")
)
