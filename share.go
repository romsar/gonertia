package gonertia

import "html/template"

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

// FlushSharedProps flushes shared props.
func (i *Inertia) FlushSharedProps() {
	i.sharedProps = make(Props)
}

// ShareTemplateData adds passed data to shared template data.
func (i *Inertia) ShareTemplateData(key string, val any) {
	i.sharedTemplateData[key] = val
}

// FlushSharedTemplateData flushes shared template data.
func (i *Inertia) FlushSharedTemplateData() {
	i.sharedTemplateData = make(templateData)
}

// ShareTemplateFunc adds passed value to the shared template func map.
func (i *Inertia) ShareTemplateFunc(key string, val any) {
	i.sharedTemplateFuncs[key] = val
}

// FlushSharedTemplateFuncs flushes the shared template func map.
func (i *Inertia) FlushSharedTemplateFuncs() {
	i.sharedTemplateFuncs = make(template.FuncMap)
}
