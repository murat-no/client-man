package main

import (
	"embed"
	"fmt"
	"sync"

	"fyne.io/fyne/v2"
)

//go:embed icons/*.svg
var embeddedIcons embed.FS

var (
	iconCache   = make(map[string]fyne.Resource)
	iconCacheMu sync.Mutex
)

// loadIconResource returns an embedded SVG resource, falling back to the supplied resource on error.
func loadIconResource(name string, fallback fyne.Resource) fyne.Resource {
	iconCacheMu.Lock()
	defer iconCacheMu.Unlock()

	if cached, ok := iconCache[name]; ok {
		return cached
	}

	iconPath := "icons/" + name + ".svg"
	data, err := embeddedIcons.ReadFile(iconPath)
	if err != nil {
		fyne.LogError(fmt.Sprintf("icon %s could not be loaded", iconPath), err)
		iconCache[name] = fallback
		return fallback
	}

	resource := fyne.NewStaticResource(name+".svg", data)
	iconCache[name] = resource
	return resource
}
