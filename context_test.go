package gonertia

import (
	"context"
	"reflect"
	"testing"
)

func TestInertia_SetTemplateData(t *testing.T) {
	t.Parallel()

	t.Run("fresh", func(t *testing.T) {
		t.Parallel()

		ctx := SetTemplateData(context.Background(), TemplateData{"foo": "bar"})

		got, ok := ctx.Value(templateDataContextKey).(TemplateData)
		if !ok {
			t.Fatal("template data from context is not `TemplateData` type")
		}

		want := TemplateData{"foo": "bar"}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("TemplateData=%#v, want=%#v", got, want)
		}
	})

	t.Run("already filled", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(context.Background(), templateDataContextKey, TemplateData{"baz": "quz", "foo": "quz"})
		ctx = SetTemplateData(ctx, TemplateData{"foo": "bar"})

		got, ok := ctx.Value(templateDataContextKey).(TemplateData)
		if !ok {
			t.Fatal("template data from context is not `TemplateData` type")
		}

		want := TemplateData{"foo": "bar"}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("TemplateData=%#v, want=%#v", got, want)
		}
	})
}

func TestInertia_SetTemplateDatum(t *testing.T) {
	t.Parallel()

	t.Run("fresh", func(t *testing.T) {
		t.Parallel()

		ctx := SetTemplateDatum(context.Background(), "foo", "bar")

		got, ok := ctx.Value(templateDataContextKey).(TemplateData)
		if !ok {
			t.Fatal("template data from context is not `TemplateData` type")
		}

		want := TemplateData{"foo": "bar"}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("TemplateData=%#v, want=%#v", got, want)
		}
	})

	t.Run("already filled", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(context.Background(), templateDataContextKey, TemplateData{"baz": "quz", "foo": "quz"})
		ctx = SetTemplateDatum(ctx, "foo", "bar")

		got, ok := ctx.Value(templateDataContextKey).(TemplateData)
		if !ok {
			t.Fatal("template data from context is not `TemplateData` type")
		}

		want := TemplateData{"foo": "bar", "baz": "quz"}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("TemplateData=%#v, want=%#v", got, want)
		}
	})
}

func Test_TemplateDataFromContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ctxData any
		want    TemplateData
	}{
		{
			name:    "nil",
			ctxData: nil,
			want:    TemplateData{},
		},
		{
			name:    "empty",
			ctxData: TemplateData{},
			want:    TemplateData{},
		},
		{
			name:    "filled",
			ctxData: TemplateData{"foo": "bar"},
			want:    TemplateData{"foo": "bar"},
		},
		{
			name:    "wrong type",
			ctxData: []string{"foo", "bar"},
			want:    TemplateData{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.WithValue(context.Background(), templateDataContextKey, tt.ctxData)

			got := TemplateDataFromContext(ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("TemplateData=%#v, want=%#v", got, tt.want)
			}
		})
	}
}

func TestInertia_SetProps(t *testing.T) {
	t.Parallel()

	t.Run("fresh", func(t *testing.T) {
		t.Parallel()

		ctx := SetProps(context.Background(), Props{"foo": "bar"})

		got, ok := ctx.Value(propsContextKey).(Props)
		if !ok {
			t.Fatal("props from context is not `Props` type")
		}

		want := Props{"foo": "bar"}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Props=%#v, want=%#v", got, want)
		}
	})

	t.Run("already filled", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(context.Background(), propsContextKey, Props{"baz": "quz", "foo": "quz"})
		ctx = SetProps(ctx, Props{"foo": "bar"})

		got, ok := ctx.Value(propsContextKey).(Props)
		if !ok {
			t.Fatal("props from context is not `Props` type")
		}

		want := Props{"foo": "bar"}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Props=%#v, want=%#v", got, want)
		}
	})
}

