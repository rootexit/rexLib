package rexDatabase

type SeoBase struct {
	// note: og设置
	OgTitle       string `json:"og_title"`
	OgDescription string `json:"og_description"`
	OgImage       string `json:"og_image"`
	OgUrl         string `json:"og_url"`
}

type SeoExtra struct {
	SeoBase
	Debug       bool   `gorm:"column:debug;comment:debug;type: boolean;default:false" json:"debug"`
	Card        string `gorm:"column:card;comment:分享卡片;type: varchar(255);" json:"card"`                    // 分享卡片
	CreatorSite string `gorm:"column:creator_site;comment:示例@xxxx;type: varchar(255);" json:"creator_site"` // 创建者站点
}
