package gonertia

import (
	"fmt"
	"net/http"
)

// Props are key-value data structure, that will
// be available in your front-end component.
type Props map[string]any

// LazyProp is a property value that will only evaluated then needed.
//
// https://inertiajs.com/partial-reloads
type LazyProp func() (any, error)

// prepareProps returns prepared and filtered props.
func (i *Inertia) prepareProps(r *http.Request, component string, props Props) (Props, error) {
	result := make(Props)

	// Add shared props to result.
	for key, val := range i.sharedProps {
		result[key] = val
	}

	// Add props from context to result.
	ctxProps, err := propsFromContext(r.Context())
	if err != nil {
		return nil, fmt.Errorf("getting props from context error: %w", err)
	}

	for key, val := range ctxProps {
		result[key] = val
	}

	// Add passed props to result.
	for key, val := range props {
		result[key] = val
	}

	// Get props keys to return. If len == 0, then return all.
	only := i.propsKeysToReturn(r, component)

	// Filter props.
	if len(only) > 0 {
		for key := range result {
			if _, ok := only[key]; !ok {
				delete(result, key)
			}
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
		val, err := resolvePropVal(val)
		if err != nil {
			return nil, fmt.Errorf("resolve prop value error: %w", err)
		}
		result[key] = val
	}

	return result, nil
}

// returns props keys to include to response.
func (i *Inertia) propsKeysToReturn(r *http.Request, component string) map[string]struct{} {
	// Partial reloads only work for visits made to the same page component.
	//
	// https://inertiajs.com/partial-reloads
	if partialComponentFromRequest(r) == component {
		return set[string](partialDataFromRequest(r))
	}

	return nil
}

// resolvePropVal resolves scalar value from prop value.
func resolvePropVal(val any) (_ any, err error) {
	if closure, ok := val.(func() (any, error)); ok {
		val, err = closure()

		if err != nil {
			return nil, fmt.Errorf("closure prop resolving error: %w", err)
		}
	} else if lazy, ok := val.(LazyProp); ok {
		val, err = lazy()

		if err != nil {
			return nil, fmt.Errorf("lazy prop resolving error: %w", err)
		}
	}

	return val, nil
}
