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
