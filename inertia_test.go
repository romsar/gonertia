package gonertia

import (
	"reflect"
	"testing"
)

var rootTemplate = `<html>
<head>{{ .inertiaHead }}</head>
<body>{{ .inertia }}</body>
</html>`

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		i, err := New(rootTemplate)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if i.rootTemplateHTML != rootTemplate {
			t.Fatalf("root template html=%s, want=%s", i.rootTemplateHTML, rootTemplate)
		}
	})

	t.Run("blank", func(t *testing.T) {
		t.Parallel()

		_, err := New("")
		if err == nil {
			t.Fatal("error expected")
		}
	})
}

func TestNewFromFile(t *testing.T) {
	t.Parallel()

	f := tmpFile(t, rootTemplate)

	i, err := NewFromFile(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if i.rootTemplateHTML != rootTemplate {
		t.Fatalf("root template html=%s, want=%s", i.rootTemplateHTML, rootTemplate)
	}
}

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

func TestInertia_ShareTemplateData(t *testing.T) {
	t.Parallel()

	type args struct {
		key string
		val any
	}
	tests := []struct {
		name         string
		templateData TemplateData
		args         args
		want         TemplateData
	}{
		{
			"add",
			TemplateData{},
			args{
				key: "foo",
				val: "bar",
			},
			TemplateData{"foo": "bar"},
		},
		{
			"replace",
			TemplateData{"foo": "zoo"},
			args{
				key: "foo",
				val: "bar",
			},
			TemplateData{"foo": "bar"},
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

func TestInertia_ShareTemplateFunc(t *testing.T) {
	t.Parallel()

	type args struct {
		key string
		val any
	}
	tests := []struct {
		name          string
		templateFuncs TemplateFuncs
		args          args
		want          TemplateFuncs
	}{
		{
			"add",
			TemplateFuncs{},
			args{
				key: "foo",
				val: "bar",
			},
			TemplateFuncs{"foo": "bar"},
		},
		{
			"replace",
			TemplateFuncs{"foo": "zoo"},
			args{
				key: "foo",
				val: "bar",
			},
			TemplateFuncs{"foo": "bar"},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			i := I(func(i *Inertia) {
				i.sharedTemplateFuncs = tt.templateFuncs
			})

			i.ShareTemplateFunc(tt.args.key, tt.args.val)

			if !reflect.DeepEqual(i.sharedTemplateFuncs, tt.want) {
				t.Fatalf("sharedTemplateFuncs=%#v, want=%#v", i.sharedTemplateFuncs, tt.want)
			}
		})
	}
}
