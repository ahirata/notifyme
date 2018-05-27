package ui

// #cgo pkg-config: gdk-3.0 gtk+-3.0
// #include <gtk/gtk.h>
// #include <gdk/gdk.h>
import "C"

import (
	"errors"
	"github.com/gotk3/gotk3/gdk"
	"sync"
	"unsafe"
)

var errNilPointer = errors.New("cgo returned unexpected nil pointer")

var mutex = &sync.Mutex{}
var memo *gdk.Rectangle

func getWorkarea(display *gdk.Display) (*gdk.Rectangle, error) {
	mutex.Lock()
	if memo != nil {
		mutex.Unlock()
		return memo, nil
	}

	monitor := C.gdk_display_get_primary_monitor((*C.GdkDisplay)(unsafe.Pointer(display.GObject)))
	if monitor == nil {
		return nil, errNilPointer
	}

	gdkRectangle := C.GdkRectangle{}
	C.gdk_monitor_get_workarea(monitor, &gdkRectangle)
	workarea := gdk.WrapRectangle(uintptr(unsafe.Pointer(&gdkRectangle)))

	memo = workarea
	mutex.Unlock()
	return workarea, nil
}
