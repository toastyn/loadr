package loadr

import (
	"errors"
	"fmt"
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

// Renderer with a set of options and components
type Renderer struct {
	opts       RendererOpts
	components []string
}

// Represents a loaded file contianing templates
type TemplateFile struct {
	patterns         []string
	templates        map[string]*template.Template
	opts             RendererOpts
	topLevelTemplate *template.Template
}

// Individual template to be rendered
type Template struct {
	template   *template.Template
	name       string
	patterns   []string
	liveReload bool
	fs         fs.FS
}

// Helper function to parse new pages
// Function will panic on error
func NewRenderer(r RendererOpts) *Renderer {
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

	return &Renderer{
		opts: r,
	}
}

// Clones and adds the relevant components
func (r Renderer) WithComponents(commonComponents ...string) *Renderer {
	r.components = append(r.components, commonComponents...)
	return &r
}

// Load in a group of files containing templates
func (r Renderer) LoadFiles(pages ...string) *TemplateFile {
	var file TemplateFile
	file.opts = r.opts
	file.patterns = r.components
	file.patterns = append(pages, file.patterns...)
	file.topLevelTemplate = template.Must(template.ParseFS(r.opts.FS, file.patterns...))

	file.templates = make(map[string]*template.Template)
	for _, subTmpl := range file.topLevelTemplate.Templates() {
		file.templates[subTmpl.Name()] = subTmpl
	}
	return &file
}

const NoTemplateName string = ""

func (f TemplateFile) Template(name string) *Template {
	var t Template
	t.fs = f.opts.FS
	t.liveReload = f.opts.LiveReload
	t.patterns = f.patterns
	t.name = name
	if name == NoTemplateName {
		t.name = f.topLevelTemplate.Name()
		t.template = f.topLevelTemplate
		return &t
	}
	var ok bool
	t.template, ok = f.templates[t.name]
	if !ok {
		panic(fmt.Sprintf("template %s not found", name))
	}
	return &t
}

var ErrTemplateNotYetParsed = errors.New("the template has not yet been parsed. Did you run LoadPages or ReloadPages?")

// Clones and reloads the existing template
func (t Template) ReloadTemplate() *Template {
	name := t.name
	t.template = template.Must(template.ParseFS(t.fs, t.patterns...))
	if t.name != NoTemplateName {
		t.template = t.template.Lookup(name)
	}
	return &t
}

// Given the data, and the liveReload flag renders the HTML
// If the live reload flag is set, the pages will first be reloaded then rendered again
func (t *Template) Render(wr io.Writer, data any) error {
	if t.liveReload {
		t = t.ReloadTemplate()
	}

	if t.template == nil {
		return ErrTemplateNotYetParsed
	}
	return t.template.Execute(wr, data)
}
