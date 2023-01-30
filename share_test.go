package gonertia

import (
	"html/template"
	"reflect"
	"testing"
)

func TestInertia_ShareProp(t *testing.T) {
	t.Parallel()

	t.Run("add value", func(t *testing.T) {
		t.Parallel()

		key := "foo"
		val := "bar"

		i := &Inertia{sharedProps: make(Props)}
		i.ShareProp(key, val)

		got := i.sharedProps[key]
		if got != val {
			t.Fatalf("got=%#v, want=%#v", got, val)
		}
	})

	t.Run("add empty value", func(t *testing.T) {
		t.Parallel()

		key := "foo"
		val := ""

		i := &Inertia{sharedProps: make(Props)}
		i.ShareProp(key, "")

		got, ok := i.sharedProps[key]
		if !ok {
			t.Fatalf("value with key %q not found", key)
		}

		if got != val {
			t.Fatalf("got=%#v, want=%#v", got, val)
		}
	})

	t.Run("replace value", func(t *testing.T) {
		t.Parallel()

		key := "foo"
		val := "zoo"

		i := &Inertia{sharedProps: Props{key: "bar"}}
		i.ShareProp(key, val)

		got := i.sharedProps[key]
		if got != val {
			t.Fatalf("got=%#v, want=%#v", got, val)
		}

		gotLen := len(i.sharedProps)
		wantLen := 1
		if gotLen != wantLen {
			t.Fatalf("got len=%#v, want len=%#v", gotLen, wantLen)
		}
	})
}

func TestInertia_SharedProps(t *testing.T) {
	t.Parallel()

	t.Run("basic test", func(t *testing.T) {
		t.Parallel()

		want := Props{"foo": "bar"}

		i := &Inertia{sharedProps: want}

		got := i.SharedProps()

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got=%#v, want=%#v", got, want)
		}
	})
}

func TestInertia_SharedProp(t *testing.T) {
	t.Parallel()

	t.Run("found", func(t *testing.T) {
		t.Parallel()

		key := "foo"
		val := "bar"

		i := &Inertia{sharedProps: Props{key: val}}

		got, ok := i.SharedProp(key)
		if !ok {
			t.Fatalf("value with key %q not found", key)
		}

		if got != val {
			t.Fatalf("got=%#v, want=%#v", got, val)
		}
	})

	t.Run("found empty", func(t *testing.T) {
		t.Parallel()

		key := "foo"
		val := ""

		i := &Inertia{sharedProps: Props{key: val}}

		got, ok := i.SharedProp(key)
		if !ok {
			t.Fatalf("value with key %q not found", key)
		}

		if got != val {
			t.Fatalf("got=%#v, want=%#v", got, val)
		}
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		key := "foo"
		val := "bar"

		i := &Inertia{sharedProps: Props{key: val}}

		got, ok := i.SharedProp("zoo")
		if ok {
			t.Fatalf("value with key %q found", key)
		}

		if got != nil {
			t.Fatalf("got=%#v, want=%#v", got, nil)
		}
	})
}

func TestInertia_FlushSharedProps(t *testing.T) {
	t.Parallel()

	t.Run("basic test", func(t *testing.T) {
		t.Parallel()

		i := &Inertia{sharedProps: Props{"foo": "bar"}}

		i.FlushSharedProps()

		got := i.sharedProps
		want := Props{}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got=%#v, want=%#v", got, want)
		}
	})

	t.Run("flush on empty", func(t *testing.T) {
		t.Parallel()

		i := &Inertia{sharedProps: Props{}}

		i.FlushSharedProps()

		got := i.sharedProps
		want := Props{}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got=%#v, want=%#v", got, want)
		}
	})
}

