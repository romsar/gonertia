package gonertia

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// Inertia is a structure that contain all logic of Inertia server adapter.
type Inertia struct {
	// url is the app URL.
	url string

	// rootTemplate is the parsed root template.
	rootTemplate *template.Template

	// rootTemplatePath is the path to root template.
	rootTemplatePath string

	// templateFS is the FS that contain root template.
	templateFS fs.FS

	// sharedProps are global props.
	sharedProps Props

	// sharedTemplateData are global template data.
	sharedTemplateData TemplateData

	// sharedTemplateFuncMap is template's function map.
	sharedTemplateFuncMap template.FuncMap

	// version is the server asset version.
	version string

	// marshallJSON is the function that can encode bytes into JSON.
	marshallJSON marshallJSON

	// containerID is id of the Inertia HTML container.
	containerID string

	// logger is the package logger.
	logger logger
}

// New initializes and returns Inertia.
func New(url, rootTemplatePath string, opts ...Option) (*Inertia, error) {
	i := &Inertia{
		url:                   url,
		rootTemplatePath:      rootTemplatePath,
		marshallJSON:          json.Marshal,
		containerID:           "app",
		logger:                log.Default(),
		sharedProps:           make(Props),
		sharedTemplateData:    make(TemplateData),
		sharedTemplateFuncMap: make(template.FuncMap),
	}

	for _, opt := range opts {
		if err := opt(i); err != nil {
			return nil, fmt.Errorf("initialize inertia error: %w", err)
		}
	}

	return i, nil
}

// marshallJSON is the function that can encode bytes into JSON.
//
// By default, this package will use json.Marshal,
// but you are free to change this behavior.
type marshallJSON func(v any) ([]byte, error)

// logger gives methods to display status of the Inertia.
//
// Sometimes it's not possible to return an error,
// so we will this messages to logger.
type logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
}

// Option is an option parameter that modifies Inertia.
type Option func(i *Inertia) error

// WithTemplateFS returns Option that will set Inertia's templateFS.
func WithTemplateFS(templateFS fs.FS) Option {
	return func(i *Inertia) error {
		i.templateFS = templateFS
		return nil
	}
}

// WithVersion returns Option that will set Inertia's version.
func WithVersion(version string) Option {
	return func(i *Inertia) error {
		i.version = version
		return nil
	}
}

// WithAssetURL returns Option that will set Inertia's version based on asset url.
func WithAssetURL(url string) Option {
	return WithVersion(md5(url))
}

// WithManifestFile returns Option that will set Inertia's version based on manifest file.
func WithManifestFile(path string) Option {
	version, err := md5File(path)
	if err == nil {
		return WithVersion(version)
	}

	return func(i *Inertia) error {
		return fmt.Errorf("calculating md5 hash of manifest file error: %w", err)
	}
}

// WithMarshalJSON returns Option that will set Inertia's marshallJSON func.
func WithMarshalJSON(f marshallJSON) Option {
	return func(i *Inertia) error {
		i.marshallJSON = f
		return nil
	}
}

// WithLogger returns Option that will set Inertia's logger.
func WithLogger(log logger) Option {
	if log == nil {
		return WithoutLogger()
	}

	return func(i *Inertia) error {
		i.logger = log
		return nil
	}
}

// WithoutLogger returns Option that will unset Inertia's logger.
// Actually set a logger with io.Discard output.
func WithoutLogger() Option {
	return func(i *Inertia) error {
		i.logger = log.New(io.Discard, "", 0)
		return nil
	}
}

// WithContainerID returns Option that will set Inertia's container id.
func WithContainerID(id string) Option {
	return func(i *Inertia) error {
		i.containerID = id
		return nil
	}
}

