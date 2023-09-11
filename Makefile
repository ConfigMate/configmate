### ConfigMate

VERSION := 0.0.0
GOOS := linux

GO_PKG := github.com/ConfigMate/configmate
GO_DEBUG_FLAGS := -gcflags="all=-N -l"
GO_FLAGS = -ldflags '-X "main.Version=$(VERSION)" -X "main.BuildDate=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")" -X "main.GitHash=$(shell git rev-parse HEAD)"'

clean:
	rm -rf bin/

bin/:
	mkdir bin/

configm: bin/
	go build $(GO_FLAGS) -o bin/configm $(GO_PKG)

configm-debug: bin/
	go build $(GO_FLAGS) $(GO_DEBUG_FLAGS) -o bin/configm $(GO_PKG)