func TestInertia_ShareTemplateData(t *testing.T) {
	t.Parallel()

	t.Run("add value", func(t *testing.T) {
		t.Parallel()

		key := "foo"
		val := "bar"

		i := &Inertia{sharedTemplateData: make(templateData)}
		i.ShareTemplateData(key, val)

		got := i.sharedTemplateData[key]
		if got != val {
			t.Fatalf("got=%#v, want=%#v", got, val)
		}
	})

	t.Run("add empty value", func(t *testing.T) {
		t.Parallel()

		key := "foo"
		val := ""

		i := &Inertia{sharedTemplateData: make(templateData)}
		i.ShareTemplateData(key, "")

		got, ok := i.sharedTemplateData[key]
		if !ok {
			t.Fatalf("value with key %q not found", key)
		}

		if got != val {
			t.Fatalf("got=%#v, want=%#v", got, val)
		}
	})

	t.Run("replace value", func(t *testing.T) {
		t.Parallel()

		key := "foo"
		val := "zoo"

		i := &Inertia{sharedTemplateData: templateData{key: "bar"}}
		i.ShareTemplateData(key, val)

		got := i.sharedTemplateData[key]
		if got != val {
			t.Fatalf("got=%#v, want=%#v", got, val)
		}

		gotLen := len(i.sharedTemplateData)
		wantLen := 1
		if gotLen != wantLen {
			t.Fatalf("got len=%#v, want len=%#v", gotLen, wantLen)
		}
	})
}

func TestInertia_FlushSharedTemplateData(t *testing.T) {
	t.Parallel()

	t.Run("basic test", func(t *testing.T) {
		t.Parallel()

		i := &Inertia{sharedTemplateData: templateData{"foo": "bar"}}

		i.FlushSharedTemplateData()

		got := i.sharedTemplateData
		want := templateData{}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got=%#v, want=%#v", got, want)
		}
	})

	t.Run("flush on empty", func(t *testing.T) {
		t.Parallel()

		i := &Inertia{sharedTemplateData: templateData{}}

		i.FlushSharedTemplateData()

		got := i.sharedTemplateData
		want := templateData{}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got=%#v, want=%#v", got, want)
		}
	})
}

func TestInertia_ShareTemplateFunc(t *testing.T) {
	t.Parallel()

	t.Run("add value", func(t *testing.T) {
		t.Parallel()

		key := "foo"
		val := "bar"

		i := &Inertia{sharedTemplateFuncMap: make(template.FuncMap)}
		i.ShareTemplateFunc(key, val)

		got := i.sharedTemplateFuncMap[key]
		if got != val {
			t.Fatalf("got=%#v, want=%#v", got, val)
		}
	})

	t.Run("add empty value", func(t *testing.T) {
		t.Parallel()

		key := "foo"
		val := ""

		i := &Inertia{sharedTemplateFuncMap: make(template.FuncMap)}
		i.ShareTemplateFunc(key, "")

		got, ok := i.sharedTemplateFuncMap[key]
		if !ok {
			t.Fatalf("value with key %q not found", key)
		}

		if got != val {
			t.Fatalf("got=%#v, want=%#v", got, val)
		}
	})

	t.Run("replace value", func(t *testing.T) {
		t.Parallel()

		key := "foo"
		val := "zoo"

		i := &Inertia{sharedTemplateFuncMap: template.FuncMap{key: func() any { return nil }}}
		i.ShareTemplateFunc(key, val)

		got := i.sharedTemplateFuncMap[key]
		if got != val {
			t.Fatalf("got=%#v, want=%#v", got, val)
		}

		gotLen := len(i.sharedTemplateFuncMap)
		wantLen := 1
		if gotLen != wantLen {
			t.Fatalf("got len=%#v, want len=%#v", gotLen, wantLen)
		}
	})
}

func TestInertia_FlushSharedTemplateFunc(t *testing.T) {
	t.Parallel()

	t.Run("basic test", func(t *testing.T) {
		t.Parallel()

		i := &Inertia{sharedTemplateFuncMap: template.FuncMap{"foo": func() any { return nil }}}

		i.FlushSharedTemplateFunc()

		got := i.sharedTemplateFuncMap
		want := template.FuncMap{}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got=%#v, want=%#v", got, want)
		}
	})

	t.Run("flush on empty", func(t *testing.T) {
		t.Parallel()

		i := &Inertia{sharedTemplateFuncMap: template.FuncMap{}}

		i.FlushSharedTemplateFunc()

		got := i.sharedTemplateFuncMap
		want := template.FuncMap{}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got=%#v, want=%#v", got, want)
		}
	})
}
