package gonertia

import (
	"context"
	"fmt"
	"net/http"
)

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

func (i *Inertia) prepareProps(r *http.Request, component string, props Props) (Props, error) {
	result := make(Props)

	// Add validation errors to the result.
	validationErrors, err := i.resolveValidationErrors(r)
	if err != nil {
		return nil, fmt.Errorf("resolve validation errors: %w", err)
	}
	result["errors"] = AlwaysProp{validationErrors}

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

func (i *Inertia) resolveValidationErrors(r *http.Request) (ValidationErrors, error) {
	// Add validation errors from storage.
	storageValidationErrors, err := i.restoreValidationErrors(r.Context())
	if err != nil {
		return nil, fmt.Errorf("getting validation errors from context: %w", err)
	}

	// ... and from context.
	ctxValidationErrors, err := ValidationErrorsFromContext(r.Context())
	if err != nil {
		return nil, fmt.Errorf("getting validation errors from context: %w", err)
	}

	validationErrors := make(ValidationErrors)
	for key, val := range storageValidationErrors {
		validationErrors[key] = val
	}
	for key, val := range ctxValidationErrors {
		validationErrors[key] = val
	}

	return ctxValidationErrors, nil
}

func (i *Inertia) restoreValidationErrors(ctx context.Context) (ValidationErrors, error) {
	if i.errorsStore == nil {
		return nil, nil
	}

	storageValidationErrors, err := i.errorsStore.Pop(ctx)
	if err != nil {
		return nil, fmt.Errorf("errors store pop: %w", err)
	}

	return storageValidationErrors, nil
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
