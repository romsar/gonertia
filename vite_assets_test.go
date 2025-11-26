package gonertia

import (
	"encoding/json"
	"html/template"
	"os"
	"strings"
	"testing"
)

func TestEntryPointHandling(t *testing.T) {
	t.Parallel()

	t.Run("template args take precedence", func(t *testing.T) {
		t.Parallel()

		hotFile := t.TempDir() + "/hot"
		if err := os.WriteFile(hotFile, []byte("http://localhost:5173"), 0o644); err != nil {
			t.Fatal(err)
		}

		i := &Inertia{}
		vi := &ViteInstance{
			Inertia: i,
			viteConfig: ViteConfig{
				HotFile:       hotFile,
				EntryPoints:   []string{"config-entry.js"}, // Should be ignored
				HotReloadPort: "http://localhost:5173",
			},
		}

		// Template args should override config
		html, err := vi.buildHotAssets("", "template-arg.js")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		htmlStr := string(html)
		if !strings.Contains(htmlStr, "template-arg.js") {
			t.Error("should use template arg entry point")
		}

		if strings.Contains(htmlStr, "config-entry.js") {
			t.Error("should not use config entry when template arg provided")
		}
	})

	t.Run("config entries used when no template args", func(t *testing.T) {
		t.Parallel()

		hotFile := t.TempDir() + "/hot"
		if err := os.WriteFile(hotFile, []byte("http://localhost:5173"), 0o644); err != nil {
			t.Fatal(err)
		}

		i := &Inertia{}
		vi := &ViteInstance{
			Inertia: i,
			viteConfig: ViteConfig{
				HotFile:       hotFile,
				EntryPoints:   []string{"app.js"},
				HotReloadPort: "http://localhost:5173",
			},
		}

		html, err := vi.buildHotAssets("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		htmlStr := string(html)
		if !strings.Contains(htmlStr, "app.js") {
			t.Error("should use config entry point")
		}
	})

	t.Run("error when no entries configured", func(t *testing.T) {
		t.Parallel()

		hotFile := t.TempDir() + "/hot"
		if err := os.WriteFile(hotFile, []byte("http://localhost:5173"), 0o644); err != nil {
			t.Fatal(err)
		}

		i := &Inertia{}
		vi := &ViteInstance{
			Inertia: i,
			viteConfig: ViteConfig{
				HotFile:       hotFile,
				HotReloadPort: "http://localhost:5173",
			},
		}

		_, err := vi.buildHotAssets("")
		if err == nil {
			t.Error("expected error when no entries configured")
		}

		if !strings.Contains(err.Error(), "no entry points configured") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("multiple template args", func(t *testing.T) {
		t.Parallel()

		hotFile := t.TempDir() + "/hot"
		if err := os.WriteFile(hotFile, []byte("http://localhost:5173"), 0o644); err != nil {
			t.Fatal(err)
		}

		i := &Inertia{}
		vi := &ViteInstance{
			Inertia: i,
			viteConfig: ViteConfig{
				HotFile:       hotFile,
				HotReloadPort: "http://localhost:5173",
			},
		}

		html, err := vi.buildHotAssets("", "app.js", "admin.js")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		htmlStr := string(html)
		if !strings.Contains(htmlStr, "app.js") {
			t.Error("should include first entry")
		}

		if !strings.Contains(htmlStr, "admin.js") {
			t.Error("should include second entry")
		}
	})
}

func TestGeneratePreloadTag(t *testing.T) {
	t.Parallel()

	t.Run("basic preload tag", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{
			viteConfig: ViteConfig{
				BuildDir: "/build/",
			},
		}

		asset := Asset{
			File: "assets/app-abc123.js",
		}

		tag := vi.generatePreloadTag(asset, "")
		tagStr := string(tag)

		if !strings.Contains(tagStr, `rel="modulepreload"`) {
			t.Error("tag should contain rel=modulepreload")
		}

		if !strings.Contains(tagStr, `href="/build/assets/app-abc123.js"`) {
			t.Error("tag should contain correct href")
		}

		if !strings.Contains(tagStr, `crossorigin`) {
			t.Error("tag should contain crossorigin")
		}
	})

	t.Run("preload tag with nonce", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{
			viteConfig: ViteConfig{
				BuildDir: "/build/",
			},
		}

		asset := Asset{
			File: "assets/app-abc123.js",
		}

		tag := vi.generatePreloadTag(asset, "test-nonce")
		tagStr := string(tag)

		if !strings.Contains(tagStr, `nonce="test-nonce"`) {
			t.Errorf("tag should contain nonce, got: %s", tagStr)
		}
	})
}

func TestGenerateStylesheetTag(t *testing.T) {
	t.Parallel()

	t.Run("basic stylesheet tag", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{
			viteConfig: ViteConfig{
				BuildDir: "/build/",
			},
		}

		asset := Asset{
			File: "assets/app-abc123.css",
		}

		tag := vi.generateStylesheetTag(asset, "")
		tagStr := string(tag)

		if !strings.Contains(tagStr, `rel="stylesheet"`) {
			t.Error("tag should contain rel=stylesheet")
		}

		if !strings.Contains(tagStr, `href="/build/assets/app-abc123.css"`) {
			t.Error("tag should contain correct href")
		}
	})

	t.Run("stylesheet tag with nonce", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{
			viteConfig: ViteConfig{
				BuildDir: "/build/",
			},
		}

		asset := Asset{
			File: "assets/app-abc123.css",
		}

		tag := vi.generateStylesheetTag(asset, "test-nonce")
		tagStr := string(tag)

		if !strings.Contains(tagStr, `nonce="test-nonce"`) {
			t.Errorf("tag should contain nonce, got: %s", tagStr)
		}
	})
}

