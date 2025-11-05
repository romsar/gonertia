package gonertia

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
)

type viteAssets struct {
	processedAssets map[string]bool
	preloadTags     []template.HTML
	stylesheetTags  []template.HTML
	scriptTags      []template.HTML
	prefetchAssets  []string
}

func (va *viteAssets) htmlStrings(tags []template.HTML) []string {
	result := make([]string, len(tags))
	for i, tag := range tags {
		result[i] = string(tag)
	}
	return result
}

func (vi *ViteInstance) isJavaScript(asset Asset) bool {
	return strings.HasSuffix(asset.File, ".js")
}

func (vi *ViteInstance) isStylesheet(asset Asset) bool {
	return strings.HasSuffix(asset.File, ".css")
}

func (vi *ViteInstance) processAsset(
	src string,
	manifest map[string]Asset,
	assets *viteAssets,
	nonce string,
) error {
	// Skip if already processed
	if assets.processedAssets[src] {
		return nil
	}

	asset, ok := manifest[src]
	if !ok {
		return fmt.Errorf("asset %q not found in manifest (%d entries available)", src, len(manifest))
	}

	// Mark as processed
	assets.processedAssets[src] = true

	// Process static dependencies first (depth-first for proper load order)
	for _, imp := range asset.Imports {
		if err := vi.processAsset(imp, manifest, assets, nonce); err != nil {
			return err
		}
	}

	// Process CSS files
	for _, cssPath := range asset.Css {
		if assets.processedAssets[cssPath] {
			continue
		}

		cssAsset, ok := manifest[cssPath]
		if !ok {
			continue
		}

		tag := vi.generateStylesheetTag(cssAsset, nonce)
		assets.stylesheetTags = append(assets.stylesheetTags, tag)
		assets.processedAssets[cssPath] = true
	}

	switch vi.viteConfig.PreloadStrategy {
	case PreloadAggressive:
		if vi.isJavaScript(asset) {
			tag := vi.generatePreloadTag(asset, nonce)
			assets.preloadTags = append(assets.preloadTags, tag)
		}

	case PreloadWaterfall:
		if !asset.IsEntry && vi.isJavaScript(asset) {
			url := vi.viteConfig.BuildDir + asset.File
			assets.prefetchAssets = append(assets.prefetchAssets, url)
		}

	case PreloadNone:
	}

	if asset.IsEntry && vi.isJavaScript(asset) {
		tag := vi.generateScriptTag(asset, nonce)
		assets.scriptTags = append(assets.scriptTags, tag)
	}

	if asset.IsEntry && vi.isStylesheet(asset) {
		tag := vi.generateStylesheetTag(asset, nonce)
		assets.stylesheetTags = append(assets.stylesheetTags, tag)
	}

	return nil
}

func (vi *ViteInstance) generatePreloadTag(asset Asset, nonce string) template.HTML {
	url := vi.viteConfig.BuildDir + asset.File

	var attrs []string
	attrs = append(attrs, `rel="modulepreload"`)
	attrs = append(attrs, fmt.Sprintf(`href="%s"`, url))
	attrs = append(attrs, `crossorigin`) // Required for modules

	// Add integrity if available
	if integrity := vi.getIntegrity(asset); integrity != "" {
		attrs = append(attrs, fmt.Sprintf(`integrity="%s"`, integrity))
	}

	// Add nonce if present
	if nonce != "" {
		attrs = append(attrs, fmt.Sprintf(`nonce="%s"`, nonce))
	}

	return template.HTML(fmt.Sprintf(`<link %s>`, strings.Join(attrs, " ")))
}

func (vi *ViteInstance) generateStylesheetTag(asset Asset, nonce string) template.HTML {
	url := vi.viteConfig.BuildDir + asset.File

	var attrs []string
	attrs = append(attrs, `rel="stylesheet"`)
	attrs = append(attrs, fmt.Sprintf(`href="%s"`, url))

	// Add integrity if available
	if integrity := vi.getIntegrity(asset); integrity != "" {
		attrs = append(attrs, fmt.Sprintf(`integrity="%s"`, integrity))
	}

	// Add nonce if present
	if nonce != "" {
		attrs = append(attrs, fmt.Sprintf(`nonce="%s"`, nonce))
	}

	return template.HTML(fmt.Sprintf(`<link %s>`, strings.Join(attrs, " ")))
}

func (vi *ViteInstance) generateScriptTag(asset Asset, nonce string) template.HTML {
	url := vi.viteConfig.BuildDir + asset.File

	var attrs []string
	attrs = append(attrs, `type="module"`)
	attrs = append(attrs, fmt.Sprintf(`src="%s"`, url))

	// Add integrity if available
	if integrity := vi.getIntegrity(asset); integrity != "" {
		attrs = append(attrs, fmt.Sprintf(`integrity="%s"`, integrity))
	}

	// Add nonce if present
	if nonce != "" {
		attrs = append(attrs, fmt.Sprintf(`nonce="%s"`, nonce))
	}

	return template.HTML(fmt.Sprintf(`<script %s></script>`, strings.Join(attrs, " ")))
}

func (vi *ViteInstance) getIntegrity(asset Asset) string {
	return asset.Integrity
}

