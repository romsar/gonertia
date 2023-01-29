package gonertia

import (
	"context"
	"fmt"
)

// contextKey represents an internal key for adding context fields.
type contextKey int

const (
	// TemplateDataContextKey is the context key for template data.
	TemplateDataContextKey = contextKey(iota + 1)

	// PropsContextKey is the context key for props.
	PropsContextKey
)

// WithTemplateData appends template data value to passed context.Context.
func (i *Inertia) WithTemplateData(ctx context.Context, key string, val any) context.Context {
	if ctxData := ctx.Value(TemplateDataContextKey); ctxData != nil {
		ctxData, ok := ctxData.(TemplateData)

		if ok {
			ctxData[key] = val
			return context.WithValue(ctx, TemplateDataContextKey, ctxData)
		}
	}

	return context.WithValue(ctx, TemplateDataContextKey, TemplateData{
		key: val,
	})
}

// WithProp appends prop value to passed context.Context.
func (i *Inertia) WithProp(ctx context.Context, key string, val any) context.Context {
	if ctxData := ctx.Value(PropsContextKey); ctxData != nil {
		ctxData, ok := ctxData.(Props)

		if ok {
			ctxData[key] = val
			return context.WithValue(ctx, PropsContextKey, ctxData)
		}
	}

	return context.WithValue(ctx, PropsContextKey, Props{
		key: val,
	})
}

// WithProps appends props values to passed context.Context.
func (i *Inertia) WithProps(ctx context.Context, props Props) context.Context {
	if ctxData := ctx.Value(PropsContextKey); ctxData != nil {
		ctxData, ok := ctxData.(Props)

		if ok {
			for key, val := range props {
				ctxData[key] = val
			}

			return context.WithValue(ctx, PropsContextKey, ctxData)
		}
	}

	return context.WithValue(ctx, PropsContextKey, props)
}

// templateDataFromContext returns template data from context.
func templateDataFromContext(ctx context.Context) (TemplateData, error) {
	ctxData := ctx.Value(TemplateDataContextKey)

	if ctxData != nil {
		data, ok := ctxData.(TemplateData)
		if !ok {
			return nil, fmt.Errorf("template data in context has invalid type")
		}

		return data, nil
	}

	return nil, nil
}

// propsFromContext returns props from context.
func propsFromContext(ctx context.Context) (Props, error) {
	ctxData := ctx.Value(PropsContextKey)

	if ctxData != nil {
		props, ok := ctxData.(Props)
		if !ok {
			return nil, fmt.Errorf("props in context have invalid type")
		}

		return props, nil
	}

	return nil, nil
}
