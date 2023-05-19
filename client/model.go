package client

type HttpHeader map[string]string
type HttpBody map[string]any

type HttpResponse struct {
	Status int
	Body   []byte
}
