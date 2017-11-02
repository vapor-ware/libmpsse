package libmpsse

import (
	"sync"
	"unsafe"
	"fmt"
)

// #cgo pkg-config: libftdi
// #cgo CFLAGS: -I/usr/local/include/mpsse
// #cgo LDFLAGS: -lmpsse -L/usr/local/lib
// #include <stdio.h>
// #include "mpsse.h"
import "C"


const (
	// MpsseOK represents the "ok" response from an MPSSE command.
	MpsseOK = 0

	// MpsseFail represents the "failed" response from an MPSSE command.
	MpsseFail = -1
)


// Mode is an integer that is used to identify the MPSSE operating
// mode. The values here match the values in the C implementation
// enum.
type Mode int

// Supported MPSSE modes.
const (
	SPI0 Mode = 1
	SPI1 Mode = 2
	SPI2 Mode = 3
	SPI3 Mode = 4
	I2C  Mode = 5
	GPIO Mode = 6
	BITBANG Mode = 7
)

// Frequency is an integer that is used to identify the clock frequency
// for the specified mode. These values match up with the frequencies
// defined in the C implementation.
type Frequency int

// Common clock rates.
const (
	OneHundredKHZ   Frequency = 100000
	FourHundredKHZ  Frequency = 400000
	OneMHZ          Frequency = 1000000
	TwoMHZ          Frequency = 2000000
	FiveMHZ         Frequency = 5000000
	SixMHZ          Frequency = 6000000
	TenMHZ          Frequency = 10000000
	TwelveMHZ       Frequency = 12000000
	FifteenMHZ      Frequency = 15000000
	ThirtyMHZ       Frequency = 30000000
	SixtyMHZ        Frequency = 60000000
)

// Endianess defines how data is clocked in and out (MSB/LSB). These values
// match up with the endianess values defined in the C implementation.
type Endianess int

// Supported endianess values.
const (
	MSB Endianess = 0x00
	LSB Endianess = 0x08
)

// I2CAck are the values used to represent ACK and NACK for I2C. These values
// match up with the I2C ACK values defined in the C implementation.
type I2CAck int

// Supported I2C ACK values.
const (
	ACK  I2CAck = 0
	NACK I2CAck = 1
)

// GPIOPin is an integer that describes a GPIO pin. These values match up with
// the values defined in the C implementation.
type GPIOPin int

// Supported GPIO pin identifiers.
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

// Iface is the FTDI interface that should be used. These values match up with
// the values defined in the C implementation.
type Iface int

// FTDI interfaces.
const (
	InterfaceAny = 0
	InterfaceA   = 1
	InterfaceB   = 2
	InterfaceC   = 3
	InterfaceD   = 4
)


// Mpsse is a struct that holds the context information for an MPSSE session.
// It holds a reference to the C context pointer that is used for all
// commands.
type Mpsse struct {
	ctx  *C.struct_mpsse_context
	open bool
	lock sync.Mutex
}


// ok is a helper function to check if the response status of an MPSSE command
// completed successfully.
func ok(status int) bool {
	return status == MpsseOK
}


// MPSSE opens and initializes the first FTDI device found.
//
// It is a wrapper for the mpsse C function:
//     struct mpsse_context *MPSSE(enum modes mode, int freq, int endianess);
func MPSSE(mode Mode, frequency Frequency, endianess Endianess) (*Mpsse, error) {

	ctx := C.MPSSE(C.enum_modes(mode), C.int(frequency), C.int(endianess))
	d := &Mpsse{ctx, true, sync.Mutex{}}

	// on success, mpsse->open will be set to 1. on failure, mpsse-open will be
	// set to 0.
	if ctx.open == 0 {
		return nil, &MpsseError{d.ErrorString()}
	}

	return d, nil
}


// Open opens a device by VID/PID.
//
// Since the C version is a wrapper around OpenIndex for index 0, we just pass
// index 0 to the OpenIndex function here.
func Open(vid int, pid int, mode Mode, frequency Frequency, endianess Endianess, iface Iface, description *string, serial *string) (*Mpsse, error) {
	return OpenIndex(vid, pid, mode, frequency, endianess, iface, description, serial, 0)
}


