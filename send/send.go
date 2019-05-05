//+build linux, arm

package main

// -g (debug) -Wall (show warnings) -Wno-unused-function (but not those related to unused functions)
// -lwiringPi (-l (load lib) wiringPi (wiringPi.so))

// #cgo CFLAGS: -g -Wall -Wno-unused-function
// #cgo LDFLAGS: -lm -lwiringPi
// #include <wiringPi.h>
import "C"
import (
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"
)
import "../protocols"

// repeat one transmission n times
const RepeatTransmit = 10

// properties of the PT2262 protocol
var pt2262 = protocols.GetPT2262Protocol()

func send(pin int, word string)  {

	// add zeros til size 24
	for len(word) <= 23 {
		word = "0" + word
	}

	for i := 0; i <= 23; i++ {

		bit := string(word[i])

		if bit == "1" {
			fmt.Print(1)
			// transmit 1
			transmit(pin, pt2262.One)
		} else if bit == "0" {
			fmt.Print(0)
			// transmit 0
			transmit(pin, pt2262.Zero)
		} else {
			log.Println("invalid word")
			syscall.Exit(-1)
		}
	}

	// transmit the sync bit at the end
	transmit(pin, pt2262.SyncFactor)
	fmt.Println("")
	C.digitalWrite(C.int(pin), C.LOW)
}

func transmit(pin int, bit protocols.HighLow)  {
	// write the value 1 (high) on the pin...
	C.digitalWrite(C.int(pin), C.HIGH)
	// ...for pulse length * bit high length microseconds
	C.delayMicroseconds(C.uint(pt2262.PulseLength * bit.High))
	// then write 0 (low)...
	C.digitalWrite(C.int(pin), C.LOW)
	// ...for pulse length * bit low length microseconds
	C.delayMicroseconds(C.uint(pt2262.PulseLength * bit.Low))
}

func main()  {
	args := os.Args
	if len(args) < 3 {
		log.Println("parameters not valid, require pin, code")
		log.Println("ex: send 0 9999")
		return
	}

	strPin := os.Args[1]
	pin, _ := strconv.Atoi(strPin)

	word := os.Args[2]

	log.Println("calling setup")
	init := C.wiringPiSetup()

	if init == -1 {
		log.Println("init failed")
		return
	}

	code, _ := strconv.Atoi(word)

	// decimal & operation on code word to convert it to binary
	// ex: 1361 -> 10101010001
	word = strconv.FormatInt(int64(code), 2)

	log.Printf("sending %d times %d\n", RepeatTransmit, code)
	for i := 0; i < RepeatTransmit; i++ {
		send(pin, word)
	}
}