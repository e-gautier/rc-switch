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
	"time"
)
import "./protocols"

// Number of maximum high/Low changes per packet.
// We can handle up to (unsigned long) => 32 bit * 2 H/L changes per bit + 2 for sync
const RCSwitchMaxChanges = 67

// MaxTimeBetweenSignals = 4300 microseconds
const MaxTimeBetweenSignals = 4300

// NReceiveTolerance 60 microseconds
const NReceiveTolerance = 60

var transmissionTime [RCSwitchMaxChanges]int
var lastTime = 0
var changeCount = 0
var repeatCount = 0

func diff(A int, B int) int {
	return int(math.Abs(float64(A) - float64(B)))
}

func receiveProtocol(count int) bool {
	// Assuming the longer pulse length is the pulse captured in transmissionTime[0]
	pt2262 := protocols.GetPT2262Protocol()
	code := 0
	var syncLengthInPulses int
	if pt2262.SyncFactor.Low > pt2262.SyncFactor.High {
		syncLengthInPulses = pt2262.SyncFactor.Low
	} else {
		syncLengthInPulses = pt2262.SyncFactor.High
	}
	delay := transmissionTime[0] / syncLengthInPulses
	delayTolerance := delay * NReceiveTolerance / 100

	var firstDataTransmissionTime int
	if pt2262.InvertedSignal {
		firstDataTransmissionTime = 2
	} else {
		firstDataTransmissionTime = 1
	}

	for i := firstDataTransmissionTime; i < count-1; i += 2 {
		// shift bit to left
		code <<= 1

		if diff(transmissionTime[i], delay*pt2262.Zero.High) < delayTolerance && diff(transmissionTime[i+1], delay*pt2262.Zero.Low) < delayTolerance {
			// zero
		} else if diff(transmissionTime[i], delay*pt2262.One.High) < delayTolerance && diff(transmissionTime[i+1], delay*pt2262.One.Low) < delayTolerance {
			// one
			code |= 1
		} else {
			// failure
			return false
		}
	}

	if count > 7 {
		// ignore short transmission to avoid noise interpretation
		fmt.Println("code: ", code)
		fmt.Println("received bit strength: ", (count-1)/2)
		fmt.Println("delay: ", delay)
		return true
	}

	return false
}

//export signalHandler
func signalHandler() {
	// time in microseconds since wiringPiSetup was called
	var timeSinceInit = int(C.micros())

	signalDuration := timeSinceInit - lastTime

	if signalDuration > MaxTimeBetweenSignals {
		// possibly a too long signal or a silence, or a gap between two signals

		if diff(signalDuration, transmissionTime[0]) < 200 {
			// if this signal is close in time duration to the first signal we got then it could be a gap between two
			// transmissions

			repeatCount++
			if repeatCount == 2 {
				receiveProtocol(changeCount)
				repeatCount = 0
			}
		}
		changeCount = 0
	}

	// prevent transmissionTime overflow
	if changeCount >= RCSwitchMaxChanges {
		changeCount = 0
		repeatCount = 0
	}

	transmissionTime[changeCount] = signalDuration
	changeCount++
	lastTime = timeSinceInit
}

func main() {
	log.Println("calling setup")
	init := C.wiringPiSetup()

	if init == -1 {
		fmt.Println("init failed")
		return
	}

	log.Println("calling handler")

	C.wiringPiISRWrapper(2, C.INT_EDGE_BOTH)

	for {
		time.Sleep(10 * time.Second)
	}
}
