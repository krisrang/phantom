FROM ubuntu:14.04
MAINTAINER krisrang "mail@rang.ee"

ENV HOME /root
ENV GOPATH /root/go
ENV PATH /root/go/bin:/usr/local/slimerjs:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games
ENV LD_LIBRARY_PATH /usr/local/lib:/usr/lib:/lib/x86_64-linux-gnu:/lib:/usr/local/slimerjs
ENV LC_ALL C
ENV INITRD No
ENV DISPLAY :99
ENV PORT 5000

RUN mkdir -p /root/go

COPY . /root/go/src/github.com/krisrang/phantom
COPY ./fonts /usr/local/share/fonts/skyltmax
COPY . /build

ADD ./multiverse.list /etc/apt/sources.list.d/multiverse.list
ADD bin/phantomjs-linux_x86_64 /usr/local/bin/phantomjs
ADD supervisord.conf /etc/supervisor/conf.d/supervisord.conf

RUN /build/build.sh

CMD supervisord
