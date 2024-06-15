package gonertia

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

// Inertia is a main Gonertia structure, which contains all the logic for being an Inertia adapter.
type Inertia struct {
	rootTemplate     *template.Template
	rootTemplateHTML string

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
func New(rootTemplate string, opts ...Option) (*Inertia, error) {
	if f, err := os.Open(rootTemplate); err == nil {
		defer f.Close()

		var bs []byte
		bs, err = io.ReadAll(f)
		if err != nil {
			return nil, fmt.Errorf("cannot read root template file %q: %w", rootTemplate, err)
		}

		rootTemplate = string(bs)
	}

	i := &Inertia{
		rootTemplateHTML:    rootTemplate,
		marshallJSON:        json.Marshal,
		containerID:         "app",
		logger:              log.New(io.Discard, "", 0),
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

// ShareProp adds passed prop to shared props.
func (i *Inertia) ShareProp(key string, val any) {
	i.sharedProps[key] = val
}

// SharedProps returns shared props.
func (i *Inertia) SharedProps() Props {
	return i.sharedProps
}

// SharedProp return the shared prop.
func (i *Inertia) SharedProp(key string) (any, bool) {
	val, ok := i.sharedProps[key]
	return val, ok
}

// ShareTemplateData adds passed data to shared template data.
func (i *Inertia) ShareTemplateData(key string, val any) {
	i.sharedTemplateData[key] = val
}

// ShareTemplateFunc adds passed value to the shared template func map.
func (i *Inertia) ShareTemplateFunc(key string, val any) {
	i.sharedTemplateFuncs[key] = val
}
