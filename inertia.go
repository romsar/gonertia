package gonertia

import (
	"context"
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

	flash FlashProvider

	ssrURL        string
	ssrHTTPClient *http.Client

	containerID    string
	version        string
	jsonMarshaller JSONMarshaller
	logger         Logger
}

// New initializes and returns Inertia.
func New(rootTemplateHTML string, opts ...Option) (*Inertia, error) {
	if rootTemplateHTML == "" {
		return nil, fmt.Errorf("blank root template")
	}

	i := &Inertia{
		rootTemplateHTML:    rootTemplateHTML,
		jsonMarshaller:      jsonDefaultMarshaller{},
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

// NewFromFile reads all bytes from the root template file and then initializes Inertia.
func NewFromFile(rootTemplatePath string, opts ...Option) (*Inertia, error) {
	bs, err := os.ReadFile(rootTemplatePath)
	if err != nil {
		return nil, fmt.Errorf("read file %q: %w", rootTemplatePath, err)
	}

	return NewFromBytes(bs, opts...)
}

// NewFromReader reads all bytes from the reader with root template html and then initializes Inertia.
func NewFromReader(rootTemplateReader io.Reader, opts ...Option) (*Inertia, error) {
	bs, err := io.ReadAll(rootTemplateReader)
	if err != nil {
		return nil, fmt.Errorf("read root template: %w", err)
	}
	if closer, ok := rootTemplateReader.(io.Closer); ok {
		_ = closer.Close()
	}

	return NewFromBytes(bs, opts...)
}

// NewFromBytes receive bytes with root template html and then initializes Inertia.
func NewFromBytes(rootTemplateBs []byte, opts ...Option) (*Inertia, error) {
	return New(string(rootTemplateBs), opts...)
}

// Logger defines an interface for debug messages.
type Logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
}

// FlashProvider defines an interface for flash data provider.
type FlashProvider interface {
	FlashErrors(ctx context.Context, errors ValidationErrors) error
	GetErrors(ctx context.Context) (ValidationErrors, error)
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
