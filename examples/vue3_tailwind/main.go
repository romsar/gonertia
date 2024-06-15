package main

import (
	"encoding/json"
	"fmt"
	inertia "github.com/romsar/gonertia"
	"log"
	"net/http"
	"os"
)

func main() {
	i := initInertia()

	mux := http.NewServeMux()

	mux.Handle("/home", i.Middleware(homeHandler(i)))
	mux.Handle("/secondary", i.Middleware(secondaryHandler(i)))
	mux.Handle("/build/", http.StripPrefix("/build/", http.FileServer(http.Dir("./public/build"))))

	http.ListenAndServe(":3000", mux)
}

func initInertia() *inertia.Inertia {
	manifestPath := "./public/build/manifest.json"

	i, err := inertia.New(
		"./resources/views/root.html",
		inertia.WithVersionFromFile(manifestPath),
	)
	if err != nil {
		log.Fatal(err)
	}

	i.ShareTemplateFunc("vite", vite(manifestPath))

	return i
}

func vite(manifestPath string) func(path string) (string, error) {
	f, err := os.Open(manifestPath)
	if err != nil {
		log.Fatalf("cannot open provided vite manifest file: %s", err)
	}
	defer f.Close()

	viteAssets := make(map[string]*struct {
		File   string `json:"file"`
		Source string `json:"src"`
	})
	err = json.NewDecoder(f).Decode(&viteAssets)
	if err != nil {
		log.Fatalf("cannot unmarshal vite manifest file to json: %s", err)
	}

	return func(path string) (string, error) {
		if val, ok := viteAssets[path]; ok {
			return val.File, nil
		}
		return "", fmt.Errorf("asset %q not found", path)
	}
}

func homeHandler(i *inertia.Inertia) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := i.Render(w, r, "Home/Index", inertia.Props{
			"text": "world",
		})

		if err != nil {
			handleServerErr(w, err)
			return
		}
	}

	return http.HandlerFunc(fn)
}

func secondaryHandler(i *inertia.Inertia) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := i.Render(w, r, "Home/Secondary")

		if err != nil {
			handleServerErr(w, err)
			return
		}
	}

	return http.HandlerFunc(fn)
}

func handleServerErr(w http.ResponseWriter, err error) {
	log.Printf("http error: %s\n", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("server error"))
}
