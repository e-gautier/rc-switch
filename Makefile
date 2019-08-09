CC=arm-linux-gnueabihf-gcc

.PHONY: help build clean
.DEFAULT_GOAL := build

help: ## help
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

build: ## build main
	env GOOS=linux GOARCH=arm CC=${CC} CC_FOR_TARGET=${CC} CGO_ENABLED=1 go build -i -v -o build/sniffer sniffer/sniffer.go
	env GOOS=linux GOARCH=arm CC=${CC} CC_FOR_TARGET=${CC} CGO_ENABLED=1 go build -i -v -o build/send send/send.go

clean: ## clean build
	rm -f build/*