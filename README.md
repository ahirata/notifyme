Notifyme
========

Notifyme is a simple [freedesktop notification server](https://developer.gnome.org/notification-spec/) implementation.

## Installation

### From sources
Make sure you have the following:

* `GOPATH` and `GOBIN` environment variables defined;
* [dep](https://github.com/golang/dep) in your `PATH`;
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
