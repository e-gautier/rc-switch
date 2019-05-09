//+build linux, arm

package main

// -g (debug) -Wall (show warnings) -Wno-unused-function (but not those related to unused functions)
// -lwiringPi (-l (load lib) wiringPi (wiringPi.so))

// #cgo CFLAGS: -g -Wall -Wno-unused-function
// #cgo LDFLAGS: -lm -lwiringPi
// #include <wiringPi.h>
// extern void signalHandler();
// static void wiringPiISRWrapper(int pin, int edgeType) {
//	wiringPiISR(pin, edgeType, &signalHandler);
// }
import "C"
import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)
import "../protocols"

// max high/low changes per frame
// We can handle up to 32 bit * 2 H/L changes per bit + 2 for sync
const MaxChangesPerFrame = 67

// NReceiveTolerance 60 microseconds TODO define
const NReceiveTolerance = 60

// properties of the PT2262 protocol
var pt2262 = protocols.GetPT2262Protocol()

// an array of successive signals which happened at a specific time
var frame [MaxChangesPerFrame]int

// globals used to build the frame with the callback handler over time
var lastTime = 0
var changeCount = 0
var repeatCount = 0

func diff(A int, B int) int {
	return int(math.Abs(float64(A) - float64(B)))
}

func decode() bool {

	// ignore short transmission to avoid noise interpretation
	if changeCount <= 7 {
		return false
	}

	// Assuming the longer pulse length is the pulse captured in frame[0]
	var binaryCode string

	var syncLengthInPulses int
	if pt2262.SyncFactor.Low > pt2262.SyncFactor.High {
		syncLengthInPulses = pt2262.SyncFactor.Low
	} else {
		syncLengthInPulses = pt2262.SyncFactor.High
	}

	delay := frame[0] / syncLengthInPulses
	delayTolerance := delay * NReceiveTolerance / 100

	var firstDataTransmissionTime int
	if pt2262.InvertedSignal {
		firstDataTransmissionTime = 2
	} else {
		firstDataTransmissionTime = 1
	}

	// iterate over the signal frame to make a 32 bit word (big endian way)
	for i := firstDataTransmissionTime; i < changeCount-1; i += 2 {
		if diff(frame[i], delay*pt2262.Zero.High) < delayTolerance && diff(frame[i+1], delay*pt2262.Zero.Low) < delayTolerance {
			// zero
			binaryCode += "0"
		} else if diff(frame[i], delay*pt2262.One.High) < delayTolerance && diff(frame[i+1], delay*pt2262.One.Low) < delayTolerance {
			// one
			binaryCode += "1"
		} else {
			// failure
			return false
		}
	}

	code, _ := strconv.ParseInt(binaryCode, 2, 64)
	fmt.Println("code: ", code)
	fmt.Println("binary code: ", binaryCode)
	fmt.Println("length: ", len(binaryCode))
	fmt.Println("received bit strength: ", (changeCount-1)/2)
	fmt.Println("delay: ", delay)
	return true
}

//export signalHandler
// signalHandler is called by wiring-pi whenever there is a up/down signal recorded on the GPIO pin
// this function try to build a frame according to TODO defined why and send it to the protocol decoder
func signalHandler() {
	// time in microseconds since wiringPiSetup was called
	var timeSinceInit = int(C.micros())

	// ____| <- time elapsed before the signal change
	signalDuration := timeSinceInit - lastTime

	// if the signal is longer than the max time allowed between signals then it's probably:
	// - a too long signal
	// - a silence
	// - a gap between two signals
	if signalDuration > pt2262.Gap {

		// if the signal duration is close to the first recorded signal duration then it's probably:
		// x a too long signal
		// x a silence
		// v a gap between two signals
		if diff(signalDuration, frame[0]) < 200 {

			// increment the repeat counter which is used to confirm since most of protocol repeat twice a signal on
			// sent
			repeatCount++

			// if it equals 2 then we assume that the signal is saved in the frame
			if repeatCount == 2 {

				// since the signal is saved to the global frame we need to decode it following a protocol
				decode()

				// reset the confirm counter
				repeatCount = 0
			}
		}

		// reset change counter
		changeCount = 0
	}

	// prevent frame overflow, if there is no success but we filled the 32 bit frame then drop it and start
	// again
	if changeCount >= MaxChangesPerFrame {
		changeCount = 0
		repeatCount = 0
	}

	// add the signal duration to the frame and increment the index
	frame[changeCount] = signalDuration
	changeCount++

	// save the last signal time, it will permit to calculate the next signal duration
	lastTime = timeSinceInit
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Println("parameters not valid, require pin")
		log.Println("./sniffer 2")
		return
	}

	pinStr := os.Args[1]
	pin, err := strconv.Atoi(pinStr)
	if err != nil {
		log.Println(err)
		return
	}

	// init wiring-pi library
	log.Println("calling setup")
	init := C.wiringPiSetup()

	if init == -1 {
		log.Println("init failed")
		return
	}

	// this wiring-pi method will call the handle callback any time a high/low change happens on the GPIO pin
	log.Println("calling handler")
	C.wiringPiISRWrapper(C.int(pin), C.INT_EDGE_BOTH)

	// sleep the main thread til SIGTERM
	for {
		time.Sleep(10 * time.Second)
	}
}
