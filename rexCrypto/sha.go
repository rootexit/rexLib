package rexCrypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
)

// note: 单向不可逆算法

func Sha256(data string) (hexStr string, baseStr string) {
	hash := sha256.Sum256([]byte(data))
	hexStr = hex.EncodeToString(hash[:])
	baseStr = base64.StdEncoding.EncodeToString(hash[:])
	return hexStr, baseStr
}

func Sha512(data string) (hexStr string, baseStr string) {
	hash := sha512.Sum512([]byte(data))
	hashHex := hex.EncodeToString(hash[:])
	hexStr = hex.EncodeToString(hash[:])
	baseStr = base64.StdEncoding.EncodeToString(hash[:])
	return hashHex, baseStr
}

func HMACSha256(data string, key string) (hexStr string, baseStr string) {
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write([]byte(data))
	hexStr = hex.EncodeToString(hash.Sum(nil))
	baseStr = base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return hexStr, baseStr
}

func HMACSha512(data string, key string) (hexStr string, baseStr string) {
	hash := hmac.New(sha512.New, []byte(key))
	hash.Write([]byte(data))
	hexStr = hex.EncodeToString(hash.Sum(nil))
	baseStr = base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return hexStr, baseStr
}
