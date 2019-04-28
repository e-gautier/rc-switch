//+build linux, arm

package main

// -g (debug) -Wall (show warnings) -Wno-unused-function (but not those related to unused functions)
// -lwiringPi (-l (load lib) wiringPi (wiringPi.so))

// #cgo CFLAGS: -g -Wall -Wno-unused-function
// #cgo LDFLAGS: -lm -lwiringPi
// #include <wiringPi.h>
import "C"
import (
	"log"
	"os"
	"strconv"
)

func send(code int)  {
	
}

func main()  {
	args := os.Args
	if len(args) < 2 {
		log.Println("parameters not valid")
	}

	log.Println("calling setup")
	init := C.wiringPiSetup()

	if init == -1 {
		log.Println("init failed")
		return
	}

	codeStr := os.Args[1]
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		log.Println(err)
		return
	}

	send(code)
}