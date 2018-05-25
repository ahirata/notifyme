package main

import (
	"fmt"
	"github.com/godbus/dbus"
	"github.com/gotk3/gotk3/glib"
	"sync/atomic"
	"time"
)

const (
	expired   = 1
	dismissed = 2
	closed    = 3
)

var reasons = map[string]uint32{
	"expired":   expired,
	"dismissed": dismissed,
	"closed":    closed,
}

// GetServerInformation returns the information on the server. Specifically, the server name, vendor, and version number
func (server *Server) GetServerInformation() (string, string, string, string, *dbus.Error) {
	fmt.Println("Received: GetServerInformation")
	return server.info.Name, server.info.Vendor, server.info.Version, server.info.SpecVersion, nil
}

// GetCapabilities returns an array of strings. Each string describes an optional capability implemented by the server
func (server *Server) GetCapabilities() ([]string, *dbus.Error) {
	fmt.Println("Received: GetCapabilities")
	return server.capabilities, nil
}

// Notify sends a notification to this notification server
func (server *Server) Notify(appName string, replacesID uint32, appIcon string, summary string, body string, actions []interface{}, hints map[string]dbus.Variant, expireTimeout int32) (uint32, *dbus.Error) {
	fmt.Printf("Received: Notify(%s, %d, %s, %s, %s, %v, %d)\n", appName, replacesID, appIcon, summary, body, actions, expireTimeout)

	var ID uint32
	if replacesID > 0 {
		ID = replacesID
	} else {
		ID = atomic.AddUint32(&server.counter, 1)
	}

	expiration := expireTimeout
	if expiration < 0 {
		expiration = server.defaultExpires
	}

	if server.mute {
		return ID, nil
	}

	notification := Notification{
		ID:            ID,
		AppName:       appName,
		ReplacesID:    replacesID,
		AppIcon:       appIcon,
		Summary:       summary,
		Body:          body,
		Actions:       actions,
		Hints:         hints,
		ExpireTimeout: expiration,
		timestamp:     time.Now(),
	}

	server.showWidget(&notification)

	if expiration > 0 {
		glib.TimeoutAdd(uint(expiration), func() {
			widget := server.store.Get(notification.ID)
			if widget != nil && widget.Notification == &notification {
				server.closeWidget(notification.ID, "expired")
			}
		})
	}

	return ID, nil
}

func (server *Server) showWidget(notification *Notification) {
	glib.IdleAdd(func() {
		widget := server.store.Remove(notification.ID)
		if widget != nil {
			widget.ReplaceNotification(notification)
			widget.Show()
			server.store.Push(widget)
		} else {
			widget, err := NotificationWidgetNew(notification, server.Outbound)
			if err != nil {
				panic(err)
			}
			server.store.Push(widget)
			widget.Show()
		}
	})
}

// CloseNotification causes a notification to be forcefully closed and removed from the user's view
func (server *Server) CloseNotification(id uint32) *dbus.Error {
	fmt.Printf("Received: CloseNotification %d\n", id)
	glib.IdleAdd(func() {
		server.closeWidget(id, "closed")
	})
	return nil
}

func (server *Server) closeWidget(id uint32, reason string) {
	widget := server.store.Remove(id)
	if widget != nil {
		widget.Close(reason)
	}
}

// CloseLastNotification closes the most recent notification. This is a non-standard message
func (server *Server) CloseLastNotification() *dbus.Error {
	fmt.Println("Received: CloseLastNotification")
	glib.IdleAdd(func() {
		if !server.store.IsEmpty() {
			widget := server.store.Pop()
			widget.Close("dismiss")
		}
	})
	return nil
}

// OpenLastNotification opens the application that sent the most recent notification. This is a non-standard message
func (server *Server) OpenLastNotification() *dbus.Error {
	fmt.Println("Received: OpenLastNotification")
	glib.IdleAdd(func() {
		if !server.store.IsEmpty() {
			widget := server.store.Pop()
			widget.OpenApp()
		}
	})
	return nil
}

// ToggleMute controls if future messages will be displayed to the user or not. This is a non-standard message
func (server *Server) ToggleMute() *dbus.Error {
	server.mute = !server.mute
	fmt.Printf("Received: ToggleMute. Is muted? %t\n", server.mute)
	return nil
}

// NotificationClosed signals a completed notification which is one that has timed out, or has been dismissed by the user.
func (server *Server) NotificationClosed(id uint32, reason uint32) {
	fmt.Printf("NotificationClosed %d: %d\n", id, reason)
	server.conn.Emit("/org/freedesktop/Notifications", "org.freedesktop.Notifications.NotificationClosed", id, reason)
}

func connect() *dbus.Conn {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	reply, err := conn.RequestName("org.freedesktop.Notifications", dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		panic("Name already taken")
	}
	fmt.Println("Connected to dbus")
	return conn
}

func (server *Server) commands() map[string]interface{} {
	methodTable := make(map[string]interface{})
	methodTable["GetServerInformation"] = server.GetServerInformation
	methodTable["GetCapabilities"] = server.GetCapabilities
	methodTable["Notify"] = server.Notify
	methodTable["CloseNotification"] = server.CloseNotification
	methodTable["CloseLastNotification"] = server.CloseLastNotification
	methodTable["OpenLastNotification"] = server.OpenLastNotification
	methodTable["ToggleMute"] = server.ToggleMute
	return methodTable
}

func closeEmitter(conn *dbus.Conn) func(uint32, uint32) {
	return func(id uint32, reason uint32) {
		conn.Emit("/org/freedesktop/Notifications", "org.freedesktop.Notifications.NotificationClosed", id, reason)
	}
}

func actionEmitter(conn *dbus.Conn) func(uint32, string) {
	return func(id uint32, action string) {
		conn.Emit("/org/freedesktop/Notifications", "org.freedesktop.Notifications.ActionInvoked", id, action)
	}
}

// Start connects the sever to d-bus to receive message
func (server *Server) start() {
	conn := connect()
	conn.ExportMethodTable(server.commands(), "/org/freedesktop/Notifications", "org.freedesktop.Notifications")
	emitClosed := closeEmitter(conn)
	emitAction := actionEmitter(conn)

	for {
		select {
		case action := <-server.Outbound:
			fmt.Println("Received action", action)
			if reason, exists := reasons[action.action]; exists {
				emitClosed(action.ID, reason)
			} else {
				emitAction(action.ID, action.action)
			}
		}
	}
}
