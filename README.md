# Gonertia

<img src="https://user-images.githubusercontent.com/27378369/215432769-35e7b0f5-29a9-41d0-ba79-ca81e624b970.png" style="width: 200px"  alt="gonertia"/>

Gonertia is a Inertia.js server-side adapter for Golang. Visit [inertiajs.com](https://inertiajs.com/) to learn more.

[![Audit Workflow](https://github.com/romsar/gonertia/actions/workflows/audit.yml/badge.svg?branch=master)](https://github.com/romsar/gonertia/actions/workflows/audit.yml?query=branch:master)
[![Go Report Card](https://goreportcard.com/badge/github.com/romsar/gonertia)](https://goreportcard.com/report/github.com/romsar/gonertia)
[![Go Reference](https://godoc.org/github.com/romsar/gonertia?status.svg)](https://pkg.go.dev/github.com/romsar/gonertia)
[![MIT license](https://img.shields.io/badge/LICENSE-MIT-orange.svg)](https://github.com/romsar/gonertia/blob/master/LICENSE)

## Introdution

Inertia allows you to create fully client-side rendered, single-page apps, without the complexity that comes with modern SPAs. It does this by leveraging existing server-side patterns that you already love.

This package based on official Laravel adapter for Inertia.js: [inertiajs/inertia-laravel](https://github.com/inertiajs/inertia-laravel).

## Roadmap
- [ ] Tests
- [ ] Validation errors
- [ ] SSR

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
    "http"
    "time"
    
    inertia "github.com/romsar/gonertia"
)

func main() {
    i, err := inertia.New(
        "https://yourwebsite.com",
        "./ui/templates/root.html",
    )
    if err != nil {
        log.Fatal(err)
    }

    // Now create your HTTP server.
    // Gonertia works well with standard http handlers,
    // but you free to use some frameworks like Gorilla Mux or Chi.
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
    <title>title</title>
    {{ .inertiaHead }}
</head>

<body>
{{ .inertia }}
<script src="/static/js/app.js"></script>
</body>
</html>
```

### More examples

#### Load root template using embed
```go
import "embed"

//go:embed templates
var templateFS embed.FS

// ...

i, err := inertia.New(
    /* ... */
    inertia.WithTemplateFS(templateFS),
)
```

#### Set asset version ([learn more](https://inertiajs.com/asset-versioning))

```go
i, err := inertia.New(
    /* ... */
    inertia.WithVersion("some-version"),
)
```

#### Set asset version by asset url

```go
i, err := inertia.New(
    /* ... */
    inertia.WithAssetURL("/static/js/1f0f8sc6.js"),
)
```

#### Set asset version by manifest file

```go
i, err := inertia.New(
    /* ... */
    inertia.WithManifestFile("./ui/manifest.json"),
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
    inertia.WithLogger(somelogger.New()),
    // or inertia.WithoutLogger(),
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
    "closure": func () (any, error) { return "prop", nil },
    "lazy": inertia.LazyProp(func () (any, error) {
        return "prop", nil
    },
}

i.Render(w, r, "Some/Page", props)
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
If request was Inertia request - user will be redirected via Inertia.js.
If not - user will be redirected via 302 status.
If response is empty - user will be redirected to the previous url.

#### Share template data ([learn more](https://inertiajs.com/responses#root-template-data))

```go
i.ShareTemplateData("title", "Home page")
```

```html
<title>{{ .title }}</title>
```

#### Share template func

```go
i.ShareTemplateFunc("trim", strings.Trim)
```

```html
<title>{{ trim "foo bar" }}</title>
```

#### Pass template data via context (in middleware)

```go
ctx := i.WithTemplateData(r.Context(), "title", "Home page")

// pass it to the next middleware or inertia.Render function via r.WithContext(ctx).
```

#### Share prop globally ([learn more](https://inertiajs.com/shared-data))

```go
i.ShareProp("name", "Roman")
```

#### Pass props via context (in middleware)

```go
ctx := i.WithProp(r.Context(), "name", "Roman")

// pass it to the next middleware or inertia.Render function via r.WithContext(ctx).
```

## License

Gonertia is released under the [MIT License](http://www.opensource.org/licenses/MIT).
