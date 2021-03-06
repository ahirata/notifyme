package ui

import (
	"github.com/ahirata/notifyme/pkg/notifyme/schema"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"strings"
)

var (
	defaultOffsetX = 10
	defaultOffsetY = 10
)

// NotificationWidget ...
type NotificationWidget struct {
	Notification *schema.Notification
	Window       *gtk.Window
	Icon         *gtk.Image
	Summary      *gtk.Label
	Body         *gtk.Label
	Actions      map[string]*gtk.Button
	Buttons      []*gtk.Button
	channel      chan schema.ActionInvoked
}

// NotificationWidgetNew ...
func NotificationWidgetNew(notification *schema.Notification, maxY int, channel chan schema.ActionInvoked) (*NotificationWidget, error) {
	var err error
	widget := NotificationWidget{Notification: notification, channel: channel}
	if widget.Window, err = gtk.WindowNew(gtk.WINDOW_POPUP); err != nil {
		return nil, err
	}
	if widget.Summary, err = gtk.LabelNew(notification.Summary); err != nil {
		return nil, err
	}
	if widget.Body, err = gtk.LabelNew(notification.Body); err != nil {
		return nil, err
	}
	if widget.Icon, err = gtk.ImageNew(); err != nil {
		return nil, err
	}
	if widget.Buttons, err = widget.createButtons(notification); err != nil {
		return nil, err
	}
	if err = widget.configure(); err != nil {
		return nil, err
	}
	if err = widget.place(maxY); err != nil {
		return nil, err
	}
	return &widget, nil
}

func (widget *NotificationWidget) createButtons(notification *schema.Notification) ([]*gtk.Button, error) {
	var buttons []*gtk.Button
	for i, j := 0, 1; j < len(notification.Actions); i, j = i+2, j+2 {
		actionID := notification.Actions[i].(string)
		actionName := notification.Actions[j].(string)

		button, err := gtk.ButtonNewWithLabel(actionName)
		if err != nil {
			return nil, err
		}
		button.Connect("button-release-event", func() {
			widget.CloseAction(actionID)
		})
		buttons = append(buttons, button)
	}
	return buttons, nil
}

func (widget *NotificationWidget) configure() error {
	configureWindow(widget.Window)
	configureSummary(widget.Summary)
	configureBody(widget.Body)
	setIcon(widget.Icon, widget.Notification)

	return widget.layout()
}

func configureWindow(window *gtk.Window) {
	window.SetName("notifyme")
	window.SetSkipTaskbarHint(true)
	window.SetSkipPagerHint(true)
	window.SetDecorated(false)
	window.SetTypeHint(gdk.WINDOW_TYPE_HINT_NOTIFICATION)
	window.SetGravity(gdk.GDK_GRAVITY_SOUTH_EAST)
	window.SetCanFocus(false)
	window.SetAcceptFocus(false)
	window.SetKeepAbove(true)
}

func configureSummary(label *gtk.Label) {
	label.SetUseMarkup(true)
	label.SetLineWrap(false)
	label.SetHAlign(gtk.ALIGN_START)
	label.SetXAlign(0)
	label.SetMaxWidthChars(45)
	label.SetEllipsize(pango.ELLIPSIZE_END)
}

func configureBody(label *gtk.Label) {
	label.SetUseMarkup(true)
	label.SetLineWrap(false)
	label.SetHAlign(gtk.ALIGN_START)
	label.SetXAlign(0)
	label.SetMaxWidthChars(45)
	label.SetEllipsize(pango.ELLIPSIZE_END)
}

