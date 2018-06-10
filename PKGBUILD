# Maintainer: Arthur Hirata<arthur_hirata@yahoo.com.br>

pkgname=notifyme-git
_gitname=notifyme
pkgver=r19.9836a8a
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

  cd "$srcdir/go/src/github.com/ahirata/${_gitname}"
  ./configure --prefix=/usr --config-dir=/usr/share/notifyme --services-dir=/usr/share/dbus-1/services
  make DESTDIR=${srcdir}/tmp
}

package() {
	cp -rT $srcdir/tmp $pkgdir
}

# vim:set ts=2 sw=2 et:
