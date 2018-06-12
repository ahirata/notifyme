package ui

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"os"
)

// StyledContainer ...
type StyledContainer interface {
	GetStyleContext() (*gtk.StyleContext, error)
}

// Container ...
type Container interface {
	Add(gtk.IWidget)
}

// LoadCSSProvider ...
func LoadCSSProvider(window *gtk.Window) {
	cssProvider := cssProviderNew()
	screen, err := window.GetScreen()
	if err != nil {
		fmt.Println("Unable get screen")
		return
	}
	gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_USER)
}

func cssProviderNew() *gtk.CssProvider {
	cssProvider, err := gtk.CssProviderNew()
	if err != nil {
		fmt.Println("Unable to create css provider")
		return nil
	}

	loadCSS("/usr/share/notifyme/themes/notifyme.css", cssProvider)
	loadCSS(os.Getenv("HOME")+"/.config/notifyme/themes/notifyme.css", cssProvider)
	loadCSS("./themes/notifyme.css", cssProvider)

	return cssProvider
}

func loadCSS(path string, cssProvider *gtk.CssProvider) {
	if _, err := os.Stat(path); err != nil {
		fmt.Println("Could not find css file: ", path)
		return
	}
	if err := cssProvider.LoadFromPath(path); err != nil {
		fmt.Println("Unable to load css file: ", path)
		return
	}
	fmt.Println("Loaded css from file", path)
}

// AddClass ...
func AddClass(container StyledContainer, class string) {
	style, err := container.GetStyleContext()
	if err != nil {
		return
	}
	style.AddClass(class)
	style.Save()
}

// AddBox ...
func AddBox(container Container, orientation gtk.Orientation, class string) (*gtk.Box, error) {
	box, err := gtk.BoxNew(orientation, 0)
	if err != nil {
		return nil, err
	}
	container.Add(box)
	AddClass(box, class)

	return box, nil
}
