package rexUtils

import "github.com/rootexit/rexLib/rexCrypto"

type ICPRecord struct {
	IsShow bool   `json:"isShow"`
	Icon   string `json:"icon"`
	Href   string `json:"href"`
	Text   string `json:"text"`
}

func DefaultICPRecord() ICPRecord {
	return ICPRecord{
		IsShow: false,
		Icon:   "https://corecdn.csvw88.com/statics/gov/icp.gif",
		Href:   "https://beian.miit.gov.cn/",
		Text:   "XICP备1996091901号-x",
	}
}

type MpsRecord struct {
	IsShow bool   `json:"isShow"`
	Icon   string `json:"icon"`
	Href   string `json:"href"`
	Text   string `json:"text"`
}

func DefaultMpsRecord() MpsRecord {
	return MpsRecord{
		IsShow: false,
		Icon:   "https://corecdn.csvw88.com/statics/gov/mps.png",
		Href:   "https://beian.mps.gov.cn/",
		Text:   "X公网安备 00000000000000号",
	}
}

func GenPassword(realPwd string) (encryptPwdBase, key, salt string) {
	keyBt, _ := rexCrypto.RandBytes(32)
	// 生成盐值
	saltBt, _ := rexCrypto.RandBytes(16)
	// 生成hmac-sha256后的密码
	_, basePwd := rexCrypto.HMACSha256(realPwd+string(saltBt), string(keyBt))
	return basePwd, key, salt
}
