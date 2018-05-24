package main

import (
	"github.com/gotk3/gotk3/gtk"
)

func notificationHandlerNew() *NotificationHandler {
	return &NotificationHandler{
		ids:      make(map[uint32]*NotificationWidget),
		Outbound: make(chan Action, 10),
		Inbound:  make(chan interface{}, 10),
	}
}

func serverNew() Server {
	server := Server{
		capabilities:   []string{"body", "actions", "body-hyperlinks", "body-markup", "body-images", "action-icons"},
		counter:        0,
		defaultExpires: 5000,
		mute:           false,
		info: ServerInformation{
			Name:        "ahirata-notification-server",
			Vendor:      "ahirata",
			Version:     "1.0",
			SpecVersion: "1.2",
		},
		Outbound: make(chan Action, 10),
		store:    WidgetStoreNew(),
	}
	return server
}

func main() {
	gtk.Init(nil)

	server := serverNew()

	go server.start()

	gtk.Main()
}
