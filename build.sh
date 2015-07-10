#!/bin/bash

set -e
set -x

export LC_ALL=C
export DEBIAN_FRONTEND=noninteractive

## Temporarily disable dpkg fsync to make building faster.
if [[ ! -e /etc/dpkg/dpkg.cfg.d/docker-apt-speedup ]]; then
  echo force-unsafe-io > /etc/dpkg/dpkg.cfg.d/docker-apt-speedup
fi

dpkg-divert --local --rename --add /sbin/initctl
ln -sf /bin/true /sbin/initctl

# apt packages
apt-get update
apt-get upgrade -y --force-yes

echo ttf-mscorefonts-installer msttcorefonts/accepted-mscorefonts-eula select true | debconf-set-selections
xargs apt-get install -y --force-yes < /build/packages.txt
fc-cache -fv

# install supervisord
apt-get install -y --no-install-recommends supervisor
mkdir -p /var/log/supervisor

# install go
wget -qO- http://golang.org/dl/go1.4.2.linux-amd64.tar.gz | tar -C /usr/local -xzf -

# install slimerjs
SLIMERJS_VERSION="0.9.6"
SLIMERJS_ARCHIVE_NAME=slimerjs-${SLIMERJS_VERSION}
SLIMERJS_BINARIES_URL=http://download.slimerjs.org/releases/${SLIMERJS_VERSION}/${SLIMERJS_ARCHIVE_NAME}-linux-x86_64.tar.bz2
curl -s $SLIMERJS_BINARIES_URL | tar xj -C /root
mv /root/$SLIMERJS_ARCHIVE_NAME /usr/local/slimerjs

# install app
cd /root/go/src/github.com/krisrang/phantom
go install

# cleanup
apt-get clean

cd /
rm -rf /build
rm -rf /tmp/* /var/tmp/*
rm -rf /var/lib/apt/lists/*
rm -rf /var/cache/apt/archives/*.deb
rm -f /etc/dpkg/dpkg.cfg.d/02apt-speedup

# remove SUID and SGID flags from all binaries
function pruned_find() {
  find / -type d \( -name dev -o -name proc \) -prune -o $@ -print
}

pruned_find -perm /u+s | xargs -r chmod u-s
pruned_find -perm /g+s | xargs -r chmod g-s

# remove non-root ownership of files
chown root:root /var/lib/libuuid

# display build summary
set +x
echo -e "\nRemaining suspicious security bits:"
(
  pruned_find ! -user root
  pruned_find -perm /u+s
  pruned_find -perm /g+s
  pruned_find -perm /+t
) | sed -u "s/^/  /"

echo -e "\nInstalled versions:"
(
  go version
  phantomjs -v
  slimerjs -v
) | sed -u "s/^/  /"

echo -e "\nSuccess!"
exit 0
