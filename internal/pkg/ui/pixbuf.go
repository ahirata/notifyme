package ui

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"strings"
)

func pixbufNewFromData(data []byte, colorspace gdk.Colorspace, hasAlpha bool, bitsPerSample, originalWidth, originalHeight desiredWidth, desiredHeight int) (*gdk.Pixbuf, error) {
	pixbuf, err := gdk.PixbufNew(colorspace, hasAlpha, bitsPerSample, width, height)
	if err != nil {
		return nil, err
	}
	pixels := pixbuf.GetPixels()
	for i := 0; i < len(pixels); i++ {
		pixels[i] = data[i]
	}

	return pixbuf.ScaleSimple(desiredWidth, desiredHeight, gdk.INTERP_BILINEAR)
}

func loadPixbufFromFile(filename string, width, height int) (*gdk.Pixbuf, error) {
	path := strings.Replace(filename, "file://", "", 1)
	return gdk.PixbufNewFromFileAtScale(path, width, height, true)
}
