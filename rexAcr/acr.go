package rexAcr

import (
	"fmt"
	"slices"
	"strings"
)

const (
	AmrTypePwd         string = "pwd"            // 密码
	AmrTypeOtpEmail    string = "otp:email"      // 邮箱验证码
	AmrTypeOtpSms      string = "otp:sms"        // 短信验证码
	AmrTypeOtpTotp     string = "otp:totp"       // TOTP验证码
	AmrTypeDeviceToken string = "device_token"   // 硬件验证器
	AmrTypeFingerprint string = "fingerprint"    // 指纹
	AmrTypeFaceId      string = "faceid"         // 人脸识别
	AmrTypeIris        string = "iris"           // 虹膜识别
	AmrTypeVoice       string = "voice"          // 语音识别
	AmrTypeWebAuthUv   string = "webauthn_uv"    // WebAuthn用户验证器
	AmrDeviceBinding   string = "device_binding" // 设备绑定
)

type (
	Level   string
	AcrTool interface {
		GetAcrLow() Level
		GetAcrMedium() Level
		GetAcrHigh() Level
		GetRank() map[Level]int
		GetLevelCombos() map[Level][][]string
		Norm(xs []string) []string
		ContainsAll(completed []string, combo []string) bool
		AchievedLevel(completed []string) (Level, []string)
		Meets(required Level, completed []string) (ok bool, achieved Level, matchedCombo []string)
	}
	defaultAcrTool struct {
		AcrLow      Level
		AcrMedium   Level
		AcrHigh     Level
		Rank        map[Level]int
		LevelCombos map[Level][][]string
	}
)

func NewDefaultAcrTool(appName string) AcrTool {
	lowLevel := Level(fmt.Sprintf("urn:%s:acr:low", appName))
	mediumLevel := Level(fmt.Sprintf("urn:%s:acr:medium", appName))
	highLevel := Level(fmt.Sprintf("urn:%s:acr:high", appName))
	return &defaultAcrTool{
		AcrLow:    lowLevel,
		AcrMedium: mediumLevel,
		AcrHigh:   highLevel,
		Rank: map[Level]int{
			lowLevel:    1,
			mediumLevel: 2,
			highLevel:   3,
		},
		LevelCombos: map[Level][][]string{
			lowLevel: {
				{AmrTypePwd},
				{AmrTypeOtpEmail},
			},
			mediumLevel: {
				{AmrTypeOtpSms},
				{AmrTypePwd, AmrDeviceBinding},
				{AmrTypePwd, AmrTypeOtpSms},
				{AmrTypePwd, AmrTypeOtpTotp},
				{AmrTypePwd, AmrTypeOtpEmail},
				{AmrTypePwd, AmrTypeDeviceToken},
			},
			highLevel: {
				{AmrTypeFingerprint, AmrDeviceBinding},
				{AmrTypeFaceId, AmrDeviceBinding},
				{AmrTypeIris, AmrDeviceBinding},
				{AmrTypeVoice, AmrDeviceBinding},
			},
		},
	}
}

func NewCustomAcrTool(appName string, lowCombos, mediumCombos, highCombos [][]string) AcrTool {
	lowLevel := Level(fmt.Sprintf("urn:%s:acr:low", appName))
	mediumLevel := Level(fmt.Sprintf("urn:%s:acr:medium", appName))
	highLevel := Level(fmt.Sprintf("urn:%s:acr:high", appName))
	return &defaultAcrTool{
		AcrLow:    lowLevel,
		AcrMedium: mediumLevel,
		AcrHigh:   highLevel,
		Rank: map[Level]int{
			lowLevel:    1,
			mediumLevel: 2,
			highLevel:   3,
		},
		LevelCombos: map[Level][][]string{
			lowLevel:    lowCombos,
			mediumLevel: mediumCombos,
			highLevel:   highCombos,
		},
	}
}

// 后期修改LevelCombos
func (d *defaultAcrTool) WithLevelCombos(lowCombos, mediumCombos, highCombos [][]string) {
	d.LevelCombos = map[Level][][]string{
		d.AcrLow:    lowCombos,
		d.AcrMedium: mediumCombos,
		d.AcrHigh:   highCombos,
	}
}

// 获取当前的Low的数据

func (d *defaultAcrTool) GetAcrLow() Level {
	return d.AcrLow
}

// 获取当前的Medium的数据

func (d *defaultAcrTool) GetAcrMedium() Level {
	return d.AcrMedium
}

// 获取当前的High的数据

func (d *defaultAcrTool) GetAcrHigh() Level {
	return d.AcrHigh
}

// 获取当前的Rank的数据

func (d *defaultAcrTool) GetRank() map[Level]int {
	return d.Rank
}

// 获取当前的LevelCombos的数据

func (d *defaultAcrTool) GetLevelCombos() map[Level][][]string {
	return d.LevelCombos
}

// 规范化：全部小写去空白，便于比较

func (d *defaultAcrTool) Norm(xs []string) []string {
	out := make([]string, 0, len(xs))
	for _, x := range xs {
		x = strings.TrimSpace(strings.ToLower(x))
		if x != "" {
			out = append(out, x)
		}
	}
	slices.Sort(out)
	out = slices.Compact(out)
	return out
}

// containsAll: completed 是否包含 combo 里的全部要素

func (d *defaultAcrTool) ContainsAll(completed []string, combo []string) bool {
	for _, need := range combo {
		if !slices.Contains(completed, need) {
			return false
		}
	}
	return true
}

// AchievedLevel 计算“已完成要素”所达到的**最高** ACR 等级；若不满足任何等级，返回空串。

func (d *defaultAcrTool) AchievedLevel(completed []string) (Level, []string) {
	c := d.Norm(completed)
	// 从高到低检查，匹配到第一条就返回（最大等级）
	order := []Level{d.AcrHigh, d.AcrMedium, d.AcrLow}
	for _, lv := range order {
		for _, combo := range d.LevelCombos[lv] {
			cb := d.Norm(combo)
			if d.ContainsAll(c, cb) {
				return lv, cb
			}
		}
	}
	return "", nil
}

// Meets 要求：只要“达到的等级” >= “所需等级”即可通过

func (d *defaultAcrTool) Meets(required Level, completed []string) (ok bool, achieved Level, matchedCombo []string) {
	achieved, matchedCombo = d.AchievedLevel(completed)
	if achieved == "" {
		return false, achieved, nil
	}
	return d.Rank[achieved] >= d.Rank[required], achieved, matchedCombo
}
