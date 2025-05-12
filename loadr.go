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

// Render templates with a specific set of options
type Renderer struct {
	opts       *RendererOpts
	components []string
	files      []*TemplateFile
}

// Files registered containing templates
type TemplateFile struct {
	patterns  []string
	templates []*Template
	opts      *RendererOpts
}

// Inividual template to be rendered
type Template struct {
	template *template.Template
	name     string
	patterns []string
	opts     *RendererOpts
}

// Create a new renderer
func NewRenderer() *Renderer {
	return &Renderer{}
}

// Clones and adds the relevant components
func (r Renderer) WithComponents(commonComponents ...string) *Renderer {
	r.components = append(r.components, commonComponents...)
	return &r
}

// Register files containing templates
func (r *Renderer) LoadFiles(pages ...string) *TemplateFile {
	var file TemplateFile
	file.opts = r.opts
	file.patterns = r.components
	file.patterns = append(pages, file.patterns...)
	r.files = append(r.files, &file)
	return &file
}

// Set the options for the renderer. Verifies the registered files and templates by loading, panics on error.
// Designed to be called in main() or later to set opts programatically
func (r *Renderer) SetOptions(opts RendererOpts) {
	if opts.StripPrefix == "" {
		opts.StripPrefix = "."
	}

	var err error
	opts.FS, err = fs.Sub(opts.FS, opts.StripPrefix)
	if err != nil {
		panic(err)
	}

	if opts.FS == nil {
		panic("filesystem not set for HTML renderer")
	}
	r.opts = &opts

	for _, f := range r.files {
		f.opts = &opts
		for _, t := range f.templates {
			t.opts = &opts
			t.loadTemplate()
		}
	}

}

const NoTemplateName string = ""

// Extract templates from a file.
// Use NoTemplateName for not defined templates.
// If the options are set this will load in the template
func (f *TemplateFile) Template(name string) *Template {
	var t Template
	t.opts = f.opts
	t.patterns = f.patterns
	t.name = name
	if t.opts != nil {
		t.loadTemplate()
	}
	f.templates = append(f.templates, &t)
	return &t
}

var ErrTemplateNotYetParsed = errors.New("the template has not yet been parsed. Did you set the renderer options?")

// Loads the Template from the files
func (t *Template) loadTemplate() {
	name := t.name
	t.template = template.Must(template.ParseFS(t.opts.FS, t.patterns...))
	if t.name != NoTemplateName {
		t.template = t.template.Lookup(name)
	}
}

// Given the data, and the liveReload flag renders the HTML
// If the live reload flag is set, the pages will first be reloaded then rendered again
func (t *Template) Render(wr io.Writer, data any) error {
	if t.template == nil || t.opts == nil {
		return ErrTemplateNotYetParsed
	}
	if t.opts.LiveReload {
		t.loadTemplate()
	}
	return t.template.Execute(wr, data)
}
