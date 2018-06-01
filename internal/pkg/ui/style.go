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
	cssProvider, err := gtk.CssProviderNew()
	if err != nil {
		fmt.Println("Unable to create css provider")
		return
	}
	err = cssProvider.LoadFromPath(os.Getenv("HOME") + "/.config/notifyme/notifyme.css")
	if err != nil {
		fmt.Println("Unable to load css file: ", os.Getenv("HOME")+"/.config/notifyme/notifyme.css")
	}
	err = cssProvider.LoadFromPath("./themes/notifyme.css")
	if err != nil {
		fmt.Println("Unable to load css file")
	}
	screen, err := window.GetScreen()
	if err != nil {
		fmt.Println("Unable get screen")
		return
	}
	gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_USER)
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
