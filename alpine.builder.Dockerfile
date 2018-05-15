FROM iron/go:dev

COPY . /mpsse
WORKDIR /mpsse

RUN apk --update --no-cache --virtual .build-dep add \
        python-dev build-base swig \
    && apk --update --no-cache add \
        libftdi1-dev \
    && make build \
    && apk del .build-dep
