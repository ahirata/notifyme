PKG=notifyme
BINARY=notifyme
SERVICES_DIR=$(HOME)/.local/share/dbus-1/services
CONFIG_DIR=$(HOME)/.config/$(PKG)
BINARY_PATH=$(GOBIN)/$(BINARY)

.PHONY: prepare stop run install

prepare:
	mkdir -p $(CONFIG_DIR)
	mkdir -p $(SERVICES_DIR)
	dep ensure

stop:
	-killall notifyme

run: stop
	go run cmd/notifyme/notifyme.go

install: prepare
	go install ./cmd/$(PKG)/$(BINARY).go
	-mv $(CONFIG_DIR)/$(PKG).css{,.old}
	-mv $(SERVICES_DIR)/org.freedesktop.Notifications.service{,.old}
	cp ./themes/$(PKG).css $(CONFIG_DIR)
	sed -e "s,\@BINARY_PATH\@,$(BINARY_PATH),g" ./init/org.freedesktop.Notifications.service > $(SERVICES_DIR)/org.freedesktop.Notifications.service
	-killall notifyme
