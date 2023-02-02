package gonertia

import (
	"fmt"
	"net/http"
)

// page is a structure that will be encoded to the JSON
// and present inside "data-page" attribute of the Inertia HTML container,
// or will be returned to the browser directly (if request was made by Inertia).
type page struct {
	// Component is the front-end component that will be rendered on users browsers.
	Component string `json:"component"`

	// Props are key-value data structure that will be available in a front-end component.
	Props Props `json:"props"`

	// URL is the current page URL.
	URL string `json:"url"`

	// Version is the asset version of the server.
	Version string `json:"version"`
}

// buildPage creates page with prepared props and other needle data.
func (i *Inertia) buildPage(r *http.Request, component string, props Props) (*page, error) {
	props, err := i.prepareProps(r, component, props)
	if err != nil {
		return nil, fmt.Errorf("prepare props error: %w", err)
	}

	return &page{
		Component: component,
		Props:     props,
		URL:       r.RequestURI,
		Version:   i.version,
	}, nil
}
