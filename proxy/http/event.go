package http

import (
	"net/http"
	"time"
)

type Interaction struct {
	Request *Request
	Response *Response
}

type Request struct {
	StartTime time.Time
	EndTime time.Time

	Method string
	Address string
	Headers http.Header
	Body []byte
}

type Response struct {
	StartTime time.Time
	EndTime time.Time

	StatusCode int
	Message string
	Headers http.Header
	Body []byte
}