package client

const (
	ErrResponseFailed = "response failed"
)

type HttpResponse struct {
	Status int
	Body   []byte
}