// Location creates redirect response.
//
// If request is Inertia request, it will set status to 409 and url will be in "X-Inertia-Location" header.
// Otherwise, it will do an HTTP redirect with specified status (default is 302).
func (i *Inertia) Location(w http.ResponseWriter, r *http.Request, url string, status ...int) {
	if IsInertiaRequest(r) {
		setInertiaLocationToResponse(w, url)
		return
	}

	redirectResponse(w, r, url, status...)
}

// ShareTemplateData adds passed data to shared template data.
func (i *Inertia) ShareTemplateData(key string, val any) {
	i.sharedTemplateData[key] = val
}

// FlushSharedTemplateData flushes shared template data.
func (i *Inertia) FlushSharedTemplateData() {
	i.sharedTemplateData = make(TemplateData)
}

// ShareTemplateFunc adds passed value to the shared template func map.
func (i *Inertia) ShareTemplateFunc(key string, val any) {
	i.sharedTemplateFuncMap[key] = val
}

// FlushSharedTemplateFunc flushes the shared template func map.
func (i *Inertia) FlushSharedTemplateFunc() {
	i.sharedTemplateFuncMap = make(template.FuncMap)
}

// Render return response with Inertia data.
//
// If request is Inertia request - it will return JSON.
// Otherwise, it will return root template.
func (i *Inertia) Render(w http.ResponseWriter, r *http.Request, component string, props Props) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("render error: %w", err)
		}
	}()

	page, err := i.buildPage(r, component, props)
	if err != nil {
		return err
	}

	if IsInertiaRequest(r) {
		return i.doInertiaResponse(w, page)
	}

	return i.doHTMLResponse(w, r, page)
}

// doInertiaResponse writes Inertia response to the response writer.
func (i *Inertia) doInertiaResponse(w http.ResponseWriter, page *page) error {
	pageJSON, err := i.marshallJSON(page)
	if err != nil {
		return fmt.Errorf("marshal page into json error: %w", err)
	}

	markAsInertiaResponse(w)
	markAsJSONResponse(w)
	setResponseStatus(w, http.StatusOK)

	if _, err := w.Write(pageJSON); err != nil {
		return fmt.Errorf("write bytes to response error: %w", err)
	}

	return nil
}

// doHTMLResponse writes HTML response to the response writer.
func (i *Inertia) doHTMLResponse(w http.ResponseWriter, r *http.Request, page *page) (err error) {
	// If root template is already created - we'll use it to save some time.
	if i.rootTemplate == nil {
		i.rootTemplate, err = i.buildRootTemplate()
		if err != nil {
			return fmt.Errorf("build root template error: %w", err)
		}
	}

	templateData, err := i.buildTemplateData(r, page)
	if err != nil {
		return fmt.Errorf("build template data error: %w", err)
	}

	markAsHTMLResponse(w)

	if err := i.rootTemplate.Execute(w, templateData); err != nil {
		return fmt.Errorf("execute root template error: %w", err)
	}

	return nil
}

// buildRootTemplate parses files or FS and then returns root template.
func (i *Inertia) buildRootTemplate() (*template.Template, error) {
	tmpl := template.New(filepath.Base(i.rootTemplatePath)).Funcs(i.sharedTemplateFuncMap)

	if i.templateFS != nil {
		return tmpl.ParseFS(i.templateFS, i.rootTemplatePath)
	}

	return tmpl.ParseFiles(i.rootTemplatePath)
}

// inertiaContainerHTML returns Inertia container HTML based on page.
func (i *Inertia) inertiaContainerHTML(pageJSON []byte) template.HTML {
	builder := new(strings.Builder)

	// It doesn't look pretty, but fast!
	builder.WriteString(`<div id="`)
	builder.WriteString(i.containerID)
	builder.WriteString(`" data-page="`)
	template.HTMLEscape(builder, pageJSON)
	builder.WriteString(`"></div>`)

	return template.HTML(builder.String())
}

// backURL returns url that will be used to redirect browser to previous page.
// At the moment it based only by Referer HTTP header.
func (i *Inertia) backURL(r *http.Request) string {
	return refererFromRequest(r)
}
