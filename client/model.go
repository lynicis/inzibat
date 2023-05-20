package client

type HttpHeader map[string]string

type HttpResponse struct {
	Status int
	Body   []byte
}
