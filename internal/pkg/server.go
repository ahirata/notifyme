package notifyme

import (
	"fmt"
	"github.com/ahirata/notifyme/internal/pkg/store"
	"github.com/ahirata/notifyme/pkg/notifyme/schema"
	"github.com/godbus/dbus"
	"sync/atomic"
	"time"
)

// Server ...
type Server struct {
	conn                     *dbus.Conn
	capabilities             []string
	counter                  uint32
	defaultExpires           int32
	mute                     bool
	info                     schema.ServerInformation
	store                    *store.WidgetStore
	NotificationClosedSignal chan schema.NotificationClosed
}

// ServerNew ...
func ServerNew() Server {
	server := Server{
		capabilities:   []string{"body", "actions", "body-hyperlinks", "body-markup"},
		counter:        0,
		defaultExpires: 5000,
		mute:           false,
		info: schema.ServerInformation{
			Name:        "notifyme",
			Vendor:      "ahirata",
			Version:     "0.0.1",
			SpecVersion: "1.2",
		},
		NotificationClosedSignal: make(chan schema.NotificationClosed, 10),
		store: store.WidgetStoreNew(make(chan schema.ActionInvoked, 10)),
	}
	return server
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

	notification := schema.Notification{
		ID:            ID,
		AppName:       appName,
		ReplacesID:    replacesID,
		AppIcon:       appIcon,
		Summary:       summary,
		Body:          body,
		Actions:       actions,
		Hints:         hints,
		ExpireTimeout: expiration,
	}

	server.store.Put(&notification)

	if expiration > 0 {
		go server.scheduleExpiration(&notification)
	}

	return ID, nil
}

func (server *Server) scheduleExpiration(notification *schema.Notification) {
	select {
	case <-time.After(time.Duration(notification.ExpireTimeout) * time.Millisecond):
		widget := server.store.Get(notification.ID)
		if widget != nil && widget.Notification == notification {
			server.store.Close(notification.ID)
			server.NotificationClosedSignal <- schema.NotificationClosed{ID: notification.ID, Reason: schema.Expired}
		}
	}
}

// CloseNotification causes a notification to be forcefully closed and removed from the user's view
func (server *Server) CloseNotification(id uint32) *dbus.Error {
	fmt.Println("Received: CloseNotification: ", id)
	server.store.Close(id)
	server.NotificationClosedSignal <- schema.NotificationClosed{ID: id, Reason: schema.Closed}
	return nil
}

// CloseLastNotification closes the most recent notification. This is a non-standard message
func (server *Server) CloseLastNotification() *dbus.Error {
	fmt.Println("Received: CloseLastNotification")
	if !server.store.IsEmpty() {
		widget := server.store.CloseLast()
		if widget != nil {
			server.NotificationClosedSignal <- schema.NotificationClosed{ID: widget.Notification.ID, Reason: schema.Dismissed}
		}
	}
	return nil
}

// OpenLastNotification opens the application that sent the most recent notification. This is a non-standard message
func (server *Server) OpenLastNotification() *dbus.Error {
	fmt.Println("Received: OpenLastNotification")
	if !server.store.IsEmpty() {
		server.store.OpenLast("default")
	}
	return nil
}

// ToggleMute controls if future messages will be displayed to the user or not. This is a non-standard message
func (server *Server) ToggleMute() *dbus.Error {
	server.mute = !server.mute
	fmt.Println("Received: ToggleMute. Is muted? ", server.mute)
	return nil
}

// NotificationClosed signals a completed notification which is one that has timed out, or has been dismissed by the user.
func (server *Server) NotificationClosed(id uint32, reason uint32) {
	fmt.Println("Received: NotificationClosed: ", id, reason)
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

func emitNotificationClosed(conn *dbus.Conn, notificationClosed schema.NotificationClosed) {
	go conn.Emit("/org/freedesktop/Notifications", "org.freedesktop.Notifications.NotificationClosed", notificationClosed.ID, notificationClosed.Reason)
}

func emitActionInvoked(conn *dbus.Conn, actionInvoked schema.ActionInvoked) {
	go conn.Emit("/org/freedesktop/Notifications", "org.freedesktop.Notifications.ActionInvoked", actionInvoked.ID, actionInvoked.ActionKey)
}

// Start connects the sever to d-bus to receive message
func (server *Server) Start() {
	conn := connect()
	conn.ExportMethodTable(server.commands(), "/org/freedesktop/Notifications", "org.freedesktop.Notifications")

	for {
		select {
		case notificationClosed := <-server.NotificationClosedSignal:
			fmt.Println("Sending NotificationClosed", notificationClosed)
			emitNotificationClosed(conn, notificationClosed)
		case actionInvoked := <-server.store.ActionInvokedSignal:
			fmt.Println("Sending ActionInvoked", actionInvoked)
			emitActionInvoked(conn, actionInvoked)
		}
	}
}