func TestGenerateScriptTag(t *testing.T) {
	t.Parallel()

	t.Run("basic script tag", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{
			viteConfig: ViteConfig{
				BuildDir: "/build/",
			},
		}

		asset := Asset{
			File: "assets/app-abc123.js",
		}

		tag := vi.generateScriptTag(asset, "")
		tagStr := string(tag)

		if !strings.Contains(tagStr, `type="module"`) {
			t.Error("tag should contain type=module")
		}

		if !strings.Contains(tagStr, `src="/build/assets/app-abc123.js"`) {
			t.Error("tag should contain correct src")
		}

		if !strings.HasPrefix(tagStr, "<script") || !strings.HasSuffix(tagStr, "</script>") {
			t.Error("tag should be properly formatted script tag")
		}
	})

	t.Run("script tag with nonce", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{
			viteConfig: ViteConfig{
				BuildDir: "/build/",
			},
		}

		asset := Asset{
			File: "assets/app-abc123.js",
		}

		tag := vi.generateScriptTag(asset, "test-nonce")
		tagStr := string(tag)

		if !strings.Contains(tagStr, `nonce="test-nonce"`) {
			t.Errorf("tag should contain nonce, got: %s", tagStr)
		}
	})
}

func TestGetIntegrity(t *testing.T) {
	t.Parallel()

	t.Run("returns integrity hash when present", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{}
		asset := Asset{
			File:      "assets/app-abc123.js",
			Integrity: "sha384-oqVuAfXRKap7fdgcCY5uykM6+R9GqQ8K/uxy9rx7HNQlGYl1kPzQho1wx4JwY8wC",
		}

		integrity := vi.getIntegrity(asset)
		expected := "sha384-oqVuAfXRKap7fdgcCY5uykM6+R9GqQ8K/uxy9rx7HNQlGYl1kPzQho1wx4JwY8wC"

		if integrity != expected {
			t.Errorf("expected %s, got %s", expected, integrity)
		}
	})

	t.Run("returns empty string when no integrity", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{}
		asset := Asset{
			File: "assets/app-abc123.js",
		}

		integrity := vi.getIntegrity(asset)
		if integrity != "" {
			t.Errorf("expected empty string, got %s", integrity)
		}
	})
}

