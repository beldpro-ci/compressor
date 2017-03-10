SHELL			:=	/bin/bash
PKGS			:=	builder

all: lint test

deps:
	glide install

lint:
	golint .

install-%:
	go install -v

test:
	gofmt -s -w .
	go test -v



.PHONY: deps install lint test clean

