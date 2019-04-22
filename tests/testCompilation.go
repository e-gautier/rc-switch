//+build linux, arm

package main

// -g (debug) -Wall (show warnings) -Wno-unused-function (but not those related to unused functions)
// -lwiringPi (-l (load lib) wiringPi (wiringPi.so))

// #cgo CFLAGS: -g -Wall -Wno-unused-function
// #cgo LDFLAGS: -lm -lwiringPi
// #include <wiringPi.h>
// extern void callback();
// static void wiringPiISRWrapper(int pin, int edgeType) {
//	wiringPiISR(pin, edgeType, &callback);
// }
import "C"
import "fmt"
import "log"

//export callback
func callback() {
	fmt.Println("compilation succeed, lib loaded")
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
}
