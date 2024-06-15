package gonertia

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// TemplateData are data that will be available in the root template.
type TemplateData map[string]any

// TemplateFuncs are functions that will be available in the root template.
type TemplateFuncs map[string]any

func (i *Inertia) buildRootTemplate() (*template.Template, error) {
	tmpl := template.New(filepath.Base(i.rootTemplatePath)).Funcs(template.FuncMap(i.sharedTemplateFuncs))

	if i.templateFS != nil {
		return tmpl.ParseFS(i.templateFS, i.rootTemplatePath)
	}

	return tmpl.ParseFiles(i.rootTemplatePath)
}

func (i *Inertia) buildTemplateData(r *http.Request, page *page) (TemplateData, error) {
	pageJSON, err := i.marshallJSON(page)
	if err != nil {
		return nil, fmt.Errorf("marshal page into json: %w", err)
	}

	// Get template data from context.
	ctxTemplateData, err := TemplateDataFromContext(r.Context())
	if err != nil {
		return nil, fmt.Errorf("getting template data from context: %w", err)
	}

	// Defaults.
	result := TemplateData{
		"inertiaHead": "", // todo reserved for SSR.
		"inertia":     i.inertiaContainerHTML(pageJSON),
	}

	// Add the shared template data to the result.
	for key, val := range i.sharedTemplateData {
		result[key] = val
	}

	// Add template data from context to the result.
	for key, val := range ctxTemplateData {
		result[key] = val
	}

	return result, nil
}
