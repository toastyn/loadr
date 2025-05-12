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
func getRenderOpts() loadr.RendererOpts {
	return loadr.RendererOpts{
		FS:          os.DirFS("."),
		StripPrefix: "pages",
		LiveReload:  true,
	}
}

// This creates a base component with predefined fragments
// using the html/template
var baseWithComponents = loadr.NewRenderer(getRenderOpts()).WithComponents("global_components.html")

// Using the base component, index.html is loaded in
// We have now created our index component where it will
// the variable should always be lowercased
var index = baseWithComponents.LoadFiles("index.html")

// From the .html file we extract the templates

// Using no template name, equivalent to calling Execute on a parsed template file
var indexPage = index.Template(loadr.NoTemplateName)

// To render only a named template
var indexContent = index.Template("content")

// The defined data that the render function takes in
type IndexData struct {
	Name    string
	Content string
}

// This is the render function where the HTML rendering and custom data model come together

// If the LiveReload option is true every time this function
// is called, it will automatically re-parse the HTML
// If false, the cached version will be used instead
// check out the benchmark (here:) to see the performance difference
func RenderIndex(w io.Writer, d IndexData) error {
	return indexPage.Render(w, d)
}

func RenderIndexContent(w io.Writer, d IndexData) error {
	return indexContent.Render(w, d)
}

// Uncomment the line below and see how the program will fail immediately on start
// var _ = baseWithComponents.LoadPages("I-do-not-exist.html")

// Bringing it all together below
func main() {
	r := http.NewServeMux()

	// The rendering is called in here
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		RenderIndex(w, IndexData{"Bob", "SomeContent"})
	})

	r.HandleFunc("/content", func(w http.ResponseWriter, r *http.Request) {
		RenderIndexContent(w, IndexData{"Bob", "SomeContent"})
	})

	fmt.Println("Listening on 8080, open http://localhost:8080/")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln(err)
	}
}
