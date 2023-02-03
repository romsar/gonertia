package gonertia

import (
	"html/template"
	"reflect"
	"testing"
)

func TestInertia_ShareProp(t *testing.T) {
	t.Parallel()

	type args struct {
		key string
		val any
	}
	tests := []struct {
		name  string
		props Props
		args  args
		want  Props
	}{
		{
			"add",
			Props{},
			args{
				key: "foo",
				val: "bar",
			},
			Props{"foo": "bar"},
		},
		{
			"replace",
			Props{"foo": "zoo"},
			args{
				key: "foo",
				val: "bar",
			},
			Props{"foo": "bar"},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			i := I(func(i *Inertia) {
				i.sharedProps = tt.props
			})

			i.ShareProp(tt.args.key, tt.args.val)

			if !reflect.DeepEqual(i.sharedProps, tt.want) {
				t.Fatalf("sharedProps=%#v, want=%#v", i.sharedProps, tt.want)
			}
		})
	}
}

func TestInertia_SharedProps(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		props Props
	}{
		{
			"empty",
			Props{},
		},
		{
			"with values",
			Props{"foo": "bar"},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			i := I(func(i *Inertia) {
				i.sharedProps = tt.props
			})

			got := i.SharedProps()

			if !reflect.DeepEqual(got, i.sharedProps) {
				t.Fatalf("sharedProps=%#v, want=%#v", got, i.sharedProps)
			}
		})
	}
}

func TestInertia_SharedProp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		props  Props
		key    string
		want   any
		wantOk bool
	}{
		{
			"empty props",
			Props{},
			"foo",
			nil,
			false,
		},
		{
			"not found",
			Props{"foo": 123},
			"bar",
			nil,
			false,
		},
		{
			"found",
			Props{"foo": 123},
			"foo",
			123,
			true,
		},
		{
			"found nil value",
			Props{"foo": nil},
			"foo",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			i := I(func(i *Inertia) {
				i.sharedProps = tt.props
			})

			got, ok := i.SharedProp(tt.key)
			if ok != tt.wantOk {
				t.Fatalf("ok=%t, want=%t", ok, tt.wantOk)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("SharedProp()=%#v, want=%#v", got, tt.want)
			}
		})
	}
}

func TestInertia_FlushSharedProps(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		props Props
	}{
		{
			"empty props",
			Props{},
		},
		{
			"non-empty props",
			Props{"foo": "bar"},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			i := I(func(i *Inertia) {
				i.sharedProps = tt.props
			})

			i.FlushSharedProps()

			if !reflect.DeepEqual(i.sharedProps, Props{}) {
				t.Fatalf("sharedProps=%#v, want=%#v", i.sharedProps, Props{})
			}
		})
	}
}

func TestInertia_ShareTemplateData(t *testing.T) {
	t.Parallel()

	type args struct {
		key string
		val any
	}
	tests := []struct {
		name         string
		templateData templateData
		args         args
		want         templateData
	}{
		{
			"add",
			templateData{},
			args{
				key: "foo",
				val: "bar",
			},
			templateData{"foo": "bar"},
		},
		{
			"replace",
			templateData{"foo": "zoo"},
			args{
				key: "foo",
				val: "bar",
			},
			templateData{"foo": "bar"},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			i := I(func(i *Inertia) {
				i.sharedTemplateData = tt.templateData
			})

			i.ShareTemplateData(tt.args.key, tt.args.val)

			if !reflect.DeepEqual(i.sharedTemplateData, tt.want) {
				t.Fatalf("sharedTemplateData=%#v, want=%#v", i.sharedTemplateData, tt.want)
			}
		})
	}
}

func TestInertia_FlushTemplateData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		templateData templateData
	}{
		{
			"empty template data",
			templateData{},
		},
		{
			"non-empty template data",
			templateData{"foo": "bar"},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			i := I(func(i *Inertia) {
				i.sharedTemplateData = tt.templateData
			})

			i.FlushSharedTemplateData()

			if !reflect.DeepEqual(i.sharedTemplateData, templateData{}) {
				t.Fatalf("sharedTemplateData=%#v, want=%#v", i.sharedTemplateData, templateData{})
			}
		})
	}
}

func TestInertia_ShareTemplateFunc(t *testing.T) {
	t.Parallel()

	type args struct {
		key string
		val any
	}
	tests := []struct {
		name    string
		funcMap template.FuncMap
		args    args
		want    template.FuncMap
	}{
		{
			"add",
			template.FuncMap{},
			args{
				key: "foo",
				val: "bar",
			},
			template.FuncMap{"foo": "bar"},
		},
		{
			"replace",
			template.FuncMap{"foo": "zoo"},
			args{
				key: "foo",
				val: "bar",
			},
			template.FuncMap{"foo": "bar"},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			i := I(func(i *Inertia) {
				i.sharedTemplateFuncs = tt.funcMap
			})

			i.ShareTemplateFunc(tt.args.key, tt.args.val)

			if !reflect.DeepEqual(i.sharedTemplateFuncs, tt.want) {
				t.Fatalf("sharedTemplateFuncs=%#v, want=%#v", i.sharedTemplateFuncs, tt.want)
			}
		})
	}
}

func TestInertia_FlushSharedTemplateFuncs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		funcMap template.FuncMap
	}{
		{
			"empty func map",
			template.FuncMap{},
		},
		{
			"non-empty func map",
			template.FuncMap{"foo": "bar"},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			i := I(func(i *Inertia) {
				i.sharedTemplateFuncs = tt.funcMap
			})

			i.FlushSharedTemplateFuncs()

			if !reflect.DeepEqual(i.sharedTemplateFuncs, template.FuncMap{}) {
				t.Fatalf("sharedTemplateFuncs=%#v, want=%#v", i.sharedTemplateFuncs, template.FuncMap{})
			}
		})
	}
}
