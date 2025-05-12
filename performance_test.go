package loadr

import (
	"bytes"
	"os"
	"testing"
)

// With live reload
func BenchmarkWithReload(b *testing.B) {

	var ropts = RendererOpts{
		FS:          os.DirFS("./_examples/basic"),
		StripPrefix: "",
		LiveReload:  true,
	}

	r := NewRenderer().WithComponents("global_components.html")
	t := r.LoadFiles("index.html").Template(NoTemplateName)
	r.SetOptions(ropts)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var bs bytes.Buffer
		t.Render(&bs, nil)
	}
}

// With parsing being cached
func BenchmarkNoReload(b *testing.B) {

	var ropts = RendererOpts{
		FS:          os.DirFS("./_examples/basic/"),
		StripPrefix: "",
		LiveReload:  false,
	}

	r := NewRenderer().WithComponents("global_components.html")
	t := r.LoadFiles("index.html").Template(NoTemplateName)
	r.SetOptions(ropts)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var bs bytes.Buffer
		t.Render(&bs, nil)
	}
}

// Same as above but with named templates
func BenchmarkNamedTemplateReload(b *testing.B) {

	var ropts = RendererOpts{
		FS:          os.DirFS("./_examples/basic"),
		StripPrefix: "",
		LiveReload:  true,
	}

	r := NewRenderer().WithComponents("global_components.html")
	t := r.LoadFiles("index.html").Template("content")
	r.SetOptions(ropts)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var bs bytes.Buffer
		t.Render(&bs, nil)
	}
}

func BenchmarkNamedTemplateNoReload(b *testing.B) {

	var ropts = RendererOpts{
		FS:          os.DirFS("./_examples/basic/"),
		StripPrefix: "",
		LiveReload:  false,
	}

	r := NewRenderer().WithComponents("global_components.html")
	t := r.LoadFiles("index.html").Template("content")
	r.SetOptions(ropts)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var bs bytes.Buffer
		t.Render(&bs, nil)
	}
}
