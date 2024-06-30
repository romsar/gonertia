# Gonertia

<img src="https://user-images.githubusercontent.com/27378369/215432769-35e7b0f5-29a9-41d0-ba79-ca81e624b970.png" style="width: 200px"  alt="gonertia"/>

Gonertia is a well-tested and zero-dependency Inertia.js server-side adapter for Golang. Visit [inertiajs.com](https://inertiajs.com/) to learn more.

[![Audit Workflow](https://github.com/romsar/gonertia/actions/workflows/audit.yml/badge.svg?branch=master)](https://github.com/romsar/gonertia/actions/workflows/audit.yml?query=branch:master)
[![Go Report Card](https://goreportcard.com/badge/github.com/romsar/gonertia)](https://goreportcard.com/report/github.com/romsar/gonertia)
[![Go Reference](https://godoc.org/github.com/romsar/gonertia?status.svg)](https://pkg.go.dev/github.com/romsar/gonertia)
[![MIT license](https://img.shields.io/badge/LICENSE-MIT-orange.svg)](https://github.com/romsar/gonertia/blob/master/LICENSE)

## Introduction

Inertia allows you to create fully client-side rendered, single-page apps, without the complexity that comes with modern SPAs. It does this by leveraging existing server-side patterns that you already love.

This package based on the official Laravel adapter for Inertia.js: [inertiajs/inertia-laravel](https://github.com/inertiajs/inertia-laravel).

## Roadmap
- [x] Tests
- [x] Helpers for testing
- [x] Helpers for validation errors
- [x] Examples
- [x] SSR

## Installation
Install using `go get` command:
```shell
go get github.com/romsar/gonertia
```

## Usage

### Basic example
Initialize Gonertia in your `main.go`:
```go
package main

import (
    "log"
    "net/http"
	
    inertia "github.com/romsar/gonertia"
)

func main() {
    i, err := inertia.New(rootHTMLString)
    // i, err := inertia.NewFromFile("resources/views/root.html")
    // i, err := inertia.NewFromReader(rootHTMLReader)
    // i, err := inertia.NewFromBytes(rootHTMLBytes)
    if err != nil {
        log.Fatal(err)
    }

    // Now create your HTTP server.
    // Gonertia works well with standard http server library,
    // but you are free to use some external routers like Gorilla Mux or Chi.
    mux := http.NewServeMux()

    mux.Handle("/home", i.Middleware(homeHandler(i)))
}

func homeHandler(i *inertia.Inertia) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        err := i.Render(w, r, "Home/Index", inertia.Props{
            "some": "data",
        })
		
        if err != nil {
            handleServerErr(w, err)
            return
        }
    }

    return http.HandlerFunc(fn)
}
```

Create `root.html` template:
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <!-- Put here your styles, meta and other stuff -->
    {{ .inertiaHead }}
</head>

<body>
{{ .inertia }}
<script type="module" src="/build/assets/app.js"></script>
</body>
</html>
```

You can also find more examples and starter kits in [examples folder](./examples) or in section above.

### More examples

#### Set asset version ([learn more](https://inertiajs.com/asset-versioning))

```go
i, err := inertia.New(
    /* ... */
    inertia.WithVersion("some-version"), // by any string
    inertia.WithVersionFromFile("./public/build/manifest.json"), // by file checksum
)
```

#### Replace standard JSON marshall function

```go
import jsoniter "github.com/json-iterator/go"

// ...

i, err := inertia.New(
    /* ... */, 
    inertia.WithMarshalJSON(jsoniter.Marshal),
)
```

#### Use your logger

```go
i, err := inertia.New(
    /* ... */
    inertia.WithLogger(), // default logger
    inertia.WithLogger(somelogger.New()), // custom logger
)
```

#### Set custom container id

```go
i, err := inertia.New(
    /* ... */
    inertia.WithContainerID("inertia"),
)
```

#### Closure and lazy props ([learn more](https://inertiajs.com/partial-reloads))

```go
props := inertia.Props{
    "regular": "prop",
    "closure": func () (any) { return "prop" },
    "closure_with_err": func () (any, error) { return "prop", nil },
    "lazy": inertia.LazyProp(func () (any, error) {
        return "prop", nil
    },
}

i.Render(w, r, "Some/Page", props)
```

Also, Gonertia have support for `always` props ([learn more](https://github.com/inertiajs/inertia-laravel/pull/627)):

```go
props := inertia.Props{
    "foo": inertia.AlwaysProp{"bar"},
}
```

#### Redirects ([learn more](https://inertiajs.com/redirects))

```go
func homeHandler(i *inertia.Inertia) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        i.Location(w, r, "/some-url")
    }

    return http.HandlerFunc(fn)
}
```

NOTES:
If response is empty - user will be redirected to the previous url, just like in Laravel official adapter.

To manually redirect back, you can use `Back` helper:

```go
i.Back(w, r)
```

#### Share template data ([learn more](https://inertiajs.com/responses#root-template-data))

```go
i.ShareTemplateData("title", "Home page")
```

```html
<h1>{{ .title }}</h1>
```

#### Share template func

```go
i.ShareTemplateFunc("trim", strings.TrimSpace)
```

```html
<h1>{{ trim " foo bar " }}</h1>
```

#### Pass template data via context (in middleware)

```go
ctx := inertia.WithTemplateData(r.Context(), "title", "Home page")

// pass it to the next middleware or inertia.Render function using r.WithContext(ctx).
```

#### Share prop globally ([learn more](https://inertiajs.com/shared-data))

```go
i.ShareProp("name", "Roman")
```

#### Pass props via context (in middleware)

```go
ctx := inertia.WithProp(r.Context(), "name", "Roman")
// or inertia.WithProps(r.Context(), inertia.Props{"name": "Roman"})

// pass it to the next middleware or inertia.Render function using r.WithContext(ctx).
```

#### Validation errors ([learn more](https://inertiajs.com/validation))

```go
ctx := inertia.WithValidationError(r.Context(), "some_field", "some error")
// or inertia.WithValidationErrors(r.Context(), inertia.ValidationErrors{"some_field": "some error"})

// pass it to the next middleware or inertia.Render function using r.WithContext(ctx).
```

#### SSR (Server Side Rendering) ([learn more](https://inertiajs.com/server-side-rendering))

To enable server side rendering you have to provide an option in place where you initialize Gonertia:

```go
i, err := inertia.New(
/* ... */
    inertia.WithSSR(), // default is http://127.0.0.1:13714
    inertia.WithSSR("http://127.0.0.1:1234"), // custom url
)
```

Also, you have to use asset bundling tools like [Vite](https://vitejs.dev/) or [Webpack](https://webpack.js.org/) (especially with [Laravel Mix](https://laravel-mix.com/)). The setup will vary depending on this choice, you can read more about it in [official docs](https://inertiajs.com/server-side-rendering) or check an [example](./examples/vue3_tailwind) that works on Vite. 

#### Testing

Of course, this package provides convenient interfaces for testing!

```go
func TestHomepage(t *testing.T) {
    body := ... // get an HTML or JSON using httptest package or real HTTP request. 
	
    // ...
	
    assertable := inertia.Assert(t, body) // from io.Reader body
    // OR
    assertable := inertia.AssertFromBytes(t, body) // from []byte body
    // OR
    assertable := inertia.AssertFromString(t, body) // from string body
	
    // now you can do assertions using assertable.Assert[...] methods:
    assertable.AssertComponent("Foo/Bar")
    assertable.AssertVersion("foo bar")
    assertable.AssertURL("https://example.com")
    assertable.AssertProps(inertia.Props{"foo": "bar"})

    // or work with the data yourself:
    assertable.Component // Foo/Bar
    assertable.Version // foo bar
    assertable.URL // https://example.com
    assertable.Props // inertia.Props{"foo": "bar"}
    assertable.Body // full response body
}
```

## More community adapters

Also, you can check one more golang adapter called [petaki/inertia-go](https://github.com/petaki/inertia-go).

Full list of community adapters is located on [inertiajs.com](https://inertiajs.com/community-adapters).

## Credits

This package is based on [inertiajs/inertia-laravel](https://github.com/inertiajs/inertia-laravel) and uses some ideas of [petaki/inertia-go](https://github.com/petaki/inertia-go).

## License

Gonertia is released under the [MIT License](http://www.opensource.org/licenses/MIT).