// OpenIndex opens a device by VID/PID/index.
//
// It is a wrapper for the mpsse C function:
//     struct mpsse_context *OpenIndex(int vid, int pid, enum modes mode, int freq, int endianess, int interface, const char *description, const char *serial, int index);
func OpenIndex(vid int, pid int, mode Mode, frequency Frequency, endianess Endianess, iface Iface, description *string, serial *string, index int) (*Mpsse, error) {

	// The description must be passed as a C char pointer. If the
	// description is nil, we will pass a null pointer.
	var descP *C.char
	if description == nil {
		descP = nil
	} else {
		descP = C.CString(*description)
		defer C.free(unsafe.Pointer(descP))
	}

	// The serial value must be passed as a C char pointer. If the
	// description is nil, we will pass a null pointer.
	var serP *C.char
	if serial == nil {
		serP = nil
	} else {
		serP = C.CString(*serial)
		defer C.free(unsafe.Pointer(serP))
	}

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

	d := &Mpsse{ctx, true, sync.Mutex{}}

	// on success, mpsse->open will be set to 1. on failure, mpsse-open will be
	// set to 0.
	if ctx.open == 0 {
		return nil, &MpsseError{d.ErrorString()}
	}

	return d, nil
}


// Close closes the device, deinitializes libftdi, and frees the MPSSE
// context pointer.
//
// It is a wrapper for the mpsse C function:
//     void Close(struct mpsse_context *mpsse);
func (m *Mpsse) Close() {
	C.Close(m.ctx)
}


// ErrorString retrieves the last error string from libftdi.
//
// It is a wrapper for the mpsse C function:
//     const char *ErrorString(struct mpsse_context *mpsse);
func (m *Mpsse) ErrorString() string {
	return C.GoString(C.ErrorString(m.ctx))
}


// SetMode sets the appropriate transmit and receive commands based on the
// requested mode and byte order.
//
// It is a wrapper for the mpsse C function:
//     int SetMode(struct mpsse_context *mpsse, int endianess);
func (m *Mpsse) SetMode(endianess Endianess) error {
	status := int(C.SetMode(m.ctx, C.int(endianess)))

	if !ok(status) {
		return &MpsseError{m.ErrorString()}
	}
	return nil
}


// EnableBitmode enables bit-wise data transfers. Must be called after
// MPSSE() / Open() / OpenIndex().
//
// It is a wrapper for the mpsse C function:
//     void EnableBitmode(struct mpsse_context *mpsse, int tf);
func (m *Mpsse) EnableBitmode(tf int) {
	C.EnableBitmode(m.ctx, C.int(tf))
}


// SetClock sets tha appropriate divisor for the desired clock frequency.
//
// It is a wrapper for the mpsse C function:
//     int SetClock(struct mpsse_context *mpsse, uint32_t freq);
func (m *Mpsse) SetClock(freq uint32) error {
	status := int(C.SetClock(m.ctx, (C.uint32_t)(freq)))

	if !ok(status) {
		return &MpsseError{m.ErrorString()}
	}
	return nil
}


// GetClock gets the currently configured clock rate.
//
// It is a wrapper for the mpsse C function:
//     int GetClock(struct mpsse_context *mpsse);
func (m *Mpsse) GetClock() int {
	return int(C.GetClock(m.ctx))
}


// GetVid returns the vendor ID of the FTDI chip.
//
// It is a wrapper for the mpsse C function:
//     int GetVid(struct mpsse_context *mpsse);
func (m *Mpsse) GetVid() int {
	return int(C.GetVid(m.ctx))
}


// GetPid returns the product ID of the FTDI chip.
//
// It is a wrapper for the mpsse C function:
//     int GetPid(struct mpsse_context *mpsse);
func (m *Mpsse) GetPid() int {
	return int(C.GetPid(m.ctx))
}


// GetDescription returns the description of the FTDI chip, if any.
//
// It is a wrapper for the mpsse C function:
//     const char *GetDescription(struct mpsse_context *mpsse);
func (m *Mpsse) GetDescription() string {
	return C.GoString(C.GetDescription(m.ctx))
}


// SetLoopback enables or disables internal loopback.
//
// It is a wrapper for the mpsse C function:
//     int SetLoopback(struct mpsse_context *mpsse, int enable);
func (m *Mpsse) SetLoopback(enable int) error {
	status := int(C.SetLoopback(m.ctx, C.int(enable)))

	if !ok(status) {
		return &MpsseError{m.ErrorString()}
	}
	return nil
}


// SetCSIdle sets the idle state of the chip select pin. CS idles high
// by default.
//
// It is a wrapper for the mpsse C function:
//     void SetCSIdle(struct mpsse_context *mpsse, int idle);
func (m *Mpsse) SetCSIdle(idle int) {
	C.SetCSIdle(m.ctx, C.int(idle))
}


// Start sends the data start condition.
//
// It is a wrapper for the mpsse C function:
//     int Start(struct mpsse_context *mpsse);
func (m *Mpsse) Start() error {
	status := int(C.Start(m.ctx))

	if !ok(status) {
		return &MpsseError{m.ErrorString()}
	}
	return nil
}


// Write sends data out via the selected serial protocol.
//
// It is a wrapper for the mpsse C function:
//     int Write(struct mpsse_context *mpsse, char *data, int size);
func (m *Mpsse) Write(data string) error {
	dataP := C.CString(data)
	defer C.free(unsafe.Pointer(dataP))

	// FIXME -- need to check that this works. not clear that len(data) gives
	// us the size that we want. maybe unsafe.Sizeof will give the int
	// size we want? but I'm also unsure about that. will need to test.
	status := int(C.Write(m.ctx, dataP, C.int(len(data))))

	if !ok(status) {
		return &MpsseError{m.ErrorString()}
	}
	return nil
}


// Stop sends the data stop condition.
//
// It is a wrapper for the mpsse C function:
//     int Stop(struct mpsse_context *mpsse);
func (m *Mpsse) Stop() error {
	status := int(C.Stop(m.ctx))

	if !ok(status) {
		return &MpsseError{m.ErrorString()}
	}
	return nil
}


// GetAck returns the last received ACK bit.
//
// It is a wrapper for the mpsse C function:
//     int GetAck(struct mpsse_context *mpsse);
func (m *Mpsse) GetAck() int {
	return int(C.GetAck(m.ctx))
}


// SetAck sets the transmitted ACK bit.
//
// It is a wrapper for the mpsse C function:
//     void SetAck(struct mpsse_context *mpsse, int ack);
func (m *Mpsse) SetAck(ack I2CAck) {
	C.SetAck(m.ctx, C.int(ack))
}


// SendAcks causes libmpsse to send ACKs after each read byte in
// I2C mode.
//
// It is a wrapper for the mpsse C function:
//     void SendAcks(struct mpsse_context *mpsse);
func (m *Mpsse) SendAcks() {
	C.SendAcks(m.ctx)
}


// SendNacks causes libmpsse to send NACKs after each read byte in
// I2C mode.
//
// It is a wrapper for the mpsse C function:
//     void SendNacks(struct mpsse_context *mpsse);
func (m *Mpsse) SendNacks() {
	C.SendNacks(m.ctx)
}


// FlushAfterRead enables or disables flushing of the FTDI chip's RX
// buffers after each read operation. Flushing is disabled by default.
//
// It is a wrapper for the mpsse C function:
//     void FlushAfterRead(struct mpsse_context *mpsse, int tf);
func (m *Mpsse) FlushAfterRead(tf int) {
	C.FlushAfterRead(m.ctx, C.int(tf))
}


// PinHigh sets the specified pin high.
//
// It is a wrapper for the mpsse C function:
//     int PinHigh(struct mpsse_context *mpsse, int pin);
func (m *Mpsse) PinHigh(pin GPIOPin) error {
	status := int(C.PinHigh(m.ctx, C.int(pin)))

	if !ok(status) {
		return &MpsseError{m.ErrorString()}
	}
	return nil
}


// PinLow sets the specified pin low.
//
// It is a wrapper for the mpsse C function:
//     int PinLow(struct mpsse_context *mpsse, int pin);
func (m *Mpsse) PinLow(pin GPIOPin) error {
	status := int(C.PinLow(m.ctx, C.int(pin)))

	if !ok(status) {
		return &MpsseError{m.ErrorString()}
	}
	return nil
}


