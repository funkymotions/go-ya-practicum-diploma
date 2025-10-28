package interfaces

import "net/http"

type RESTClient interface {
	Get(path string) (*http.Response, error)
}
