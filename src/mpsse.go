package src

import (
	"sync"
	"fmt"
	"unsafe"
)

// #cgo pkg-config: libftdi
// #include <stdio.h>
// #include "mpsse.h"
import "C"


// modes
type Mode int
const (
	SPI0 Mode = 1
	SPI1 Mode = 2
	SPI2 Mode = 3
	SPI3 Mode = 4
	I2C  Mode = 5
	GPIO Mode = 6
	BITBANG Mode = 7
)

// frequencies
type Frequency int
const (
	ONE_HUNDRED_KHZ  Frequency = 100000
	FOUR_HUNDRED_KHZ Frequency = 400000
	ONE_MHZ          Frequency = 1000000
	TWO_MHZ          Frequency = 2000000
	FIVE_MHZ         Frequency = 5000000
	SIX_MHZ          Frequency = 6000000
	TEN_MHZ          Frequency = 10000000
	TWELVE_MHZ       Frequency = 12000000
	FIFTEEN_MHZ      Frequency = 15000000
	THIRTY_MHZ       Frequency = 30000000
	SIXTY_MHZ        Frequency = 60000000
)

// endianess
type Endianess int
const (
	MSB Endianess = 0x00
	LSB Endianess = 0x08
)

// i2c ack
type I2C_ACK int
const (
	ACK  I2C_ACK = 0
	NACK I2C_ACK = 1
)

// gpio pins
type GPIOPin int
const (
	GPIOL0 GPIOPin = 0
	GPIOL1 GPIOPin = 1
	GPIOL2 GPIOPin = 2
	GPIOL3 GPIOPin = 3
	GPIOH0 GPIOPin = 4
	GPIOH1 GPIOPin = 5
	GPIOH2 GPIOPin = 6
	GPIOH3 GPIOPin = 7
	GPIOH4 GPIOPin = 8
	GPIOH5 GPIOPin = 9
	GPIOH6 GPIOPin = 10
	GPIOH7 GPIOPin = 11
)

// ftdi interfaces
type Iface int
const (
	INTERFACE_ANY = 0
	INTERFACE_A   = 1
	INTERFACE_B   = 2
	INTERFACE_C   = 3
	INTERFACE_D   = 4
)


type Mpsse struct {
	ctx  *C.struct_mpsse_context
	open bool
	lock sync.Mutex
}




// Once we have the functionality it, this could be setup in a way that is
// similar to the python lib, namely have an MPSSE struct. A New() fn would
// act as the initializer (e.g. the current MPSSE fn) and then the actions
// (read, write, ack, etc) would be part of the struct - would make it a
// little easier to use, e.g. Write(vec, "\xA1") would be vec.Write("\xA1")




// could probably build in error checking to most of these methods too!



// TODO - USED -- actually, we don't really need this. we just need to define
//  open and possibly perform some of the same logic that MPSSE does. we don't
//  need this because it opens the first device, but we want to open a specific
//  device. we should still implement this, but I don't think it will be used
//  for our plugin.
//
// MPSSE opens and initializes the first FTDI device found.
//
// It is a wrapper for the mpsse C function:
//     struct mpsse_context *MPSSE(enum modes mode, int freq, int endianess);
func MPSSE(mode Mode, frequency Frequency, endianess Endianess) (*Mpsse, error) {

	ctx := C.MPSSE(C.enum_modes(mode), C.int(frequency), C.int(endianess))
	// FIXME - should check if ok
	d := &Mpsse{ctx, true, sync.Mutex{}}

	fmt.Printf("device: %+v", d)
	return d, nil
}

// TODO - USED
//   since the C version is just a wrapper around OpenIndex for idx 0, don't wrap the C fn here, just
//   use idx 0 with the wrapped OpenIndex fn.
func Open(vid int, pid int, mode Mode, frequency Frequency, endianess Endianess, iface Iface, description *string, serial *string) (*Mpsse, error) {
	return OpenIndex(vid, pid, mode, frequency, endianess, iface, description, serial, 0)
}


// FIXME - not used by us (well, it could be .. Open() is just a wrapper for this at idx 0)
// It is a wrapper for the mpsse C function:
//     struct mpsse_context *OpenIndex(int vid, int pid, enum modes mode, int freq, int endianess, int interface, const char *description, const char *serial, int index);
func OpenIndex(vid int, pid int, mode Mode, frequency Frequency, endianess Endianess, iface Iface, description *string, serial *string, index int) (*Mpsse, error) {

	descP := C.CString(*description)
	defer C.free(unsafe.Pointer(descP))

	serP := C.CString(*serial)
	defer C.free(unsafe.Pointer(serP))

	ctx := C.OpenIndex(
		C.int(vid),
		C.int(pid),
		C.enum_modes(mode),
		C.int(frequency),
		C.int(endianess),
		C.int(iface),
		descP,
		serP,
		C.int(index),
	)

	// FIXME - should check if ctx ok

	d := &Mpsse{ctx, true, sync.Mutex{}}
	return d, nil
}


// TODO - USED
// It is a wrapper for the mpsse C function:
//     void Close(struct mpsse_context *mpsse);
func (m *Mpsse) Close() {
	C.Close(unsafe.Pointer(m.ctx))
}


// TODO - USED
// It is a wrapper for the mpsse C function:
//     const char *ErrorString(struct mpsse_context *mpsse);
func (m *Mpsse) ErrorString() string {
	return C.GoString(C.ErrorString(unsafe.Pointer(m.ctx)))
}


// It is a wrapper for the mpsse C function:
//     int SetMode(struct mpsse_context *mpsse, int endianess);
func SetMode() {}


// It is a wrapper for the mpsse C function:
//     void EnableBitmode(struct mpsse_context *mpsse, int tf);
func EnableBitmode() {}


