package gonertia

import (
	"context"
	"reflect"
	"testing"
)

func TestInertia_WithTemplateData(t *testing.T) {
	t.Parallel()

	ctx := I().WithTemplateData(context.Background(), "foo", "bar")

	got, ok := ctx.Value(templateDataContextKey).(templateData)
	if !ok {
		t.Fatal("template data from context is not templateData")
	}

	want := templateData{"foo": "bar"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("temlateData=%#v, want=%#v", got, want)
	}
}

func TestInertia_WithTemplateProp(t *testing.T) {
	t.Parallel()

	ctx := I().WithProp(context.Background(), "foo", "bar")

	got, ok := ctx.Value(propsContextKey).(Props)
	if !ok {
		t.Fatal("props from context is not templateData")
	}

	want := Props{"foo": "bar"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Props=%#v, want=%#v", got, want)
	}
}

func TestInertia_WithTemplateProps(t *testing.T) {
	t.Parallel()

	ctx := I().WithProps(context.Background(), Props{"foo": "bar"})

	got, ok := ctx.Value(propsContextKey).(Props)
	if !ok {
		t.Fatal("props from context is not templateData")
	}

	want := Props{"foo": "bar"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Props=%#v, want=%#v", got, want)
	}
}

func Test_TemplateDataFromContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ctxData any
		want    templateData
		wantErr bool
	}{
		{
			name:    "nil",
			ctxData: nil,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "empty",
			ctxData: templateData{},
			want:    templateData{},
			wantErr: false,
		},
		{
			name:    "filled",
			ctxData: templateData{"foo": "bar"},
			want:    templateData{"foo": "bar"},
			wantErr: false,
		},
		{
			name:    "wrong type",
			ctxData: []string{"foo", "bar"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.WithValue(context.Background(), templateDataContextKey, tt.ctxData)

			got, err := TemplateDataFromContext(ctx)
			if tt.wantErr && err == nil {
				t.Fatal("error expected")
			} else if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %#v", err)
			} else if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("templateData=%#v, want=%#v", got, tt.want)
			}
		})
	}
}

func Test_PropsFromContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ctxData any
		want    Props
		wantErr bool
	}{
		{
			name:    "nil",
			ctxData: nil,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "empty",
			ctxData: Props{},
			want:    Props{},
			wantErr: false,
		},
		{
			name:    "filled",
			ctxData: Props{"foo": "bar"},
			want:    Props{"foo": "bar"},
			wantErr: false,
		},
		{
			name:    "wrong type",
			ctxData: []string{"foo", "bar"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.WithValue(context.Background(), propsContextKey, tt.ctxData)

			got, err := PropsFromContext(ctx)
			if tt.wantErr && err == nil {
				t.Fatal("error expected")
			} else if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %#v", err)
			} else if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Props=%#v, want=%#v", got, tt.want)
			}
		})
	}
}
