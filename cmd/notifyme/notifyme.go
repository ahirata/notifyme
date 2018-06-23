package main

import (
	"flag"
	"github.com/ahirata/notifyme/internal/pkg/server"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	kill := flag.Bool("k", false, "kill notifyme")
	flag.Parse()

	if *kill {
		dbusHandler := notifyme.DbusHandler{}
		dbusHandler.KillServer()
		return
	}

	gtk.Init(nil)

	server := notifyme.ServerNew()

	go server.Start()

	gtk.Main()
}
