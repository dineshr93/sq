ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
OS := $(go version | cut -d' ' -f 4 | cut -d'/' -f 1)
ARCH := $(go version | cut -d' ' -f 4 | cut -d'/' -f 2)
BINARY_NAME :=  $(shell basename $(CURDIR))


ifeq ($(OS),Windows_NT)
	BINARY_NAME := ${BINARY_NAME}.exe
endif

ifeq ($(OS),Windows_NT)
	SHELL := powershell.exe
	.SHELLFLAGS := -NoProfile -Command
	RM_F_CMD = Remove-Item -erroraction silentlycontinue -Force
    RM_RF_CMD = ${RM_F_CMD} -Recurse
	HELP_CMD = Select-String "^[a-zA-Z_-]+:.*?\#\# .*$$" "./Makefile" | Foreach-Object { $$_data = $$_.matches -split ":.*?\#\# "; $$obj = New-Object PSCustomObject; Add-Member -InputObject $$obj -NotePropertyName ('Command') -NotePropertyValue $$_data[0]; Add-Member -InputObject $$obj -NotePropertyName ('Description') -NotePropertyValue $$_data[1]; $$obj } | Format-Table -HideTableHeaders @{Expression={ $$e = [char]27; "$$e[36m$$($$_.Command)$${e}[0m" }}, Description
else
	SHELL := bash
	RM_F_CMD = rm -f
	RM_RF_CMD = ${RM_F_CMD} -r
	HELP_CMD = grep -E '^[a-zA-Z_-]+:.*?\#\# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?\#\# "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
endif


.DEFAULT_GOAL := help
.PHONY: init add build test compile help

init: $(PROJ_DIR) ## cobra init  and creates this makefile and performs initial build (provide name of project ex: make init n=modulename)
add: $@ ## cobra add new commands only (first level) For second level use cobra add 2ndlevelCMD -p 1stlevelCmd (provide name of command make add c=commandname)
build: $@ ## Build only for this platform
test: $@ ## Performs build and does generatedbinary -h
compile: $@ ## Generating binary for every OS and Platform
copy: $@ ## Copy binary to desired environment path (please modify)

init:
	echo "make init n=modulename"
	mkdir -p ${n} && cd ${n} && go mod init github.com/dineshr93/${n} && cobra init --author "Dinesh Ravi dineshr93@gmail.com" --license Apache-2.0 --viper && go mod tidy && go build -o ./bin/${n} main.go && ./bin/${n} -h && cp ${ROOT_DIR}/Makefile .

add:
	echo "make add c=commandname"
	cobra add ${c}
	go build -o ./bin/${BINARY_NAME} main.go
	./bin/${BINARY_NAME} -h

build:
	echo "make build"
	go build -o ./bin/${BINARY_NAME} main.go
	cp ./bin/${BINARY_NAME} .

compile:
	echo "Generating binary for every OS and Platform"
	GOOS=freebsd GOARCH=386 go build -o bin/${BINARY_NAME}-freebsd-386 main.go
	GOOS=linux GOARCH=386 go build -o bin/${BINARY_NAME}-linux-386 main.go
	GOOS=windows GOARCH=386 go build -o bin/${BINARY_NAME}-windows-386.exe main.go
	GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME}-linux-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/${BINARY_NAME}-windows-amd64.exe main.go

test: build
	./bin/${BINARY_NAME} -h
copy: test
	cp bin/${BINARY_NAME} /sw/bin/

help: ## Show this help
	@${HELP_CMD}
