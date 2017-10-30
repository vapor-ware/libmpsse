FROM ubuntu:16.04

RUN apt-get update
RUN apt-get install -y --no-install-recommends \
      libftdi-dev make swig build-essential golang-go python-dev pkg-config


WORKDIR /mpsse


