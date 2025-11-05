package gonertia

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
)

type cspConfig struct {
	nonceGenerator func() string
	policy         string
	nonceKey       string
}

// CSPOption configures CSP middleware.
type CSPOption func(*cspConfig)

// WithCSPNonceGenerator sets a custom nonce generator function.
func WithCSPNonceGenerator(gen func() string) CSPOption {
	return func(c *cspConfig) {
		c.nonceGenerator = gen
	}
}

// WithCSPPolicy sets a custom CSP policy template.
// Use {{nonce}} placeholders where nonces should be inserted.
func WithCSPPolicy(policy string) CSPOption {
	return func(c *cspConfig) {
		c.policy = policy
	}
}

// WithCSPNonceKey sets the template data key for the nonce.
func WithCSPNonceKey(key string) CSPOption {
	return func(c *cspConfig) {
		c.nonceKey = key
	}
}

// CSPMiddleware returns middleware that generates per-request CSP nonces.
// Nonces are added to the Content-Security-Policy header and made available
// in templates via the configured key (default: "csp_nonce").
func (vi *ViteInstance) CSPMiddleware(opts ...CSPOption) func(http.Handler) http.Handler {
	config := &cspConfig{
		nonceGenerator: generateCryptoNonce,
		policy: "script-src 'nonce-{{nonce}}' 'strict-dynamic'; " +
			"style-src 'nonce-{{nonce}}'; " +
			"font-src 'self' data:; " +
			"img-src 'self' data: https:; " +
			"object-src 'none'; " +
			"base-uri 'none'",
		nonceKey: "csp_nonce",
	}

	for _, opt := range opts {
		opt(config)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nonce := config.nonceGenerator()
			newCSP := strings.ReplaceAll(config.policy, "{{nonce}}", nonce)

			existing := w.Header().Get("Content-Security-Policy")
			if existing != "" {
				w.Header().Set("Content-Security-Policy", existing+"; "+newCSP)
			} else {
				w.Header().Set("Content-Security-Policy", newCSP)
			}

			ctx := SetTemplateDatum(r.Context(), config.nonceKey, nonce)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func generateCryptoNonce() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}
