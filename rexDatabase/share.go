package rexDatabase

type SeoShareSetting struct {
	Debug         bool   `gorm:"column:debug;comment:debug;type: boolean;default:false" json:"debug"`
	OgTitle       string `gorm:"column:og_title;comment:og标题;type: varchar(255);" json:"og_title"`
	OgDescription string `gorm:"column:og_description;comment:og描述;type: varchar(255);" json:"og_description"`
	OgImage       string `gorm:"column:og_image;comment:og图片;type: varchar(255);" json:"og_image"`
	Url           string `gorm:"column:url;comment:分享链接;type: varchar(255);" json:"url"`                      // 分享链接
	Card          string `gorm:"column:card;comment:分享卡片;type: varchar(255);" json:"card"`                    // 分享卡片
	CreatorSite   string `gorm:"column:creator_site;comment:示例@xxxx;type: varchar(255);" json:"creator_site"` // 创建者站点
}
