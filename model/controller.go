package model

// mime types for file response
const (
	mimeContentTypeText       = "text/plain"
	mimeContentTypePdf        = "application/pdf"
	mimeContentTypeDoc        = "application/zip"
	mimeContentTypePowerPoint = "application/vnd.ms-powerpoint"
	mimeContentTypePng        = "image/png"
	mimeContentTypeJpeg       = "image/jpeg"
	mimeContentTypeJpg        = "image/jpg"
)

const (
	MimeTypeText = "text/plain"
	MimeTypeJSON = "application/json"
	MimeTypeZip  = "zip"
)

const (
	HeaderContentType = "Content-Type"
)

var ContentTypes = []string{
	mimeContentTypeText,
	mimeContentTypePdf,
	mimeContentTypeDoc,
	mimeContentTypePowerPoint,
	mimeContentTypePng,
	mimeContentTypeJpeg,
	mimeContentTypeJpg,
}

// MimeTypeToExt maps mime type to relative file extension
var MimeTypeToExt = map[string]string{
	mimeContentTypePdf:        "pdf",
	MimeTypeZip:               "zip",
	MimeTypeText:              "txt",
	mimeContentTypeDoc:        "docx",
	mimeContentTypePowerPoint: "pptx",
	mimeContentTypePng:        "png",
	mimeContentTypeJpeg:       "jpeg",
	mimeContentTypeJpg:        "jpg",
}

// CtxKey represents a key (string) for retrieving struct saved in request context.
type CtxKey struct {
	Key string
}

var (
	ClaimsCtxKey = CtxKey{Key: "ClaimsCtxKey"}
	CtxKeyBody   = CtxKey{Key: "body"}   // key used for retrieving body from request context.
	CtxKeyParams = CtxKey{Key: "params"} // key used for retrieving path and query params from request context.
	CtxKeyID     = CtxKey{Key: "id"}     // key used for retrieving the call id from request context.
)
