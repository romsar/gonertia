// Vite integration for Inertia.
package gonertia

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path"
	"strings"
)

// PreloadStrategy determines how dependencies are loaded.
type PreloadStrategy string

const (
	// PreloadNone loads only the main entry point with no preloading.
	PreloadNone PreloadStrategy = "none"
	// PreloadAggressive preloads all static dependencies immediately.
	PreloadAggressive PreloadStrategy = "aggressive"
	// PreloadWaterfall loads dependencies in controlled batches after page load.
	PreloadWaterfall PreloadStrategy = "waterfall"
)

// NonceGenerator generates CSP nonces.
type NonceGenerator func() string

// ViteConfig holds Vite configuration.
type ViteConfig struct {
	HotFile          string
	BuildManifest    string
	FallbackManifest string
	BuildDir         string
	HotReloadPort    string
	EmbedFS          fs.FS // Optional fs.FS (embed.FS or os.DirFS) for production builds
	UseEmbedFS       bool  // Whether to use embed.FS for manifest loading

	// Asset management configuration
	Nonce             string          // Static CSP nonce
	NonceGenerator    NonceGenerator  // Dynamic nonce generator (called once)
	IntegrityKey      *string         // Manifest key for SRI hashes (nil = disabled)
	EntryPoints       []string        // Entry points for asset generation
	PreloadStrategy   PreloadStrategy // How to handle dependency loading
	PreloadConcurrent int             // Concurrent prefetch count for waterfall
}

// ViteInstance wraps Inertia with Vite functionality.
type ViteInstance struct {
	*Inertia
	viteConfig ViteConfig
}

// ViteOption configures Vite.
type ViteOption func(*ViteConfig)

// WithHotFile sets the hot reload file path.
func WithHotFile(path string) ViteOption {
	return func(c *ViteConfig) {
		c.HotFile = path
	}
}

// WithBuildManifest sets the build manifest path.
func WithBuildManifest(path string) ViteOption {
	return func(c *ViteConfig) {
		c.BuildManifest = path
	}
}

// WithFallbackManifest sets the fallback manifest path.
func WithFallbackManifest(path string) ViteOption {
	return func(c *ViteConfig) {
		c.FallbackManifest = path
	}
}

// WithBuildDir sets the build directory.
func WithBuildDir(dir string) ViteOption {
	return func(c *ViteConfig) {
		c.BuildDir = dir
	}
}

// WithHotReloadPort sets the hot reload port.
func WithHotReloadPort(port string) ViteOption {
	return func(c *ViteConfig) {
		c.HotReloadPort = port
	}
}

// WithNonce sets a static CSP nonce.
func WithNonce(nonce string) ViteOption {
	return func(c *ViteConfig) {
		c.Nonce = nonce
	}
}

// WithAutoNonce generates a cryptographically secure nonce automatically.
func WithAutoNonce() ViteOption {
	return func(c *ViteConfig) {
		c.Nonce = generateCryptoNonce()
	}
}

// WithNonceGenerator sets a custom nonce generator function.
func WithNonceGenerator(gen NonceGenerator) ViteOption {
	return func(c *ViteConfig) {
		c.NonceGenerator = gen
	}
}

// WithIntegrity enables SubResource Integrity with the default manifest key "integrity".
func WithIntegrity() ViteOption {
	return func(c *ViteConfig) {
		key := "integrity"
		c.IntegrityKey = &key
	}
}

// WithIntegrityKey enables SubResource Integrity with a custom manifest key.
func WithIntegrityKey(key string) ViteOption {
	return func(c *ViteConfig) {
		c.IntegrityKey = &key
	}
}

// WithEntryPoints explicitly sets entry points to load.
func WithEntryPoints(entries ...string) ViteOption {
	return func(c *ViteConfig) {
		c.EntryPoints = entries
	}
}

// WithoutPreloading disables preloading. Browser handles module discovery naturally.
// This is the default behavior.
func WithoutPreloading() ViteOption {
	return withPreloadStrategy(PreloadNone, 0)
}

// WithAggressivePreload enables aggressive preloading of all dependencies.
// All JavaScript imports are preloaded immediately using modulepreload.
func WithAggressivePreload() ViteOption {
	return withPreloadStrategy(PreloadAggressive, 0)
}

// WithWaterfallPreload enables batched prefetch with concurrency control.
// Assets are loaded after page load in batches. concurrent controls how many
// assets load in parallel (default: 3 if not specified or <= 0).
func WithWaterfallPreload(concurrent int) ViteOption {
	return withPreloadStrategy(PreloadWaterfall, concurrent)
}