func TestInertia_SetProp(t *testing.T) {
	t.Parallel()

	t.Run("fresh", func(t *testing.T) {
		t.Parallel()

		ctx := SetProp(context.Background(), "foo", "bar")

		got, ok := ctx.Value(propsContextKey).(Props)
		if !ok {
			t.Fatal("props from context is not `Props` type")
		}

		want := Props{"foo": "bar"}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Props=%#v, want=%#v", got, want)
		}
	})

	t.Run("already filled", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(context.Background(), propsContextKey, Props{"baz": "quz", "foo": "quz"})
		ctx = SetProp(ctx, "foo", "bar")

		got, ok := ctx.Value(propsContextKey).(Props)
		if !ok {
			t.Fatal("props from context is not `Props` type")
		}

		want := Props{"foo": "bar", "baz": "quz"}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("props=%#v, want=%#v", got, want)
		}
	})
}

func Test_PropsFromContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ctxData any
		want    Props
	}{
		{
			name:    "nil",
			ctxData: nil,
			want:    Props{},
		},
		{
			name:    "empty",
			ctxData: Props{},
			want:    Props{},
		},
		{
			name:    "filled",
			ctxData: Props{"foo": "bar"},
			want:    Props{"foo": "bar"},
		},
		{
			name:    "wrong type",
			ctxData: []string{"foo", "bar"},
			want:    Props{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.WithValue(context.Background(), propsContextKey, tt.ctxData)

			got := PropsFromContext(ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Props=%#v, want=%#v", got, tt.want)
			}
		})
	}
}

func TestInertia_SetValidationErrors(t *testing.T) {
	t.Parallel()

	t.Run("fresh", func(t *testing.T) {
		t.Parallel()

		ctx := SetValidationErrors(context.Background(), ValidationErrors{"foo": "bar"})

		got, ok := ctx.Value(validationErrorsContextKey).(ValidationErrors)
		if !ok {
			t.Fatal("validation errors from context is not `ValidationErrors` type")
		}

		want := ValidationErrors{"foo": "bar"}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("ValidationErrors=%#v, want=%#v", got, want)
		}
	})

	t.Run("already filled", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(context.Background(), validationErrorsContextKey, ValidationErrors{"baz": "quz", "foo": "quz"})
		ctx = SetValidationErrors(ctx, ValidationErrors{"foo": "bar"})

		got, ok := ctx.Value(validationErrorsContextKey).(ValidationErrors)
		if !ok {
			t.Fatal("validation errors from context is not `ValidationErrors` type")
		}

		want := ValidationErrors{"foo": "bar"}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("ValidationErrors=%#v, want=%#v", got, want)
		}
	})
}

func TestInertia_AddValidationErrors(t *testing.T) {
	t.Parallel()

	t.Run("fresh", func(t *testing.T) {
		t.Parallel()

		ctx := AddValidationErrors(context.Background(), ValidationErrors{"foo": "bar"})

		got, ok := ctx.Value(validationErrorsContextKey).(ValidationErrors)
		if !ok {
			t.Fatal("validation errors from context is not `ValidationErrors` type")
		}

		want := ValidationErrors{"foo": "bar"}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("ValidationErrors=%#v, want=%#v", got, want)
		}
	})

	t.Run("already filled", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(context.Background(), validationErrorsContextKey, ValidationErrors{"baz": "quz", "foo": "quz"})
		ctx = AddValidationErrors(ctx, ValidationErrors{"foo": "bar"})

		got, ok := ctx.Value(validationErrorsContextKey).(ValidationErrors)
		if !ok {
			t.Fatal("validation errors from context is not `ValidationErrors` type")
		}

		want := ValidationErrors{"baz": "quz", "foo": "bar"}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("ValidationErrors=%#v, want=%#v", got, want)
		}
	})
}

