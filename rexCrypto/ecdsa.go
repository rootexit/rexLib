package rexCrypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"time"
)

var (
	ErrorInvalidECDSAPrivateKeyPEMFormat = errors.New("invalid public key PEM format")
	ErrorInvalidECDSAPublicKeyPEMFormat  = errors.New("invalid public key PEM format")
	ErrorPublicKeyNotECDSA               = errors.New("public key not ECDSA type")
)

// **解析 ECDSA 公钥 **
func ParseECDSAPrivateKeyFromPEM(privateKeyPEM string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, ErrorInvalidECDSAPrivateKeyPEMFormat
	}

	// **解析 ECDSA 私钥**
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func ParseECDSAPublicKeyFromCert(certPEM string) (*ecdsa.PublicKey, error) {
	// **解析 PEM 证书**
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, ErrorInvalidECDSAPublicKeyPEMFormat
	}
	// **解析 X.509 证书**
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	// **获取公钥**
	ecdsaPubKey, ok := cert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrorPublicKeyNotECDSA
	}
	return ecdsaPubKey, nil
}

// note: 生成 ECDSA 公私钥
func ECDSAGenerateKeys(curve elliptic.Curve) (privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, err error) {
	switch curve.Params() {
	case elliptic.P224().Params():
		return nil, nil, fmt.Errorf("curve must be P256 or P384 or P521")
	}
	// **1. 生成 ECDSA 私钥**
	privateKey, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// note: 生成 ECC 证书
func ECDSAGenerateECCCertificate(curve elliptic.Curve, subject pkix.Name) (keyPem, certPem []byte, err error) {
	privateKey, publicKey, err := ECDSAGenerateKeys(curve)
	if err != nil {
		return nil, nil, err
	}

	nowTime := time.Now()
	// 创建证书模板
	serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               subject,
		NotBefore:             nowTime,
		NotAfter:              nowTime.Add(365 * 24 * time.Hour), // 证书有效期 1 年
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	// 生成自签名证书
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	// 编码证书 & 私钥到 PEM 格式**
	certPEMBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	keyPEMBytes, _ := x509.MarshalECPrivateKey(privateKey)
	keyPEMBytes = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyPEMBytes})

	return keyPEMBytes, certPEMBytes, nil
}
