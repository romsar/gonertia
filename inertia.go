package gonertia

import (
	"encoding/json"
	"fmt"
	"html/template"
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
	sharedTemplateData templateData

	// sharedTemplateFuncs is template's function map.
	sharedTemplateFuncs template.FuncMap

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
		url:                 url,
		rootTemplatePath:    rootTemplatePath,
		marshallJSON:        json.Marshal,
		containerID:         "app",
		logger:              log.Default(),
		sharedProps:         make(Props),
		sharedTemplateData:  make(templateData),
		sharedTemplateFuncs: make(template.FuncMap),
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
// so we will send those messages to logger.
type logger interface {
	Printf(format string, v ...any)
	Println(v ...any)
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

// Render return response with Inertia data.
//
// If request is Inertia request - it will return JSON.
// Otherwise, it will return root template.
func (i *Inertia) Render(w http.ResponseWriter, r *http.Request, component string, props ...Props) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("render error: %w", err)
		}
	}()

	page, err := i.buildPage(r, component, firstOr[Props](props, nil))
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
	tmpl := template.New(filepath.Base(i.rootTemplatePath)).Funcs(i.sharedTemplateFuncs)

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
// At the moment it based only by the "Referer" HTTP header.
func (i *Inertia) backURL(r *http.Request) string {
	return refererFromRequest(r)
}
