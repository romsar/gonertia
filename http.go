package gonertia

import (
	"net/http"
	"strings"
)

// IsInertiaRequest returns true if this request is Inertia request.
func IsInertiaRequest(r *http.Request) bool {
	return r.Header.Get("X-Inertia") != ""
}

// markAsInertiaResponse tells browser that it's Inertia response.
func markAsInertiaResponse(w http.ResponseWriter) {
	w.Header().Set("X-Inertia", "true")
}

// setInertiaVaryToResponse sets the "Vary" response header to "X-Inertia".
func setInertiaVaryToResponse(w http.ResponseWriter) {
	w.Header().Set("Vary", "X-Inertia")
}

// setInertiaLocationToResponse sets Inertia location header to response.
func setInertiaLocationToResponse(w http.ResponseWriter, url string) {
	w.Header().Set("X-Inertia-Location", url)
	setResponseStatus(w, http.StatusConflict)
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

// setResponseStatus sets status code for response.
func setResponseStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

// redirectResponse creates redirect response.
func redirectResponse(w http.ResponseWriter, r *http.Request, url string, status ...int) {
	http.Redirect(w, r, url, firstOr[int](status, http.StatusFound))
}

// markAsJSONResponse tells browser that response is the JSON response.
func markAsJSONResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

// markAsHTMLResponse tells browser that response is the HTML response.
func markAsHTMLResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
}

// isSeeOtherRedirectMethod returns HTTP methods that
// need to have 303 response status, instead of 302.
func isSeeOtherRedirectMethod(method string) bool {
	return method == http.MethodPost || method == http.MethodPatch || method == http.MethodPut
}

// refererFromRequest returns Referer header from request.
func refererFromRequest(r *http.Request) string {
	return r.Referer()
}
