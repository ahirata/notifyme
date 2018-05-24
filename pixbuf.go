package main

import (
	"github.com/gotk3/gotk3/gdk"
)

func pixbufNewFromData(data []byte, colorspace gdk.Colorspace, hasAlpha bool, bitsPerSample, width, height int) (*gdk.Pixbuf, error) {
	pixbuf, err := gdk.PixbufNew(colorspace, hasAlpha, bitsPerSample, width, height)
	if err != nil {
		return nil, err
	}
	pixels := pixbuf.GetPixels()
	for i := 0; i < len(pixels); i++ {
		pixels[i] = data[i]
	}
	return pixbuf, nil
}
