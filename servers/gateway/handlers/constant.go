package handlers

const (
	headerAccessControlAllowOrigin   = "Access-Control-Allow-Origin"
	headerAccessControlAllowMethods  = "Access-Control-Allow-Methods"
	headerAccessControlAllowHeaders  = "Access-Control-Allow-Headers"
	headerAccessControlExposeHeaders = "Access-Control-Expose-Headers"
	headerAccessControlMaxAge        = "Access-Control-Max-Age"
	contentTypeJSON                  = "application/json"
	contentTypeText                  = "text/plain"
	contentTypeHTML                  = "text/html"
	contentType                      = "Content-Type"
	contentAuth                      = "Authorization"
	charsetUTF8                      = "charset=utf-8"
	contentTypeJSONUTF8              = contentTypeJSON + "; " + charsetUTF8
	contentTypeTextUTF8              = contentTypeText + "; " + charsetUTF8
)
