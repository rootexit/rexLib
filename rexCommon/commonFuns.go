package rexCommon

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"google.golang.org/grpc/peer"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

/**
 * @Author joker
 * @Description //TODO 常用函数库
 * @Date 2020-7-12 17:08:54
 **/

// 注意client 本身是连接池，不要每次请求时创建client
var (
	HttpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
)

func GetScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	if scheme := r.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	return "http"
}

func ExtractCNAndENAndNum(s string) []string {
	// 匹配：中文、英文大写、小写
	re := regexp.MustCompile("[a-zA-Z0-9\u4e00-\u9fa5]+")
	matches := re.FindAllString(s, -1)
	// 去重 + 去空
	seen := make(map[string]struct{})
	var result []string
	for _, part := range matches {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if _, exists := seen[part]; !exists {
			seen[part] = struct{}{}
			result = append(result, part)
		}
	}
	return result
}

func SplitBySymbolsAndDedup(s string) []string {
	// 匹配所有非字母数字（包括空格、符号、标点等）
	re := regexp.MustCompile(`[\P{L}\P{N}]+`)
	parts := re.Split(s, -1)

	// 去重 + 去空
	seen := make(map[string]struct{})
	var result []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if _, exists := seen[part]; !exists {
			seen[part] = struct{}{}
			result = append(result, part)
		}
	}
	return result
}

func SplitAndDedup(s string) []string {
	// 使用正则表达式分割：空格、中文逗号、英文逗号
	re := regexp.MustCompile(`[ ,，]+`)
	parts := re.Split(s, -1)

	// 去重 + 去空
	seen := make(map[string]struct{})
	var result []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if _, exists := seen[part]; !exists {
			seen[part] = struct{}{}
			result = append(result, part)
		}
	}
	return result
}

func ExtractPath(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("URL parse error: %w", err)
	}
	return parsedURL.Path, nil
}

func ExtractScheme(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("URL parse error: %w", err)
	}
	return parsedURL.Scheme, nil
}

func ExtractDomain(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("URL parse error: %w", err)
	}
	return parsedURL.Hostname(), nil
}

func GenerateClientId() string {
	return GenerateRandHex(8)
}

func GenerateClientSecretHex() string {
	return GenerateRandHex(16)
}

func GenerateRandHex(num int) string {
	b := make([]byte, num) // 32字节
	rand.Read(b)
	return hex.EncodeToString(b)
}

func PrintMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Printf("Alloc: %v KB\n", m.Alloc/1024)           // 当前已分配的堆内存
	log.Printf("TotalAlloc: %v KB\n", m.TotalAlloc/1024) // 分配的堆内存总量
	log.Printf("Sys: %v KB\n", m.Sys/1024)               // 从系统申请的内存
	log.Printf("NumGC: %v\n", m.NumGC)                   // 垃圾回收的次数
}

func FlatToNested(input map[string]interface{}) map[string]interface{} {
	nested := make(map[string]interface{})
	for key, value := range input {
		parts := strings.Split(key, ".")
		current := nested
		for i, part := range parts {
			if i == len(parts)-1 {
				current[part] = value
			} else {
				if _, ok := current[part]; !ok {
					current[part] = make(map[string]interface{})
				}
				current = current[part].(map[string]interface{})
			}
		}
	}
	return nested
}

func MaskPhoneWithRegex(phone string) string {
	// 匹配区号的正则：支持 "+数字 " 格式
	re := regexp.MustCompile(`^\+[\d]+\s`)
	matches := re.FindString(phone)
	if matches != "" {
		// 提取区号和号码
		countryCode := matches
		actualNumber := phone[len(matches):]
		if len(actualNumber) < 7 {
			return "Invalid phone number"
		}
		// 替换中间部分为星号
		return countryCode + actualNumber[:3] + "****" + actualNumber[len(actualNumber)-4:]
	}
	// 如果没有区号，按普通号码处理
	if len(phone) < 7 {
		return "Invalid phone number"
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

func MaskPhoneDynamic(phone string) string {
	length := len(phone)
	if length < 7 {
		return "Invalid phone number"
	}
	// 保留前三位和后四位，中间用星号替代
	return phone[:3] + "****" + phone[length-4:]
}

func CalculateAgeTime(birthday time.Time) int {
	now := time.Now()
	age := now.Year() - birthday.Year()

	if now.YearDay() < birthday.YearDay() {
		age--
	}
	return age
}

func CalculateAge(birthday string) int {
	if birthday == "0000-00-00" || birthday == "" || len(birthday) != 10 {
		return 0
	}
	birthdayTime, _ := time.Parse(time.DateOnly, birthday)
	now := time.Now()
	age := now.Year() - birthdayTime.Year()

	if now.YearDay() < birthdayTime.YearDay() {
		age--
	}
	return age
}

func BindAndCheck(ctx *gin.Context, data interface{}) error {
	if err := ctx.ShouldBindJSON(data); err != nil {
		return errors.New(fmt.Sprintf("bindjson err%s", err))
	}
	// 校验数据
	validate := validator.New()
	if err := validate.Struct(data); err != nil {
		return errors.New(fmt.Sprintf("validator err%s", err))
	}
	return nil
}

func RandInt64(min, max int64) int {
	if min >= max || min == 0 || max == 0 {
		return int(max)
	}
	return int(rand.Int63n(max-min) + min)
}

func RandInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return int(max)
	}
	return rand.Intn(max-min) + min
}

