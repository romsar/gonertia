package gonertia

import (
	"net/http"
	"strings"
)

const (
	headerInertia                 = "X-Inertia"
	headerInertiaLocation         = "X-Inertia-Location"
	headerInertiaPartialData      = "X-Inertia-Partial-Data"
	headerInertiaPartialExcept    = "X-Inertia-Partial-Except"
	headerInertiaPartialComponent = "X-Inertia-Partial-Component"
	headerInertiaVersion          = "X-Inertia-Version"
	headerInertiaReset            = "X-Inertia-Reset"
	headerVary                    = "Vary"
	headerContentType             = "Content-Type"
)

// IsInertiaRequest returns true if the request is an Inertia request.
func IsInertiaRequest(r *http.Request) bool {
	return r.Header.Get(headerInertia) != ""
}

func setInertiaInResponse(w http.ResponseWriter) {
	w.Header().Set(headerInertia, "true")
}

func deleteInertiaInResponse(w http.ResponseWriter) {
	w.Header().Del(headerInertia)
}

func setInertiaVaryInResponse(w http.ResponseWriter) {
	w.Header().Set(headerVary, headerInertia)
}

func deleteVaryInResponse(w http.ResponseWriter) {
	w.Header().Del(headerVary)
}

func setInertiaLocationInResponse(w http.ResponseWriter, url string) {
	w.Header().Set(headerInertiaLocation, url)
}

func setResponseStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func onlyFromRequest(r *http.Request) []string {
	header := r.Header.Get(headerInertiaPartialData)
	if header == "" {
		return nil
	}

	return strings.Split(header, ",")
}

func exceptFromRequest(r *http.Request) []string {
	header := r.Header.Get(headerInertiaPartialExcept)
	if header == "" {
		return nil
	}

	return strings.Split(header, ",")
}

func resetFromRequest(r *http.Request) []string {
	header := r.Header.Get(headerInertiaReset)
	if header == "" {
		return nil
	}

	return strings.Split(header, ",")
}

func partialComponentFromRequest(r *http.Request) string {
	return r.Header.Get(headerInertiaPartialComponent)
}

func inertiaVersionFromRequest(r *http.Request) string {
	return r.Header.Get(headerInertiaVersion)
}

func redirectResponse(w http.ResponseWriter, r *http.Request, url string, status ...int) {
	http.Redirect(w, r, url, firstOr[int](status, http.StatusFound))
}

func setJSONResponse(w http.ResponseWriter) {
	w.Header().Set(headerContentType, "application/json")
}

func setJSONRequest(r *http.Request) {
	r.Header.Set(headerContentType, "application/json")
}

func setHTMLResponse(w http.ResponseWriter) {
	w.Header().Set(headerContentType, "text/html")
}

func isSeeOtherRedirectMethod(method string) bool {
	return method == http.MethodPut || method == http.MethodPatch || method == http.MethodDelete
}

func refererFromRequest(r *http.Request) string {
	return r.Referer()
}
