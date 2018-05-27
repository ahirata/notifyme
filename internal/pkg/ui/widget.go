package ui

import (
	"github.com/ahirata/notifyme/pkg/notifyme/schema"
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
	Notification *schema.Notification
	Window       *gtk.Window
	Icon         *gtk.Image
	Summary      *gtk.Label
	Body         *gtk.Label
	Actions      map[string]*gtk.Button
	Buttons      []*gtk.Button
	channel      chan schema.ActionInvoked // TODO - maybe this should not be here
}

// NotificationWidgetNew ...
func NotificationWidgetNew(notification *schema.Notification, position int, channel chan schema.ActionInvoked) (*NotificationWidget, error) {
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
	if widget.Icon, err = loadImage(notification); err != nil {
		return nil, err
	}
	if err = widget.createButtons(notification); err != nil {
		return nil, err
	}
	if err = widget.configure(); err != nil {
		return nil, err
	}
	if err = widget.move(position); err != nil {
		return nil, err
	}
	return &widget, nil
}

func (widget *NotificationWidget) createButtons(notification *schema.Notification) error {
	for i, j := 0, 1; j < len(notification.Actions); i, j = i+2, j+2 {
		actionID := notification.Actions[i].(string)
		actionName := notification.Actions[j].(string)

		button, err := gtk.ButtonNewWithLabel(actionName)
		if err != nil {
			return err
		}
		button.Connect("button-release-event", func() {
			widget.CloseAction(actionID)
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

func loadImage(notification *schema.Notification) (*gtk.Image, error) {
	var image *gtk.Image
	var err error
	if strings.HasPrefix(notification.AppIcon, "file://") {
		path := strings.Replace(notification.AppIcon, "file://", "", 1)
		pixbuf, err := gdk.PixbufNewFromFileAtScale(path, 64, 64, true)
		if err != nil {
			return nil, err
		}
		if image, err = gtk.ImageNewFromPixbuf(pixbuf); err != nil {
			return nil, err
		}
	} else if image, err = gtk.ImageNewFromIconName(notification.AppIcon, gtk.ICON_SIZE_DIALOG); err != nil {
		return nil, err
	} else if image, err = gtk.ImageNewFromPixbuf(pixbufNew(notification)); err != nil {
		return nil, err
	}
	return image, err
}

func pixbufNew(notification *schema.Notification) *gdk.Pixbuf {
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

func (widget *NotificationWidget) move(position int) error {
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

	widget.Window.Connect("size-allocate", func() {
		width := max(widget.Window.GetAllocatedWidth(), 500)
		height := max(widget.Window.GetAllocatedHeight(), 129)

		widget.Window.SetSizeRequest(width, height)
		widget.Window.Move(workarea.GetX()+workarea.GetWidth()-width-offsetX, workarea.GetY()+workarea.GetHeight()-height-offsetY-position)
	})

	return nil
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

// ReplaceNotification replaces the image, summary and body of the notification with same ID
func (widget *NotificationWidget) ReplaceNotification(notification *schema.Notification) {
	widget.Icon.Clear()
	if image, err := loadImage(notification); err == nil {
		widget.Icon.SetFromPixbuf(image.GetPixbuf())
	}
	widget.Summary.SetLabel(notification.Summary)
	widget.Body.SetLabel(notification.Body)
	widget.Notification = notification
}

// OpenApp opens the app that sent the notification
func (widget *NotificationWidget) OpenApp() {
	cmd := exec.Command("wmctrl", "-xa", widget.Notification.AppName)
	cmd.Start()
}

// Close closes the widget
func (widget *NotificationWidget) Close() {
	widget.Window.Destroy()
}

// CloseAction ...
func (widget *NotificationWidget) CloseAction(actionKey string) {
	widget.Close()
	widget.channel <- schema.ActionInvoked{ID: widget.Notification.ID, ActionKey: actionKey}
}

// Show shows the widget
func (widget *NotificationWidget) Show() {
	widget.Window.ShowAll()
	widget.Window.Present()
}