// SetDirection sets ths input/output direction of all pins. For use in
// BITBANG mode only.
//
// It is a wrapper for the mpsse C function:
//     int SetDirection(struct mpsse_context *mpsse, uint8_t direction);
func (m *Mpsse) SetDirection(direction uint8) error {
	status := int(C.SetDirection(m.ctx, (C.uint8_t)(direction)))

	if !ok(status) {
		return &MpsseError{m.ErrorString()}
	}
	return nil
}


// WriteBits performs a bit-wise write of up to 8 bits at a time.
//
// It is a wrapper for the mpsse C function:
//     int WriteBits(struct mpsse_context *mpsse, char bits, int size);
func (m *Mpsse) WriteBits() {}


// ReadBits performs a bit-wise read of up to 8 bits.
//
// It is a wrapper for the mpsse C function:
//     char ReadBits(struct mpsse_context *mpsse, int size);
func (m *Mpsse) ReadBits() {}


// WritePins sets the input/output value of all pins. For use in BITBANG
// mode only.
//
// It is a wrapper for the mpsse C function:
//     int WritePins(struct mpsse_context *mpsse, uint8_t data);
func (m *Mpsse) WritePins(data uint8) error {
	status := int(C.WritePins(m.ctx, (C.uint8_t)(data)))

	if !ok(status) {
		return &MpsseError{m.ErrorString()}
	}
	return nil
}


// ReadPins reads the state of the chip's pins. For use in BITBANG mode
// only.
//
// It is a wrapper for the mpsse C function:
//     int ReadPins(struct mpsse_context *mpsse);
func (m *Mpsse) ReadPins() int {
	return int(C.ReadPins(m.ctx))
}


// PinState checks if a specific pin is high or low. For use in BITBANG
// mode only.
//
// It is a wrapper for the mpsse C function:
//     int PinState(struct mpsse_context *mpsse, int pin, int state);
func (m *Mpsse) PinState(pin, state int) int {
	return int(C.PinState(m.ctx, C.int(pin), C.int(state)))
}


// Tristate places all I/O pins into a tristate mode.
//
// It is a wrapper for the mpsse C function:
//     int Tristate(struct mpsse_context *mpsse);
func (m *Mpsse) Tristate() error {
	status := int(C.Tristate(m.ctx))

	if !ok(status) {
		return &MpsseError{m.ErrorString()}
	}
	return nil
}


// Version returns the libmpsse version number.
//
// It is a wrapper for the mpsse C function:
//     char Version(void);
func Version() {}


// Read reads data over the selected serial protocol.
//
// It is a wrapper for the mpsse C function:
//     char *Read(struct mpsse_context *mpsse, int size);
func (m *Mpsse) Read(size int) string {
	resp := C.Read(m.ctx, C.int(size))
	fmt.Printf("libmpsse >> read response C: %#v\n", resp)
	fmt.Printf("libmpsse >> read response Go: %#v\n", C.GoString(resp))
	return C.GoString(resp)
}


// Transfer reads and writes data over the selected serial protocol
// (SPI only).
//
// It is a wrapper for the mpsse C function:
//     char *Transfer(struct mpsse_context *mpsse, char *data, int size);
func (m *Mpsse) Transfer() {}


// FastWrite is a function for performing fast writes in MPSSE.
//
// It is a wrapper for the mpsse C function:
//     int FastWrite(struct mpsse_context *mpsse, char *data, int size);
func (m *Mpsse) FastWrite() {}


// FastRead is a function for performing fast reads in MPSSE.
//
// It is a wrapper for the mpsse C function:
//     int FastRead(struct mpsse_context *mpsse, char *data, int size);
func (m *Mpsse) FastRead() {}


// FastTransfer is a function to perform fast transfers in MPSSE.
//
// It is a wrapper for the mpsse C function:
//     int FastTransfer(struct mpsse_context *mpsse, char *wdata, char *rdata, int size);
func (m *Mpsse) FastTransfer() {}
