# Maintainer: Arthur Hirata<arthur_hirata@yahoo.com.br>

pkgname=notifyme-git
_gitname=notifyme
pkgver=r19.c92b649
pkgrel=1
pkgdesc="Freedesktop notification server implementation"
arch=('x86_64')
url="https://github.com/ahirata/${_gitname}"
license=('MIT')
makedepends=('go' 'git' 'dep')
options=('!strip' '!emptydirs')
source=("git://github.com/ahirata/${_gitname}.git")
sha256sums=('SKIP')

pkgver() {
  cd "${srcdir}/${_gitname}"
  echo "r$(git rev-list --count HEAD).$(git rev-parse --short HEAD)"
}

prepare() {
  mkdir -p "$srcdir/go/src/github.com/ahirata"
  ln -sf "$srcdir/${_gitname}" "$srcdir/go/src/github.com/ahirata/${_gitname}"
}

build() {
  export GOPATH="$srcdir/go"
  export DBUS_SERVICES="/usr/share/dbus-1/services"
  export CONFIG_DIR="/usr/share/notifyme"
  export prefix="/usr"

  cd "$srcdir/go/src/github.com/ahirata/${_gitname}"
  make
}

package() {
  export GOPATH="$srcdir/go"
  export DBUS_SERVICES="/usr/share/dbus-1/services"
  export CONFIG_DIR="/usr/share/notifyme"
  export prefix="/usr"

  cd "$srcdir/go/src/github.com/ahirata/${_gitname}"
  make DESTDIR="$pkgdir" install
}

# vim:set ts=2 sw=2 et:
