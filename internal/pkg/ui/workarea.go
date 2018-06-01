package ui

// #cgo pkg-config: gdk-3.0 gtk+-3.0
// #include <gtk/gtk.h>
// #include <gdk/gdk.h>
import "C"

import (
	"errors"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"unsafe"
)

func getWorkarea(window *gtk.Window) (*gdk.Rectangle, error) {
	screen, err := window.GetScreen()
	if err != nil {
		return nil, err
	}

	display, err := screen.GetDisplay()
	if err != nil {
		return nil, err
	}

	monitor := C.gdk_display_get_primary_monitor((*C.GdkDisplay)(unsafe.Pointer(display.GObject)))
	if monitor == nil {
		return nil, errors.New("cgo returned unexpected nil pointer")
	}

	gdkRectangle := C.GdkRectangle{}
	C.gdk_monitor_get_workarea(monitor, &gdkRectangle)
	workarea := gdk.WrapRectangle(uintptr(unsafe.Pointer(&gdkRectangle)))
	return workarea, nil
}
