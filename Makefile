# ConfigMate

## Do not modify; this is maintained by a Github Action. 
VERSION := $(shell grep 'current_version =' .bumpversion.cfg | sed 's/^[[:space:]]*current_version = //')

GOOS := linux
ifeq ($(OS), Windows_NT)
	GOOS := windows
	EXT := ".exe"
endif

GO_PKG := github.com/ConfigMate/configmate
GO_DEBUG_FLAGS := -gcflags="all=-N -l"
GO_FLAGS = -ldflags '-X "main.Version=$(VERSION)" -X "main.BuildDate=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")" -X "main.GitHash=$(shell git rev-parse HEAD)"'

clean: clean_parsers clean_mocks
	rm -rf bin/

bin/:
	mkdir bin/

configm: generate_parsers bin/
	go build $(GO_FLAGS) -o bin/configm$(EXT) $(GO_PKG)

configm_debug: generate_parsers bin/
	go build $(GO_FLAGS) $(GO_DEBUG_FLAGS) -o "bin/configm$(EXT)" $(GO_PKG)

generate_parsers:
	chmod +x ./scripts/generate_parsers.sh
	./scripts/generate_parsers.sh

clean_parsers:
	rm -rf lib/
	rm -rf parsers/gen/

## Testing
test: mocks
	go test ./...

mocks:
	chmod +x ./scripts/generate_mocks.sh
	./scripts/generate_mocks.sh

clean_mocks:
	find ./ -type f -name '*_mocks.go' -exec rm {} +