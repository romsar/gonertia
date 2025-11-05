package gonertia

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCSPMiddlewareDefault(t *testing.T) {
	t.Parallel()

	t.Run("default configuration", func(t *testing.T) {
		t.Parallel()

		i := &Inertia{}
		vi := &ViteInstance{Inertia: i}

		handler := vi.CSPMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		csp := rec.Header().Get("Content-Security-Policy")
		if csp == "" {
			t.Fatal("CSP header not set")
		}

		if !strings.Contains(csp, "script-src 'nonce-") {
			t.Error("CSP should contain script-src with nonce")
		}

		if !strings.Contains(csp, "style-src 'nonce-") {
			t.Error("CSP should contain style-src with nonce")
		}

		if !strings.Contains(csp, "object-src 'none'") {
			t.Error("CSP should contain object-src 'none'")
		}
	})
}

func TestCSPMiddlewareCustom(t *testing.T) {
	t.Parallel()

	t.Run("custom policy", func(t *testing.T) {
		t.Parallel()

		i := &Inertia{}
		vi := &ViteInstance{Inertia: i}

		handler := vi.CSPMiddleware(
			WithCSPPolicy("script-src 'nonce-%s'; default-src 'self'"),
		)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		csp := rec.Header().Get("Content-Security-Policy")
		if !strings.Contains(csp, "default-src 'self'") {
			t.Error("CSP should contain custom policy")
		}
	})

	t.Run("custom nonce generator", func(t *testing.T) {
		t.Parallel()

		i := &Inertia{}
		vi := &ViteInstance{Inertia: i}

		called := false
		customGen := func() string {
			called = true
			return "custom-test-nonce"
		}

		handler := vi.CSPMiddleware(
			WithCSPNonceGenerator(customGen),
		)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if !called {
			t.Error("custom nonce generator should be called")
		}

		csp := rec.Header().Get("Content-Security-Policy")
		if !strings.Contains(csp, "custom-test-nonce") {
			t.Error("CSP should contain custom nonce")
		}
	})

	t.Run("custom nonce key", func(t *testing.T) {
		t.Parallel()

		i := &Inertia{}
		vi := &ViteInstance{Inertia: i}

		var capturedNonce string
		handler := vi.CSPMiddleware(
			WithCSPNonceKey("my_custom_nonce"),
		)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := TemplateDataFromContext(r.Context())
			if val, ok := data["my_custom_nonce"]; ok {
				capturedNonce = val.(string)
			}
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if capturedNonce == "" {
			t.Error("custom nonce key should be set in template data")
		}
	})
}

func TestCSPMiddlewareMerge(t *testing.T) {
	t.Parallel()

	t.Run("merges with existing CSP header", func(t *testing.T) {
		t.Parallel()

		i := &Inertia{}
		vi := &ViteInstance{Inertia: i}

		handler := vi.CSPMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		rec.Header().Set("Content-Security-Policy", "frame-ancestors 'none'")

		handler.ServeHTTP(rec, req)

		csp := rec.Header().Get("Content-Security-Policy")
		if !strings.Contains(csp, "frame-ancestors 'none'") {
			t.Error("should preserve existing CSP directive")
		}

		if !strings.Contains(csp, "script-src 'nonce-") {
			t.Error("should add new CSP directives")
		}
	})
}

func TestCSPMiddlewareContext(t *testing.T) {
	t.Parallel()

	t.Run("nonce is available in context", func(t *testing.T) {
		t.Parallel()

		i := &Inertia{}
		vi := &ViteInstance{Inertia: i}

		var contextNonce string
		handler := vi.CSPMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := TemplateDataFromContext(r.Context())
			if val, ok := data["csp_nonce"]; ok {
				contextNonce = val.(string)
			}
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if contextNonce == "" {
			t.Error("nonce should be available in context")
		}

		csp := rec.Header().Get("Content-Security-Policy")
		if !strings.Contains(csp, contextNonce) {
			t.Error("nonce in context should match nonce in CSP header")
		}
	})

	t.Run("generates unique nonces per request", func(t *testing.T) {
		t.Parallel()

		i := &Inertia{}
		vi := &ViteInstance{Inertia: i}

		var nonce1, nonce2 string
		handler := vi.CSPMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req1 := httptest.NewRequest(http.MethodGet, "/", nil)
		rec1 := httptest.NewRecorder()
		handler.ServeHTTP(rec1, req1)
		csp1 := rec1.Header().Get("Content-Security-Policy")
		parts1 := strings.Split(csp1, "'nonce-")
		if len(parts1) > 1 {
			nonce1 = strings.Split(parts1[1], "'")[0]
		}

		req2 := httptest.NewRequest(http.MethodGet, "/", nil)
		rec2 := httptest.NewRecorder()
		handler.ServeHTTP(rec2, req2)
		csp2 := rec2.Header().Get("Content-Security-Policy")
		parts2 := strings.Split(csp2, "'nonce-")
		if len(parts2) > 1 {
			nonce2 = strings.Split(parts2[1], "'")[0]
		}

		if nonce1 == "" || nonce2 == "" {
			t.Fatal("nonces should be generated")
		}

		if nonce1 == nonce2 {
			t.Error("nonces should be unique per request")
		}
	})
}

func TestCSPOptions(t *testing.T) {
	t.Parallel()

	t.Run("WithCSPPolicy", func(t *testing.T) {
		t.Parallel()

		cfg := &cspConfig{}
		WithCSPPolicy("custom policy")(cfg)

		if cfg.policy != "custom policy" {
			t.Errorf("expected 'custom policy', got %s", cfg.policy)
		}
	})

	t.Run("WithCSPNonceGenerator", func(t *testing.T) {
		t.Parallel()

		customGen := func() string { return "test" }
		cfg := &cspConfig{}
		WithCSPNonceGenerator(customGen)(cfg)

		if cfg.nonceGenerator == nil {
			t.Error("nonceGenerator should be set")
		}

		if cfg.nonceGenerator() != "test" {
			t.Error("nonceGenerator should return expected value")
		}
	})

	t.Run("WithCSPNonceKey", func(t *testing.T) {
		t.Parallel()

		cfg := &cspConfig{}
		WithCSPNonceKey("custom_key")(cfg)

		if cfg.nonceKey != "custom_key" {
			t.Errorf("expected 'custom_key', got %s", cfg.nonceKey)
		}
	})
}
