package main

func (store *WidgetStore) add(widget *NotificationWidget) {
	store.ids[widget.Notification.ID] = widget
	store.widgets = append(store.widgets, widget)
}

func (store *WidgetStore) remove(widget *NotificationWidget) {
	store := delete(store.ids, widget.Notification.ID)
	result := store.widgets[:0]
	for _, notification := range store.widgets {
		if notification.Notification.ID != widget.Notification.ID {
			result = append(result, notification)
		}
	}
	store.widgets = result
}

// WidgetStoreNew ...
func WidgetStoreNew() *WidgetStore {
	return &WidgetStore{ids: make(map[uint32]*NotificationWidget)}
}