func setIcon(icon *gtk.Image, notification *schema.Notification) {
	icon.Clear()
	if imageData, found := notification.ImageData(); found {
		icon.SetFromPixbuf(pixbufNewFromImageData(&imageData))
	} else if imagePath, found := notification.ImagePath(); found {
		icon.SetFromPixbuf(loadPixbufFromFile(imagePath, 64, 64))
	} else if strings.HasPrefix(notification.AppIcon, "file://") {
		icon.SetFromPixbuf(loadPixbufFromFile(notification.AppIcon, 64, 64))
	} else if notification.AppIcon != "" {
		icon.SetFromIconName(notification.AppIcon, gtk.ICON_SIZE_DIALOG)
	} else if iconData, found := notification.IconData(); found {
		icon.SetFromPixbuf(pixbufNewFromImageData(&iconData))
	}
	return
}

func loadImageFromFile(filename string, width, height int) (*gtk.Image, error) {
	return gtk.ImageNewFromPixbuf(loadPixbufFromFile(filename, width, height))
}

func pixbufNewFromImageData(imageData *schema.ImageData) *gdk.Pixbuf {
	pixbuf, err := pixbufNewFromData(imageData.Data, gdk.COLORSPACE_RGB, imageData.HasAlpha, int(imageData.BitsPerSample), int(imageData.Width), int(imageData.Height), 64, 64)
	if err != nil {
		return nil
	}

	return pixbuf
}

func (widget *NotificationWidget) layout() error {
	LoadCSSProvider(widget.Window)

	AddClass(widget.Window, "notifyme")
	AddClass(widget.Summary, "summary")
	AddClass(widget.Body, "body")

	vbox, err := AddBox(widget.Window, gtk.ORIENTATION_VERTICAL, "main")
	if err != nil {
		return err
	}

	content, err := AddBox(vbox, gtk.ORIENTATION_HORIZONTAL, "content")
	if err != nil {
		return (err)
	}
	content.Add(widget.Icon)

	textBox, err := AddBox(content, gtk.ORIENTATION_VERTICAL, "message")
	if err != nil {
		return err
	}
	textBox.Add(widget.Summary)
	textBox.Add(widget.Body)

	actions, err := AddBox(vbox, gtk.ORIENTATION_HORIZONTAL, "actions")
	if err != nil {
		return err
	}
	actions.SetHAlign(gtk.ALIGN_END)
	actions.SetVAlign(gtk.ALIGN_END)
	for _, button := range widget.Buttons {
		actions.Add(button)
	}

	return nil
}

func (widget *NotificationWidget) place(maxY int) error {
	workarea, err := getWorkarea(widget.Window)
	if err != nil {
		panic(err)
	}

	widget.Window.Connect("size-allocate", func() {
		positionX := widget.getPositionX(workarea)
		positionY := widget.getPositionY(workarea, maxY)

		widget.Window.Move(positionX, positionY)
	})

	return nil
}

func (widget *NotificationWidget) getPositionX(workarea *gdk.Rectangle) int {
	width := widget.Window.GetAllocatedWidth()
	return workarea.GetX() + workarea.GetWidth() - width - defaultOffsetX
}

func (widget *NotificationWidget) getPositionY(workarea *gdk.Rectangle, maxY int) int {
	height := widget.Window.GetAllocatedHeight()
	positionY := workarea.GetY() + workarea.GetHeight()
	if maxY < positionY {
		positionY = maxY
	}
	return positionY - height - defaultOffsetY
}

// ReplaceNotification replaces the image, summary and body of the notification with same ID
func (widget *NotificationWidget) ReplaceNotification(notification *schema.Notification) {
	setIcon(widget.Icon, notification)
	widget.Summary.SetLabel(notification.Summary)
	widget.Body.SetLabel(notification.Body)
	widget.Notification = notification
}

// Close closes the widget
func (widget *NotificationWidget) Close() {
	widget.Window.Destroy()
}

// CloseAction ...
func (widget *NotificationWidget) CloseAction(actionKey string) {
	widget.channel <- schema.ActionInvoked{ID: widget.Notification.ID, ActionKey: actionKey}
	widget.Close()
}

// Show shows the widget
func (widget *NotificationWidget) Show() {
	widget.Window.ShowAll()
}
