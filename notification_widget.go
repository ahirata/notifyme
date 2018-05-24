package main

import (
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"os/exec"
	"strings"
)

var (
	offsetX = 20
	offsetY = 20
)

// NotificationWidget ...
type NotificationWidget struct {
	Notification *Notification
	Window       *gtk.Window
	Icon         *gtk.Image
	Summary      *gtk.Label
	Body         *gtk.Label
	Actions      map[string]*gtk.Button
	Buttons      []*gtk.Button
	channel      chan Action // TODO - maybe this should not be here
}

// NotificationWidgetNew ...
func NotificationWidgetNew(notificationWidget *NotificationWidget) (*NotificationWidget, error) {
	var err error
	notification := notificationWidget.Notification

	if notificationWidget.Window, err = gtk.WindowNew(gtk.WINDOW_POPUP); err != nil {
		return nil, err
	}
	if strings.HasPrefix(notification.AppIcon, "file://") {
		if notificationWidget.Icon, err = gtk.ImageNewFromFile(strings.Replace(notification.AppIcon, "file://", "", 1)); err != nil {
			return nil, err
		}
	} else if notificationWidget.Icon, err = gtk.ImageNewFromIconName(notification.AppIcon, gtk.ICON_SIZE_DIALOG); err != nil {
		fmt.Println(err)
		return nil, err
	} else if notificationWidget.Icon, err = gtk.ImageNewFromPixbuf(pixbufNew(notification)); err != nil {
		return nil, err
	}
	if notificationWidget.Summary, err = gtk.LabelNew(notification.Summary); err != nil {
		return nil, err
	}
	if notificationWidget.Body, err = gtk.LabelNew(notification.Body); err != nil {
		return nil, err
	}
	if err = notificationWidget.createButtons(notification); err != nil {
		return nil, err
	}
	if err = notificationWidget.configure(); err != nil {
		return nil, err
	}
	if err = notificationWidget.move(); err != nil {
		return nil, err
	}
	return notificationWidget, nil
}

func (widget *NotificationWidget) createButtons(notification *Notification) error {
	for i, j := 0, 1; j < len(notification.Actions); i, j = i+2, j+2 {
		actionID := notification.Actions[i].(string)
		actionName := notification.Actions[j].(string)

		button, err := gtk.ButtonNewWithLabel(actionName)
		if err != nil {
			return err
		}
		button.Connect("button-release-event", func() {
			widget.channel <- Action{ID: notification.ID, action: actionID}
		})
		widget.Buttons = append(widget.Buttons, button)
	}
	return nil
}

func (widget *NotificationWidget) configure() error {
	configureWindow(widget.Window)
	configureSummary(widget.Summary)
	configureBody(widget.Body)

	return widget.layout()
}

func configureWindow(window *gtk.Window) {
	window.SetSkipTaskbarHint(true)
	window.SetDecorated(false)
	window.SetTypeHint(gdk.WINDOW_TYPE_HINT_NOTIFICATION)
	window.SetGravity(gdk.GDK_GRAVITY_SOUTH_EAST)
	window.SetCanFocus(false)
	window.SetAcceptFocus(false)
	window.SetKeepAbove(true)
}

func configureSummary(label *gtk.Label) {
	label.SetUseMarkup(true)
	label.SetHAlign(gtk.ALIGN_START)
	label.SetXAlign(0)
	label.SetMaxWidthChars(40)
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

func pixbufNew(notification *Notification) *gdk.Pixbuf {
	imageData, exists := notification.ImageData()
	if !exists {
		return nil
	}

	pixbuf, err := pixbufNewFromData(imageData.Data, gdk.COLORSPACE_RGB, imageData.HasAlpha, int(imageData.BitsPerSample), int(imageData.Width), int(imageData.Height))
	if err != nil {
		return nil
	}

	scaled, err := pixbuf.ScaleSimple(64, 64, gdk.INTERP_BILINEAR)
	if err != nil {
		return nil
	}
	return scaled
}

func (widget *NotificationWidget) layout() error {
	AddClass(widget.Window, "notificationd")
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

func (widget *NotificationWidget) move() error {
	screen, err := widget.Window.GetScreen()
	if err != nil {
		return err
	}

	display, err := screen.GetDisplay()
	if err != nil {
		return err
	}

	workarea, err := getWorkarea(display)
	if err != nil {
		panic(err)
	}
	width, height := 500, 129
	widget.Window.SetSizeRequest(width, height)

	widget.Window.Move(workarea.GetX()+workarea.GetWidth()-width-offsetX, workarea.GetY()+workarea.GetHeight()-height-offsetY)

	return nil
}

// ReplaceNotification replaces the image, summary and body of the notification with same ID
func (widget *NotificationWidget) ReplaceNotification(notification *Notification) {
	widget.Icon.Clear()
	widget.Icon.SetFromPixbuf(pixbufNew(notification))
	widget.Summary.SetLabel(notification.Summary)
	widget.Body.SetLabel(notification.Body)
	widget.Notification = notification
}

// OpenApp opens the app that sent the notification
func (widget *NotificationWidget) OpenApp() {
	cmd := exec.Command("wmctrl", "-xa", widget.Notification.AppName)
	cmd.Start()
	widget.channel <- Action{ID: widget.Notification.ID, action: "default"}
	widget.Close("dismiss")
}

// Close closes the widget
func (widget *NotificationWidget) Close(reason string) {
	if widget.Window != nil {
		widget.Window.Destroy()
		widget.channel <- Action{ID: widget.Notification.ID, action: reason}
	}
}

// Show shows the widget
func (widget *NotificationWidget) Show() {
	widget.Window.ShowAll()
}
