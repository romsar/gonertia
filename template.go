package gonertia

import (
	"fmt"
	"net/http"
)

// templateData is a map with values that
// will be available in the root template.
type templateData map[string]any

// buildTemplateData returns sharedProps based the page.
func (i *Inertia) buildTemplateData(r *http.Request, page *page) (templateData, error) {
	pageJSON, err := i.marshallJSON(page)
	if err != nil {
		return nil, fmt.Errorf("marshal page into json error: %w", err)
	}

	// Get template data from context.
	ctxTemplateData, err := templateDataFromContext(r.Context())
	if err != nil {
		return nil, fmt.Errorf("getting template data from context error: %w", err)
	}

	result := templateData{
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
