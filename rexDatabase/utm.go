package rexDatabase

type Utm struct {
	UtmSource   string `gorm:"column:utm_source;comment:流量来源;type:varchar(255);" json:"utm_source"`          // 流量来源,google、wechat、facebook
	UtmMedium   string `gorm:"column:utm_medium;comment:渠道类型;type:varchar(255);" json:"utm_medium"`          // 渠道类型,cpc、email、social
	UtmCampaign string `gorm:"column:utm_campaign;comment:活动名称;type:varchar(255);" json:"utm_campaign"`      // 活动名称,black_friday_2026
	UtmTerm     string `gorm:"column:utm_term;comment:广告关键词（主要用于搜索广告）;type:varchar(255);" json:"utm_term"`   // 广告关键词,running+shoes
	UtmContent  string `gorm:"column:utm_content;comment:区分同一广告的不同素材;type:varchar(255);" json:"utm_content"` // 区分同一广告的不同素材,banner_a、button_red
	UtmId       string `gorm:"column:utm_id;comment:广告ID;type:varchar(255);" json:"utm_id"`                  // 广告ID,123456
}

type Attribution struct {
	UTM Utm `json:"utm"` // 标准化

	GCLID string `json:"gclid,omitempty"` //Google Ads

	FBCLID string `json:"fbclid,omitempty"` // Facebook

	MSCLKID string `json:"msclkid,omitempty"` // Microsoft Ads

	TTCLID string `json:"ttclid,omitempty"` // TikTok

	WBRAID string `json:"wbraid,omitempty"` // Google iOS App 流量归因

	GBRAID string `json:"gbraid,omitempty"` // Google Web 流量归因（受隐私限制场景）

	RDTCID string `json:"rdt_cid,omitempty"` // Reddit

	TWCLID string `json:"twclid,omitempty"` // X(Twitter)
}
