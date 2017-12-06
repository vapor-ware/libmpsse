FROM ubuntu:16.04

ADD . /mpsse
WORKDIR /mpsse

RUN set -ex \
  && buildDeps='python-dev golang-go build-essential swig make' \
  && apt-get update \
  && apt-get install -y --no-install-recommends \
    libftdi-dev pkg-config $buildDeps \
  && make build \
  && rm -rf /var/lib/apt/lists/* \
  && apt-get purge -y --auto-remove $buildDeps \
  && apt-get autoremove -y \
  && apt-get clean


