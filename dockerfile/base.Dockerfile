FROM vaporio/foundation:bionic

ADD . /mpsse
WORKDIR /mpsse

RUN set -ex \
 && buildDeps='python-dev build-essential swig make' \
 && apt-get update \
 && apt-get install -y --no-install-recommends \
   libftdi-dev pkg-config $buildDeps \
 && make install \
 && rm -rf /var/lib/apt/lists/* \
 && apt-get purge -y --auto-remove $buildDeps \
 && apt-get autoremove -y \
 && apt-get clean