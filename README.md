# rc-switch

## requirements
- raspberry pi with GPIO and antenna receiver plugged on GPIO 2
- Go 1.12
- wiring pi lib https://github.com/WiringPi/WiringPi
- raspberry pi toolchain https://github.com/raspberrypi/tools
- Wiringpi installed on the target raspberry pi http://wiringpi.com/download-and-install/
## setup cross compile
get the toolchain
```bash
cd ~/Download && git clone https://github.com/raspberrypi/tools
```
add compiler to path
```bash
PATH=$PATH:~/Downloads/tools/arm-bcm2708/gcc-linaro-arm-linux-gnueabihf-raspbian-x64/bin/
```
build wiring pi lib (command has to be ran on the target!)
```bash
sudo apt install git && cd /tmp && git clone https://github.com/WiringPi/WiringPi && cd WiringPi && ./build
```
repatriate the compiled .so(s) to the compiler
```bash
scp <user@pi>:/usr/local/lib/libwiringPi* ~/Download/tools/arm-bcm2708/gcc-linaro-arm-linux-gnueabihf-raspbian-x64/arm-linux-gnueabihf/lib/
```
cross the hell compile to arm
```bash
env GOOS=linux GOARCH=arm CC=arm-linux-gnueabihf-gcc CC_FOR_TARGET=arm-linux-gnueabihf-gcc CGO_ENABLED=1 go build -i -v -o build/rc-switch main.go
```
send binary to target
```bash
scp build/rc-switch <user@pi>:/tmp
```
run
```bash
cd /tmp && ./rc-switch
```
### tests
easy way to find out if compilation is well made and lib is loaded in the wrapper:
```bash
env GOOS=linux GOARCH=arm CC=arm-linux-gnueabihf-gcc CC_FOR_TARGET=arm-linux-gnueabihf-gcc CGO_ENABLED=1 go build -i -v -o build/test tests/testCompilation.go
```
```bash
./build/test 
2019/04/22 15:06:16 calling setup
2019/04/22 15:06:16 calling handler
compilation succeed, lib loaded
```

### sources
- http://rfelektronik.se/manuals/Datasheets/HX2262.pdf
- https://github.com/sui77/rc-switch
- http://wiringpi.com/reference
- https://golang.org/cmd/cgo/