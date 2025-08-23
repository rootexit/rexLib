package rexCrypto

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

// 字节(Byte): 通常将可表示常用英文字符8位二进制称为一字节。
// 一个英文字母(不分大小写)占一个字节的空间，一个中文汉字占两个字节的空间．
// 符号：英文标点2占一个字节，中文标点占两个字节．
// 1字节(Byte） = 8位(bit)
// 比特(Bit)，亦称二进制位
// 比特指二进制中的一位，是二进制最小信息单位。
// 1比特就是1位

const (
	Bits16Len    BitLen = 2
	Bits32Len    BitLen = 4
	Bits64Len    BitLen = 8
	Bits128Len   BitLen = 16
	Bits256Len   BitLen = 32
	Bits512Len   BitLen = 64
	Bits1024Len  BitLen = 128
	Bits2048Len  BitLen = 256
	Bits4096Len  BitLen = 512
	Bits8192Len  BitLen = 1024
	Bits16384Len BitLen = 2048
)

type (
	BitLen   int
	randTool interface {
		GetAnyBtLen(num int) BitLen
		RandInt(min, max int) int
		RandInt64(min, max int64) int64
		RandLowerString(stringLen int) string
		RandStringRunes(n int) string
		GenValidateCode(width int) string
		RandBytesHexNoErr(btLen BitLen) string
		RandBytesHex(btLen BitLen) (string, error)
		RandBytesBase(btLen BitLen) (string, error)
		RandBytesBaseNoErr(btLen BitLen) string
		RandBytesUrlBase(btLen BitLen) (string, error)
		RandBytesUrlBaseNoErr(btLen BitLen) string
		RandBytes(btLen BitLen) ([]byte, error)
		RandBytesNoErr(btLen BitLen) []byte
	}
	defaultRand struct {
	}
)

func NewRand() randTool {
	return &defaultRand{}
}

// note: 输入比特的长度，返回字节位数
func (r *defaultRand) GetAnyBtLen(BytesNum int) BitLen {
	if BytesNum%8 != 0 {
		return BitLen(BytesNum / 8)
	}
	return BitLen(BytesNum / 8)
}

// RandInt64 返回 [min, max) 之间的随机 int64
func (r *defaultRand) RandInt64(min, max int64) int64 {
	if min >= max {
		return max
	}
	// big.NewInt 参数必须 > 0
	nBig, err := rand.Int(rand.Reader, big.NewInt(max-min))
	if err != nil {
		panic(fmt.Errorf("crypto/rand failed: %w", err))
	}
	return nBig.Int64() + min
}

// RandInt 返回 [min, max) 之间的随机 int
func (r *defaultRand) RandInt(min, max int) int {
	if min >= max {
		return max
	}
	diff := int64(max - min)
	nBig, err := rand.Int(rand.Reader, big.NewInt(diff))
	if err != nil {
		panic(fmt.Errorf("crypto/rand failed: %w", err))
	}
	return int(nBig.Int64()) + min
}

// 生成指定宽度的数字验证码
func (r *defaultRand) GenValidateCode(width int) string {
	if width <= 0 {
		return ""
	}

	var sb strings.Builder
	for i := 0; i < width; i++ {
		// 生成 [0,10) 的随机数
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			panic(fmt.Errorf("crypto/rand failed: %w", err))
		}
		sb.WriteByte('0' + byte(n.Int64())) // 转为字符 '0'...'9'
	}
	return sb.String()
}

// RandStringRunes 返回一个指定长度的随机字符串（包含数字+大小写字母）
func (r *defaultRand) RandStringRunes(n int) string {
	const letterRunes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]rune, n)
	for i := range b {
		// 生成 [0, len(letterRunes)) 的安全随机数
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
		if err != nil {
			panic(err) // 理论上不会发生
		}
		b[i] = rune(letterRunes[num.Int64()])
	}
	return string(b)
}

func (r *defaultRand) RandLowerString(stringLen int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, stringLen)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			panic(err)
		}
		b[i] = letters[num.Int64()]
	}
	return string(b)
}

func (r *defaultRand) RandBytesHexNoErr(btLen BitLen) string {
	return hex.EncodeToString(r.RandBytesNoErr(btLen))
}

func (r *defaultRand) RandBytesHex(btLen BitLen) (string, error) {
	tmp, err := r.RandBytes(btLen)
	return hex.EncodeToString(tmp), err
}

func (r *defaultRand) RandBytesBase(btLen BitLen) (string, error) {
	tmp, err := r.RandBytes(btLen)
	return base64.StdEncoding.EncodeToString(tmp), err
}

func (r *defaultRand) RandBytesBaseNoErr(btLen BitLen) string {
	return base64.StdEncoding.EncodeToString(r.RandBytesNoErr(btLen))
}

func (r *defaultRand) RandBytesUrlBase(btLen BitLen) (string, error) {
	tmp, err := r.RandBytes(btLen)
	return base64.URLEncoding.EncodeToString(tmp), err
}

func (r *defaultRand) RandBytesUrlBaseNoErr(btLen BitLen) string {
	return base64.URLEncoding.EncodeToString(r.RandBytesNoErr(btLen))
}

func (r *defaultRand) RandBytes(btLen BitLen) ([]byte, error) {
	tmp := make([]byte, btLen) // 32 字节
	_, err := rand.Read(tmp)
	if err != nil {
		return nil, err
	}
	return tmp, err
}

func (r *defaultRand) RandBytesNoErr(btLen BitLen) []byte {
	bt, _ := r.RandBytes(btLen)
	return bt
}
