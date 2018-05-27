package schema

import "github.com/godbus/dbus"

// Reason codes
const (
	Expired   = 1
	Dismissed = 2
	Closed    = 3
)

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

// ServerInformation ...
type ServerInformation struct {
	Name        string
	Vendor      string
	Version     string
	SpecVersion string
}

// ActionInvoked ...
type ActionInvoked struct {
	ID        uint32
	ActionKey string
}

// NotificationClosed ...
type NotificationClosed struct {
	ID     uint32
	Reason uint32
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
