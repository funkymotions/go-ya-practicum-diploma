package utils

import "net/http"

type HTTPContentType string

const (
	ContentTypeJSON  HTTPContentType = "application/json"
	ContentTypePlain HTTPContentType = "text/plain"
)

func CheckRequestContentType(r *http.Request, validContentType HTTPContentType) bool {
	contentType := r.Header.Get("Content-Type")
	return HTTPContentType(contentType) == validContentType
}
