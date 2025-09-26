package vite

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/romsar/gonertia/v2"
)

const rootTemplate = `<html><head>{{ .inertiaHead }}</head><body>{{ .inertia }}</body></html>`

func TestNew(t *testing.T) {
	i, err := gonertia.New(rootTemplate)
	if err != nil {
		t.Fatalf("failed to create inertia: %v", err)
	}

	vi, err := New(i)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	if vi == nil {
		t.Fatal("New() returned nil")
		return
	}

	if vi.Inertia != i {
		t.Error("Instance should embed the provided Inertia instance")
	}
}

func TestNewWithOptions(t *testing.T) {
	i, err := gonertia.New(rootTemplate)
	if err != nil {
		t.Fatalf("failed to create inertia: %v", err)
	}

	vi, err := New(i,
		WithHotFile("custom/hot"),
		WithBuildManifest("custom/manifest.json"),
		WithBuildDir("/custom/"),
		WithHotReloadPort("//localhost:3000"),
	)
	if err != nil {
		t.Fatalf("New() with options failed: %v", err)
	}

	config := vi.config
	if config.HotFile != "custom/hot" {
		t.Errorf("HotFile = %q, want %q", config.HotFile, "custom/hot")
	}
	if config.BuildManifest != "custom/manifest.json" {
		t.Errorf("BuildManifest = %q, want %q", config.BuildManifest, "custom/manifest.json")
	}
	if config.BuildDir != "/custom/" {
		t.Errorf("BuildDir = %q, want %q", config.BuildDir, "/custom/")
	}
	if config.HotReloadPort != "//localhost:3000" {
		t.Errorf("HotReloadPort = %q, want %q", config.HotReloadPort, "//localhost:3000")
	}
}

func TestIsHotReload(t *testing.T) {
	tmpDir := t.TempDir()
	hotFile := filepath.Join(tmpDir, "hot")

	i, _ := gonertia.New(rootTemplate)
	vi, _ := New(i, WithHotFile(hotFile))

	// No hot file - should be bundled mode
	if vi.isHotReload() {
		t.Error("isHotReload() = true, want false when hot file doesn't exist")
	}

	// Create hot file - should be hot reload mode
	if err := os.WriteFile(hotFile, []byte("//localhost:5173"), 0o644); err != nil {
		t.Fatalf("failed to create hot file: %v", err)
	}

	if !vi.isHotReload() {
		t.Error("isHotReload() = false, want true when hot file exists")
	}
}

func TestReadHotReloadURL(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		createFile bool
		expected   string
	}{
		{"no file", "", false, "//localhost:5173"},
		{"empty file", "", true, "//localhost:5173"},
		{"simple url", "//localhost:3000", true, "//localhost:3000"},
		{"http url", "http://localhost:3000", true, "//localhost:3000"},
		{"https url", "https://localhost:3000", true, "//localhost:3000"},
		{"whitespace", "  //localhost:3000  ", true, "//localhost:3000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			hotFile := filepath.Join(tmpDir, "hot")

			i, _ := gonertia.New(rootTemplate)
			vi, _ := New(i, WithHotFile(hotFile))

			if tt.createFile {
				if err := os.WriteFile(hotFile, []byte(tt.content), 0o644); err != nil {
					t.Fatalf("failed to create hot file: %v", err)
				}
			}

			url := vi.readHotReloadURL()

			if url != tt.expected {
				t.Errorf("readHotReloadURL() = %q, want %q", url, tt.expected)
			}
		})
	}
}

func TestHotReloadResolver(t *testing.T) {
	tmpDir := t.TempDir()
	hotFile := filepath.Join(tmpDir, "hot")

	// Create hot file with custom URL
	if err := os.WriteFile(hotFile, []byte("//localhost:3000"), 0o644); err != nil {
		t.Fatalf("failed to create hot file: %v", err)
	}

	i, _ := gonertia.New(rootTemplate)
	vi, _ := New(i, WithHotFile(hotFile))

	resolver := vi.hotReloadResolver()

	tests := []struct {
		asset    string
		expected string
	}{
		{"", "//localhost:3000"},
		{"app.js", "//localhost:3000/app.js"},
		{"/app.js", "//localhost:3000/app.js"},
		{"@vite/client", "//localhost:3000/@vite/client"},
	}

	for _, tt := range tests {
		t.Run(tt.asset, func(t *testing.T) {
			url, err := resolver(tt.asset)
			if err != nil {
				t.Fatalf("resolver(%q) error: %v", tt.asset, err)
			}
			if url != tt.expected {
				t.Errorf("resolver(%q) = %q, want %q", tt.asset, url, tt.expected)
			}
		})
	}
}

