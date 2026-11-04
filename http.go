package gonertia

import (
	"net/http"
	"strings"
)

// Headers used by the Inertia protocol, so you don't have to retype them
// when you deal with Inertia requests yourself (custom middleware, tests…).
//
// https://inertiajs.com/the-protocol
const (
	HeaderInertia                          = "X-Inertia"
	HeaderInertiaLocation                  = "X-Inertia-Location"
	HeaderInertiaPartialData               = "X-Inertia-Partial-Data"
	HeaderInertiaPartialExcept             = "X-Inertia-Partial-Except"
	HeaderInertiaPartialComponent          = "X-Inertia-Partial-Component"
	HeaderInertiaVersion                   = "X-Inertia-Version"
	HeaderInertiaReset                     = "X-Inertia-Reset"
	HeaderInertiaInfiniteScrollMergeIntent = "X-Inertia-Infinite-Scroll-Merge-Intent"
)

const (
	headerVary        = "Vary"
	headerContentType = "Content-Type"
)

// IsInertiaRequest returns true if the request is an Inertia request.
func IsInertiaRequest(r *http.Request) bool {
	return r.Header.Get(HeaderInertia) != ""
}

func setInertiaInResponse(w http.ResponseWriter) {
	w.Header().Set(HeaderInertia, "true")
}

func deleteInertiaInResponse(w http.ResponseWriter) {
	w.Header().Del(HeaderInertia)
}

func setInertiaVaryInResponse(w http.ResponseWriter) {
	w.Header().Set(headerVary, HeaderInertia)
}

func deleteVaryInResponse(w http.ResponseWriter) {
	w.Header().Del(headerVary)
}

func setInertiaLocationInResponse(w http.ResponseWriter, url string) {
	w.Header().Set(HeaderInertiaLocation, url)
}

func setResponseStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func onlyFromRequest(r *http.Request) []string {
	header := r.Header.Get(HeaderInertiaPartialData)
	if header == "" {
		return nil
	}

	return strings.Split(header, ",")
}

func exceptFromRequest(r *http.Request) []string {
	header := r.Header.Get(HeaderInertiaPartialExcept)
	if header == "" {
		return nil
	}

	return strings.Split(header, ",")
}

func resetFromRequest(r *http.Request) []string {
	header := r.Header.Get(HeaderInertiaReset)
	if header == "" {
		return nil
	}

	return strings.Split(header, ",")
}

func partialComponentFromRequest(r *http.Request) string {
	return r.Header.Get(HeaderInertiaPartialComponent)
}

func inertiaVersionFromRequest(r *http.Request) string {
	return r.Header.Get(HeaderInertiaVersion)
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

func infiniteScrollMergeIntentFromRequest(r *http.Request) string {
	return r.Header.Get(HeaderInertiaInfiniteScrollMergeIntent)
}
