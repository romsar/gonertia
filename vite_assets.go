// Asset generation logic for Vite integration.
package gonertia

import (
	"fmt"
	"html/template"
	"strings"
)

// viteAssets tracks generated tags during asset processing.
type viteAssets struct {
	processedAssets map[string]bool // Track processed to avoid duplicates
	preloadTags     []template.HTML // <link rel="modulepreload">
	stylesheetTags  []template.HTML // <link rel="stylesheet">
	scriptTags      []template.HTML // <script type="module">
	prefetchAssets  []string        // Assets for waterfall prefetch
}

// htmlStrings converts template.HTML slice to string slice.
func (va *viteAssets) htmlStrings(tags []template.HTML) []string {
	result := make([]string, len(tags))
	for i, tag := range tags {
		result[i] = string(tag)
	}
	return result
}

// resolveNonce gets the nonce to use for this ViteInstance.
func (vi *ViteInstance) resolveNonce() string {
	// 1. Use static nonce if set
	if vi.viteConfig.Nonce != "" {
		return vi.viteConfig.Nonce
	}

	// 2. Call generator if provided
	if vi.viteConfig.NonceGenerator != nil {
		return vi.viteConfig.NonceGenerator()
	}

	// 3. No nonce
	return ""
}

// processAsset recursively processes an asset and its dependencies.
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

	// Get asset from manifest
	asset, ok := manifest[src]
	if !ok {
		return fmt.Errorf("asset %q not found in manifest", src)
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

	// Handle based on preload strategy
	switch vi.viteConfig.PreloadStrategy {
	case PreloadAggressive:
		// Add modulepreload for all JavaScript imports
		if strings.HasSuffix(asset.File, ".js") {
			tag := vi.generatePreloadTag(asset, nonce)
			assets.preloadTags = append(assets.preloadTags, tag)
		}

	case PreloadWaterfall:
		// Add to prefetch list (will be loaded after page load)
		// Skip the entry itself (it's in scriptTags)
		if !asset.IsEntry && strings.HasSuffix(asset.File, ".js") {
			url := vi.viteConfig.BuildDir + asset.File
			assets.prefetchAssets = append(assets.prefetchAssets, url)
		}

	case PreloadNone:
		// Do nothing - browser handles discovery
	}

	// Add script tag ONLY for JavaScript entry points
	if asset.IsEntry && strings.HasSuffix(asset.File, ".js") {
		tag := vi.generateScriptTag(asset, nonce)
		assets.scriptTags = append(assets.scriptTags, tag)
	}

	// Add stylesheet tag for CSS entry points
	if asset.IsEntry && strings.HasSuffix(asset.File, ".css") {
		tag := vi.generateStylesheetTag(asset, nonce)
		assets.stylesheetTags = append(assets.stylesheetTags, tag)
	}

	return nil
}

// generatePreloadTag generates a modulepreload link tag.
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

// generateStylesheetTag generates a stylesheet link tag.
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

// generateScriptTag generates a module script tag.
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

// getIntegrity retrieves integrity hash from asset.
// Returns the integrity hash if present in the manifest (e.g., from vite-plugin-manifest-sri).
// The hash is used for SubResource Integrity (SRI) to verify asset integrity.
func (vi *ViteInstance) getIntegrity(asset Asset) string {
	return asset.Integrity
}

// generateAllAssets is the main template function that outputs all Vite assets.
// It accepts variadic entry points: {{ viteAssets "app.js" "admin.js" }}
// If no args provided, uses vi.viteConfig.EntryPoints from config.
// It routes to hot reload or production mode based on environment.
func (vi *ViteInstance) generateAllAssets(entries ...string) (template.HTML, error) {
	if vi.isHotReload() {
		return vi.generateHotAssets(entries...)
	}
	return vi.generateProductionAssets(entries...)
}

// generateHotAssets generates tags for development mode with HMR support.
func (vi *ViteInstance) generateHotAssets(entries ...string) (template.HTML, error) {
	// Use provided entries or fall back to config
	if len(entries) == 0 {
		entries = vi.viteConfig.EntryPoints
	}

	// Error if no entries configured
	if len(entries) == 0 {
		return "", fmt.Errorf("no entry points configured: use WithEntryPoints() or provide template args")
	}

	hotURL := vi.readHotReloadURL()

	// Pre-allocate: 1 for Vite client + entry points
	tags := make([]string, 0, 1+len(entries))

	// Add Vite HMR client
	tags = append(tags, fmt.Sprintf(
		`<script type="module" src="%s/@vite/client"></script>`,
		hotURL,
	))

	// Add each entry point
	for _, entry := range entries {
		tags = append(tags, fmt.Sprintf(
			`<script type="module" src="%s/%s"></script>`,
			hotURL, entry,
		))
	}

	return template.HTML(strings.Join(tags, "\n")), nil
}

// generateProductionAssets generates optimized tags for production.
func (vi *ViteInstance) generateProductionAssets(entries ...string) (template.HTML, error) {
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

	// Resolve nonce once for all tags
	nonce := vi.resolveNonce()

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

// generateWaterfallScript generates an inline script for batched prefetching.
func (vi *ViteInstance) generateWaterfallScript(prefetchAssets []string, nonce string) template.HTML {
	if len(prefetchAssets) == 0 {
		return ""
	}

	// Default concurrent to 3 if not set or invalid
	concurrent := vi.viteConfig.PreloadConcurrent
	if concurrent <= 0 {
		concurrent = 3
	}

	// Build JavaScript assets array
	var assetsJSON strings.Builder
	assetsJSON.WriteString("[")
	for i, url := range prefetchAssets {
		if i > 0 {
			assetsJSON.WriteString(",")
		}
		// Escape quotes in URL
		escapedURL := strings.ReplaceAll(url, `"`, `\"`)
		assetsJSON.WriteString(fmt.Sprintf(`{href:"%s",rel:"prefetch"}`, escapedURL))
	}
	assetsJSON.WriteString("]")

	// Build nonce attribute
	nonceAttr := ""
	if nonce != "" {
		nonceAttr = fmt.Sprintf(` nonce="%s"`, nonce)
	}

	// Generate the waterfall script (minified)
	script := fmt.Sprintf(`<script%s>
window.addEventListener('load',()=>setTimeout(()=>{
const loadNext=(assets,count)=>{
if(count>assets.length)count=assets.length
const fragment=new DocumentFragment
while(count>0){
const cfg=assets.shift()
const link=document.createElement('link')
link.href=cfg.href
link.rel=cfg.rel
fragment.append(link)
count--
if(assets.length){
link.onload=()=>loadNext(assets,1)
link.onerror=()=>loadNext(assets,1)
}
}
document.head.append(fragment)
}
loadNext(%s,%d)
}))
</script>`, nonceAttr, assetsJSON.String(), concurrent)

	return template.HTML(script)
}
