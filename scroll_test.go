package gonertia

import (
	"net/http"
	"testing"
)

func TestScrollProp(t *testing.T) {
	t.Parallel()

	t.Run("basic scroll prop", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/home")
		asInertiaRequest(r)

		items := []string{"item1", "item2", "item3"}
		err := I().Render(w, r, "Some/Component", Props{
			"items": Scroll(items),
			"foo":   "bar",
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		assertable := AssertFromString(t, w.Body.String())
		assertable.AssertProps(Props{
			"items": map[string]any{
				"data": []any{"item1", "item2", "item3"},
			},
			"foo":    "bar",
			"errors": map[string]any{},
		})
		// ScrollProp appends to wrapper "data" by default
		assertable.AssertMergeProps([]string{"items.data"})
	})

	t.Run("scroll prop with append intent", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/home")
		asInertiaRequest(r)
		r.Header.Set(headerInertiaInfiniteScrollMergeIntent, "append")

		items := []string{"item1", "item2", "item3"}
		err := I().Render(w, r, "Some/Component", Props{
			"items": Scroll(items),
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		assertable := AssertFromString(t, w.Body.String())
		assertable.AssertProps(Props{
			"items": map[string]any{
				"data": []any{"item1", "item2", "item3"},
			},
			"errors": map[string]any{},
		})
		assertable.AssertMergeProps([]string{"items.data"})
	})

	t.Run("scroll prop with prepend intent", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/home")
		asInertiaRequest(r)
		r.Header.Set(headerInertiaInfiniteScrollMergeIntent, "prepend")

		items := []string{"item1", "item2", "item3"}
		err := I().Render(w, r, "Some/Component", Props{
			"items": Scroll(items),
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		assertable := AssertFromString(t, w.Body.String())
		assertable.AssertProps(Props{
			"items": map[string]any{
				"data": []any{"item1", "item2", "item3"},
			},
			"errors": map[string]any{},
		})
		assertable.AssertPrependProps([]string{"items.data"})
	})

	t.Run("scroll prop with custom wrapper", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/home")
		asInertiaRequest(r)
		r.Header.Set(headerInertiaInfiniteScrollMergeIntent, "append")

		items := []string{"item1", "item2", "item3"}
		err := I().Render(w, r, "Some/Component", Props{
			"items": Scroll(items, WithWrapper("results")),
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		assertable := AssertFromString(t, w.Body.String())
		assertable.AssertMergeProps([]string{"items.results"})
	})

	t.Run("scroll prop with metadata", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/home")
		asInertiaRequest(r)

		metadata := ScrollMetadata{
			PageName:     "page",
			PreviousPage: 1,
			NextPage:     3,
			CurrentPage:  2,
		}

		err := I().Render(w, r, "Some/Component", Props{
			"items": Scroll([]string{"item1", "item2"}, WithMetadata(metadata)),
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		assertable := AssertFromString(t, w.Body.String())
		assertable.AssertScrollProps(map[string]map[string]any{
			"items": {
				"pageName":     "page",
				"previousPage": float64(1),
				"nextPage":     float64(3),
				"currentPage":  float64(2),
				"reset":        false,
			},
		})
	})

	t.Run("scroll prop with reset", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/home")
		asInertiaRequest(r)
		withReset(r, []string{"items"})

		err := I().Render(w, r, "Some/Component", Props{
			"items": Scroll([]string{"item1", "item2"}),
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		assertable := AssertFromString(t, w.Body.String())
		assertable.AssertScrollProps(map[string]map[string]any{
			"items": {
				"pageName":     "page",
				"previousPage": nil,
				"nextPage":     nil,
				"currentPage":  nil,
				"reset":        true,
			},
		})
		// Reset props should not be in merge props
		if len(assertable.MergeProps) > 0 {
			t.Fatalf("expected no merge props, got %v", assertable.MergeProps)
		}
	})

	t.Run("auto-wraps data when not already wrapped", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/home")
		asInertiaRequest(r)

		items := []string{"item1", "item2", "item3"}
		err := I().Render(w, r, "Some/Component", Props{
			"items": Scroll(items),
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		assertable := AssertFromString(t, w.Body.String())
		// Data should be wrapped in "data" key
		assertable.AssertProps(Props{
			"items": map[string]any{
				"data": []any{"item1", "item2", "item3"},
			},
			"errors": map[string]any{},
		})
	})

	t.Run("does not double-wrap already wrapped data", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/home")
		asInertiaRequest(r)

		// Data already wrapped in "data" key
		wrappedData := map[string]any{
			"data": []string{"item1", "item2", "item3"},
		}
		err := I().Render(w, r, "Some/Component", Props{
			"items": Scroll(wrappedData),
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		assertable := AssertFromString(t, w.Body.String())
		// Should not double-wrap
		assertable.AssertProps(Props{
			"items": map[string]any{
				"data": []any{"item1", "item2", "item3"},
			},
			"errors": map[string]any{},
		})
	})

	t.Run("auto-wraps with custom wrapper", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/home")
		asInertiaRequest(r)

		items := []string{"item1", "item2", "item3"}
		err := I().Render(w, r, "Some/Component", Props{
			"items": Scroll(items, WithWrapper("results")),
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		assertable := AssertFromString(t, w.Body.String())
		// Data should be wrapped in "results" key
		assertable.AssertProps(Props{
			"items": map[string]any{
				"results": []any{"item1", "item2", "item3"},
			},
			"errors": map[string]any{},
		})
	})

	t.Run("does not double-wrap with custom wrapper", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/home")
		asInertiaRequest(r)

		// Data already wrapped in "results" key
		wrappedData := map[string]any{
			"results": []string{"item1", "item2", "item3"},
		}
		err := I().Render(w, r, "Some/Component", Props{
			"items": Scroll(wrappedData, WithWrapper("results")),
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		assertable := AssertFromString(t, w.Body.String())
		// Should not double-wrap
		assertable.AssertProps(Props{
			"items": map[string]any{
				"results": []any{"item1", "item2", "item3"},
			},
			"errors": map[string]any{},
		})
	})
}

func TestAdvancedMergeProps(t *testing.T) {
	t.Parallel()

	t.Run("deep merge props", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/home")
		asInertiaRequest(r)

		err := I().Render(w, r, "Some/Component", Props{
			"settings": DeepMerge(map[string]any{
				"theme": "dark",
				"lang":  "en",
			}),
			"foo": "bar",
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		assertable := AssertFromString(t, w.Body.String())
		assertable.AssertDeepMergeProps([]string{"settings"})
	})

	t.Run("prepend props at root", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/home")
		asInertiaRequest(r)

		err := I().Render(w, r, "Some/Component", Props{
			"items": Merge([]int{1, 2, 3}).Prepend(),
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		assertable := AssertFromString(t, w.Body.String())
		assertable.AssertPrependProps([]string{"items"})
	})

	t.Run("append props at paths", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/home")
		asInertiaRequest(r)

		err := I().Render(w, r, "Some/Component", Props{
			"data": Merge(map[string]any{
				"items": []int{1, 2, 3},
			}).Append("items"),
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		assertable := AssertFromString(t, w.Body.String())
		assertable.AssertMergeProps([]string{"data.items"})
	})

	t.Run("prepend props at paths", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/home")
		asInertiaRequest(r)

		err := I().Render(w, r, "Some/Component", Props{
			"data": Merge(map[string]any{
				"items": []int{1, 2, 3},
			}).Prepend("items"),
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		assertable := AssertFromString(t, w.Body.String())
		assertable.AssertPrependProps([]string{"data.items"})
	})

	t.Run("match on keys", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/home")
		asInertiaRequest(r)

		err := I().Render(w, r, "Some/Component", Props{
			"users": Merge([]map[string]any{
				{"id": 1, "name": "Alice"},
				{"id": 2, "name": "Bob"},
			}).MatchOn("id"),
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		assertable := AssertFromString(t, w.Body.String())
		assertable.AssertMatchPropsOn([]string{"users.id"})
	})

	t.Run("multiple match on keys", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/home")
		asInertiaRequest(r)

		err := I().Render(w, r, "Some/Component", Props{
			"users": Merge([]map[string]any{
				{"id": 1, "email": "alice@example.com"},
			}).MatchOn("id", "email"),
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		assertable := AssertFromString(t, w.Body.String())
		assertable.AssertMatchPropsOn([]string{"users.id", "users.email"})
	})
}

func TestScrollMetadata(t *testing.T) {
	t.Parallel()

	t.Run("implements ProvidesScrollMetadata", func(t *testing.T) {
		t.Parallel()

		metadata := ScrollMetadata{
			PageName:     "page",
			PreviousPage: 1,
			NextPage:     3,
			CurrentPage:  2,
		}

		var _ ProvidesScrollMetadata = metadata

		if metadata.GetPageName() != "page" {
			t.Fatalf("expected page name to be 'page', got %s", metadata.GetPageName())
		}
		if metadata.GetPreviousPage() != 1 {
			t.Fatalf("expected previous page to be 1, got %v", metadata.GetPreviousPage())
		}
		if metadata.GetNextPage() != 3 {
			t.Fatalf("expected next page to be 3, got %v", metadata.GetNextPage())
		}
		if metadata.GetCurrentPage() != 2 {
			t.Fatalf("expected current page to be 2, got %v", metadata.GetCurrentPage())
		}
	})

	t.Run("custom metadata provider", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/home")
		asInertiaRequest(r)

		// Custom metadata provider
		customMetadata := ScrollMetadata{
			PageName:     "cursor",
			PreviousPage: "prev-token",
			NextPage:     "next-token",
			CurrentPage:  "current-token",
		}

		err := I().Render(w, r, "Some/Component", Props{
			"items": Scroll([]string{"item1", "item2"}, WithMetadata(customMetadata)),
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		assertable := AssertFromString(t, w.Body.String())
		assertable.AssertScrollProps(map[string]map[string]any{
			"items": {
				"pageName":     "cursor",
				"previousPage": "prev-token",
				"nextPage":     "next-token",
				"currentPage":  "current-token",
				"reset":        false,
			},
		})
	})

	t.Run("metadata function", func(t *testing.T) {
		t.Parallel()

		w, r := requestMock(http.MethodGet, "/home")
		asInertiaRequest(r)

		type paginatedData struct {
			Data        []string
			CurrentPage int
			NextPage    *int
			PrevPage    *int
		}

		next := 3
		prev := 1
		data := paginatedData{
			Data:        []string{"item1", "item2"},
			CurrentPage: 2,
			NextPage:    &next,
			PrevPage:    &prev,
		}

		err := I().Render(w, r, "Some/Component", Props{
			"items": Scroll(data, WithMetadataFunc(func(val any) ProvidesScrollMetadata {
				d := val.(paginatedData)
				return ScrollMetadata{
					PageName:     "page",
					CurrentPage:  d.CurrentPage,
					NextPage:     d.NextPage,
					PreviousPage: d.PrevPage,
				}
			})),
		})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		assertable := AssertFromString(t, w.Body.String())
		assertable.AssertScrollProps(map[string]map[string]any{
			"items": {
				"pageName":     "page",
				"previousPage": float64(1),
				"nextPage":     float64(3),
				"currentPage":  float64(2),
				"reset":        false,
			},
		})
	})
}
