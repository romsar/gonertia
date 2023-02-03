package gonertia

import (
	"net/http"
	"strings"
)

// IsInertiaRequest returns true if the request is an Inertia request.
func IsInertiaRequest(r *http.Request) bool {
	return r.Header.Get("X-Inertia") != ""
}

// setInertiaInResponse sets Inertia header in the response.
func setInertiaInResponse(w http.ResponseWriter) {
	w.Header().Set("X-Inertia", "true")
}

// setInertiaVaryInResponse sets Inertia Vary header in the response.
func setInertiaVaryInResponse(w http.ResponseWriter) {
	w.Header().Set("Vary", "X-Inertia")
}

// setInertiaLocationInResponse sets Inertia location header in the response.
func setInertiaLocationInResponse(w http.ResponseWriter, url string) {
	w.Header().Set("X-Inertia-Location", url)
	setResponseStatus(w, http.StatusConflict)
}

// setResponseStatus sets status code in the response.
func setResponseStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

// partialDataFromRequest returns Inertia partial data from request headers.
func partialDataFromRequest(r *http.Request) []string {
	header := r.Header.Get("X-Inertia-Partial-Data")
	if header == "" {
		return nil
	}

	return strings.Split(header, ",")
}

// partialComponentFromRequest returns Inertia partial component from request headers.
func partialComponentFromRequest(r *http.Request) string {
	return r.Header.Get("X-Inertia-Partial-Component")
}

// setResponseStatus returns Inertia version from request headers.
func inertiaVersionFromRequest(r *http.Request) string {
	return r.Header.Get("X-Inertia-Version")
}

// redirectResponse creates redirect response.
func redirectResponse(w http.ResponseWriter, r *http.Request, url string, status ...int) {
	http.Redirect(w, r, url, firstOr[int](status, http.StatusFound))
}

// markAsJSONResponse sets the type of response content in JSON format.
func markAsJSONResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

// markAsHTMLResponse sets the type of response content in HTML format.
func markAsHTMLResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
}

// isSeeOtherRedirectMethod returns HTTP methods that
// must have 303 response status, instead of 302.
func isSeeOtherRedirectMethod(method string) bool {
	return method == http.MethodPost || method == http.MethodPatch || method == http.MethodPut
}

// refererFromRequest returns referer header from request.
func refererFromRequest(r *http.Request) string {
	return r.Referer()
}
