package rexCrypto

import (
	"encoding/base64"
	"encoding/hex"
)

// note: 单向不可逆算法

func Sha256(data string) (hexStr string, baseStr string) {
	hash := NewHash().Sha256([]byte(data))
	hexStr = hex.EncodeToString(hash[:])
	baseStr = base64.StdEncoding.EncodeToString(hash[:])
	return hexStr, baseStr
}

func Sha512(data string) (hexStr string, baseStr string) {
	hash := NewHash().Sha512([]byte(data))
	hashHex := hex.EncodeToString(hash[:])
	hexStr = hex.EncodeToString(hash[:])
	baseStr = base64.StdEncoding.EncodeToString(hash[:])
	return hashHex, baseStr
}

func HMACSha256(data string, key string) (hexStr string, baseStr string) {
	hash := NewHash().HMACSha256([]byte(data), []byte(key))
	hexStr = hex.EncodeToString(hash)
	baseStr = base64.StdEncoding.EncodeToString(hash)
	return hexStr, baseStr
}

func HMACSha512(data string, key string) (hexStr string, baseStr string) {
	hash := NewHash().HMACSha256([]byte(data), []byte(key))
	hexStr = hex.EncodeToString(hash)
	baseStr = base64.StdEncoding.EncodeToString(hash)
	return hexStr, baseStr
}
