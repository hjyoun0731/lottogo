default: help

help:
	@echo 'usage: make [target] ...'
	@echo ''
	@echo 'targets:'
	@egrep '^(.+)\:\ .*##\ (.+)' ${MAKEFILE_LIST} | sed 's/:.*##/#/' | column -t -c 2 -s '#'

all:
	make clean
	make build 

clean:
	go clean

build:	dummy
	go build api/*.go
	go build cmd/*.go
	cp ./server ./build/docker/bin/	

dummy:
