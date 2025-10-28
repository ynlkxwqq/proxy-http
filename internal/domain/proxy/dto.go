package proxy

import (
	"errors"
	"net/http"
	"net/url"
)

type Request struct {
	Method        string            `json:"method" example:"GET" validate:"required"`
	URL           string            `json:"url" example:"http://google.com" validate:"required"`
	Headers       map[string]string `json:"headers" example:"Content-Type:application/json,Authorization:Bearer token" validate:"required"`
	ParsedURL     *url.URL          `swaggerignore:"true"`
	ParsedHeaders http.Header       `swaggerignore:"true"`
} // @name Request

type Response struct {
	RequestID  string            `json:"id" example:"1"`
	StatusCode int               `json:"status" example:"200"`
	Headers    map[string]string `json:"headers" example:"Content-Type:application/json"`
	Length     int64             `json:"length" example:"100"`
} // @name Response

func (req *Request) Bind(r *http.Request) (err error) {
	if req.Method == "" {
		return errors.New("method is required")
	}

	if req.URL == "" {
		return errors.New("url is required")
	}

	if len(req.Headers) == 0 {
		return errors.New("headers is required")
	}

	parsedURL, err := url.ParseRequestURI(req.URL)
	if err != nil {
		return errors.New("invalid url")
	}

	req.ParsedURL = parsedURL

	parsedHeaders := make(http.Header)
	for header, values := range req.Headers {
		parsedHeaders.Add(header, values)
	}

	req.ParsedHeaders = parsedHeaders

	return
}
