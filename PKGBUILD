# Maintainer: Tom McKeesick <tmck01@gmail.com>
pkgname=pokesay-bin
pkgver=0.18.1
pkgrel=1
pkgdesc="Print pokemon in the CLI! An adaptation of the classic 'cowsay'"
arch=('x86_64')
url="https://github.com/tmck-code/pokesay"
license=('BSD-3-Clause')
depends=()
source=("https://github.com/tmck-code/pokesay/releases/download/v$pkgver/pokesay-$pkgver-linux-amd64.tar.gz")
sha256sums=('SKIP')

package() {
    install -Dm755 "pokesay-$pkgver-linux-amd64" "$pkgdir/usr/bin/pokesay"
    install -Dm644 "pokesay.1" "$pkgdir/usr/share/man/man1/pokesay.1"
    install -Dm644 "LICENSE" "$pkgdir/usr/share/licenses/pokesay/LICENSE"
}
