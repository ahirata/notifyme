package main

import (
	"github.com/gotk3/gotk3/gtk"
)

// StyledContainer ...
type StyledContainer interface {
	GetStyleContext() (*gtk.StyleContext, error)
}

// Container ...
type Container interface {
	Add(gtk.IWidget)
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
