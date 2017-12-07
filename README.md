libmpsse
========

Open source library for SPI/I2C control via FTDI chips


## Mac (Darwin)

Get some dependencies
```
brew install libftdi
brew install libusb
```

Then, from the `src` directory
```
make distclean
./configure --disable-python
make
make install
```