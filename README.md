Notifyme
========

Notifyme is a simple [freedesktop notification server](https://developer.gnome.org/notification-spec/) implementation written in Go.

## Installation

### From sources
Make sure you have the following:

* `GOPATH` and `$GOBIN` defined;
* [dep](https://github.com/golang/dep)
* No other notification service registered and running (eg. notify-osd, xfce-notifyd, dunst).

```
make
make install
```

### Arch Linux
Build and install from PKGBUILD:
```
cd build/archlinux
makepkg -i
```
