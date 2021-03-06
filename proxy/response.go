package proxy

import (
	"io"
	"mime"
	"net/http"
)

type ResponseReader struct {
	io.Reader
	r *http.Response
}

func newResponseReader(r *http.Response) *ResponseReader {
	return &ResponseReader{
		Reader: r.Body,
		r:      r,
	}
}

func (r *ResponseReader) ContentType() string {
	cth := r.Header().Get("Content-Type")
	ct, _, _ := mime.ParseMediaType(cth)
	return ct
}

func (r *ResponseReader) Header() http.Header {
	return r.r.Header
}

func (r *ResponseReader) Request() *http.Request {
	return r.r.Request
}

type ResponseWriter struct {
	io.Writer
	rw          http.ResponseWriter
	statusCode  int
	headersDone bool
}

func newResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		Writer: w,
		rw:     w,
	}
}

func (w *ResponseWriter) takeHeaders(r *ResponseReader) {
	for k, v := range r.Header() {
		for _, v := range v {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(r.r.StatusCode)
}

func (w *ResponseWriter) Header() http.Header {
	return w.rw.Header()
}

func (w *ResponseWriter) WriteHeader(s int) {
	w.statusCode = s
}

func (w *ResponseWriter) ReadFrom(r io.Reader) error {
	w.flushHeaders()
	_, err := io.Copy(w.Writer, r)
	return err
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.flushHeaders()
	return w.Writer.Write(b)
}

func (w *ResponseWriter) flushHeaders() {
	if w.headersDone {
		return
	}
	w.rw.WriteHeader(w.statusCode)
	w.headersDone = true
}

func (w *ResponseWriter) setChunked() {
	w.Header().Del("Content-Length")
}
