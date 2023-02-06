package gonertia

import (
	"bytes"
	"io"
	"net/http"
)

// Middleware returns Inertia middleware handler.
//
// All of your handlers that can be handled by
// the Inertia should be under this middleware.
func (i *Inertia) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set header Vary to "X-Inertia".
		//
		// https://github.com/inertiajs/inertia-laravel/pull/404
		setInertiaVaryInResponse(w)

		if !IsInertiaRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		// Now we know that this request was made by Inertia.
		//
		// But there is one problem:
		// http.ResponseWriter has no methods for getting the response status code and response content.
		// So, we have to create our own response writer wrapper, that will contain that info.
		//
		// It's not critical that we will have a byte buffer, because we
		// know that Inertia response always in JSON format and actually not very big.
		w2 := buildInertiaResponseWrapper(w)

		// Now put our response writer wrapper to other handlers.
		next.ServeHTTP(w2, r)

		// Our response writer wrapper does have all needle data! Yuppy!
		//
		// Don't forget to copy all data to the original
		// response writer before end!
		defer func() {
			i.copyWrapperResponse(w, w2)
		}()

		// Determines what to do when the Inertia asset version has changed.
		// By default, we'll initiate a client-side location visit to force an update.
		//
		// https://inertiajs.com/asset-versioning
		if r.Method == http.MethodGet && inertiaVersionFromRequest(r) != i.version {
			setInertiaLocationInResponse(w2, r.URL.RequestURI())
			return
		}

		// Determines what to do when an Inertia action returned empty response.
		// By default, we will redirect the user back to where he came from.
		if w2.StatusCode() == http.StatusOK && w2.IsEmpty() {
			backURL := i.backURL(r)

			if backURL != "" {
				setInertiaLocationInResponse(w2, backURL)
				return
			}
		}

		// The POST, PUT and PATCH requests cannot have the 302 code status.
		// Let's set the status code to the 303 instead.
		//
		// https://inertiajs.com/redirects#303-response-code
		if w2.StatusCode() == http.StatusFound && isSeeOtherRedirectMethod(r.Method) {
			setResponseStatus(w2, http.StatusSeeOther)
		}
	})
}

func (i *Inertia) copyWrapperResponse(dst http.ResponseWriter, src *inertiaResponseWrapper) {
	i.copyWrapperHeaders(dst, src)
	i.copyWrapperStatusCode(dst, src)
	i.copyWrapperBuffer(dst, src)
}

func (i *Inertia) copyWrapperBuffer(dst http.ResponseWriter, src *inertiaResponseWrapper) {
	if _, err := io.Copy(dst, src.buf); err != nil {
		i.logger.Printf("cannot copy inertia response buffer to writer: %s", err)
	}
}

func (i *Inertia) copyWrapperStatusCode(dst http.ResponseWriter, src *inertiaResponseWrapper) {
	dst.WriteHeader(src.statusCode)
}

func (i *Inertia) copyWrapperHeaders(dst http.ResponseWriter, src *inertiaResponseWrapper) {
	for key, headers := range src.header {
		dst.Header().Del(key)

		for _, header := range headers {
			dst.Header().Add(key, header)
		}
	}
}

type inertiaResponseWrapper struct {
	statusCode int
	buf        *bytes.Buffer
	header     http.Header
}

var _ http.ResponseWriter = (*inertiaResponseWrapper)(nil)

func (w *inertiaResponseWrapper) StatusCode() int {
	return w.statusCode
}

func (w *inertiaResponseWrapper) IsEmpty() bool {
	return w.buf.Len() == 0
}

func (w *inertiaResponseWrapper) Header() http.Header {
	return w.header
}

func (w *inertiaResponseWrapper) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}

func (w *inertiaResponseWrapper) WriteHeader(code int) {
	w.statusCode = code
}

func buildInertiaResponseWrapper(w http.ResponseWriter) *inertiaResponseWrapper {
	w2 := &inertiaResponseWrapper{
		statusCode: http.StatusOK,
		buf:        bytes.NewBuffer(nil),
		header:     w.Header(),
	}

	// In some situations, we can pass a http.ResponseWriter,
	// that also implements this interface.
	if val, ok := w.(interface{ StatusCode() int }); ok {
		w2.statusCode = val.StatusCode()
	}

	return w2
}
