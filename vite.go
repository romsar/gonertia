// Vite integration for Inertia.
package gonertia

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"
)

// ViteConfig holds Vite configuration.
type ViteConfig struct {
	HotFile          string
	BuildManifest    string
	FallbackManifest string
	BuildDir         string
	HotReloadPort    string
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

// NewWithVite creates a Vite instance with the given Inertia instance.
func NewWithVite(i *Inertia, opts ...ViteOption) (*ViteInstance, error) {
	config := ViteConfig{
		HotFile:          "public/hot",
		BuildManifest:    "public/build/manifest.json",
		FallbackManifest: "public/build/.vite/manifest.json",
		BuildDir:         "/build/",
		HotReloadPort:    "//localhost:5173",
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

func (vi *ViteInstance) setup() error {
	hotReload := vi.isHotReload()

	if err := vi.ShareTemplateFunc("vite", vi.assetResolver(hotReload)); err != nil {
		return fmt.Errorf("share vite function: %w", err)
	}

	if err := vi.ShareTemplateFunc("viteReactRefresh", vi.reactRefreshHelper(hotReload)); err != nil {
		return fmt.Errorf("share vite react refresh function: %w", err)
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
	File string `json:"file"`
	Src  string `json:"src"`
}
