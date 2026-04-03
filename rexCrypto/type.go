package rexCrypto

type KeyType string

const (
	KeyTypeRSA KeyType = "RSA" // note: RSA密钥
	KeyTypeEC  KeyType = "EC"  // note: EC密钥
	KeyTypeOKP KeyType = "OKP"
	KeyTypeOct KeyType = "OCT" // note: 对称密钥, Octet Key
)

type SignatureAlgorithm string

const (
	// ES256 returns an object representing ECDSA signature algorithm using P-256 curve and SHA-256.

	SignatureAlgorithmES256 SignatureAlgorithm = "ES256"

	// ES256K returns an object representing ECDSA signature algorithm using secp256k1 curve and SHA-256.

	SignatureAlgorithmES256K SignatureAlgorithm = "ES256K"

	// ES384 returns an object representing ECDSA signature algorithm using P-384 curve and SHA-384.

	SignatureAlgorithmES384 SignatureAlgorithm = "ES384"

	// ES512 returns an object representing ECDSA signature algorithm using P-521 curve and SHA-512.

	SignatureAlgorithmES512 SignatureAlgorithm = "ES512"

	// EdDSA returns an object representing EdDSA signature algorithms.

	SignatureAlgorithmEdDSA SignatureAlgorithm = "EdDSA"

	// HS256 returns an object representing HMAC signature algorithm using SHA-256.

	SignatureAlgorithmHS256 SignatureAlgorithm = "HS256"

	// HS384 returns an object representing HMAC signature algorithm using SHA-384.

	SignatureAlgorithmHS384 SignatureAlgorithm = "HS384"

	// HS512 returns an object representing HMAC signature algorithm using SHA-512.

	SignatureAlgorithmHS512 SignatureAlgorithm = "HS512"

	// NoSignature returns an object representing the lack of a signature algorithm. Using this value specifies that the content should not be signed, which you should avoid doing.

	SignatureAlgorithmNoSignature SignatureAlgorithm = "none"

	// PS256 returns an object representing RSASSA-PSS signature algorithm using SHA-256 and MGF1-SHA256.

	SignatureAlgorithmPS256 SignatureAlgorithm = "PS256"

	// PS384 returns an object representing RSASSA-PSS signature algorithm using SHA-384 and MGF1-SHA384.

	SignatureAlgorithmPS384 SignatureAlgorithm = "PS384"

	// PS512 returns an object representing RSASSA-PSS signature algorithm using SHA-512 and MGF1-SHA512.

	SignatureAlgorithmPS512 SignatureAlgorithm = "PS512"

	// RS256 returns an object representing RSASSA-PKCS-v1.5 signature algorithm using SHA-256.

	SignatureAlgorithmRS256 SignatureAlgorithm = "RS256"

	// RS384 returns an object representing RSASSA-PKCS-v1.5 signature algorithm using SHA-384.

	SignatureAlgorithmRS384 SignatureAlgorithm = "RS384"

	// RS512 returns an object representing RSASSA-PKCS-v1.5 signature algorithm using SHA-512.

	SignatureAlgorithmRS512 SignatureAlgorithm = "RS512"
)

const (

	// note: RSA-2048, RSA-3072, RSA-4096
	RSA2048 = "RSA-2048"
	RSA3072 = "RSA-3072"
	RSA4096 = "RSA-4096"

	// note: RSA-2048, RSA-3072, RSA-4096
	ECP224 = "EC-P224"
	ECP256 = "EC-P256"
	ECP384 = "EC-P384"
	ECP521 = "EC-P521"

	// note: RSA-SHA-256
	RsaSha256 = "RSA-SHA-256"
	// note: RSA-SHA-512
	RsaSha512 = "RSA-SHA-512"

	// note: SHA-256, SHA-512

	SHA256 = "SHA-256"
	SHA512 = "SHA-512"

	// note: HMAC-SHA-256, HMAC-SHA-512
	HMACSHA256 = "HMAC-SHA-256"
	HMACSHA512 = "HMAC-SHA-512"

	// note: AES-128-GCM, AES-192-GCM,AES-256-GCM
	AesGCM128 = "AES-128-GCM"
	AesGCM192 = "AES-192-GCM"
	AesGCM256 = "AES-256-GCM"

	// note: AES-128-CBC, AES-192-CBC, AES-256-CBC
	AesCBC128 = "AES-128-CBC"
	AesCBC192 = "AES-192-CBC"
	AesCBC256 = "AES-256-CBC"

	// note: AES-128-CCM, AES-192-CCM, AES-256-CCM
	AesCCM128 = "AES-128-CCM"
	AesCCM192 = "AES-192-CCM"
	AesCCM256 = "AES-256-CCM"

	// note: AES-128-CTR, AES-192-CTR, AES-256-CTR
	AesCTR128 = "AES-128-CTR"
	AesCTR192 = "AES-192-CTR"
	AesCTR256 = "AES-256-CTR"
)
