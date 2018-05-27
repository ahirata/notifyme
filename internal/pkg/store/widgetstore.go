package store

import "github.com/ahirata/notifyme/internal/pkg/ui"
import "github.com/ahirata/notifyme/pkg/notifyme/schema"
import "github.com/gotk3/gotk3/glib"

// WidgetStore ...
type WidgetStore struct {
	widgets             []*ui.NotificationWidget
	ActionInvokedSignal chan schema.ActionInvoked
}

// WidgetStoreNew ...
func WidgetStoreNew(actionInvokedSignal chan schema.ActionInvoked) *WidgetStore {
	return &WidgetStore{ActionInvokedSignal: actionInvokedSignal}
}

func (store *WidgetStore) push(widget *ui.NotificationWidget) {
	store.widgets = append(store.widgets, widget)
}

func (store *WidgetStore) pop() *ui.NotificationWidget {
	last := len(store.widgets) - 1
	widget, array := store.widgets[last], store.widgets[:last]
	store.widgets = array
	return widget
}

func (store *WidgetStore) remove(id uint32) *ui.NotificationWidget {
	filtered := store.widgets[:0]
	var removed *ui.NotificationWidget
	for _, widget := range store.widgets {
		if widget.Notification.ID != id {
			filtered = append(filtered, widget)
		} else {
			removed = widget
		}
	}
	store.widgets = filtered
	return removed
}

func (store *WidgetStore) getLast(fn func(widget *ui.NotificationWidget)) *ui.NotificationWidget {
	if !store.IsEmpty() {
		result := make(chan *ui.NotificationWidget)
		glib.IdleAdd(func() {
			widget := store.pop()
			fn(widget)
			result <- widget
		})
		return <-result
	}
	return nil
}

func (store *WidgetStore) currentHeight() int {
	var height int
	for _, widget := range store.widgets {
		height = height + widget.Window.GetAllocatedHeight() + 10
	}
	return height
}

// Get ...
func (store *WidgetStore) Get(id uint32) *ui.NotificationWidget {
	for _, widget := range store.widgets {
		if widget.Notification.ID == id {
			return widget
		}
	}
	return nil
}

// Close ...
func (store *WidgetStore) Close(id uint32) {
	glib.IdleAdd(func() {
		removed := store.remove(id)
		if removed != nil {
			removed.Close()
		}
	})
}

// IsEmpty ...
func (store *WidgetStore) IsEmpty() bool {
	return len(store.widgets) == 0
}

// Put ...
func (store *WidgetStore) Put(notification *schema.Notification) {
	glib.IdleAdd(func() {
		widget := store.Get(notification.ID)
		if widget != nil {
			widget.ReplaceNotification(notification)
		} else {
			widget, err := ui.NotificationWidgetNew(notification, store.currentHeight(), store.ActionInvokedSignal)
			if err != nil {
				panic(err)
			}
			store.push(widget)
			widget.Show()
		}
	})
}

// CloseLast ...
func (store *WidgetStore) CloseLast() *ui.NotificationWidget {
	return store.getLast(func(widget *ui.NotificationWidget) {
		widget.Close()
	})
}

// OpenLast ...
func (store *WidgetStore) OpenLast(actionKey string) *ui.NotificationWidget {
	return store.getLast(func(widget *ui.NotificationWidget) {
		widget.OpenApp()
		widget.CloseAction(actionKey)
	})
}