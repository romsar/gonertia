// Package vite provides Vite integration for Inertia.
package vite

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"

	inertia "github.com/romsar/gonertia/v2"
)

// Config holds Vite configuration.
type Config struct {
	HotFile          string
	BuildManifest    string
	FallbackManifest string
	BuildDir         string
	HotReloadPort    string
}

// Instance wraps Inertia with Vite functionality.
type Instance struct {
	*inertia.Inertia
	config Config
}

// Option configures Vite.
type Option func(*Config)

// WithHotFile sets the hot reload file path.
func WithHotFile(path string) Option {
	return func(c *Config) {
		c.HotFile = path
	}
}

// WithBuildManifest sets the build manifest path.
func WithBuildManifest(path string) Option {
	return func(c *Config) {
		c.BuildManifest = path
	}
}

// WithFallbackManifest sets the fallback manifest path.
func WithFallbackManifest(path string) Option {
	return func(c *Config) {
		c.FallbackManifest = path
	}
}

// WithBuildDir sets the build directory.
func WithBuildDir(dir string) Option {
	return func(c *Config) {
		c.BuildDir = dir
	}
}

// WithHotReloadPort sets the hot reload port.
func WithHotReloadPort(port string) Option {
	return func(c *Config) {
		c.HotReloadPort = port
	}
}

// New creates a Vite instance with the given Inertia instance.
func New(i *inertia.Inertia, opts ...Option) (*Instance, error) {
	config := Config{
		HotFile:          "public/hot",
		BuildManifest:    "public/build/manifest.json",
		FallbackManifest: "public/build/.vite/manifest.json",
		BuildDir:         "/build/",
		HotReloadPort:    "//localhost:5173",
	}

	for _, opt := range opts {
		opt(&config)
	}

	vi := &Instance{
		Inertia: i,
		config:  config,
	}

	if err := vi.setup(); err != nil {
		return nil, fmt.Errorf("setup vite: %w", err)
	}

	return vi, nil
}

func (vi *Instance) setup() error {
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

func (vi *Instance) isHotReload() bool {
	_, err := os.Stat(vi.config.HotFile)
	return err == nil
}

func (vi *Instance) assetResolver(hotReload bool) func(string) (string, error) {
	if hotReload {
		return vi.hotReloadResolver()
	}
	return vi.bundledResolver()
}

func (vi *Instance) hotReloadResolver() func(string) (string, error) {
	return func(asset string) (string, error) {
		url, _ := vi.readHotReloadURL()
		if asset != "" && !strings.HasPrefix(asset, "/") {
			asset = "/" + asset
		}
		return url + asset, nil
	}
}

func (vi *Instance) readHotReloadURL() (string, error) {
	content, err := os.ReadFile(vi.config.HotFile)
	if err != nil {
		return vi.config.HotReloadPort, nil
	}

	url := strings.TrimSpace(string(content))
	if url == "" {
		return vi.config.HotReloadPort, nil
	}

	if strings.HasPrefix(url, "http://") {
		return "//" + url[7:], nil
	}
	if strings.HasPrefix(url, "https://") {
		return "//" + url[8:], nil
	}

	return url, nil
}

func (vi *Instance) bundledResolver() func(string) (string, error) {
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
		return path.Join(vi.config.BuildDir, entry.File), nil
	}
}

func (vi *Instance) loadManifest() (map[string]Asset, error) {
	manifestPath, err := vi.findManifest()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("open manifest: %w", err)
	}
	defer file.Close()

	var manifest map[string]Asset
	if err := json.NewDecoder(file).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("decode manifest: %w", err)
	}

	return manifest, nil
}

func (vi *Instance) reactRefreshHelper(hotReload bool) func() template.HTML {
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

func (vi *Instance) findManifest() (string, error) {
	if _, err := os.Stat(vi.config.BuildManifest); err == nil {
		return vi.config.BuildManifest, nil
	}

	if _, err := os.Stat(vi.config.FallbackManifest); err == nil {
		if err := os.Rename(vi.config.FallbackManifest, vi.config.BuildManifest); err != nil {
			return "", fmt.Errorf("move manifest: %w", err)
		}
		return vi.config.BuildManifest, nil
	}

	return "", fmt.Errorf("manifest not found")
}

// Asset represents a Vite manifest entry.
type Asset struct {
	File string `json:"file"`
	Src  string `json:"src"`
}