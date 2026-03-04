package rexDatabase

type ContentType string

// note: 内容类型，1->H5,2->外部链接,3->小程序,4->视频文章,5->PDF文件,6->快讯;

const (
	ContentTypeH5      ContentType = "h5"
	ContentTypeLink    ContentType = "external_link"
	ContentTypeMiniApp ContentType = "mini_app"
	ContentTypeVideo   ContentType = "video"
	ContentTypePdf     ContentType = "pdf"
	ContentTypeNews    ContentType = "express"
)
