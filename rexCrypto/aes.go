package rexCrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	AES128KeyLen = 16
	AES192KeyLen = 24
	AES256KeyLen = 32
	AESGCMIvLen  = 12
	AESCCMIvLen  = 24
	AEsCBCIvLen  = 16
	AESCTRIvLen  = 16
)

func GenAESKeyAndIv(KeyLen int, IvLen int) (keyBase, ivBase string, err error) {
	key := make([]byte, KeyLen)
	_, err = rand.Read(key)
	if err != nil {
		return "", "", err
	}
	iv := make([]byte, IvLen)
	_, err = rand.Read(iv)
	if err != nil {
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(key), base64.StdEncoding.EncodeToString(iv), nil
}

// **AES-GCM 加密**
func AESEncryptByGCM(plainText []byte, keyBase, ivBase string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(keyBase)
	if err != nil {
		return "", err
	}
	iv, err := base64.StdEncoding.DecodeString(ivBase)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	// 加密 + 认证
	cipherText := aesGCM.Seal(nil, iv, plainText, nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func AESEncryptByGCMBt(plainText, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	// 加密 + 认证
	cipherText := aesGCM.Seal(nil, iv, plainText, nil)
	return cipherText, nil
}

// **AES-GCM 解密**
func AESDecryptByGCM(cipherTextBase64 string, keyBase, ivBase string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(keyBase)
	if err != nil {
		return nil, err
	}
	iv, err := base64.StdEncoding.DecodeString(ivBase)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	cipherText, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return nil, err
	}
	plainText, err := aesGCM.Open(nil, iv, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

// **AES-GCM 解密**
func AESDecryptByGCMBt(cipherTextBt, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plainText, err := aesGCM.Open(nil, iv, cipherTextBt, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

// **AES-CCM 加密**
func AESEncryptByCCM(plainText []byte, keyBase, ivBase string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(keyBase)
	if err != nil {
		return "", err
	}
	iv, err := base64.StdEncoding.DecodeString(ivBase)
	if err != nil {
		return "", err
	}
	aesCCM, err := chacha20poly1305.NewX(key)
	if err != nil {
		return "", err
	}
	cipherText := aesCCM.Seal(nil, iv, plainText, nil)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// **AES-CCM 解密**
func AESDecryptByCCM(cipherTextBase64 string, keyBase, ivBase string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(keyBase)
	if err != nil {
		return nil, err
	}
	iv, err := base64.StdEncoding.DecodeString(ivBase)
	if err != nil {
		return nil, err
	}
	aesCCM, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}
	cipherText, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return nil, err
	}
	plainText, err := aesCCM.Open(nil, iv, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

// **PKCS7 填充**
func CBCPkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// **PKCS7 去填充**
func CBCPkcs7Unpad(data []byte) []byte {
	length := len(data)
	unpad := int(data[length-1])
	return data[:(length - unpad)]
}

// **AES-CBC 加密**
func AESEncryptByCBC(plainText []byte, keyBase, ivBase string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(keyBase)
	if err != nil {
		return "", err
	}
	iv, err := base64.StdEncoding.DecodeString(ivBase)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// PKCS7 填充
	plainText = CBCPkcs7Pad(plainText, aes.BlockSize)

	cipherText := make([]byte, len(plainText))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, plainText)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// **AES-CBC 解密**
func AESDecryptByCBC(cipherTextBase64 string, keyBase, ivBase string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(keyBase)
	if err != nil {
		return nil, err
	}
	iv, err := base64.StdEncoding.DecodeString(ivBase)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Base64 解码密文
	cipherText, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return nil, err
	}

	// CBC 解密
	plainText := make([]byte, len(cipherText))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plainText, cipherText)

	// PKCS7 去填充
	plainText = CBCPkcs7Unpad(plainText)

	return plainText, nil
}

// **AES-CTR 加密**
func AESEncryptByCTR(plainText []byte, keyBase, ivBase string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(keyBase)
	if err != nil {
		return "", err
	}
	iv, err := base64.StdEncoding.DecodeString(ivBase)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, iv)
	cipherText := make([]byte, len(plainText))
	stream.XORKeyStream(cipherText, plainText)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// **AES-CTR 解密**
func AESDecryptByCTR(cipherTextBase64 string, keyBase, ivBase string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(keyBase)
	if err != nil {
		return nil, err
	}
	iv, err := base64.StdEncoding.DecodeString(ivBase)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)
	cipherText, _ := base64.StdEncoding.DecodeString(cipherTextBase64)
	plainText := make([]byte, len(cipherText))
	stream.XORKeyStream(plainText, cipherText)
	return plainText, nil
}
