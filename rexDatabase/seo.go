package rexDatabase

type Seo struct {
	// note: og设置
	OgTitle       string `json:"og_title"`
	OgDescription string `json:"og_description"`
	OgImage       string `json:"og_image"`
	OgUrl         string `json:"og_url"`
}
