package gonertia

import (
	"context"
	"fmt"
)

type contextKey int

const (
	templateDataContextKey = contextKey(iota + 1)
	propsContextKey
)

// WithTemplateData appends template data value to the passed context.Context.
func (i *Inertia) WithTemplateData(ctx context.Context, key string, val any) context.Context {
	if ctxData := ctx.Value(templateDataContextKey); ctxData != nil {
		ctxData, ok := ctxData.(TemplateData)

		if ok {
			ctxData[key] = val
			return context.WithValue(ctx, templateDataContextKey, ctxData)
		}
	}

	return context.WithValue(ctx, templateDataContextKey, TemplateData{
		key: val,
	})
}

// WithProp appends prop value to the passed context.Context.
func (i *Inertia) WithProp(ctx context.Context, key string, val any) context.Context {
	if ctxData := ctx.Value(propsContextKey); ctxData != nil {
		ctxData, ok := ctxData.(Props)

		if ok {
			ctxData[key] = val
			return context.WithValue(ctx, propsContextKey, ctxData)
		}
	}

	return context.WithValue(ctx, propsContextKey, Props{
		key: val,
	})
}

// WithProps appends props values to the passed context.Context.
func (i *Inertia) WithProps(ctx context.Context, props Props) context.Context {
	if ctxData := ctx.Value(propsContextKey); ctxData != nil {
		ctxData, ok := ctxData.(Props)

		if ok {
			for key, val := range props {
				ctxData[key] = val
			}

			return context.WithValue(ctx, propsContextKey, ctxData)
		}
	}

	return context.WithValue(ctx, propsContextKey, props)
}

// TemplateDataFromContext returns template data from the context.
func TemplateDataFromContext(ctx context.Context) (TemplateData, error) {
	ctxData := ctx.Value(templateDataContextKey)

	if ctxData != nil {
		data, ok := ctxData.(TemplateData)
		if !ok {
			return nil, fmt.Errorf("template data in the context has invalid type")
		}

		return data, nil
	}

	return nil, nil
}

// PropsFromContext returns props from the context.
func PropsFromContext(ctx context.Context) (Props, error) {
	ctxData := ctx.Value(propsContextKey)

	if ctxData != nil {
		props, ok := ctxData.(Props)
		if !ok {
			return nil, fmt.Errorf("props in the context have invalid type")
		}

		return props, nil
	}

	return nil, nil
}
