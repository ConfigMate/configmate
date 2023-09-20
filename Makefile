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


# Setup the ANTLR parser
ANTLR_VERSION := 4.13.1
GRAMMAR_DIR := ./grammar
ANTLR4_OUTPUT_DIR := ./gen
ANTLR4_JAR := ./lib/antlr-$(ANTLR_VERSION)-complete.jar
ANTLR4_JAR_OUTPUT_DIR := ./lib

./lib/antlr-*-complete.jar:
	chmod +x ./scripts/download_antlr_jar.sh
	./scripts/download_antlr_jar.sh $(ANTLR_VERSION) $(ANTLR4_JAR) $(ANTLR4_JAR_OUTPUT_DIR)

generate-parser: ./lib/antlr-*-complete.jar
	find $(GRAMMAR_DIR) -name '*.g4' -exec sh -c 'java -Xmx500M -cp $(ANTLR4_JAR):$(CLASSPATH) org.antlr.v4.Tool -Dlanguage=Go -o $(ANTLR4_OUTPUT_DIR) "{}"' \;

clean-parser:
	rm -rf lib/antlr-4.13.1-complete.jar
	rm -rf gen/

clean:
	rm -rf bin/

bin/:
	mkdir bin/

configm: bin/
	go build $(GO_FLAGS) -o bin/configm$(EXT) $(GO_PKG)

configm-debug: bin/
	go build $(GO_FLAGS) $(GO_DEBUG_FLAGS) -o "bin/configm$(EXT)" $(GO_PKG)

## Testing
test: mocks
	go test ./...

mocks:
	chmod +x ./scripts/generate_mocks.sh
	./scripts/generate_mocks.sh

clean_mocks:
	find ./ -type d -name 'mocks' -exec rm -rf {} +