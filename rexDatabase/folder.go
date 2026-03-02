package rexDatabase

const (
	DefaultUploadDirTypeImage  = "images"
	DefaultUploadDirVideo      = "videos"
	DefaultUploadDirAudio      = "audios"
	DefaultUploadDirDocument   = "documents"
	DefaultUploadDirPDF        = "pdfs"
	DefaultUploadDirCompressed = "compressed"
	DefaultUploadDirLinks      = "links"
	DefaultUploadDirOther      = "other"
)

type CreatorType string

const (
	CreatorTypeSystem CreatorType = "system"
	CreatorTypeTenant CreatorType = "tenant"
)