// withPreloadStrategy sets the preload strategy and concurrency.
// This is the internal implementation used by public helpers.
func withPreloadStrategy(strategy PreloadStrategy, concurrent int) ViteOption {
	return func(c *ViteConfig) {
		c.PreloadStrategy = strategy
		c.PreloadConcurrent = concurrent
	}
}

// generateCryptoNonce creates a cryptographically secure random nonce.
func generateCryptoNonce() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

// NewVite creates a Vite instance with the given Inertia instance.
// This uses the file system for hot reload detection and manifest loading.
func NewVite(i *Inertia, opts ...ViteOption) (*ViteInstance, error) {
	config := ViteConfig{
		HotFile:          "public/hot",
		BuildManifest:    "public/build/manifest.json",
		FallbackManifest: "public/build/.vite/manifest.json",
		BuildDir:         "/build/",
		HotReloadPort:    "//localhost:5173",
		UseEmbedFS:       false,

		// Default configuration
		PreloadStrategy:   PreloadNone,
		PreloadConcurrent: 3,
	}

	for _, opt := range opts {
		opt(&config)
	}

	vi := &ViteInstance{
		Inertia:    i,
		viteConfig: config,
	}

	if err := vi.setup(); err != nil {
		return nil, fmt.Errorf("setup vite: %w", err)
	}

	return vi, nil
}

// NewViteFromFS creates a Vite instance using fs.FS (embed.FS or os.DirFS) for production builds.
// It checks for hot reload file first (for development), otherwise uses the provided fs.FS.
func NewViteFromFS(i *Inertia, embedFS fs.FS, opts ...ViteOption) (*ViteInstance, error) {
	config := ViteConfig{
		HotFile:          "public/hot",
		BuildManifest:    "public/build/manifest.json",
		FallbackManifest: "public/build/.vite/manifest.json",
		BuildDir:         "/build/",
		HotReloadPort:    "//localhost:5173",
		EmbedFS:          embedFS,
		UseEmbedFS:       true,

		// Default configuration
		PreloadStrategy:   PreloadNone,
		PreloadConcurrent: 3,
	}

	for _, opt := range opts {
		opt(&config)
	}

	vi := &ViteInstance{
		Inertia:    i,
		viteConfig: config,
	}

	if err := vi.setup(); err != nil {
		return nil, fmt.Errorf("setup vite: %w", err)
	}

	return vi, nil
}

// NewWithVite is deprecated. Use NewVite instead.
func NewWithVite(i *Inertia, opts ...ViteOption) (*ViteInstance, error) {
	return NewVite(i, opts...)
}

func (vi *ViteInstance) setup() error {
	hotReload := vi.isHotReload()

	// Existing template functions (backward compatibility)
	if err := vi.ShareTemplateFunc("vite", vi.assetResolver(hotReload)); err != nil {
		return fmt.Errorf("share vite function: %w", err)
	}

	if err := vi.ShareTemplateFunc("viteReactRefresh", vi.reactRefreshHelper(hotReload)); err != nil {
		return fmt.Errorf("share vite react refresh function: %w", err)
	}

	if err := vi.ShareTemplateFunc("viteRefresh", vi.refreshHelper(hotReload)); err != nil {
		return fmt.Errorf("share vite refresh function: %w", err)
	}

	if err := vi.ShareTemplateFunc("viteAssets", vi.generateAllAssets); err != nil {
		return fmt.Errorf("share viteAssets function: %w", err)
	}

	vi.ShareTemplateData("hmr", hotReload)
	return nil
}

func (vi *ViteInstance) isHotReload() bool {
	_, err := os.Stat(vi.viteConfig.HotFile)
	return err == nil
}

func (vi *ViteInstance) assetResolver(hotReload bool) func(string) (string, error) {
	if hotReload {
		return vi.hotReloadResolver()
	}
	return vi.bundledResolver()
}

func (vi *ViteInstance) hotReloadResolver() func(string) (string, error) {
	return func(asset string) (string, error) {
		url := vi.readHotReloadURL()
		if asset != "" && !strings.HasPrefix(asset, "/") {
			asset = "/" + asset
		}
		return url + asset, nil
	}
}

func (vi *ViteInstance) readHotReloadURL() string {
	content, err := os.ReadFile(vi.viteConfig.HotFile)
	if err != nil {
		return vi.viteConfig.HotReloadPort
	}

	url := strings.TrimSpace(string(content))
	if url == "" {
		return vi.viteConfig.HotReloadPort
	}

	if strings.HasPrefix(url, "http://") {
		return "//" + url[7:]
	}
	if strings.HasPrefix(url, "https://") {
		return "//" + url[8:]
	}

	return url
}

