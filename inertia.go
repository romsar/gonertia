package gonertia

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

// Inertia is a main Gonertia structure, which contains all the logic for being an Inertia adapter.
type Inertia struct {
	templateFS       fs.FS
	rootTemplate     *template.Template
	rootTemplatePath string

	sharedProps         Props
	sharedTemplateData  TemplateData
	sharedTemplateFuncs TemplateFuncs

	ssrURL        string
	ssrHTTPClient *http.Client

	containerID  string
	version      string
	marshallJSON marshallJSON
	logger       logger
}

// New initializes and returns Inertia.
func New(rootTemplatePath string, opts ...Option) (*Inertia, error) {
	i := &Inertia{
		rootTemplatePath:    rootTemplatePath,
		marshallJSON:        json.Marshal,
		containerID:         "app",
		logger:              log.Default(),
		sharedProps:         make(Props),
		sharedTemplateData:  make(TemplateData),
		sharedTemplateFuncs: make(TemplateFuncs),
	}

	for _, opt := range opts {
		if err := opt(i); err != nil {
			return nil, fmt.Errorf("initialize inertia: %w", err)
		}
	}

	return i, nil
}

type marshallJSON func(v any) ([]byte, error)

// Sometimes it's not possible to return an error,
// so we will send those messages to the logger.
type logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
}
