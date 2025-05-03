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

	r := NewRenderer(ropts).WithComponents("global_components.html").LoadPages("index.html")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var bs bytes.Buffer
		r.Render(&bs, nil)
	}
}

// With parsing being cached
func BenchmarkNoReload(b *testing.B) {

	var ropts = RendererOpts{
		FS:          os.DirFS("./_examples/basic/"),
		StripPrefix: "",
		LiveReload:  false,
	}

	r := NewRenderer(ropts).WithComponents("global_components.html").LoadPages("index.html")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var bs bytes.Buffer
		r.Render(&bs, nil)
	}
}
