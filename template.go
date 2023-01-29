package gonertia

import (
	"fmt"
	"net/http"
)

// TemplateData is a map with values that
// will be available in the root template.
type TemplateData map[string]any

// buildTemplateData returns sharedProps based on page.
func (i *Inertia) buildTemplateData(r *http.Request, page *page) (TemplateData, error) {
	pageJSON, err := i.marshallJSON(page)
	if err != nil {
		return nil, fmt.Errorf("marshal page into json error: %w", err)
	}

	// Get template data from context.
	ctxTemplateData, err := templateDataFromContext(r.Context())
	if err != nil {
		return nil, fmt.Errorf("getting template data from context error: %w", err)
	}

	data := TemplateData{
		"inertiaHead": "", // reserved for SSR.
		"inertia":     i.inertiaContainerHTML(pageJSON),
	}

	// Add shared template data.
	for key, val := range i.sharedTemplateData {
		data[key] = val
	}

	// Add template data from context.
	for key, val := range ctxTemplateData {
		data[key] = val
	}

	return data, nil
}
