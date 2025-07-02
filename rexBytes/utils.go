package rexBytes

import (
	"fmt"
	"github.com/rootexit/rexLib/rexCommon"
	"sort"
	"strconv"
)

func BytesCombineBytes(high, low byte) int64 {
	highStr := fmt.Sprintf("0x%x", high)
	lowStr := fmt.Sprintf("0x%x", low)
	highNum, err := strconv.ParseInt(highStr, 16, 64) // base=16, bitSize=64
	if err != nil {
		fmt.Println("SixteenStr2int64 Error:", err)
		return 0
	}
	lowNum, err := strconv.ParseInt(lowStr, 16, 64) // base=16, bitSize=64
	if err != nil {
		fmt.Println("SixteenStr2int64 Error:", err)
		return 0
	}
	// 左移高位 8 位，使用按位或操作合并低位
	result := (highNum << 8) | lowNum
	return result
}

func SixteenStr2int64(param byte) int64 {
	hexStr := fmt.Sprintf("0x%x", param)
	// 将 16 进制字符串解析为数字
	number, err := strconv.ParseInt(hexStr, 16, 64) // base=16, bitSize=64
	if err != nil {
		fmt.Println("SixteenStr2int64 Error:", err)
		return 0
	}
	return number
}

func CombineBytes(high, low int) int {
	// 左移高位 8 位，使用按位或操作合并低位
	result := (high << 8) | low
	return result
}

func Ten2Sixteen(param byte) int {
	return rexCommon.Str2Int(fmt.Sprintf("0x%x", param))
}

func Map2Bytes(data map[string]uint8) []byte {
	// note: 转换成buffer
	var result []byte
	keys := make([]int, 0, len(data))
	for k := range data {
		key, _ := strconv.Atoi(k) // 转换键为整数
		keys = append(keys, key)
	}
	sort.Ints(keys) // 按键升序排列

	for _, k := range keys {
		result = append(result, data[strconv.Itoa(k)]) // 将值转换为 byte 类型
	}
	return result
}

func Ten2sixteen2uint(param byte) uint {
	return rexCommon.Str2Uint(fmt.Sprintf("%x", param))
}

func SplitBytes(data []byte, segmentSize int) [][]byte {
	// 计算能分割的段数
	totalSegments := len(data) / segmentSize

	// 创建结果存储的二维切片
	result := make([][]byte, 0, totalSegments)

	// 按段切割数据
	for i := 0; i < totalSegments; i++ {
		start := i * segmentSize
		end := start + segmentSize
		result = append(result, data[start:end])
	}
	// 保留多余部分
	//if len(data)%segmentSize != 0 {
	//	remainderStart := totalSegments * segmentSize
	//	result = append(result, data[remainderStart:])
	//}
	return result
}