func DelFilelist(path string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return nil //f为空 错误不为空 错误是文件不存在 可以忽略
		}
		if f.IsDir() {
			fmt.Printf("文件夹 继续递归 %s  \n", path)
			DelFilelist(path)
		} else {
			err := os.RemoveAll(path)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("删除文件 %s  \n", path)
			return nil
		}
		//	println(path)
		return nil
	})
	if err != nil {
		fmt.Printf("walk 错误 err: %v\n", err)
	}
}

// RandStringRunes 返回随机字符串
func RandStringRunes(n int) string {
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func UploadFile(url string, params map[string]string, nameField, fileName string, file io.Reader) ([]byte, error) {
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)

	formFile, err := writer.CreateFormFile(nameField, fileName)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(formFile, file)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	//req.Header.Set("Content-Type","multipart/form-data")
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func DistributeFile(url string, params map[string]string, nameField, path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)

	formFile, err := writer.CreateFormFile(nameField, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(formFile, file)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		writer.WriteField(key, val)
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	//req.Header.Set("Content-Type","multipart/form-data")
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// Replace 根据替换表执行批量替换
func Replace(table map[string]string, s string) string {
	for key, value := range table {
		s = strings.Replace(s, key, value, -1)
	}
	return s
}

func StringToBool(string2 string) bool {
	//todo :string to bool
	b, _ := strconv.ParseBool(string2)
	return b
}
func BoolToString(bool2 bool) string {
	//todo :bool to string
	sBool := strconv.FormatBool(bool2) //方法1
	return sBool
}

// Int2Str int类型转string类型
func Int2Str(inter int) string {
	string := strconv.Itoa(inter)
	return string
}

// Int2Str int64类型转string类型 精确到后2位 精确到后4位
func Int642Str(inter int64) string {

	string := strconv.FormatInt(inter, 10)
	return string
}

func Int64Str(inter int64) string {
	string := strconv.FormatInt(inter, 10)
	return string
}

func Str2Float64(in string) float64 {
	//num, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", num), 64)
	float, _ := strconv.ParseFloat(in, 64)
	return float
}

// Str2Int string类型转Int类型
func Str2Int(inter string) int {
	int, err := strconv.Atoi(inter)
	if err != nil {
		log.Println("err", err)
	}
	return int
}

// Str2Int string类型转Int类型
func Str2Uint(inter string) uint {
	uint64, _ := strconv.ParseUint(inter, 10, 64)
	return uint(uint64)
}

func Str2Uint32(inter string) uint32 {
	uint64, _ := strconv.ParseUint(inter, 10, 64)
	return uint32(uint64)
}

// Str2Int64 string类型转Int64类型
func Str2Int64(inter string) int64 {
	int64, _ := strconv.ParseInt(inter, 10, 64)
	return int64
}

func Arr2Str(strings []string) string {
	b, _ := json.Marshal(strings)
	return fmt.Sprintf("%s", b)
}

func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func GetClientIP(ctx context.Context) (string, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("[getClinetIP] invoke FromContext() failed")
	}
	if pr.Addr == net.Addr(nil) {
		return "", fmt.Errorf("[getClientIP] peer.Addr is nil")
	}
	addSlice := strings.Split(pr.Addr.String(), ":")
	return addSlice[0], nil
}

// 进行Sha1编码
func Sha1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// TODO: 获取当月的最后第一天或者最后一天
func ReturnSpecifyMonth(year, month int) (time.Time, time.Time) {
	//currentYear, currentMonth, _ := now.Date()
	now := time.Now()
	currentLocation := now.Location()

	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, 0)
	return firstOfMonth, lastOfMonth
}

// TODO: 获取当年的最后第一天或者最后一天
func ReturnSpecifyYear(year int) (time.Time, time.Time) {
	//currentYear, currentMonth, _ := now.Date()
	now := time.Now()
	currentLocation := now.Location()

	firstOfMonth := time.Date(year, 1, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 12, 0)
	return firstOfMonth, lastOfMonth
}

// note: 返回wechat编码

type WeChatShareConfig struct {
	AppId     string   `json:"AppId"`
	Timestamp int64    `json:"Timestamp"`
	NonceStr  string   `json:"NonceStr"`
	Debug     bool     `json:"Debug"`
	JsApiList []string `json:"JsApiList"`
	Signature string   `json:"Signature"`
}

func GetWeChatShareConfig(debug bool, ticket, shareLink, appid string, JsApiList []string) WeChatShareConfig {
	noncestr := RandStringRunes(16)
	timestamp := time.Now().Unix()
	tempUrl := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		ticket,
		noncestr,
		timestamp,
		shareLink)
	sha1URI := Sha1(tempUrl)
	weChatShareConfig := WeChatShareConfig{
		AppId:     appid,
		Timestamp: timestamp,
		NonceStr:  noncestr,
		Debug:     debug,
		JsApiList: JsApiList,
		Signature: sha1URI,
	}
	return weChatShareConfig
}

// Slice2Str string类型转Int64类型
func SliceInt2Str(inters interface{}) string {
	temp := fmt.Sprintf("%d", inters)
	temp = temp[1 : len(temp)-1]
	temp = strings.Replace(temp, " ", ",", -1)
	return temp
}
