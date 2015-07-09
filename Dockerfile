FROM ubuntu:14.04
MAINTAINER krisrang "mail@rang.ee"

ENV HOME /root
ENV GOPATH /root/go
ENV PATH /root/go/bin:/usr/local/slimerjs:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games
ENV LD_LIBRARY_PATH /usr/local/lib:/usr/lib:/lib:/usr/local/slimerjs
ENV DEBIAN_FRONTEND noninteractive
ENV LC_ALL C
ENV INITRD No
ENV DISPLAY :99
ENV PORT 5000

RUN mkdir -p /root/go

ADD ./multiverse.list /etc/apt/sources.list.d/multiverse.list
ADD . /build
ADD xvfb_init /etc/init.d/xvfb
ADD xvfb-daemon-run /usr/bin/xvfb-daemon-run
ADD bin/phantomjs-linux_x86_64 /usr/local/bin/phantomjs
ADD supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY . /root/go/src/github.com/krisrang/phantom
RUN /build/build.sh

EXPOSE 5000
CMD supervisord