// It is a wrapper for the mpsse C function:
//     int SetClock(struct mpsse_context *mpsse, uint32_t freq);
func SetClock() {}


// It is a wrapper for the mpsse C function:
//     int GetClock(struct mpsse_context *mpsse);
func GetClock() {}


// It is a wrapper for the mpsse C function:
//     int GetVid(struct mpsse_context *mpsse);
func GetVid() {}


// It is a wrapper for the mpsse C function:
//     int GetPid(struct mpsse_context *mpsse);
func GetPid() {}


// It is a wrapper for the mpsse C function:
//     const char *GetDescription(struct mpsse_context *mpsse);
func GetDescription() {}


// It is a wrapper for the mpsse C function:
//     int SetLoopback(struct mpsse_context *mpsse, int enable);
func SetLoopback() {}


// It is a wrapper for the mpsse C function:
//     void SetCSIdle(struct mpsse_context *mpsse, int idle);
func SetCSIdle() {}


// TODO - USED
// It is a wrapper for the mpsse C function:
//     int Start(struct mpsse_context *mpsse);
func (m *Mpsse) Start() int {
	return int(C.Start(unsafe.Pointer(m.ctx)))
}


// TODO - USED
// It is a wrapper for the mpsse C function:
//     int Write(struct mpsse_context *mpsse, char *data, int size);
func (m *Mpsse) Write(data string) int {
	dataP := C.CString(data)
	defer C.free(unsafe.Pointer(dataP))

	// FIXME -- need to check that this works. not clear that len(data) gives
	// us the size that we want. maybe unsafe.Sizeof will give the int
	// size we want? but I'm also unsure about that. will need to test.
	return int(C.Write(unsafe.Pointer(m.ctx), dataP, C.int(len(data))))
}


// TODO - USED
// It is a wrapper for the mpsse C function:
//     int Stop(struct mpsse_context *mpsse);
func (m *Mpsse) Stop() int {
	return int(C.Stop(unsafe.Pointer(m.ctx)))
}


// TODO - USED
// It is a wrapper for the mpsse C function:
//     int GetAck(struct mpsse_context *mpsse);
func (m *Mpsse) GetAck() int {
	return int(C.GetAck(unsafe.Pointer(m.ctx)))
}


// It is a wrapper for the mpsse C function:
//     void SetAck(struct mpsse_context *mpsse, int ack);
func (m *Mpsse) SetAck(ack I2C_ACK) {
	C.SetAck(unsafe.Pointer(m.ctx), C.int(ack))
}


// TODO - USED
// It is a wrapper for the mpsse C function:
//     void SendAcks(struct mpsse_context *mpsse);
func (m *Mpsse) SendAcks() {
	C.SendAcks(unsafe.Pointer(m.ctx))
}


// TODO -- USED
// It is a wrapper for the mpsse C function:
//     void SendNacks(struct mpsse_context *mpsse);
func (m *Mpsse) SendNacks() {
	C.SendNacks(unsafe.Pointer(m.ctx))
}


// It is a wrapper for the mpsse C function:
//     void FlushAfterRead(struct mpsse_context *mpsse, int tf);
func FlushAfterRead() {}


// TODO -- USED
// It is a wrapper for the mpsse C function:
//     int PinHigh(struct mpsse_context *mpsse, int pin);
func (m *Mpsse) PinHigh(pin GPIOPin) int {
	return int(C.PinHigh(unsafe.Pointer(m.ctx), C.int(pin)))
}


// It is a wrapper for the mpsse C function:
//     int PinLow(struct mpsse_context *mpsse, int pin);
func (m *Mpsse) PinLow(pin GPIOPin) int {
	return int(C.PinLow(unsafe.Pointer(m.ctx), C.int(pin)))
}


// It is a wrapper for the mpsse C function:
//     int SetDirection(struct mpsse_context *mpsse, uint8_t direction);
func SetDirection() {}


// It is a wrapper for the mpsse C function:
//     int WriteBits(struct mpsse_context *mpsse, char bits, int size);
func WriteBits() {}


// It is a wrapper for the mpsse C function:
//     char ReadBits(struct mpsse_context *mpsse, int size);
func ReadBits() {}


// It is a wrapper for the mpsse C function:
//     int WritePins(struct mpsse_context *mpsse, uint8_t data);
func WritePins() {}


// It is a wrapper for the mpsse C function:
//     int ReadPins(struct mpsse_context *mpsse);
func ReadPins() {}


// It is a wrapper for the mpsse C function:
//     int PinState(struct mpsse_context *mpsse, int pin, int state);
func PinState() {}


// It is a wrapper for the mpsse C function:
//     int Tristate(struct mpsse_context *mpsse);
func Tristate() {}


// It is a wrapper for the mpsse C function:
//     char Version(void);
func Version() {}


// TODO - USED
// It is a wrapper for the mpsse C function:
//     char *Read(struct mpsse_context *mpsse, int size);
func (m *Mpsse) Read(size int) string {
	return C.GoString(C.Read(unsafe.Pointer(m.ctx), C.int(size)))
}


// It is a wrapper for the mpsse C function:
//     char *Transfer(struct mpsse_context *mpsse, char *data, int size);
func Transfer() {}


// It is a wrapper for the mpsse C function:
//     int FastWrite(struct mpsse_context *mpsse, char *data, int size);
func FastWrite() {}


// It is a wrapper for the mpsse C function:
//     int FastRead(struct mpsse_context *mpsse, char *data, int size);
func FastRead() {}


// It is a wrapper for the mpsse C function:
//     int FastTransfer(struct mpsse_context *mpsse, char *wdata, char *rdata, int size);
func FastTransfer() {}