func TestSRIInTags(t *testing.T) {
	t.Parallel()

	t.Run("preload tag includes integrity attribute", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{
			viteConfig: ViteConfig{
				BuildDir: "/build/",
			},
		}

		asset := Asset{
			File:      "assets/app-abc123.js",
			Integrity: "sha384-oqVuAfXRKap7fdgcCY5uykM6+R9GqQ8K/uxy9rx7HNQlGYl1kPzQho1wx4JwY8wC",
		}

		tag := vi.generatePreloadTag(asset, "")
		tagStr := string(tag)

		if !strings.Contains(tagStr, `integrity="sha384-oqVuAfXRKap7fdgcCY5uykM6+R9GqQ8K/uxy9rx7HNQlGYl1kPzQho1wx4JwY8wC"`) {
			t.Errorf("tag should contain integrity attribute, got: %s", tagStr)
		}
	})

	t.Run("stylesheet tag includes integrity attribute", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{
			viteConfig: ViteConfig{
				BuildDir: "/build/",
			},
		}

		asset := Asset{
			File:      "assets/app-abc123.css",
			Integrity: "sha384-xyzABC123",
		}

		tag := vi.generateStylesheetTag(asset, "")
		tagStr := string(tag)

		if !strings.Contains(tagStr, `integrity="sha384-xyzABC123"`) {
			t.Errorf("tag should contain integrity attribute, got: %s", tagStr)
		}
	})

	t.Run("script tag includes integrity attribute", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{
			viteConfig: ViteConfig{
				BuildDir: "/build/",
			},
		}

		asset := Asset{
			File:      "assets/app-abc123.js",
			Integrity: "sha384-xyzABC123",
		}

		tag := vi.generateScriptTag(asset, "")
		tagStr := string(tag)

		if !strings.Contains(tagStr, `integrity="sha384-xyzABC123"`) {
			t.Errorf("tag should contain integrity attribute, got: %s", tagStr)
		}
	})

	t.Run("tags work with both nonce and integrity", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{
			viteConfig: ViteConfig{
				BuildDir: "/build/",
			},
		}

		asset := Asset{
			File:      "assets/app-abc123.js",
			Integrity: "sha384-xyzABC123",
		}

		tag := vi.generateScriptTag(asset, "test-nonce")
		tagStr := string(tag)

		if !strings.Contains(tagStr, `integrity="sha384-xyzABC123"`) {
			t.Error("tag should contain integrity attribute")
		}

		if !strings.Contains(tagStr, `nonce="test-nonce"`) {
			t.Error("tag should contain nonce attribute")
		}
	})
}

func TestGenerateHotAssets(t *testing.T) {
	t.Parallel()

	t.Run("generates hot reload tags", func(t *testing.T) {
		t.Parallel()

		// Create temp hot file
		hotFile := t.TempDir() + "/hot"
		if err := os.WriteFile(hotFile, []byte("//localhost:5173"), 0o644); err != nil {
			t.Fatal(err)
		}

		i := &Inertia{}
		vi := &ViteInstance{
			Inertia: i,
			viteConfig: ViteConfig{
				HotFile:       hotFile,
				EntryPoints:   []string{"resources/js/app.tsx"},
				HotReloadPort: "//localhost:5173",
			},
		}

		html, err := vi.buildHotAssets("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		htmlStr := string(html)

		// Should contain Vite client
		if !strings.Contains(htmlStr, "/@vite/client") {
			t.Error("should contain Vite client script")
		}

		// Should contain entry point
		if !strings.Contains(htmlStr, "resources/js/app.tsx") {
			t.Error("should contain entry point script")
		}

		// Should be script tags
		if !strings.Contains(htmlStr, `<script type="module"`) {
			t.Error("should contain script tags")
		}
	})
}

func TestGenerateWaterfallScript(t *testing.T) {
	t.Parallel()

	t.Run("generates waterfall script", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{
			viteConfig: ViteConfig{
				PreloadConcurrent: 3,
			},
		}

		assets := []string{
			"/build/assets/chunk1.js",
			"/build/assets/chunk2.js",
			"/build/assets/chunk3.js",
		}

		script := vi.generateWaterfallScript(assets, "")
		scriptStr := string(script)

		// Should be a script tag
		if !strings.HasPrefix(scriptStr, "<script") || !strings.HasSuffix(scriptStr, "</script>") {
			t.Error("should be a script tag")
		}

		// Should contain load event listener
		if !strings.Contains(scriptStr, "window.addEventListener('load'") {
			t.Error("should contain load event listener")
		}

		// Should contain loadNext function
		if !strings.Contains(scriptStr, "loadNext") {
			t.Error("should contain loadNext function")
		}

		// Should contain all assets
		for _, asset := range assets {
			if !strings.Contains(scriptStr, asset) {
				t.Errorf("should contain asset: %s", asset)
			}
		}

		// Should contain concurrent value
		if !strings.Contains(scriptStr, "loadNext(assets, 3)") {
			t.Error("should contain concurrent value")
		}
	})

	t.Run("includes nonce in script tag", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{
			viteConfig: ViteConfig{
				PreloadConcurrent: 2,
			},
		}

		assets := []string{"/build/assets/chunk1.js"}

		script := vi.generateWaterfallScript(assets, "test-nonce-456")
		scriptStr := string(script)

		if !strings.Contains(scriptStr, `nonce="test-nonce-456"`) {
			t.Error("should contain nonce attribute")
		}
	})

	t.Run("returns empty for no assets", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{
			viteConfig: ViteConfig{},
		}

		script := vi.generateWaterfallScript([]string{}, "")
		if script != "" {
			t.Error("should return empty for no assets")
		}
	})

	t.Run("escapes quotes in URLs", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{
			viteConfig: ViteConfig{
				PreloadConcurrent: 1,
			},
		}

		assets := []string{`/build/assets/"quoted".js`}

		script := vi.generateWaterfallScript(assets, "")
		scriptStr := string(script)

		// Should escape the quotes
		if !strings.Contains(scriptStr, `\"quoted\"`) {
			t.Error("should escape quotes in URLs")
		}
	})
}