func TestInertia_SetValidationError(t *testing.T) {
	t.Parallel()

	t.Run("fresh", func(t *testing.T) {
		t.Parallel()

		ctx := SetValidationError(context.Background(), "foo", "bar")

		got, ok := ctx.Value(validationErrorsContextKey).(ValidationErrors)
		if !ok {
			t.Fatal("validation errors from context is not `ValidationErrors` type")
		}

		want := ValidationErrors{"foo": "bar"}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("ValidationErrors=%#v, want=%#v", got, want)
		}
	})

	t.Run("already filled", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(context.Background(), validationErrorsContextKey, ValidationErrors{"baz": "quz", "foo": "quz"})
		ctx = SetValidationError(ctx, "foo", "bar")

		got, ok := ctx.Value(validationErrorsContextKey).(ValidationErrors)
		if !ok {
			t.Fatal("validation errors from context is not `ValidationErrors` type")
		}

		want := ValidationErrors{"foo": "bar", "baz": "quz"}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("ValidationErrors=%#v, want=%#v", got, want)
		}
	})
}

func Test_ValidationErrorsFromContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ctxData any
		want    ValidationErrors
	}{
		{
			name:    "nil",
			ctxData: nil,
			want:    ValidationErrors{},
		},
		{
			name:    "empty",
			ctxData: ValidationErrors{},
			want:    ValidationErrors{},
		},
		{
			name:    "filled",
			ctxData: ValidationErrors{"foo": "bar"},
			want:    ValidationErrors{"foo": "bar"},
		},
		{
			name:    "filled with nested",
			ctxData: ValidationErrors{"foo": ValidationErrors{"abc": "123"}},
			want:    ValidationErrors{"foo": ValidationErrors{"abc": "123"}},
		},
		{
			name:    "wrong type",
			ctxData: []string{"foo", "bar"},
			want:    ValidationErrors{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.WithValue(context.Background(), validationErrorsContextKey, tt.ctxData)

			got := ValidationErrorsFromContext(ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ValidationErrors=%#v, want=%#v", got, tt.want)
			}
		})
	}
}

func TestInertia_SetEncryptHistory(t *testing.T) {
	t.Parallel()

	ctx := SetEncryptHistory(context.Background())

	got, ok := ctx.Value(encryptHistoryContextKey).(bool)
	if !ok {
		t.Fatal("encrypt history from context is not `bool` type")
	}

	want := true

	if got != want {
		t.Fatalf("encryptHistory=%t, want=%t", got, want)
	}
}

func Test_EncryptHistoryFromContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ctxData any
		want    bool
	}{
		{
			name:    "nil",
			ctxData: nil,
			want:    false,
		},
		{
			name:    "false",
			ctxData: false,
			want:    false,
		},
		{
			name:    "true",
			ctxData: true,
			want:    true,
		},
		{
			name:    "wrong type",
			ctxData: []string{"foo", "bar"},
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.WithValue(context.Background(), encryptHistoryContextKey, tt.ctxData)

			got, _ := EncryptHistoryFromContext(ctx)
			if got != tt.want {
				t.Fatalf("encryptHistory=%t, want=%t", got, tt.want)
			}
		})
	}
}

func TestInertia_ClearHistory(t *testing.T) {
	t.Parallel()

	ctx := ClearHistory(context.Background())

	got, ok := ctx.Value(clearHistoryContextKey).(bool)
	if !ok {
		t.Fatal("clear history from context is not `bool` type")
	}

	want := true

	if got != want {
		t.Fatalf("clearHistory=%t, want=%t", got, want)
	}
}

func Test_ClearHistoryFromContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ctxData any
		want    bool
	}{
		{
			name:    "nil",
			ctxData: nil,
			want:    false,
		},
		{
			name:    "false",
			ctxData: false,
			want:    false,
		},
		{
			name:    "true",
			ctxData: true,
			want:    true,
		},
		{
			name:    "wrong type",
			ctxData: []string{"foo", "bar"},
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.WithValue(context.Background(), clearHistoryContextKey, tt.ctxData)

			got := ClearHistoryFromContext(ctx)
			if got != tt.want {
				t.Fatalf("clearHistory=%t, want=%t", got, tt.want)
			}
		})
	}
}
