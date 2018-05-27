PKG=notifyme
BINARY=notifyme
SERVICES_DIR=$(HOME)/.local/share/dbus-1/services
CONFIG_DIR=$(HOME)/.config/$(PKG)
BINARY_PATH=$(GOBIN)/$(BINARY)

pre_install:
	mkdir -p $(CONFIG_DIR)
	mkdir -p $(SERVICES_DIR)

depends:
	dep ensure

run:
	-killall notifyme
	go run cmd/notifyme/notifyme.go

install: depends pre_install
	go install ./cmd/$(PKG)/$(BINARY).go
	cp ./themes/$(PKG).css $(CONFIG_DIR)
	sed -e "s,\@BINARY_PATH\@,$(BINARY_PATH),g" ./init/org.freedesktop.Notifications.service > $(SERVICES_DIR)/org.freedesktop.Notifications.service
