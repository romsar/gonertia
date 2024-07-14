package gonertia

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// TemplateData are data that will be available in the root template.
type TemplateData map[string]any

// TemplateFuncs are functions that will be available in the root template.
type TemplateFuncs map[string]any

// Props are the data that will be transferred
// and will be available in the front-end component.
type Props map[string]any

// LazyProp is a property value that will only evaluated then needed.
//
// https://inertiajs.com/partial-reloads
type LazyProp func() (any, error)

func (p LazyProp) TryProp() (any, error) {
	return p()
}

// AlwaysProp is a property value that will always evaluated.
//
// https://github.com/inertiajs/inertia-laravel/pull/627
type AlwaysProp struct {
	Value any
}

func (p AlwaysProp) Prop() any {
	return p.Value
}

// Proper is an interface for custom type, which provides property, that will be resolved.
type Proper interface {
	Prop() any
}

// TryProper is an interface for custom type, which provides property and error, that will be resolved.
type TryProper interface {
	TryProp() (any, error)
}

// ValidationErrors are messages, that will be stored in the "errors" prop.
type ValidationErrors map[string]any

// Location creates redirect response.
//
// If request was made by Inertia - sets status to 409 and url will be in "X-Inertia-Location" header.
// Otherwise, it will do an HTTP redirect with specified status (default is 302 for GET, 303 for POST/PUT/PATCH).
func (i *Inertia) Location(w http.ResponseWriter, r *http.Request, url string, status ...int) {
	if IsInertiaRequest(r) {
		setInertiaLocationInResponse(w, url)
		return
	}

	redirectResponse(w, r, url, status...)
}

// Back creates redirect response to the previous url.
func (i *Inertia) Back(w http.ResponseWriter, r *http.Request, status ...int) {
	i.Location(w, r, i.backURL(r), status...)
}

func (i *Inertia) backURL(r *http.Request) string {
	// At the moment, it based only on the "Referer" HTTP header.
	return refererFromRequest(r)
}

// Render returns response with Inertia data.
//
// If request was made by Inertia - it will return data in JSON format.
// Otherwise, it will return HTML with root template.
//
// If SSR is enabled, pre-renders JavaScript and return HTML (https://inertiajs.com/server-side-rendering).
func (i *Inertia) Render(w http.ResponseWriter, r *http.Request, component string, props ...Props) (err error) {
	page, err := i.buildPage(r, component, firstOr[Props](props, nil))
	if err != nil {
		return fmt.Errorf("build page: %w", err)
	}

	if IsInertiaRequest(r) {
		if err = i.doInertiaResponse(w, page); err != nil {
			return fmt.Errorf("inertia response: %w", err)
		}

		return
	}

	if err = i.doHTMLResponse(w, r, page); err != nil {
		return fmt.Errorf("html response: %w", err)
	}

	return nil
}

type page struct {
	Component string `json:"component"`
	Props     Props  `json:"props"`
	URL       string `json:"url"`
	Version   string `json:"version"`
}

func (i *Inertia) buildPage(r *http.Request, component string, props Props) (*page, error) {
	props, err := i.prepareProps(r, component, props)
	if err != nil {
		return nil, fmt.Errorf("prepare props: %w", err)
	}

	return &page{
		Component: component,
		Props:     props,
		URL:       r.RequestURI,
		Version:   i.version,
	}, nil
}

func (i *Inertia) prepareProps(r *http.Request, component string, props Props) (Props, error) {
	result := make(Props)

	// Add validation errors from context.
	ctxValidationErrors, err := ValidationErrorsFromContext(r.Context())
	if err != nil {
		return nil, fmt.Errorf("getting validation errors from context: %w", err)
	}
	result["errors"] = AlwaysProp{ctxValidationErrors}

	// Add shared props to the result.
	for key, val := range i.sharedProps {
		result[key] = val
	}

	// Add props from context to the result.
	ctxProps, err := PropsFromContext(r.Context())
	if err != nil {
		return nil, fmt.Errorf("getting props from context: %w", err)
	}

	for key, val := range ctxProps {
		result[key] = val
	}

	// Add passed props to the result.
	for key, val := range props {
		result[key] = val
	}

	{
		// Only (include keys) and except (exclude keys) logic.
		only, except := i.getOnlyAndExcept(r, component)

		// Filter props.
		if len(only) > 0 {
			for key, val := range result {
				if _, ok := only[key]; ok {
					continue
				}
				if _, ok := val.(AlwaysProp); ok {
					continue
				}

				delete(result, key)
			}
		} else {
			for key, val := range result {
				if _, ok := val.(LazyProp); ok {
					delete(result, key)
				}
			}
		}
		for key := range except {
			delete(result, key)
		}
	}

	// Resolve props values.
	for key, val := range result {
		val, err = resolvePropVal(val)
		if err != nil {
			return nil, fmt.Errorf("resolve prop value: %w", err)
		}
		result[key] = val
	}

	return result, nil
}

func (i *Inertia) getOnlyAndExcept(r *http.Request, component string) (only, except map[string]struct{}) {
	// Partial reloads only work for visits made to the same page component.
	//
	// https://inertiajs.com/partial-reloads
	if partialComponentFromRequest(r) != component {
		return nil, nil
	}

	return setOf[string](onlyFromRequest(r)), setOf[string](exceptFromRequest(r))
}

