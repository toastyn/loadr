package core

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/nesbyte/loadr/livereload"
	"github.com/nesbyte/loadr/registry"
)

func NewTemplate[T, U any](tc *TemplateContext[T], pattern string, data U) *Templ[T, U] {
	t := Templ[T, U]{tc: tc, data: data, usePattern: pattern}

	registry.Add(&t)

	return &t
}

type Templ[T, U any] struct {
	t          *template.Template
	tc         *TemplateContext[T]
	data       U
	usePattern string
}

var ErrNoBaseOrPatternFound = errors.New("no basetemplate nor patterns have been provided")

type LoadingError struct {
	BaseTemplates []string
	WithTemplates []string
	UsePattern    string
	Err           error
}

func (e *LoadingError) Error() string {
	return fmt.Sprintf("basetemplates %q with templates %q and template pattern %q failed: %s", e.BaseTemplates, strings.Join(e.WithTemplates, ", "), e.UsePattern, e.Err.Error())
}

func (e *LoadingError) Unwrap() error {
	return e.Err
}

func newLoadingError[T, U any](t *Templ[T, U], err error) error {
	return &LoadingError{t.tc.baseTemplates, t.tc.withTemplates, t.usePattern, err}
}

var ErrNoConfigProvided = errors.New("no config provided")
var ErrTemplateParse = errors.New("template parse error")
var ErrInvalidTemplateData = errors.New("invalid template data")

// Base data used to define the data passed in to the
// template
type BaseData[T any, U any] struct {
	B T // BaseData passed in on every Render() call
	D U // Data passed in explicitly by the Render(data) call
}

// Loads, validates and registers the template.
// This should rarely be called directly
func (t *Templ[T, U]) Load() error {
	// Immeditately run on load
	if t.tc.onLoad != nil {
		err := t.tc.onLoad()
		if err != nil {
			return err
		}
	}

	if t.tc.config == nil {
		return ErrNoConfigProvided
	}

	patterns := []string{}
	patterns = append(patterns, t.tc.baseTemplates...)
	patterns = append(patterns, t.tc.withTemplates...)

	if len(patterns) == 0 {
		return newLoadingError(t, ErrNoBaseOrPatternFound)
	}

	// Parse and cache the template
	var err error
	t.t, err = template.New("").Funcs(t.tc.funcMap).ParseFS(t.tc.config.FS, patterns...)
	if err != nil {
		return newLoadingError(t, fmt.Errorf("%w: %v", ErrTemplateParse, err))
	}

	// Try to execute the template using the sample data provided
	bs := []byte{}
	w := bytes.NewBuffer(bs)
	err = t.t.ExecuteTemplate(w, t.usePattern, BaseData[T, U]{B: *t.tc.baseData, D: t.data})
	if err != nil {
		return newLoadingError(t, fmt.Errorf("%w has a .B or .D prefix been included for the field?: %v", ErrInvalidTemplateData, err))
	}

	return nil
}

// Renders the template to a writer with the base data
// and data of the loaded type.
// The data injected into a struct is of the form:
//
//	{
//			B: any // Base data
//			D: any // Data as passed in through the Render
//	}
//
// Even if no base data has been provided, the template will be provided
// in the above form.
//
// If live reloading is enabled, JS is injected at the end of the body.
func (t *Templ[T, U]) Render(w io.Writer, data U) {
	d := BaseData[T, U]{B: *t.tc.baseData, D: data}

	// In production rendering is short and simple
	if !registry.LiveReload() {
		err := t.t.ExecuteTemplate(w, t.usePattern, d)
		if err != nil {
			panic(&LoadingError{t.tc.baseTemplates, t.tc.withTemplates, t.usePattern, fmt.Errorf("execute template error in render %s", err)})
		}
		return
	}

	// Reload the component
	err := t.Load()
	if err != nil {
		w.Write([]byte(registry.JSToInject()))

		livereload.LiveReloadCustomErrorHandler(err)
		return
	}

	// Capture the output to a buffer
	var buf bytes.Buffer

	err = t.t.ExecuteTemplate(&buf, t.usePattern, d)
	if err != nil {
		panic(&LoadingError{t.tc.baseTemplates, t.tc.withTemplates, t.usePattern, fmt.Errorf("execute template error in render %s", err)})
	}

	html := buf.String()
	idx := strings.LastIndex(strings.ToLower(html), "</body>")
	if idx != -1 {
		html = html[:idx] + registry.JSToInject() + html[idx:]
	}

	w.Write([]byte(html))
}
