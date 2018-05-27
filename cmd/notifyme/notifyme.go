package main

import (
	"github.com/ahirata/notifyme/internal/pkg"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)

	server := notifyme.ServerNew()

	go server.Start()

	gtk.Main()
}
