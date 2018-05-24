package main

import "github.com/godbus/dbus"

// ServerInformation ...
type ServerInformation struct {
	Name        string
	Vendor      string
	Version     string
	SpecVersion string
}

// Server ...
type Server struct {
	conn           *dbus.Conn
	capabilities   []string
	counter        uint32
	defaultExpires int32
	mute           bool
	info           ServerInformation
	store          *WidgetStore
	Outbound       chan Action
}

// WidgetStore ...
type WidgetStore struct {
	widgets []*NotificationWidget
	ids     map[uint32]*NotificationWidget
}

// NotificationHandler ...
type NotificationHandler struct {
	widgets  []*NotificationWidget
	ids      map[uint32]*NotificationWidget
	Outbound chan Action
	Inbound  chan interface{}
}

// Notification ...
type Notification struct {
	ID            uint32
	AppName       string
	ReplacesID    uint32
	AppIcon       string
	Summary       string
	Body          string
	Actions       []interface{}
	Hints         map[string]dbus.Variant
	ExpireTimeout int32
}

// Action ...
type Action struct {
	ID     uint32
	action string
}

// ImageData ...
type ImageData struct {
	Width         int32
	Height        int32
	RowStride     int32
	HasAlpha      bool
	BitsPerSample int32
	Channels      int32
	Data          []byte
}

// ImageData reads the image bytes from the Hints
func (notification *Notification) ImageData() (ImageData, bool) {
	hints := notification.Hints
	variant, found := hints["image-data"]
	if !found {
		return ImageData{}, found
	}

	var imageData []interface{}
	err := dbus.Store([]interface{}{variant}, &imageData)
	if err != nil {
		panic(err)
	}

	image := ImageData{
		Width:         imageData[0].(int32),
		Height:        imageData[1].(int32),
		RowStride:     imageData[2].(int32),
		HasAlpha:      imageData[3].(bool),
		BitsPerSample: imageData[4].(int32),
		Channels:      imageData[5].(int32),
		Data:          imageData[6].([]byte),
	}
	return image, true
}
