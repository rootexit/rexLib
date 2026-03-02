package rexDatabase

import "strings"

type IsAccelerate int32

// note: 是否开启传输加速,1->否,2-是
const (
	TypeIsAccelerateNo IsAccelerate = iota + 1
	TypeIsAccelerateYes
)

type UploadStatus int8

// note: 状态,1->已创建未上传,2->已生成签名,3->已上传完毕
const (
	FileStatusCreated UploadStatus = iota + 1
	FileStatusSign
	FileStatusUploaded
)

type UploadType string

// 上传类型,1->对象存储前端直传,2->服务器接口直传,3->外部添加
const (
	UploadTypeObjectStorage UploadType = "ObjectStorage"
	UploadTypeServer        UploadType = "ServerUpload"
	UploadTypeExternalAdd   UploadType = "ExternalAdd"
)

type ArchiveType string

const (
	ArchiveTypeImage        ArchiveType = "image"
	ArchiveTypeVideo        ArchiveType = "video"
	ArchiveTypeAudio        ArchiveType = "audio"
	ArchiveTypeDocument     ArchiveType = "document"
	ArchiveTypePDF          ArchiveType = "pdf"
	ArchiveTypeCompressed   ArchiveType = "compressed"
	ArchiveTypeExternalLink ArchiveType = "external_link"
	ArchiveTypeOther        ArchiveType = "other"
)

const (
	ArchiveDefaultNameImage        string = "图片"
	ArchiveDefaultNameVideo        string = "视频"
	ArchiveDefaultNameAudio        string = "音频"
	ArchiveDefaultNameDocument     string = "文档"
	ArchiveDefaultNamePDF          string = "PDF"
	ArchiveDefaultNameCompressed   string = "压缩包"
	ArchiveDefaultNameExternalLink string = "外部链接"
	ArchiveDefaultNameOther        string = "其他"
)

func DefaultFileCover(archiveType ArchiveType) string {
	switch archiveType {
	case ArchiveTypeImage:
		return DefaultCoverArchiveTypeImage
	case ArchiveTypeVideo:
		return DefaultCoverArchiveTypeVideo
	case ArchiveTypeAudio:
		return DefaultCoverArchiveTypeAudio
	case ArchiveTypeDocument:
		return DefaultCoverArchiveTypeDocument
	case ArchiveTypePDF:
		return DefaultCoverArchiveTypePDF
	case ArchiveTypeCompressed:
		return DefaultCoverArchiveTypeCompressed
	case ArchiveTypeExternalLink:
		return DefaultCoverArchiveTypeExternalLink
	case ArchiveTypeOther:
		return DefaultCoverArchiveTypeOther
	default:
		return DefaultCoverArchiveTypeOther
	}
}

const (
	DefaultCoverArchiveTypeImage        = "https://cdn.lilsite.com/default_folder_cover/file_coll_default_images_cover.png"
	DefaultCoverArchiveTypeVideo        = "https://cdn.lilsite.com/default_folder_cover/file_coll_default_videos_cover.png"
	DefaultCoverArchiveTypeAudio        = "https://cdn.lilsite.com/default_folder_cover/file_coll_default_media_cover.png"
	DefaultCoverArchiveTypeDocument     = "https://cdn.lilsite.com/default_folder_cover/file_coll_default_docs_cover.png"
	DefaultCoverArchiveTypePDF          = "https://cdn.lilsite.com/default_folder_cover/file_coll_default_pdf_cover.png"
	DefaultCoverArchiveTypeCompressed   = "https://cdn.lilsite.com/default_folder_cover/file_coll_default_zip_cover.png"
	DefaultCoverArchiveTypeExternalLink = "https://cdn.lilsite.com/default_folder_cover/file_coll_default_link_cover.png"
	DefaultCoverArchiveTypeOther        = "https://cdn.lilsite.com/default_folder_cover/file_coll_default_others_cover.png"
)

func FormatArchiveType(uploadType UploadType, fileExt string) (archiveType ArchiveType) {
	fileLowerExt := strings.ToLower(fileExt)
	switch fileLowerExt {
	case ".png", ".jpg", ".jpeg", ".webp", ".bmp", ".tif", ".tiff", ".raw", ".heic", ".heif", ".ai", ".svg", ".eps", ".ico", ".swf":
		archiveType = ArchiveTypeImage
	case ".mp4", ".avi", ".mkv", ".mov", ".mxf", ".wmv", ".flv", ".webm", ".ts", ".3gp", ".f4v", ".m3u8":
		archiveType = ArchiveTypeVideo
	case ".mp3", ".wav", ".aiff", ".pcm", ".flac", ".m4a", ".ape", ".aac", ".ogg", ".wma", ".opus":
		archiveType = ArchiveTypeAudio
	case ".doc", ".docx", ".odt", ".rtf", ".txt", ".xps", ".xls", ".xlsx", ".ppt", ".pptx", ".ods", ".md", ".tex", ".csv":
		archiveType = ArchiveTypeDocument
	case ".pdf":
		archiveType = ArchiveTypePDF
	case ".zip", ".rar", ".7z", ".gz", ".bz2", ".xz", ".tar", ".iso", ".tar.gz", ".tgz", ".tar.bz2", ".tar.xz", ".tar.zst", ".deb", ".rpm":
		archiveType = ArchiveTypeCompressed
	default:
		if uploadType == UploadTypeExternalAdd {
			archiveType = ArchiveTypeExternalLink
		} else {
			archiveType = ArchiveTypeOther
		}
		break
	}
	return archiveType
}
