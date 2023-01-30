package gonertia

import (
	"context"
	"fmt"
)

// contextKey represents an internal key for adding context fields.
type contextKey int

const (
	// templateDataContextKey is the context key for template data.
	templateDataContextKey = contextKey(iota + 1)

	// propsContextKey is the context key for props.
	propsContextKey
)

// WithTemplateData appends template data value to passed context.Context.
func (i *Inertia) WithTemplateData(ctx context.Context, key string, val any) context.Context {
	if ctxData := ctx.Value(templateDataContextKey); ctxData != nil {
		ctxData, ok := ctxData.(templateData)

		if ok {
			ctxData[key] = val
			return context.WithValue(ctx, templateDataContextKey, ctxData)
		}
	}

	return context.WithValue(ctx, templateDataContextKey, templateData{
		key: val,
	})
}

// WithProp appends prop value to passed context.Context.
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

// WithProps appends props values to passed context.Context.
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

// templateDataFromContext returns template data from context.
func templateDataFromContext(ctx context.Context) (templateData, error) {
	ctxData := ctx.Value(templateDataContextKey)

	if ctxData != nil {
		data, ok := ctxData.(templateData)
		if !ok {
			return nil, fmt.Errorf("template data in context has invalid type")
		}

		return data, nil
	}

	return nil, nil
}

// propsFromContext returns props from context.
func propsFromContext(ctx context.Context) (Props, error) {
	ctxData := ctx.Value(propsContextKey)

	if ctxData != nil {
		props, ok := ctxData.(Props)
		if !ok {
			return nil, fmt.Errorf("props in context have invalid type")
		}

		return props, nil
	}

	return nil, nil
}
