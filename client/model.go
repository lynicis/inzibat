package client

const ErrorResponseFailed = "response failed"

type HttpResponse struct {
	Status int
	Body   []byte
}
