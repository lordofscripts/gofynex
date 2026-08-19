
# Set Compiler to either of: go, gocolor OR gopretty
MODE=gopretty
# Include the GOFLAGS from .env into go build in a seamless way
include .env
export

# Compilers
GO=go
GOCOLOR=~/go/bin/colorgo
#GOPRETTY="$(HOME)/go/bin/gofilter"
GOPRETTY="$(HOME)/go/bin/gocolor"
ifeq ($(MODE),gocolor)
        GO=$(GOCOLOR)
endif
GO_TAGS=-tags x11

# - Packagers only
PKG_SEMANTIC_VERSION=$(shell sed -n '/>>>BEGIN/,/>>>END/p' version.go > /tmp/mainversion_gfx.go && go run /tmp/mainversion_gfx.go short)
PKG_FULL_VERSION:= $(patsubst v%,%,$(PKG_SEMANTIC_VERSION))
PKG_PUBLIC_NAME=$(shell grep -m 1 'module' go.mod | sed -E 's/^module\s+//p')
PKG_NAME=goapp
PKG_PREFERRED_ARCH="amd64"
PKG_REVISION=1
PKG_FULLNAME=${PKG_NAME}_${PKG_FULL_VERSION}-${PKG_REVISION}-${PKG_PREFERRED_ARCH}

# Sources
CMD_CLI=./cmd

# Outputs
BIN=./bin
BIN_OUT_CLI=$(BIN)/gofynex_demo


all: demo

# ---------------------------------------------------
demo:
ifeq ($(MODE),gopretty)
	$(GO) build $(GO_TAGS) -v -o $(BIN_OUT_CLI) $(CMD_CLI)/*.go 2>&1 | $(GOPRETTY) -color -width 75 -version
else
	$(GO) build $(GO_TAGS) -v -o $(BIN_OUT_CLI) $(CMD_CLI)/*.go
endif

# ---------------------------------------------------

# Publish package info to GO Package Repository
proxy:
	GOPROXY=proxy.golang.org go list -m $(PKG_PUBLIC_NAME)@v$(PKG_FULL_VERSION)

update:
	go get -u all

clean:
	go clean

testall:
	go test ./...

testfull:
	go test -v tests/*_test.go

version:
	@echo $(PKG_FULL_VERSION)

name:
	@echo $(PKG_FULLNAME)	