func TestProcessAsset(t *testing.T) {
	t.Parallel()

	t.Run("processes entry asset with PreloadNone", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{
			viteConfig: ViteConfig{
				BuildDir:        "/build/",
				PreloadStrategy: PreloadNone,
			},
		}

		manifest := map[string]Asset{
			"resources/js/app.tsx": {
				File:    "assets/app-abc.js",
				IsEntry: true,
			},
		}

		assets := &viteAssets{
			processedAssets: make(map[string]bool),
		}

		err := vi.processAsset("resources/js/app.tsx", manifest, assets, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should have script tag for entry
		if len(assets.scriptTags) != 1 {
			t.Errorf("expected 1 script tag, got %d", len(assets.scriptTags))
		}

		// Should have no preload tags (PreloadNone)
		if len(assets.preloadTags) != 0 {
			t.Errorf("expected 0 preload tags, got %d", len(assets.preloadTags))
		}
	})

	t.Run("processes entry asset with PreloadAggressive", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{
			viteConfig: ViteConfig{
				BuildDir:        "/build/",
				PreloadStrategy: PreloadAggressive,
			},
		}

		manifest := map[string]Asset{
			"resources/js/app.tsx": {
				File:    "assets/app-abc.js",
				IsEntry: true,
				Imports: []string{"resources/js/helper.ts"},
			},
			"resources/js/helper.ts": {
				File: "assets/helper-def.js",
			},
		}

		assets := &viteAssets{
			processedAssets: make(map[string]bool),
		}

		err := vi.processAsset("resources/js/app.tsx", manifest, assets, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should have preload tags for both entry and import
		if len(assets.preloadTags) != 2 {
			t.Errorf("expected 2 preload tags, got %d", len(assets.preloadTags))
		}

		// Should have script tag only for entry
		if len(assets.scriptTags) != 1 {
			t.Errorf("expected 1 script tag, got %d", len(assets.scriptTags))
		}
	})

	t.Run("processes CSS files", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{
			viteConfig: ViteConfig{
				BuildDir:        "/build/",
				PreloadStrategy: PreloadNone,
			},
		}

		manifest := map[string]Asset{
			"resources/js/app.tsx": {
				File:    "assets/app-abc.js",
				IsEntry: true,
				Css:     []string{"resources/css/app.css"},
			},
			"resources/css/app.css": {
				File: "assets/app-xyz.css",
			},
		}

		assets := &viteAssets{
			processedAssets: make(map[string]bool),
		}

		err := vi.processAsset("resources/js/app.tsx", manifest, assets, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should have stylesheet tag
		if len(assets.stylesheetTags) != 1 {
			t.Errorf("expected 1 stylesheet tag, got %d", len(assets.stylesheetTags))
		}
	})

	t.Run("prevents duplicate processing", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{
			viteConfig: ViteConfig{
				BuildDir:        "/build/",
				PreloadStrategy: PreloadAggressive,
			},
		}

		manifest := map[string]Asset{
			"resources/js/app.tsx": {
				File:    "assets/app-abc.js",
				IsEntry: true,
				Imports: []string{"resources/js/shared.ts"},
			},
			"resources/js/admin.tsx": {
				File:    "assets/admin-def.js",
				IsEntry: true,
				Imports: []string{"resources/js/shared.ts"},
			},
			"resources/js/shared.ts": {
				File: "assets/shared-ghi.js",
			},
		}

		assets := &viteAssets{
			processedAssets: make(map[string]bool),
		}

		// Process both entries
		_ = vi.processAsset("resources/js/app.tsx", manifest, assets, "")
		_ = vi.processAsset("resources/js/admin.tsx", manifest, assets, "")

		// Shared module should only be preloaded once
		sharedCount := 0
		for _, tag := range assets.preloadTags {
			if strings.Contains(string(tag), "shared-ghi.js") {
				sharedCount++
			}
		}

		if sharedCount != 1 {
			t.Errorf("expected shared module to be preloaded once, got %d times", sharedCount)
		}
	})

	t.Run("error for missing asset", func(t *testing.T) {
		t.Parallel()

		vi := &ViteInstance{
			viteConfig: ViteConfig{},
		}

		manifest := map[string]Asset{}

		assets := &viteAssets{
			processedAssets: make(map[string]bool),
		}

		err := vi.processAsset("missing.js", manifest, assets, "")
		if err == nil {
			t.Error("expected error for missing asset")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("unexpected error message: %v", err)
		}
	})
}

