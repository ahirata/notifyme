package store

import (
	"github.com/ahirata/notifyme/internal/pkg/ui"
	"math"
)

// WidgetStore holds the existig widgets
type WidgetStore struct {
	widgets []*ui.NotificationWidget
}

// Push adds a Widget to the widget list
func (store *WidgetStore) Push(widget *ui.NotificationWidget) {
	store.widgets = append(store.widgets, widget)
}

// Pop removes the most recent added widget
func (store *WidgetStore) Pop() *ui.NotificationWidget {
	last := len(store.widgets) - 1
	widget, array := store.widgets[last], store.widgets[:last]
	store.widgets = array
	return widget
}

// Remove removes a widget based on its id regardless of its position in the list
func (store *WidgetStore) Remove(id uint32) *ui.NotificationWidget {
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

// MinY returns the smaller screen position Y among all widgets
func (store *WidgetStore) MinY() int {
	minY := math.MaxInt32
	for _, widget := range store.widgets {
		_, y := widget.Window.GetPosition()
		if 0 < y && y < minY {
			minY = y
		}
	}
	return minY
}

// Get retrieves the widget by id
func (store *WidgetStore) Get(id uint32) *ui.NotificationWidget {
	for _, widget := range store.widgets {
		if widget.Notification.ID == id {
			return widget
		}
	}
	return nil
}

// IsEmpty returns true if there are no widgets, false otherwise
func (store *WidgetStore) IsEmpty() bool {
	return len(store.widgets) == 0
}