// generateAllAssets is the main template function that outputs all Vite assets.
// It accepts variadic entry points: {{ viteAssets "app.js" "admin.js" }}
// If no args provided, uses vi.viteConfig.EntryPoints from config.
// It routes to hot reload or production mode based on environment.
func (vi *ViteInstance) generateAllAssets(entries ...string) (template.HTML, error) {
	if vi.isHotReload() {
		return vi.buildHotAssets("", entries...)
	}
	return vi.buildAssets("", entries...)
}

// generateAllAssetsWithNonce is the template function for CSP-enabled asset output.
// It accepts a nonce as the first parameter followed by variadic entry points.
// Usage: {{ viteAssetsWithNonce .csp_nonce "app.js" "admin.js" }}
// The nonce is added to all script, style, and preload tags for Content Security Policy.
func (vi *ViteInstance) generateAllAssetsWithNonce(nonce string, entries ...string) (template.HTML, error) {
	if vi.isHotReload() {
		return vi.buildHotAssets(nonce, entries...)
	}
	return vi.buildAssets(nonce, entries...)
}

func (vi *ViteInstance) buildHotAssets(nonce string, entries ...string) (template.HTML, error) {
	// Use provided entries or fall back to config
	if len(entries) == 0 {
		entries = vi.viteConfig.EntryPoints
	}

	// Error if no entries configured
	if len(entries) == 0 {
		return "", fmt.Errorf("no entry points configured: use WithEntryPoints() or provide template args")
	}

	hotURL := vi.readHotReloadURL()

	tags := make([]string, 0, 1+len(entries))

	nonceAttr := ""
	if nonce != "" {
		nonceAttr = fmt.Sprintf(` nonce="%s"`, nonce)
	}

	tags = append(tags, fmt.Sprintf(
		`<script type="module"%s src="%s/@vite/client"></script>`,
		nonceAttr, hotURL,
	))

	for _, entry := range entries {
		tags = append(tags, fmt.Sprintf(
			`<script type="module"%s src="%s/%s"></script>`,
			nonceAttr, hotURL, entry,
		))
	}

	return template.HTML(strings.Join(tags, "\n")), nil
}

func (vi *ViteInstance) buildAssets(nonce string, entries ...string) (template.HTML, error) {
	// Use provided entries or fall back to config
	if len(entries) == 0 {
		entries = vi.viteConfig.EntryPoints
	}

	// Error if no entries configured
	if len(entries) == 0 {
		return "", fmt.Errorf("no entry points configured: use WithEntryPoints() or provide template args")
	}

	// Load manifest
	manifest, err := vi.loadManifest()
	if err != nil {
		return "", fmt.Errorf("load manifest: %w", err)
	}

	// Initialize asset tracker
	assets := &viteAssets{
		processedAssets: make(map[string]bool),
		prefetchAssets:  []string{},
	}

	// Process each entry and its dependencies
	for _, entry := range entries {
		if err := vi.processAsset(entry, manifest, assets, nonce); err != nil {
			return "", fmt.Errorf("process asset %q: %w", entry, err)
		}
	}

	// Build final output
	var allTags []string

	// 1. Modulepreload tags (for aggressive strategy)
	allTags = append(allTags, assets.htmlStrings(assets.preloadTags)...)

	// 2. Stylesheet tags
	allTags = append(allTags, assets.htmlStrings(assets.stylesheetTags)...)

	// 3. Script tags (entry points)
	allTags = append(allTags, assets.htmlStrings(assets.scriptTags)...)

	// 4. Waterfall prefetch script (if strategy enabled)
	if vi.viteConfig.PreloadStrategy == PreloadWaterfall && len(assets.prefetchAssets) > 0 {
		script := vi.generateWaterfallScript(assets.prefetchAssets, nonce)
		allTags = append(allTags, string(script))
	}

	return template.HTML(strings.Join(allTags, "\n")), nil
}

func (vi *ViteInstance) generateWaterfallScript(prefetchAssets []string, nonce string) template.HTML {
	if len(prefetchAssets) == 0 {
		return ""
	}

	// Default concurrent to 3 if not set or invalid
	concurrent := vi.viteConfig.PreloadConcurrent
	if concurrent <= 0 {
		concurrent = 3
	}

	type asset struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	}

	assets := make([]asset, len(prefetchAssets))
	for i, url := range prefetchAssets {
		assets[i] = asset{Href: url, Rel: "prefetch"}
	}

	assetsJSON, _ := json.Marshal(assets)

	nonceAttr := ""
	if nonce != "" {
		nonceAttr = fmt.Sprintf(` nonce="%s"`, nonce)
	}

	script := fmt.Sprintf(`<script%s>
window.addEventListener('load', () => {
  const assets = %s;
  const loadNext = (remaining, count) => {
    if (count > remaining.length) count = remaining.length;
    const fragment = new DocumentFragment();
    while (count > 0) {
      const cfg = remaining.shift();
      const link = document.createElement('link');
      link.href = cfg.href;
      link.rel = cfg.rel;
      fragment.append(link);
      count--;
      if (remaining.length) {
        link.onload = () => loadNext(remaining, 1);
        link.onerror = () => loadNext(remaining, 1);
      }
    }
    document.head.append(fragment);
  };
  loadNext(assets, %d);
});
</script>`, nonceAttr, string(assetsJSON), concurrent)

	return template.HTML(script)
}
