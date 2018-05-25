package main

// Push ...
func (store *WidgetStore) Push(widget *NotificationWidget) {
	store.widgets = append(store.widgets, widget)
}

// Get ...
func (store *WidgetStore) Get(id uint32) *NotificationWidget {
	for _, widget := range store.widgets {
		if widget.Notification.ID == id {
			return widget
		}
	}
	return nil
}

// Remove ...
func (store *WidgetStore) Remove(id uint32) *NotificationWidget {
	filtered := store.widgets[:0]
	var removed *NotificationWidget
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

// IsEmpty ...
func (store *WidgetStore) IsEmpty() bool {
	return len(store.widgets) == 0
}

// Pop ...
func (store *WidgetStore) Pop() *NotificationWidget {
	last := len(store.widgets) - 1
	elem, array := store.widgets[last], store.widgets[:last]
	store.widgets = array
	return elem
}

// WidgetStoreNew ...
func WidgetStoreNew() *WidgetStore {
	return &WidgetStore{ids: make(map[uint32]*NotificationWidget)}
}
