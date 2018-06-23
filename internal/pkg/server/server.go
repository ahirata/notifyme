package notifyme

import (
	"fmt"
	"github.com/ahirata/notifyme/internal/pkg/store"
	"github.com/ahirata/notifyme/internal/pkg/ui"
	"github.com/ahirata/notifyme/pkg/notifyme/schema"
	"github.com/godbus/dbus"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"sync/atomic"
	"time"
)

// Server ...
type Server struct {
	conn                     *dbus.Conn
	capabilities             []string
	counter                  uint32
	defaultTimeout           int32
	mute                     bool
	info                     schema.ServerInformation
	store                    store.WidgetStore
	NotificationClosedSignal chan schema.NotificationClosed
	ActionInvokedSignal      chan schema.ActionInvoked
}

// ServerNew ...
func ServerNew() Server {
	server := Server{
		capabilities:   []string{"body", "actions", "body-hyperlinks", "body-markup"},
		counter:        0,
		defaultTimeout: 10000,
		mute:           false,
		info: schema.ServerInformation{
			Name:        "notifyme",
			Vendor:      "ahirata",
			Version:     "0.0.1",
			SpecVersion: "1.2",
		},
		NotificationClosedSignal: make(chan schema.NotificationClosed, 10),
		ActionInvokedSignal:      make(chan schema.ActionInvoked, 10),
		store:                    store.WidgetStore{},
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

	notification := schema.Notification{
		ID:            server.notificationID(replacesID),
		AppName:       appName,
		ReplacesID:    replacesID,
		AppIcon:       appIcon,
		Summary:       summary,
		Body:          body,
		Actions:       actions,
		Hints:         hints,
		ExpireTimeout: server.notificationTimeout(expireTimeout),
	}

	if server.mute {
		return notification.ID, nil
	}

	glib.IdleAdd(func() {
		widget := server.store.Get(notification.ID)
		if widget != nil {
			widget.ReplaceNotification(&notification)
			return
		}

		widget, err := ui.NotificationWidgetNew(&notification, server.store.MinY(), server.ActionInvokedSignal)
		if err != nil {
			fmt.Println("Error building widget", err)
		}
		server.store.Push(widget)
		widget.Show()
	})

	if notification.ExpireTimeout > 0 {
		go server.scheduleExpiration(&notification)
	}

	return notification.ID, nil
}

func (server *Server) notificationID(replacesID uint32) uint32 {
	if replacesID > 0 {
		return replacesID
	}
	return atomic.AddUint32(&server.counter, 1)
}

func (server *Server) notificationTimeout(requestedTimeout int32) int32 {
	if requestedTimeout < 0 {
		return server.defaultTimeout
	}
	return requestedTimeout
}

func (server *Server) scheduleExpiration(notification *schema.Notification) {
	select {
	case <-time.After(time.Duration(notification.ExpireTimeout) * time.Millisecond):
		glib.IdleAdd(func() {
			if widget := server.store.Get(notification.ID); widget == nil || widget.Notification != notification {
				return
			}
			if removed := server.store.Remove(notification.ID); removed != nil {
				removed.Close()
			}

			server.NotificationClosedSignal <- schema.NotificationClosed{ID: notification.ID, Reason: schema.Expired}
		})
	}
}

// CloseNotification causes a notification to be forcefully closed and removed from the user's view
func (server *Server) CloseNotification(id uint32) *dbus.Error {
	fmt.Println("Received: CloseNotification: ", id)
	glib.IdleAdd(func() {
		if removed := server.store.Remove(id); removed != nil {
			removed.Close()
		}
		server.NotificationClosedSignal <- schema.NotificationClosed{ID: id, Reason: schema.Closed}
	})
	return nil
}

// CloseLastNotification closes the most recent notification. This is a non-standard message
func (server *Server) CloseLastNotification() *dbus.Error {
	fmt.Println("Received: CloseLastNotification")
	glib.IdleAdd(func() {
		if server.store.IsEmpty() {
			return
		}

		widget := server.store.Pop()
		widget.Close()
		server.NotificationClosedSignal <- schema.NotificationClosed{ID: widget.Notification.ID, Reason: schema.Dismissed}
	})
	return nil
}

// OpenLastNotification opens the application that sent the most recent notification. This is a non-standard message
func (server *Server) OpenLastNotification() *dbus.Error {
	fmt.Println("Received: OpenLastNotification")
	glib.IdleAdd(func() {
		if server.store.IsEmpty() {
			return
		}

		widget := server.store.Pop()
		widget.CloseAction("default")
		server.NotificationClosedSignal <- schema.NotificationClosed{ID: widget.Notification.ID, Reason: schema.Dismissed}
	})
	return nil
}

// ToggleMute controls if future messages will be displayed to the user or not. This is a non-standard message
func (server *Server) ToggleMute() *dbus.Error {
	server.mute = !server.mute
	fmt.Println("Received: ToggleMute. Is muted? ", server.mute)
	return nil
}

// Kill kills the notification server
func (server *Server) Kill() *dbus.Error {
	gtk.MainQuit()
	return nil
}

// Start connects the sever to d-bus to receive messages
func (server *Server) Start() {
	handler := DbusHandlerNew(server.commands())

	for {
		select {
		case notificationClosed := <-server.NotificationClosedSignal:
			fmt.Println("Sending NotificationClosed", notificationClosed)
			go handler.EmitNotificationClosed(notificationClosed)
		case actionInvoked := <-server.ActionInvokedSignal:
			fmt.Println("Sending ActionInvoked", actionInvoked)
			go handler.EmitActionInvoked(actionInvoked)
		}
	}
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
	methodTable["Kill"] = server.Kill
	return methodTable
}
