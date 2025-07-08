package registry

import (
	"sync"
)

type Loader interface {
	Load() error
}

type registry struct {
	loaders    map[Loader]struct{}
	liveReload bool   // If true, sets the Templ to reload on every Render() call
	jsToInject string // JS to inject at the end of the body
}

var store = registry{loaders: make(map[Loader]struct{})}

var mu sync.Mutex

// Adds a BaseRender and it's pattern to the register
func Add(l Loader) {
	_, ok := store.loaders[l]
	mu.Lock()
	if !ok {
		store.loaders[l] = struct{}{}
	}
	mu.Unlock()
}

// Enables or disables live reloading
func SetLiveReload(enabled bool) {
	store.liveReload = enabled
}

// Checks if live reloading is enabled
func LiveReload() bool {
	return store.liveReload
}

func SetJSToInject(b []byte) {
	store.jsToInject = string(b)
}

func JSToInject() string {
	return store.jsToInject
}

// Prepares the templates by loading and validating them
func LoadTemplates() error {

	for loader := range store.loaders {
		err := loader.Load()
		if err != nil {
			return err
		}
	}

	return nil
}

// Should not be used unless you know what you are doing.
//
// Resests the store, this is helpful for tests to reset
// the store if incorrect templates are purposfully provided.
// WARNING: If Reset() is used directly in application logic
// this can remove existing templates, allow Load to silently
// pass and create runtime panics.
func Reset() {
	store = registry{loaders: make(map[Loader]struct{})}
}