func TestGenerateCryptoNonce(t *testing.T) {
	t.Parallel()

	nonce1 := generateCryptoNonce()
	nonce2 := generateCryptoNonce()

	// Should generate non-empty nonces
	if nonce1 == "" || nonce2 == "" {
		t.Error("should generate non-empty nonces")
	}

	// Should generate different nonces
	if nonce1 == nonce2 {
		t.Error("should generate different nonces")
	}

	// Should be base64 encoded (no special characters except +, /, =)
	for _, c := range nonce1 {
		valid := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '+' || c == '/' || c == '='
		if !valid {
			t.Errorf("nonce contains invalid character: %c", c)
		}
	}
}

func TestViteAssetsHTMLStrings(t *testing.T) {
	t.Parallel()

	va := &viteAssets{}

	tags := []template.HTML{
		template.HTML("<link href='a'>"),
		template.HTML("<link href='b'>"),
		template.HTML("<link href='c'>"),
	}

	strings := va.htmlStrings(tags)

	if len(strings) != 3 {
		t.Errorf("expected 3 strings, got %d", len(strings))
	}

	if strings[0] != "<link href='a'>" {
		t.Errorf("unexpected string: %s", strings[0])
	}
}

func TestGenerateAllAssets(t *testing.T) {
	t.Parallel()

	t.Run("routes to hot reload mode", func(t *testing.T) {
		t.Parallel()

		hotFile := t.TempDir() + "/hot"
		if err := os.WriteFile(hotFile, []byte("http://localhost:5173"), 0o644); err != nil {
			t.Fatal(err)
		}

		i := &Inertia{}
		vi := &ViteInstance{
			Inertia: i,
			viteConfig: ViteConfig{
				HotFile:       hotFile,
				EntryPoints:   []string{"app.js"},
				HotReloadPort: "http://localhost:5173",
			},
		}

		html, err := vi.generateAllAssets()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		htmlStr := string(html)
		if !strings.Contains(htmlStr, "@vite/client") {
			t.Error("should include vite client in hot reload mode")
		}

		if !strings.Contains(htmlStr, "app.js") {
			t.Error("should include entry point in hot reload mode")
		}
	})

	t.Run("routes to production mode", func(t *testing.T) {
		t.Parallel()

		manifestPath := t.TempDir() + "/manifest.json"
		manifest := map[string]Asset{
			"app.js": {
				File:    "assets/app-abc123.js",
				IsEntry: true,
			},
		}
		data, _ := json.Marshal(manifest)
		if err := os.WriteFile(manifestPath, data, 0o644); err != nil {
			t.Fatal(err)
		}

		i := &Inertia{}
		vi := &ViteInstance{
			Inertia: i,
			viteConfig: ViteConfig{
				HotFile:       t.TempDir() + "/nonexistent",
				BuildManifest: manifestPath,
				BuildDir:      "/build/",
				EntryPoints:   []string{"app.js"},
			},
		}

		html, err := vi.generateAllAssets()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		htmlStr := string(html)
		if !strings.Contains(htmlStr, "/build/assets/app-abc123.js") {
			t.Error("should include built asset in production mode")
		}
	})
}

