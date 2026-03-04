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

type VisibleType string

// note: 文章可见类型，1->不可见,2->所有成员可见,4->仅订阅用户可见,5->部分人可见,6->部分人不可见,7->所有人可见
const (
	VisibleTypeInvisible     VisibleType = "private"
	VisibleTypeAllMember     VisibleType = "all_member"
	VisibleTypeSubscribeOnly VisibleType = "subscribe_only"
	VisibleTypePartly        VisibleType = "partly"
	VisibleTypePartlyNot     VisibleType = "partly_not"
	VisibleTypePublic        VisibleType = "public"
)

type AuthorType string

// note: 作者类型, 1->独立作者, 2->联合创作
const (
	AuthorTypeIndependent AuthorType = "independent"
	AuthorTypeJoint       AuthorType = "joint"
)

type MiniAppType string

// 如果是小程序类型,1->微信小程序,2->支付宝小程序,3->抖音小程序

const (
	MiniAppTypeNone   MiniAppType = "none"
	MiniAppTypeWeChat MiniAppType = "wechat"
	MiniAppTypeAlipay MiniAppType = "alipay"
	MiniAppTypeDouyin MiniAppType = "douyin"
)

type CollectionType string

// note: 类型,1->站内文章合集,2->跳转链接;
const (
	CollectionTypeContent CollectionType = "collection"
	CollectionTypeLink    CollectionType = "link"
)
