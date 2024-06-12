package gonertia

import (
	"net/http"
	"strings"
)

// IsInertiaRequest returns true if the request is an Inertia request.
func IsInertiaRequest(r *http.Request) bool {
	return r.Header.Get("X-Inertia") != ""
}

func setInertiaInResponse(w http.ResponseWriter) {
	w.Header().Set("X-Inertia", "true")
}

func setInertiaVaryInResponse(w http.ResponseWriter) {
	w.Header().Set("Vary", "X-Inertia")
}

func setInertiaLocationInResponse(w http.ResponseWriter, url string) {
	w.Header().Set("X-Inertia-Location", url)
	setResponseStatus(w, http.StatusConflict)
}

func setResponseStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func onlyFromRequest(r *http.Request) []string {
	header := r.Header.Get("X-Inertia-Partial-Data")
	if header == "" {
		return nil
	}

	return strings.Split(header, ",")
}

func exceptFromRequest(r *http.Request) []string {
	header := r.Header.Get("X-Inertia-Partial-Except")
	if header == "" {
		return nil
	}

	return strings.Split(header, ",")
}

func partialComponentFromRequest(r *http.Request) string {
	return r.Header.Get("X-Inertia-Partial-Component")
}

func inertiaVersionFromRequest(r *http.Request) string {
	return r.Header.Get("X-Inertia-Version")
}

func redirectResponse(w http.ResponseWriter, r *http.Request, url string, status ...int) {
	http.Redirect(w, r, url, firstOr[int](status, http.StatusFound))
}

func setJSONResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func setHTMLResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
}

func isSeeOtherRedirectMethod(method string) bool {
	return method == http.MethodPost || method == http.MethodPatch || method == http.MethodPut
}

func refererFromRequest(r *http.Request) string {
	return r.Referer()
}
