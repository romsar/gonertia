package gonertia

import (
	"bytes"
	"io"
	"net/http"
)

// Middleware returns http.Handler with Inertia support.
//
// All of your handlers that need to be handled by the Inertia
// should be under this middleware.
func (i *Inertia) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set header Vary to "X-Inertia".
		//
		// https://github.com/inertiajs/inertia-laravel/pull/404
		setInertiaVaryToResponse(w)

		// If request is not Inertia request, we can move next.
		if !IsInertiaRequest(r) {
			next.ServeHTTP(w, r)

			return
		}

		// Now we know that it's Inertia request.
		//
		// But there is one problem:
		// http.ResponseWriter doesn't have methods for getting status code and response content.
		// So we have to create our own response writer, that will contain that info.
		//
		// It's not critical that we now have a byte buffer, because we
		// know that Inertia response has JSON format and usually not very big.
		w2 := buildInertiaResponseWriter(w)

		// Now put our response writer to other handlers.
		next.ServeHTTP(w2, r)

		// Now, our response writer does have all needle data! Yuppy!
		//
		// Don't forget to copy all data to the original
		// response writer before end!
		defer func() {
			i.copyHeaders(w, w2)
			i.copyStatusCode(w, w2)
			i.copyBuffer(w, w2)
		}()

		// Determines what to do when the Inertia asset version has changed.
		// By default, we'll initiate a client-side location visit to force an update.
		//
		// https://inertiajs.com/asset-versioning
		if r.Method == http.MethodGet && inertiaVersionFromRequest(r) != i.version {
			i.Location(w2, r, i.url+r.RequestURI)
			return
		}

		// Determines what to do when an Inertia action returned with no response.
		// By default, we'll redirect the user back to where they came from.
		if w2.StatusCode() == http.StatusOK && w2.IsEmpty() {
			backURL := i.backURL(r)

			if backURL != "" {
				redirectResponse(w, r, backURL)
				return
			}
		}

		// The POST, PUT and PATCH requests cannot have the status 302.
		// Let's set the status code to 303 instead.
		//
		// https://inertiajs.com/redirects#303-response-code
		if w2.StatusCode() == http.StatusFound && isSeeOtherRedirectMethod(w2.Method()) {
			setResponseStatus(w2, http.StatusSeeOther)
		}
	})
}

// copyBuffer copying source bytes buf into destination bytes buffer.
func (i *Inertia) copyBuffer(dst http.ResponseWriter, src *inertiaResponseWriter) {
	if _, err := io.Copy(dst, src.buf); err != nil {
		i.logger.Printf("cannot copy inertia response buffer to writer: %s", err)
	}
}

// copyStatusCode copying source status code into destination status code.
func (i *Inertia) copyStatusCode(dst http.ResponseWriter, src *inertiaResponseWriter) {
	dst.WriteHeader(src.statusCode)
}

// copyHeaders copying source header into destination header.
func (i *Inertia) copyHeaders(dst http.ResponseWriter, src *inertiaResponseWriter) {
	for key, headers := range src.header {
		for _, header := range headers {
			dst.Header().Add(key, header)
		}
	}
}

// inertiaResponseWriter is the implementation of http.ResponseWriter,
// that have response body buffer and status code that will be return to client.
type inertiaResponseWriter struct {
	method     string
	statusCode int
	buf        *bytes.Buffer
	header     http.Header
}

var _ http.ResponseWriter = (*inertiaResponseWriter)(nil)

// Method returns HTTP method of response.
func (w *inertiaResponseWriter) Method() string {
	return w.method
}

// StatusCode returns HTTP status code of response.
func (w *inertiaResponseWriter) StatusCode() int {
	return w.statusCode
}

// IsEmpty returns true is response body is empty.
func (w *inertiaResponseWriter) IsEmpty() bool {
	return w.buf.Len() == 0
}

// Header returns response headers.
func (w *inertiaResponseWriter) Header() http.Header {
	return w.header
}

// Write pushes some bytes to response body.
func (w *inertiaResponseWriter) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

// WriteHeader sets the status code of the response.
func (w *inertiaResponseWriter) WriteHeader(code int) {
	w.statusCode = code
}

// buildInertiaResponseWriter initializes inertiaResponseWriter.
func buildInertiaResponseWriter(w http.ResponseWriter) *inertiaResponseWriter {
	w2 := &inertiaResponseWriter{
		statusCode: http.StatusOK,
		buf:        bytes.NewBuffer(nil),
		header:     w.Header(),
	}

	// In some situations, we can pass a http.ResponseWriter,
	// that also implements these interfaces.
	if val, ok := w.(interface{ StatusCode() int }); ok {
		w2.statusCode = val.StatusCode()
	}
	if val, ok := w.(interface{ Header() http.Header }); ok {
		w2.header = val.Header()
	}
	if val, ok := w.(interface{ Buffer() *bytes.Buffer }); ok {
		w2.buf = val.Buffer()
	}

	return w2
}
