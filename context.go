package gonertia

import (
	"context"
	"fmt"
)

type contextKey int

const (
	templateDataContextKey = contextKey(iota + 1)
	propsContextKey
	validationErrorsContextKey
)

// WithTemplateData appends template data value to the passed context.Context.
func WithTemplateData(ctx context.Context, key string, val any) context.Context {
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

	return TemplateData{}, nil
}

// WithProp appends prop value to the passed context.Context.
func WithProp(ctx context.Context, key string, val any) context.Context {
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
func WithProps(ctx context.Context, props Props) context.Context {
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

	return Props{}, nil
}

// WithValidationError appends validation error to the passed context.Context.
func WithValidationError(ctx context.Context, key string, msg any) context.Context {
	if ctxData := ctx.Value(validationErrorsContextKey); ctxData != nil {
		ctxData, ok := ctxData.(ValidationErrors)

		if ok {
			ctxData[key] = msg
			return context.WithValue(ctx, validationErrorsContextKey, ctxData)
		}
	}

	return context.WithValue(ctx, validationErrorsContextKey, ValidationErrors{
		key: msg,
	})
}

// WithValidationErrors appends validation errors to the passed context.Context.
func WithValidationErrors(ctx context.Context, errors ValidationErrors) context.Context {
	if ctxData := ctx.Value(validationErrorsContextKey); ctxData != nil {
		ctxData, ok := ctxData.(ValidationErrors)

		if ok {
			for key, msg := range errors {
				ctxData[key] = msg
			}

			return context.WithValue(ctx, validationErrorsContextKey, ctxData)
		}
	}

	return context.WithValue(ctx, validationErrorsContextKey, errors)
}

// ValidationErrorsFromContext returns validation errors from the context.
func ValidationErrorsFromContext(ctx context.Context) (ValidationErrors, error) {
	ctxData := ctx.Value(validationErrorsContextKey)

	if ctxData != nil {
		validationErrors, ok := ctxData.(ValidationErrors)
		if !ok {
			return nil, fmt.Errorf("validation errors in the context have invalid type")
		}

		return validationErrors, nil
	}

	return ValidationErrors{}, nil
}
