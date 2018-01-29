libmpsse
========

Open source Go wrapper library for SPI/I2C control via FTDI chips based on [devttys0/libmpsse](https://github.com/devttys0/libmpsse).


## Setting up on Mac (Darwin)

### Installing Go
If you don't have `go` installed you can install it via HomeBrew,
```
$ brew install go
$ mkdir -p $HOME/go
```

You must then set your GOPATH - this can be done by addint the following exports
to your `~/.bashrc` file
```
export GOPATH=${HOME}/work
export GOVERSION=$(brew list go | head -n 1 | cut -d '/' -f 6)
export GOROOT=$(brew --prefix)/Cellar/go/${GOVERSION}/libexec
export PATH=${GOPATH}/bin:$PATH
```

Source the new environment
```
$ source ~/.bashrc
```

And check that it worked
```
$ go env
```

### Installing libmpsse
First you will need to get this repo. This can either be done by cloning
the repo locally, or via `go get`, which will clone it into your GOPATH
under `$GOPATH/src/github.com/vapor-ware/libmpsse`
```
$ github.com/vapor-ware/libmpsse
```

Once you have libmpsse locally, `cd` to the repo.

#### Getting dependencies
Before the source can be built, you need to install some dependencies:
* libftdi
* libusb

This can be done via HomeBrew
```
brew install libusb libftdi
```

The versions of the packages can be found using `brew info [FORMULA...]`
Currently, we are using:
- `libusb`:  stable, 1.0.21 (bottled), HEAD
- `libftdi`: stable 1.4 (bottled)

#### Installing
Then, simply `make build` from the project root. You can see the Makefile
target for `build` for more information on how the source is being built
and installed.

#### Uninstalling
To uninstall libmpsse, from the repo root:
```
$ cd src
$ make uninstall
```