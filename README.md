# rc-switch
Not really rc-switch but more like a test I made of the CGO wrapper on a cross-compiled C arm lib.
## quick start
```bash
wget https://git.io/fjnXH && \
unzip send_sniffer_0.0.1.zip
```
```bash
./send 0 9999
./sniffer 2
```
## TODO
- more doc
- more experiment
- real auto tests
## requirements
- raspberry pi with GPIO,
- antenna receiver plugged on GPIO 2
- antenna transmitter plugged on GPIO 0
- Go 1.12
- wiring pi lib https://github.com/WiringPi/WiringPi
- raspberry pi toolchain https://github.com/raspberrypi/tools
- Wiring-pi lib installed on the target raspberry pi http://wiringpi.com/download-and-install/
## setup cross compile
get the toolchain
```bash
git clone https://github.com/raspberrypi/tools
```
add toolchain compilers to path
```bash
PATH=$PATH:~/tools/arm-bcm2708/gcc-linaro-arm-linux-gnueabihf-raspbian-x64/bin/
```
build wiring pi lib (**command has to be ran on the target**)
```bash
apt install git && \
cd /tmp && \
git clone https://github.com/WiringPi/WiringPi && \
cd WiringPi && \
./build
```
repatriate the compiled .so(s) to the compiler
```bash
scp <user@pi>:/usr/local/lib/libwiringPi* ~/tools/arm-bcm2708/gcc-linaro-arm-linux-gnueabihf-raspbian-x64/arm-linux-gnueabihf/lib/
```
cross the hell compile to arm
```bash
make
```
send binary to target
```bash
scp build/send build/sniffer <user@pi>:/tmp
```
run (target)
```bash
./sniffer 2
./send 0 9999
```
## tests
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

## sources
- http://rfelektronik.se/manuals/Datasheets/HX2262.pdf
- https://github.com/sui77/rc-switch
- http://wiringpi.com/reference
- https://golang.org/cmd/cgo/