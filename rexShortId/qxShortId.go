package rexShortId

import (
	"crypto/sha256"
	"fmt"
	"math/big"
)

// Base62 字符表
var charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// 编码大整数为 base62 字符串
func base62EncodeBigInt(n *big.Int) string {
	if n.Cmp(big.NewInt(0)) == 0 {
		return "0"
	}
	base := big.NewInt(62)
	result := ""
	zero := big.NewInt(0)
	rem := new(big.Int)
	for n.Cmp(zero) > 0 {
		n, rem = new(big.Int).DivMod(n, base, rem)
		result = string(charset[rem.Int64()]) + result
	}
	return result
}

// note: 保持向下兼容，以前的 GenerateShortID 函数, 使用ShortId替代
// Deprecated: Use NewFunction instead.
func GenerateShortID(biz string, mid int64, salt string) string {
	// 拼接数据
	raw := fmt.Sprintf("%s_%d_%s", biz, mid, salt)

	// 哈希
	hash := sha256.Sum256([]byte(raw))

	// 取前 8-16 字节转换为整数（防止过长）
	hashSlice := hash[:12] // 可调整长度
	bigInt := new(big.Int).SetBytes(hashSlice)

	// 编码为 base62
	return base62EncodeBigInt(bigInt)
}

func ShortId(raw string) string {
	// 哈希
	hash := sha256.Sum256([]byte(raw))

	// 取前 8-16 字节转换为整数（防止过长）
	hashSlice := hash[:12] // 可调整长度
	bigInt := new(big.Int).SetBytes(hashSlice)

	// 编码为 base62
	return base62EncodeBigInt(bigInt)
}
