package main

// #include <stdio.h>
// #include <errno.h>
import "C"
import "fmt"

func main() {
	fmt.Print("test")
	C.printf("test2")
}
