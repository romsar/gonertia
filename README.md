# Inertia.js Go Adapter

<img src="https://user-images.githubusercontent.com/27378369/215351322-26995373-ca54-4dfa-b000-59908bbf4d5a.png" style="width: 200px" />

This is a Inertia.js server-side adapter for Golang. Visit [inertiajs.com](https://inertiajs.com/) to learn more.

## Introdution

This package based on [petaki/inertia-go](https://github.com/petaki/inertia-go), but with more striving for similarity [inertia-laravel](https://github.com/inertiajs/inertia-laravel):
- Middleware (redirect back if empty/303 redirect status for post/put/patch).
- Lazy and closure props support.
- Redirect with `inertia.Location` function works for inertia and non-inertia requests.
- Template directives `inertia` and `inertiaHead` just like in [inertia-laravel](https://github.com/inertiajs/inertia-laravel).
- Asset versioning by asset url or manifest file.
- Some other differences like testing coverage or more complex examples.

The purpose of developing this package is to implement all the functions that are available in the official Laravel package, as well as their maximum similar implementation.

## Installation
Install using `go get` command:
```shell
go get github.com/romsar/gonertia
```

## Usage

### Basic example
main.go
```go
package main

import (
    inertia "github.com/romsar/gonertia"
    "log"
    "http"
    "time"
)

func main() {
    i, err := inertia.New(
        "https://yourwebsite.com",
        "./root.html",
    )
    if err != nil {
        log.Fatal(err)
    }

    mux := http.NewServeMux()

    mux.Handle("/home", i.Middleware(homeHandler(i)))
}

func homeHandler(i *inertia.Inertia) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        props := inertia.Props{
            "some": "data",
        }
		
        if err := i.Render(w, r, "Index", props); err != nil {
           handleServerErr(w, err)
        }
    }

    return http.HandlerFunc(fn)
}
```

root.html
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

#### Load root template using embed

```go
import "embed"

//go:embed templates
var templateFS embed.FS

// ...

i, err := inertia.New(
    "https://yourwebsite.com",
    "templates/root.html",
    inertia.WithTemplateFS(templateFS),
)
```

#### With asset version

```go
i, err := inertia.New(
    "https://yourwebsite.com",
    "templates/root.html",
    inertia.WithVersion("some-version"),
)
```

or

```go
i, err := inertia.New(
    "https://yourwebsite.com",
    "templates/root.html",
    inertia.WithAssetURL("/static/js/1f0f8sc6.js"),
)
```

or

```go
i, err := inertia.New(
    "https://yourwebsite.com",
    "templates/root.html",
    inertia.WithManifestFile("./ui/manifest.json"),
)
```

#### With your marshal JSON function

```go
i, err := inertia.New(
    "https://yourwebsite.com",
    "templates/root.html",
    inertia.WithMarshalJSON(jsoniter.Marshal),
)
```

#### With your logger

```go
i, err := inertia.New(
    "https://yourwebsite.com",
    "templates/root.html",
    inertia.WithLogger(somelogger.New()),
)
```

or disable logging at all:

```go
i, err := inertia.New(
    "https://yourwebsite.com",
    "templates/root.html",
    inertia.WithoutLogger(),
)
```

#### With custom container id

```go
i, err := inertia.New(
    "https://yourwebsite.com",
    "templates/root.html",
    inertia.WithContainerID("inertia"),
)
```

#### Closure and lazy props ([learn more](https://inertiajs.com/partial-reloads))

```go
func homeHandler(i *inertia.Inertia) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        props := inertia.Props{
            "time": func () (any, error) { return time.Now().String(), nil },
            "lazy": inertia.LazyProp(func () (any, error) {
                return "lazy data", nil
            }),
        }
		
        if err := i.Render(w, r, "Index", props); err != nil {
           handleServerErr(w, err)
        }
    }

    return http.HandlerFunc(fn)
}
```

#### Redirects

```go
func homeHandler(i *inertia.Inertia) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        i.Location(w, r, "/some-url")
    }

    return http.HandlerFunc(fn)
}
```

If response is empty - user will be redirected to the previous url.

#### Share template data

```go
i, err := inertia.New(/* ... */)

i.ShareTemplateData("title", "Home page")
```

```html
<title>{{ .title }}</title>
```

#### Share template func

```go
i, err := inertia.New(/* ... */)

i.ShareTemplateFunc("trim", strings.Trim)
```

```html
<title>{{ trim " foo bar " }}</title>
```

#### Share prop globally

```go
i, err := inertia.New(/* ... */)

i.ShareProp("name", "Roman")
```

#### Pass template data via context (in middleware/handler)

```go
ctx := i.WithTemplateData(r.Context(), "title", "Home page")

// pass it to the next middleware or inertia.Render function. 
```

#### Pass props via context (in middleware/handler)

```go
ctx := i.WithProp(r.Context(), "name", "Roman")

// pass it to the next middleware or inertia.Render function. 
```

## Roadmap
- [ ] Validation errors
- [ ] SSR