func (vi *ViteInstance) bundledResolver() func(string) (string, error) {
	manifest, err := vi.loadManifest()
	if err != nil {
		return func(string) (string, error) {
			return "", fmt.Errorf("manifest error: %w", err)
		}
	}

	return func(asset string) (string, error) {
		entry, exists := manifest[asset]
		if !exists {
			return "", fmt.Errorf("asset %q not found", asset)
		}
		return path.Join(vi.viteConfig.BuildDir, entry.File), nil
	}
}

func (vi *ViteInstance) loadManifest() (map[string]Asset, error) {
	// If using embed.FS and not in hot reload mode, load from embed.FS
	if vi.viteConfig.UseEmbedFS && !vi.isHotReload() {
		return vi.loadManifestFromEmbed()
	}

	return vi.loadManifestFromFS()
}

func (vi *ViteInstance) loadManifestFromFS() (map[string]Asset, error) {
	manifestPath, err := vi.findManifest()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("open manifest: %w", err)
	}
	defer func() { _ = file.Close() }()

	var manifest map[string]Asset
	if err := json.NewDecoder(file).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("decode manifest: %w", err)
	}

	return manifest, nil
}

func (vi *ViteInstance) loadManifestFromEmbed() (map[string]Asset, error) {
	manifestPath, err := vi.findManifestInEmbed()
	if err != nil {
		return nil, err
	}

	file, err := vi.viteConfig.EmbedFS.Open(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("open manifest from embed: %w", err)
	}
	defer func() { _ = file.Close() }()

	var manifest map[string]Asset
	if err := json.NewDecoder(file).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("decode manifest: %w", err)
	}

	return manifest, nil
}

func (vi *ViteInstance) findManifestInEmbed() (string, error) {
	// Check primary manifest path
	if _, err := fs.Stat(vi.viteConfig.EmbedFS, vi.viteConfig.BuildManifest); err == nil {
		return vi.viteConfig.BuildManifest, nil
	}

	// Check fallback manifest path
	if _, err := fs.Stat(vi.viteConfig.EmbedFS, vi.viteConfig.FallbackManifest); err == nil {
		return vi.viteConfig.FallbackManifest, nil
	}

	return "", fmt.Errorf("manifest not found in embed.FS")
}

func (vi *ViteInstance) reactRefreshHelper(hotReload bool) func() template.HTML {
	return func() template.HTML {
		if !hotReload {
			return template.HTML("") // No React Refresh in production
		}

		viteClientURL, _ := vi.assetResolver(hotReload)("@vite/client")
		reactRefreshURL, _ := vi.assetResolver(hotReload)("@react-refresh")

		html := fmt.Sprintf(`<script type="module" src="%s"></script>
<script type="module">
    import RefreshRuntime from "%s"
    RefreshRuntime.injectIntoGlobalHook(window)
    window.$RefreshReg$ = () => {}
    window.$RefreshSig$ = () => (type) => type
    window.__vite_plugin_react_preamble_installed__ = true
</script>`, viteClientURL, reactRefreshURL)

		return template.HTML(html)
	}
}

func (vi *ViteInstance) refreshHelper(hotReload bool) func() template.HTML {
	return func() template.HTML {
		if !hotReload {
			return template.HTML("") // No refresh in production
		}

		viteClientURL, _ := vi.assetResolver(hotReload)("@vite/client")

		html := fmt.Sprintf(`<script type="module" src="%s"></script>`, viteClientURL)

		return template.HTML(html)
	}
}

func (vi *ViteInstance) findManifest() (string, error) {
	if _, err := os.Stat(vi.viteConfig.BuildManifest); err == nil {
		return vi.viteConfig.BuildManifest, nil
	}

	if _, err := os.Stat(vi.viteConfig.FallbackManifest); err == nil {
		if err := os.Rename(vi.viteConfig.FallbackManifest, vi.viteConfig.BuildManifest); err != nil {
			return "", fmt.Errorf("move manifest: %w", err)
		}
		return vi.viteConfig.BuildManifest, nil
	}

	return "", fmt.Errorf("manifest not found")
}

// Asset represents a Vite manifest entry.
type Asset struct {
	File           string   `json:"file"`                     // Hashed output file
	Src            string   `json:"src,omitempty"`            // Source file path
	Name           string   `json:"name,omitempty"`           // Asset name
	IsEntry        bool     `json:"isEntry,omitempty"`        // Main entry point
	IsDynamicEntry bool     `json:"isDynamicEntry,omitempty"` // Code-split chunk
	Imports        []string `json:"imports,omitempty"`        // Static dependencies
	DynamicImports []string `json:"dynamicImports,omitempty"` // Lazy imports
	Css            []string `json:"css,omitempty"`            // Associated CSS
	Integrity      string   `json:"integrity,omitempty"`      // SRI hash (e.g., from vite-plugin-manifest-sri)
}
