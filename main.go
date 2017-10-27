package main

import (
	"fmt"
	"./src"
	"time"
	"os"
	"encoding/binary"
	"errors"
)

const (
	READ_ADDRESS = "\xE3"
	WRITE_ADDRESS = "\xE2"

	READ_REGISTER = "\x6B"
	WRITE_REGISTER = "\x6A"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func unpackWord(reading []byte) uint32 {
	return binary.BigEndian.Uint32(reading)
}


func convertTempReading(reading []byte) (float32, error) {
	raw := unpackWord(reading)
	if raw == 0xFFFF {
		return -1, errors.New("No thermistor plugged in.")
	}

	slope := []float32{-0.07347, -0.07835, -0.10895, -0.15663, -0.25263, -0.37143, -0.52632}
	x1 := []float32{631, 382, 248, 161, 111, 74, 54}
	y1 := []float32{18, 38, 53, 67, 80, 94, 105}

	raw &= 0x3FF

	rawf := float32(raw)

	var temperature float32
	if rawf >= x1[0] {
		// Region 7
		temperature = slope[0] * (rawf - x1[0]) + y1[0]

	} else if x1[1] <= rawf && rawf <= x1[0] - 1 {
		// Region 6
		temperature = slope[1]*(rawf-x1[1]) + y1[1]

	} else if x1[2] <= rawf && rawf <= x1[1] - 1 {
		// Region 5
		temperature = slope[2]*(rawf-x1[2]) + y1[2]

	} else if x1[3] <= rawf && rawf <= x1[2] - 1 {
		// Region 4
		temperature = slope[3]*(rawf-x1[3]) + y1[3]

	} else if x1[4] <= rawf && rawf <= x1[3] - 1 {
		// Region 3
		temperature = slope[4]*(rawf-x1[4]) + y1[4]

	} else if x1[5] <= rawf && rawf <= x1[4] - 1 {
		// Region 2
		temperature = slope[5]*(rawf-x1[5]) + y1[5]

	} else if x1[6] <= rawf && rawf <= x1[5] - 1 {
		// Region 1
		temperature = slope[6]*(rawf-x1[6]) + y1[6]

	} else {
		// Hit max temperature of the thermistor
		temperature = 105.0
	}

	return temperature, nil
}


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

	check(vec.PinHigh(src.GPIOL0))
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


	// now we need to convert the readings.
	r := []byte(adReading)
	for i := 0; i < thermCount; i++ {
		index := i * 2
		sub := r[index:index+2]
		reading, err := convertTempReading(sub)
		check(err)

		fmt.Printf("Thermistor %v:\t%v C\n", i, reading)
	}

	fmt.Printf("results: %v", adReading)
}
