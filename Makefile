PACKAGE=notifyme
BINARY=notifyme

srcdir ?= .
prefix ?= "$(GOPATH)"
bindir ?= "$(prefix)/bin"

BINARY_PATH="$(bindir)/$(BINARY)"
BUILD_DIR="./tmp"

XDG_DATA_HOME ?= "$(HOME)/.local/share"
DBUS_SERVICES ?= "$(XDG_DATA_HOME)/dbus-1/services"

XDG_CONFIG_HOME ?= "$(HOME)/.config"
CONFIG_DIR ?= "$(XDG_CONFIG_HOME)/notifyme"

.PHONY: prepare build stop run install

all: build

prepare:
	dep ensure

build: prepare
	mkdir -p "$(BUILD_DIR)$(DBUS_SERVICES)" "$(BUILD_DIR)$(CONFIG_DIR)"
	go build -o "$(BUILD_DIR)$(BINARY_PATH)" "$(srcdir)/cmd/$(PACKAGE)/$(BINARY).go"
	sed -e "s,\@BINARY_PATH\@,$(BINARY_PATH),g" "$(srcdir)/init/org.freedesktop.Notifications.service" > "$(BUILD_DIR)$(DBUS_SERVICES)/org.freedesktop.Notifications.service"
	cp -r "$(srcdir)/themes" "$(BUILD_DIR)$(CONFIG_DIR)"

stop:
	-killall notifyme

run: stop
	go run cmd/notifyme/notifyme.go

install:
	install -Dm644 "$(BUILD_DIR)$(DBUS_SERVICES)/org.freedesktop.Notifications.service" "$(DESTDIR)$(DBUS_SERVICES)/org.freedesktop.Notifications.service"
	install -Dm644 "$(BUILD_DIR)$(CONFIG_DIR)/themes/notifyme.css" "$(DESTDIR)$(CONFIG_DIR)/themes/notifyme.css"
	install -Dm755 "$(BUILD_DIR)$(BINARY_PATH)" "$(DESTDIR)$(BINARY_PATH)"
