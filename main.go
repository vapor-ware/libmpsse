package main

import (
	"fmt"
	"./src"
	"time"
	"os"
)

const (
	READ_ADDRESS = "\xE3"
	WRITE_ADDRESS = "\xE2"

	READ_REGISTER = "\x6B"
	WRITE_REGISTER = "\x6A"
)


func main() {
	fmt.Println("Starting test script")
	thermCount := 12

	vec, err := src.SimpleOpen(0x0403, 0x6011, src.I2C, src.ONE_HUNDRED_KHZ, src.MSB, src.INTERFACE_A)
	if err != nil {
		panic(err)
	}

	gpio, err := src.SimpleOpen(0x0403, 0x6011, src.I2C, src.ONE_HUNDRED_KHZ, src.MSB, src.INTERFACE_B)
	if err != nil {
		panic(err)
	}

	vec.PinHigh(src.GPIOL0)
	time.Sleep(1 * time.Millisecond)

	vec.Start()
	vec.Write(READ_ADDRESS)

	var adReading string
	if vec.GetAck() == int(src.ACK) {

		// if we got an ack, the slave is there
		vec.SendNacks()
		vec.Read(1)

		vec.SendAcks()
		vec.Stop()

		// set channel to 3 for MAX11608
		vec.Start()
		vec.Write("\xE2\x08")
		vec.Stop()

		// verify channel was set
		vec.Start()
		vec.Write(READ_ADDRESS)
		vec.SendNacks()
		vec.Read(1)
		vec.Stop()

		// configure max116xx
		vec.SendAcks()
		vec.Start()

		vec.Write("\x6A\xD2\x0F")

		// initiating a read starts the conversion
		vec.Start()
		vec.Write(READ_REGISTER)

		time.Sleep(10 * time.Millisecond)

		adReading = vec.Read(thermCount * 2)
	} else {
		fmt.Println("No ACK from thermistors")
		vec.Stop()
		vec.Close()
		gpio.Close()
		os.Exit(1)
	}

	vec.Stop()
	vec.Close()
	gpio.Close()

	fmt.Printf("results: %v", adReading)
}
