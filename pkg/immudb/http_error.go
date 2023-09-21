package immudb

import (
	"fmt"
	"io"
	"net/http"
)

type HttpError struct {
	method, path, status, body string
	statusCode                 int
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("%s %s failed - status: %s - body: %s", e.method, e.path, e.status, e.body)
}

func HttpStatusCode(err error) int {
	if e, ok := err.(*HttpError); ok {
		return e.statusCode
	}
	return -1
}

func newHttpError(method, path string, resp *http.Response) error {
	return &HttpError{
		method:     method,
		path:       path,
		status:     resp.Status,
		body:       readAll(resp.Body),
		statusCode: resp.StatusCode,
	}
}

func readAll(data io.Reader) string {
	if data, err := io.ReadAll(data); err == nil {
		return string(data)
	} else {
		return fmt.Sprintf("<error: %s>", err)
	}
}
