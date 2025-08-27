package rexCrypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
)

type (
	HashTool interface {
		Sha256(data []byte) []byte
		Sha512(data []byte) []byte
		HMACSha256(data []byte, key []byte) []byte
		HMACSha512(data []byte, key []byte) []byte
	}
	defaultHash struct {
	}
)

func NewHash() HashTool {
	return &defaultHash{}
}

func (d *defaultHash) Sha256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func (d *defaultHash) Sha512(data []byte) []byte {
	hash := sha512.Sum512(data)
	return hash[:]
}

func (d *defaultHash) HMACSha256(data []byte, key []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

func (d *defaultHash) HMACSha512(data []byte, key []byte) []byte {
	hash := hmac.New(sha512.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}