func TestBundledResolver(t *testing.T) {
	tmpDir := t.TempDir()
	manifestFile := filepath.Join(tmpDir, "manifest.json")

	// Create manifest file
	manifest := `{
		"app.js": {"file": "assets/app.abc123.js"},
		"main.css": {"file": "assets/main.def456.css"}
	}`
	if err := os.WriteFile(manifestFile, []byte(manifest), 0o644); err != nil {
		t.Fatalf("failed to create manifest: %v", err)
	}

	i, _ := gonertia.New(rootTemplate)
	vi, _ := New(i, WithBuildManifest(manifestFile), WithBuildDir("/build/"))

	resolver := vi.bundledResolver()

	tests := []struct {
		asset       string
		expected    string
		expectError bool
	}{
		{"app.js", "/build/assets/app.abc123.js", false},
		{"main.css", "/build/assets/main.def456.css", false},
		{"missing.js", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.asset, func(t *testing.T) {
			url, err := resolver(tt.asset)

			if tt.expectError {
				if err == nil {
					t.Errorf("resolver(%q) expected error but got none", tt.asset)
				}
				return
			}

			if err != nil {
				t.Fatalf("resolver(%q) unexpected error: %v", tt.asset, err)
			}
			if url != tt.expected {
				t.Errorf("resolver(%q) = %q, want %q", tt.asset, url, tt.expected)
			}
		})
	}
}

func TestFindManifest(t *testing.T) {
	tmpDir := t.TempDir()
	buildManifest := filepath.Join(tmpDir, "manifest.json")
	fallbackManifest := filepath.Join(tmpDir, "fallback.json")

	tests := []struct {
		name           string
		createBuild    bool
		createFallback bool
		expectError    bool
		expectMove     bool
	}{
		{"build exists", true, false, false, false},
		{"fallback exists", false, true, false, true},
		{"both exist", true, true, false, false}, // build takes priority
		{"neither exists", false, false, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vi := setupManifestTest(t, buildManifest, fallbackManifest, tt.createBuild, tt.createFallback)
			path, err := vi.findManifest()

			if tt.expectError {
				assertManifestError(t, err)
				return
			}

			assertManifestSuccess(t, err, path, buildManifest)
			if tt.expectMove {
				assertManifestMove(t, buildManifest, fallbackManifest)
			}
		})
	}
}

func setupManifestTest(t *testing.T, buildManifest, fallbackManifest string, createBuild, createFallback bool) *Instance {
	_ = os.Remove(buildManifest)
	_ = os.Remove(fallbackManifest)

	if createBuild {
		_ = os.WriteFile(buildManifest, []byte("{}"), 0o644)
	}
	if createFallback {
		_ = os.WriteFile(fallbackManifest, []byte("{}"), 0o644)
	}

	i, _ := gonertia.New(rootTemplate)
	vi, _ := New(i,
		WithBuildManifest(buildManifest),
		WithFallbackManifest(fallbackManifest),
	)
	return vi
}

func assertManifestError(t *testing.T, err error) {
	if err == nil {
		t.Error("findManifest() expected error but got none")
	}
}

func assertManifestSuccess(t *testing.T, err error, path, buildManifest string) {
	if err != nil {
		t.Fatalf("findManifest() unexpected error: %v", err)
	}
	if path != buildManifest {
		t.Errorf("findManifest() = %q, want %q", path, buildManifest)
	}
}

func assertManifestMove(t *testing.T, buildManifest, fallbackManifest string) {
	if _, err := os.Stat(buildManifest); err != nil {
		t.Error("fallback manifest should have been moved to build location")
	}
	if _, err := os.Stat(fallbackManifest); err == nil {
		t.Error("fallback manifest should have been moved away")
	}
}

func TestAssetResolverIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	hotFile := filepath.Join(tmpDir, "hot")
	manifestFile := filepath.Join(tmpDir, "manifest.json")

	// Test hot reload mode
	t.Run("hot reload mode", func(t *testing.T) {
		if err := os.WriteFile(hotFile, []byte("//localhost:3000"), 0o644); err != nil {
			t.Fatalf("failed to create hot file: %v", err)
		}

		i, _ := gonertia.New(rootTemplate)
		vi, _ := New(i, WithHotFile(hotFile))

		resolver := vi.assetResolver(vi.isHotReload())
		url, err := resolver("app.js")
		if err != nil {
			t.Fatalf("resolver error: %v", err)
		}

		expected := "//localhost:3000/app.js"
		if url != expected {
			t.Errorf("hot reload resolver = %q, want %q", url, expected)
		}
	})

	// Test bundled mode
	t.Run("bundled mode", func(t *testing.T) {
		_ = os.Remove(hotFile) // Remove hot file to trigger bundled mode

		manifest := `{"app.js": {"file": "assets/app.abc123.js"}}`
		if err := os.WriteFile(manifestFile, []byte(manifest), 0o644); err != nil {
			t.Fatalf("failed to create manifest: %v", err)
		}

		i, _ := gonertia.New(rootTemplate)
		vi, _ := New(i,
			WithHotFile(hotFile),
			WithBuildManifest(manifestFile),
		)

		resolver := vi.assetResolver(vi.isHotReload())
		url, err := resolver("app.js")
		if err != nil {
			t.Fatalf("resolver error: %v", err)
		}

		expected := "/build/assets/app.abc123.js"
		if url != expected {
			t.Errorf("bundled resolver = %q, want %q", url, expected)
		}
	})
}

