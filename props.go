package gonertia

import (
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

// AlwaysProp is a property value that will always evaluated.
//
// https://github.com/inertiajs/inertia-laravel/pull/627
type AlwaysProp func() any

func (i *Inertia) prepareProps(r *http.Request, component string, props Props) (Props, error) {
	result := make(Props)

	// Add validation errors from context.
	ctxValidationErrors, err := ValidationErrorsFromContext(r.Context())
	if err != nil {
		return nil, fmt.Errorf("getting validation errors from context: %w", err)
	}
	result["errors"] = AlwaysProp(func() any { return ctxValidationErrors })

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

	// Get props keys to return. If len == 0, then return all.
	only := i.propsKeysToReturn(r, component)

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

func (i *Inertia) propsKeysToReturn(r *http.Request, component string) map[string]struct{} {
	// Partial reloads only work for visits made to the same page component.
	//
	// https://inertiajs.com/partial-reloads
	if partialComponentFromRequest(r) == component {
		return setOf[string](partialDataFromRequest(r))
	}

	return nil
}

func resolvePropVal(val any) (_ any, err error) {
	switch typed := val.(type) {
	case func() any:
		return typed(), nil
	case AlwaysProp:
		return typed(), nil
	case func() (any, error):
		val, err = typed()
		if err != nil {
			return nil, fmt.Errorf("closure prop resolving: %w", err)
		}
	case LazyProp:
		val, err = typed()
		if err != nil {
			return nil, fmt.Errorf("lazy prop resolving: %w", err)
		}
	}

	return val, nil
}