func resolvePropVal(val any) (_ any, err error) {
	switch typed := val.(type) {
	case func() any:
		return typed(), nil
	case func() (any, error):
		val, err = typed()
		if err != nil {
			return nil, fmt.Errorf("closure prop resolving: %w", err)
		}
	case Proper:
		return typed.Prop(), nil
	case TryProper:
		return typed.TryProp()
	}

	return val, nil
}

func (i *Inertia) doInertiaResponse(w http.ResponseWriter, page *page) error {
	pageJSON, err := i.jsonMarshaller.Marshal(page)
	if err != nil {
		return fmt.Errorf("json marshal page into json: %w", err)
	}

	setInertiaInResponse(w)
	setJSONResponse(w)
	setResponseStatus(w, http.StatusOK)

	if _, err = w.Write(pageJSON); err != nil {
		return fmt.Errorf("write bytes to response: %w", err)
	}

	return nil
}

func (i *Inertia) doHTMLResponse(w http.ResponseWriter, r *http.Request, page *page) (err error) {
	// If root template is already created - we'll use it to save some time.
	if i.rootTemplate == nil {
		i.rootTemplate, err = i.buildRootTemplate()
		if err != nil {
			return fmt.Errorf("build root template: %w", err)
		}
	}

	templateData, err := i.buildTemplateData(r, page)
	if err != nil {
		return fmt.Errorf("build template data: %w", err)
	}

	setHTMLResponse(w)

	if err = i.rootTemplate.Execute(w, templateData); err != nil {
		return fmt.Errorf("execute root template: %w", err)
	}

	return nil
}

func (i *Inertia) buildRootTemplate() (*template.Template, error) {
	tmpl := template.New("").Funcs(template.FuncMap(i.sharedTemplateFuncs))
	return tmpl.Parse(i.rootTemplateHTML)
}

func (i *Inertia) buildTemplateData(r *http.Request, page *page) (TemplateData, error) {
	// Get template data from context.
	ctxTemplateData, err := TemplateDataFromContext(r.Context())
	if err != nil {
		return nil, fmt.Errorf("getting template data from context: %w", err)
	}

	// Defaults.
	inertia, inertiaHead, err := i.buildInertiaHTML(page)
	if err != nil {
		return nil, fmt.Errorf("build intertia html: %w", err)
	}
	templateData := TemplateData{
		"inertia":     inertia,
		"inertiaHead": inertiaHead,
	}

	// Add the shared template data to the result.
	for key, val := range i.sharedTemplateData {
		templateData[key] = val
	}

	// Add template data from context to the result.
	for key, val := range ctxTemplateData {
		templateData[key] = val
	}

	return templateData, nil
}

func (i *Inertia) buildInertiaHTML(page *page) (inertia, inertiaHead template.HTML, _ error) {
	pageJSON, err := i.jsonMarshaller.Marshal(page)
	if err != nil {
		return "", "", fmt.Errorf("json marshal page into json: %w", err)
	}

	if i.isSSREnabled() {
		inertia, inertiaHead, err = i.htmlContainerSSR(pageJSON)
		if err == nil {
			return inertia, inertiaHead, nil
		}

		i.logger.Printf("ssr rendering error: %s", err)
	}

	return i.htmlContainer(pageJSON)
}

func (i *Inertia) isSSREnabled() bool {
	return i.ssrURL != "" && i.ssrHTTPClient != nil
}

// htmlContainerSSR will send request with json marshaled page payload to ssr render endpoint.
// That endpoint will return head and body html, which will be returned and then rendered.
func (i *Inertia) htmlContainerSSR(pageJSON []byte) (inertia, inertiaHead template.HTML, _ error) {
	url := i.prepareSSRURL()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(pageJSON))
	if err != nil {
		return "", "", fmt.Errorf("new http request: %w", err)
	}
	setJSONRequest(req)

	resp, err := i.ssrHTTPClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("execute http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return "", "", fmt.Errorf("invalid response status code: %d", resp.StatusCode)
	}

	var ssr struct {
		Head []string `json:"head"`
		Body string   `json:"body"`
	}
	err = i.jsonMarshaller.Decode(resp.Body, &ssr)
	if err != nil {
		return "", "", fmt.Errorf("json decode ssr render response: %w", err)
	}

	inertia = template.HTML(ssr.Body)
	inertiaHead = template.HTML(strings.Join(ssr.Head, "\n"))

	return inertia, inertiaHead, nil
}

func (i *Inertia) prepareSSRURL() string {
	return strings.ReplaceAll(i.ssrURL, "/render", "") + "/render"
}

func (i *Inertia) htmlContainer(pageJSON []byte) (inertia, _ template.HTML, _ error) {
	builder := new(strings.Builder)

	// It doesn't look pretty, but fast!
	builder.WriteString(`<div id="`)
	builder.WriteString(i.containerID)
	builder.WriteString(`" data-page="`)
	template.HTMLEscape(builder, pageJSON)
	builder.WriteString(`"></div>`)

	return template.HTML(builder.String()), "", nil
}