func TestGenerateProductionAssets(t *testing.T) {
	t.Parallel()

	t.Run("generates production tags with dependencies", func(t *testing.T) {
		t.Parallel()

		manifestPath := t.TempDir() + "/manifest.json"
		manifest := map[string]Asset{
			"app.js": {
				File:    "assets/app-abc.js",
				IsEntry: true,
				Imports: []string{"vendor.js"},
				Css:     []string{"app.css"},
			},
			"vendor.js": {
				File: "assets/vendor-def.js",
			},
			"app.css": {
				File: "assets/app-xyz.css",
			},
		}
		data, _ := json.Marshal(manifest)
		if err := os.WriteFile(manifestPath, data, 0o644); err != nil {
			t.Fatal(err)
		}

		i := &Inertia{}
		vi := &ViteInstance{
			Inertia: i,
			viteConfig: ViteConfig{
				BuildManifest: manifestPath,
				BuildDir:      "/build/",
				EntryPoints:   []string{"app.js"},
			},
		}

		html, err := vi.buildAssets("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		htmlStr := string(html)

		// Should include CSS
		if !strings.Contains(htmlStr, "/build/assets/app-xyz.css") {
			t.Error("should include CSS file")
		}

		// Should include entry script
		if !strings.Contains(htmlStr, "/build/assets/app-abc.js") {
			t.Error("should include entry script")
		}
	})

	t.Run("error when manifest not found", func(t *testing.T) {
		t.Parallel()

		i := &Inertia{}
		vi := &ViteInstance{
			Inertia: i,
			viteConfig: ViteConfig{
				BuildManifest: "/nonexistent/manifest.json",
				EntryPoints:   []string{"app.js"},
			},
		}

		_, err := vi.buildAssets("")
		if err == nil {
			t.Error("should return error when manifest not found")
		}
	})
}

func TestViteOptionFunctions(t *testing.T) {
	t.Parallel()

	t.Run("WithIntegrity", func(t *testing.T) {
		t.Parallel()

		cfg := &ViteConfig{}
		WithIntegrity()(cfg)

		if cfg.IntegrityKey == nil {
			t.Error("IntegrityKey should be set")
		}

		if *cfg.IntegrityKey != "integrity" {
			t.Errorf("expected integrity, got %s", *cfg.IntegrityKey)
		}
	})

	t.Run("WithIntegrityKey", func(t *testing.T) {
		t.Parallel()

		cfg := &ViteConfig{}
		WithIntegrityKey("custom-key")(cfg)

		if cfg.IntegrityKey == nil {
			t.Error("IntegrityKey should be set")
		}

		if *cfg.IntegrityKey != "custom-key" {
			t.Errorf("expected custom-key, got %s", *cfg.IntegrityKey)
		}
	})

	t.Run("WithEntryPoints", func(t *testing.T) {
		t.Parallel()

		cfg := &ViteConfig{}
		WithEntryPoints("app.js", "admin.js")(cfg)

		if len(cfg.EntryPoints) != 2 {
			t.Errorf("expected 2 entry points, got %d", len(cfg.EntryPoints))
		}

		if cfg.EntryPoints[0] != "app.js" {
			t.Errorf("expected app.js, got %s", cfg.EntryPoints[0])
		}
	})

	t.Run("WithoutPreloading", func(t *testing.T) {
		t.Parallel()

		cfg := &ViteConfig{}
		WithoutPreloading()(cfg)

		if cfg.PreloadStrategy != PreloadNone {
			t.Errorf("expected PreloadNone, got %v", cfg.PreloadStrategy)
		}
	})

	t.Run("WithAggressivePreload", func(t *testing.T) {
		t.Parallel()

		cfg := &ViteConfig{}
		WithAggressivePreload()(cfg)

		if cfg.PreloadStrategy != PreloadAggressive {
			t.Errorf("expected PreloadAggressive, got %v", cfg.PreloadStrategy)
		}
	})

	t.Run("WithWaterfallPreload", func(t *testing.T) {
		t.Parallel()

		cfg := &ViteConfig{}
		WithWaterfallPreload(5)(cfg)

		if cfg.PreloadStrategy != PreloadWaterfall {
			t.Errorf("expected PreloadWaterfall, got %v", cfg.PreloadStrategy)
		}

		if cfg.PreloadConcurrent != 5 {
			t.Errorf("expected 5, got %d", cfg.PreloadConcurrent)
		}
	})
}
