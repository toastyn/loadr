package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/nesbyte/loadr"
)

// Sets the configuration for the renderer
var ropts = loadr.RendererOpts{
	FS:          os.DirFS("."),
	StripPrefix: "",
	// setting livereload to true allows the app to be built once, and then you can update the HTML and simply refresh the page
	LiveReload: true,
}

// This creates a base component with predefined fragments
// using the html/template
var baseWithComponents = loadr.NewRenderer(ropts).WithComponents("global_components.html")

// Using the base component, index.html is loaded in
// We have now created our index component where it will
// the variable should always be lowercased
var index = baseWithComponents.LoadPages("index.html")

// The defined data that the render function takes in
type IndexData struct {
	Name string
}

// This is the render function where the HTML rendering and custom data model come together

// If the LiveReload option is true every time this function
// is called, it will automatically re-parse the HTML
// If false, the cached version will be used instead
// check out the benchmark (here:) to see the performance difference
func RenderIndex(w io.Writer, d IndexData) error {
	return index.Render(w, d)
}

// Uncomment the line below and see how the program will fail immediately on start
// var _ = baseWithComponents.LoadPages("I-do-not-exist.html")

// Bringing it all together below
func main() {
	r := http.NewServeMux()

	// The rendering is called in here
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		RenderIndex(w, IndexData{"Bob"})
	})

	fmt.Println("Listening on 8080, open http://localhost:8080/")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln(err)
	}
}
