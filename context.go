package gonertia

import (
	"context"
)

type contextKey int

const (
	templateDataContextKey = contextKey(iota + 1)
	propsContextKey
	validationErrorsContextKey
	encryptHistoryContextKey
	clearHistoryContextKey
)

// SetTemplateData sets template data to the passed context.
func SetTemplateData(ctx context.Context, templateData TemplateData) context.Context {
	return context.WithValue(ctx, templateDataContextKey, templateData)
}

// SetTemplateDatum sets single template data item to the passed context.
func SetTemplateDatum(ctx context.Context, key string, val any) context.Context {
	templateData := TemplateDataFromContext(ctx)
	templateData[key] = val
	return SetTemplateData(ctx, templateData)
}

// TemplateDataFromContext returns template data from the context.
func TemplateDataFromContext(ctx context.Context) TemplateData {
	templateData, ok := ctx.Value(templateDataContextKey).(TemplateData)
	if ok {
		return templateData
	}
	return TemplateData{}
}

// SetProps sets props values to the passed context.
func SetProps(ctx context.Context, props Props) context.Context {
	return context.WithValue(ctx, propsContextKey, props)
}

// SetProp sets prop value to the passed context.
func SetProp(ctx context.Context, key string, val any) context.Context {
	props := PropsFromContext(ctx)
	props[key] = val
	return SetProps(ctx, props)
}

// PropsFromContext returns props from the context.
func PropsFromContext(ctx context.Context) Props {
	props, ok := ctx.Value(propsContextKey).(Props)
	if ok {
		return props
	}
	return Props{}
}

// SetValidationErrors sets validation errors to the passed context.
func SetValidationErrors(ctx context.Context, errors ValidationErrors) context.Context {
	return context.WithValue(ctx, validationErrorsContextKey, errors)
}

// AddValidationErrors appends validation errors to the passed context.
func AddValidationErrors(ctx context.Context, errors ValidationErrors) context.Context {
	validationErrors := ValidationErrorsFromContext(ctx)
	for key, val := range errors {
		validationErrors[key] = val
	}
	return SetValidationErrors(ctx, validationErrors)
}

// SetValidationError sets validation error to the passed context.
func SetValidationError(ctx context.Context, key string, msg string) context.Context {
	validationErrors := ValidationErrorsFromContext(ctx)
	validationErrors[key] = msg
	return SetValidationErrors(ctx, validationErrors)
}

// ValidationErrorsFromContext returns validation errors from the context.
func ValidationErrorsFromContext(ctx context.Context) ValidationErrors {
	validationErrors, ok := ctx.Value(validationErrorsContextKey).(ValidationErrors)
	if ok {
		return validationErrors
	}
	return ValidationErrors{}
}

// SetEncryptHistory enables or disables history encryption.
func SetEncryptHistory(ctx context.Context, encrypt ...bool) context.Context {
	return context.WithValue(ctx, encryptHistoryContextKey, firstOr[bool](encrypt, true))
}

// EncryptHistoryFromContext returns history encryption value from the context.
func EncryptHistoryFromContext(ctx context.Context) (bool, bool) {
	encryptHistory, ok := ctx.Value(encryptHistoryContextKey).(bool)
	return encryptHistory, ok
}

// ClearHistory cleaning history state.
func ClearHistory(ctx context.Context) context.Context {
	return context.WithValue(ctx, clearHistoryContextKey, true)
}

// ClearHistoryFromContext returns clear history value from the context.
func ClearHistoryFromContext(ctx context.Context) bool {
	clearHistory, ok := ctx.Value(clearHistoryContextKey).(bool)
	if ok {
		return clearHistory
	}
	return false
}
