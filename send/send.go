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
)
import "../protocols"

// repeat one transmission n times
const RepeatTransmit = 10

// properties of the PT2262 protocol
var pt2262 = protocols.GetPT2262Protocol()

func length(word int) (length int) {
	for word != 0 {
		word /= 10
		length++
	}

	return length
}

func send(code int)  {
	wordLength := 24
	for i := wordLength - 1; i >= 0; i-- {

		// decimal & operation on code word to convert it to binary
		// ex: 1361 -> 10101010001
		if code & (1 << uint(i)) == 1 {
			// transmit 1
			transmit(pt2262.One)
			fmt.Print(1)
		} else {
			// transmit 0
			transmit(pt2262.Zero)
			fmt.Print(0)
		}
	}

	// transmit the sync bit at the end
	transmit(pt2262.SyncFactor)
}

func transmit(bit protocols.HighLow)  {

}

func main()  {
	args := os.Args
	if len(args) < 2 {
		log.Println("parameters not valid, require code")
		log.Println("send 9999")
	}

	codeStr := os.Args[1]
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("calling setup")
	init := C.wiringPiSetup()

	if init == -1 {
		log.Println("init failed")
		return
	}

	for i := 0; i < RepeatTransmit; i++ {
		send(code)
	}
}