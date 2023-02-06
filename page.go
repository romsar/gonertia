package gonertia

import (
	"fmt"
	"net/http"
)

type page struct {
	Component string `json:"component"`
	Props     Props  `json:"props"`
	URL       string `json:"url"`
	Version   string `json:"version"`
}

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
