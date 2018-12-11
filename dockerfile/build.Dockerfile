FROM vaporio/golang:1.11

ADD . /go/src/github.com/vapor-ware/libmpsse

RUN set -ex \
 && apt-get update \
 && apt-get install -y --no-install-recommends libftdi-dev \
 && cd /go/src/github.com/vapor-ware/libmpsse ; make install \
 && rm -rf /var/lib/apt/lists/*

# This isn't strictly necessary, but it validates that things are
# correct at image build-time.
RUN go build -v github.com/vapor-ware/libmpsse
