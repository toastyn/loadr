package loadr

import (
	"errors"
	"html/template"
	"io"
	"io/fs"
)

// Set the fields to configure the Renderer
type RendererOpts struct {
	FS          fs.FS  // Sets the FS for the renderer
	StripPrefix string // used to strip the prefix from FS
	LiveReload  bool   // Enables/disables template reloading on every render
}

// Composes the std template and keeps track of the existing patterns
type TemplateRender struct {
	template    *template.Template
	patterns    []string
	liveReload  bool
	fs          fs.FS
	stripPrefix string
}

// Helper function to parse new pages
// Function will panic on error
func NewRenderer(r RendererOpts) *TemplateRender {

	if r.StripPrefix == "" {
		r.StripPrefix = "."
	}

	var err error
	r.FS, err = fs.Sub(r.FS, r.StripPrefix)
	if err != nil {
		panic(err)
	}

	if r.FS == nil {
		panic("filesystem not set for HTML renderer")
	}

	return &TemplateRender{
		fs:          r.FS,
		stripPrefix: r.StripPrefix,
		liveReload:  r.LiveReload,
	}
}

// Clones and loads in the pages to be rendered out
func (t TemplateRender) LoadPages(pages ...string) *TemplateRender {
	t.patterns = append(pages, t.patterns...) // Order is important here
	t.template = template.Must(template.ParseFS(t.fs, t.patterns...))
	return &t
}

// Clones and reloads the existing pages
func (t TemplateRender) ReloadPages() *TemplateRender {
	t.template = template.Must(template.ParseFS(t.fs, t.patterns...))
	return &t
}

// Clones and adds the relevant components
func (t TemplateRender) WithComponents(commonComponents ...string) *TemplateRender {
	t.patterns = append(t.patterns, commonComponents...)
	return &t
}

var ErrTemplateNotYetParsed = errors.New("the template has not yet been parsed. Did you run LoadPages or ReloadPages?")

// Given the data, and the liveReload flag renders the HTML
// If the live reload flag is true, the pages will first be reloaded then rendered again
func (t *TemplateRender) Render(wr io.Writer, data any) error {
	if t.liveReload {
		t = t.ReloadPages()
	}

	if t.template == nil {
		return ErrTemplateNotYetParsed
	}

	return t.template.Execute(wr, data)
}