func TestSetupAddsViteFunction(t *testing.T) {
	i, _ := gonertia.New(rootTemplate)
	vi, err := New(i)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// The vite function should be available in templates
	// We can't easily test this without accessing internals,
	// but we can verify setup doesn't error
	if err := vi.setup(); err != nil {
		t.Errorf("setup() failed: %v", err)
	}
}

func TestInvalidManifest(t *testing.T) {
	tmpDir := t.TempDir()
	manifestFile := filepath.Join(tmpDir, "manifest.json")

	// Create invalid JSON
	if err := os.WriteFile(manifestFile, []byte("invalid json"), 0o644); err != nil {
		t.Fatalf("failed to create invalid manifest: %v", err)
	}

	i, _ := gonertia.New(rootTemplate)
	vi, _ := New(i, WithBuildManifest(manifestFile))

	_, err := vi.loadManifest()
	if err == nil {
		t.Error("loadManifest() should fail with invalid JSON")
	}
	if !strings.Contains(err.Error(), "decode manifest") {
		t.Errorf("expected decode error, got: %v", err)
	}
}

func TestViteReactRefresh(t *testing.T) {
	tmpDir := t.TempDir()
	hotFile := filepath.Join(tmpDir, "hot")

	i, _ := gonertia.New(rootTemplate)

	tests := []struct {
		name        string
		createHot   bool
		expectEmpty bool
	}{
		{
			name:        "hot reload mode - should output React refresh setup",
			createHot:   true,
			expectEmpty: false,
		},
		{
			name:        "production mode - should output empty string",
			createHot:   false,
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vi := setupViteInstance(t, i, hotFile, tt.createHot)
			resultStr := getReactRefreshResult(vi)

			if tt.expectEmpty {
				assertEmptyResult(t, resultStr)
			} else {
				assertValidReactRefresh(t, resultStr)
			}
		})
	}
}

func setupViteInstance(t *testing.T, i *gonertia.Inertia, hotFile string, createHot bool) *Instance {
	_ = os.Remove(hotFile)

	if createHot {
		if err := os.WriteFile(hotFile, []byte("//localhost:5173"), 0o644); err != nil {
			t.Fatalf("failed to create hot file: %v", err)
		}
	}

	vi, err := New(i, WithHotFile(hotFile))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	return vi
}

func getReactRefreshResult(vi *Instance) string {
	helper := vi.reactRefreshHelper(vi.isHotReload())
	result := helper()
	return string(result)
}

func assertEmptyResult(t *testing.T, resultStr string) {
	if resultStr != "" {
		t.Errorf("expected empty string in production mode, got: %q", resultStr)
	}
}

func assertValidReactRefresh(t *testing.T, resultStr string) {
	if resultStr == "" {
		t.Error("expected non-empty React refresh setup in development mode")
		return
	}

	expectedParts := []string{
		"@vite/client",
		"@react-refresh",
		"RefreshRuntime.injectIntoGlobalHook",
		"window.$RefreshReg$",
		"window.$RefreshSig$",
		"__vite_plugin_react_preamble_installed__",
	}

	for _, part := range expectedParts {
		if !strings.Contains(resultStr, part) {
			t.Errorf("React refresh setup missing expected part: %q", part)
		}
	}

	if !strings.Contains(resultStr, "<script type=\"module\"") {
		t.Error("React refresh setup should contain script tags")
	}
}

func TestViteReactRefreshIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	hotFile := filepath.Join(tmpDir, "hot")

	// Create hot file for development mode
	if err := os.WriteFile(hotFile, []byte("//localhost:3000"), 0o644); err != nil {
		t.Fatalf("failed to create hot file: %v", err)
	}

	i, _ := gonertia.New(rootTemplate)
	vi, err := New(i, WithHotFile(hotFile))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	// The viteReactRefresh function should be available in shared template functions
	// We can't directly test this without accessing internals, but we can verify
	// the function doesn't error during setup
	if err := vi.setup(); err != nil {
		t.Errorf("setup() with React refresh failed: %v", err)
	}

	// Test that URLs are correctly resolved in the helper
	helper := vi.reactRefreshHelper(true) // Force development mode
	result := helper()
	resultStr := string(result)

	// Should contain the correct dev server URLs
	if !strings.Contains(resultStr, "//localhost:3000/@vite/client") {
		t.Error("React refresh setup should contain correct Vite client URL")
	}
	if !strings.Contains(resultStr, "//localhost:3000/@react-refresh") {
		t.Error("React refresh setup should contain correct React refresh URL")
	}
}
