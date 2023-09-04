# Makefile

# TODO in case need run with args 
# $(GO) run $(GOFLAGS) $(MAIN_FILE)

# Define the name of your main Go file
MAIN_FILE = main.go

# Define the Go compiler and flags
GO = go
GOFLAGS = -v

# Define targets and commands
.PHONY: all build run test

all: build

build:
	$(GO) build $(GOFLAGS) -o myapp $(MAIN_FILE)

run:
	$(GO) run $(MAIN_FILE)

test:
	$(GO) test $(GOFLAGS) ./...

