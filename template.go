package gonertia

import (
	"fmt"
	"net/http"
)

// TemplateData are the data that will be transferred
// and will be available in the root template.
type TemplateData map[string]any

func (i *Inertia) buildTemplateData(r *http.Request, page *page) (TemplateData, error) {
	pageJSON, err := i.marshallJSON(page)
	if err != nil {
		return nil, fmt.Errorf("marshal page into json error: %w", err)
	}

	// Get template data from context.
	ctxTemplateData, err := TemplateDataFromContext(r.Context())
	if err != nil {
		return nil, fmt.Errorf("getting template data from context error: %w", err)
	}

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
