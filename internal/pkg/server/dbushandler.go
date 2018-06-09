package notifyme

import (
	"fmt"
	"github.com/ahirata/notifyme/pkg/notifyme/schema"
	"github.com/godbus/dbus"
)

const (
	objectPath               = "/org/freedesktop/Notifications"
	serviceInterface         = "org.freedesktop.Notifications"
	actionInvokedSignal      = serviceInterface + ".ActionInvoked"
	notificationClosedSignal = serviceInterface + ".NotificationClosed"
)

// DbusHandler type struct
type DbusHandler struct {
	conn *dbus.Conn
}

// DbusHandlerNew connects to dbus exporting the commands on methodTable
func DbusHandlerNew(methodTable map[string]interface{}) *DbusHandler {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	reply, err := conn.RequestName(serviceInterface, dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		panic("Name already taken")
	}
	fmt.Println("Connected to dbus")
	conn.ExportMethodTable(methodTable, objectPath, serviceInterface)

	return &DbusHandler{
		conn: conn,
	}
}

// EmitNotificationClosed emits the NotificationClosed signal
func (handler *DbusHandler) EmitNotificationClosed(notificationClosed schema.NotificationClosed) {
	handler.conn.Emit(objectPath, notificationClosedSignal, notificationClosed.ID, notificationClosed.Reason)
}

// EmitActionInvoked emits the ActionInvoked signa
func (handler *DbusHandler) EmitActionInvoked(actionInvoked schema.ActionInvoked) {
	handler.conn.Emit(objectPath, actionInvokedSignal, actionInvoked.ID, actionInvoked.ActionKey)
}
