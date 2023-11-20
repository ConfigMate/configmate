#!/bin/bash

# list of package paths to generate mocks for
PACKAGES=(
    "analyzer"
    "analyzer/check"
    "analyzer/spec"
    "analyzer/types"
    "files"
    "parsers"
    "server"
)

for pkg in "${PACKAGES[@]}"; do
    PKG_DIR="${pkg}"
    PKG_NAME="${pkg##*/}"
    
    # Loop over .go files in the current package directory
    for go_file in "${PKG_DIR}"/*.go; do
        # Find interfaces in the current .go file
        while IFS=: read src_file line; do
            # Get the base name of the source file without the extension
            src_base=$(basename "${src_file}" .go)

            # Set the destination filename to {originalfile}_mocks.go
            dst_file="${PKG_DIR}/${src_base}_mocks.go"

            # Generate a mock for all interfaces in the source file
            mockgen \
                -source="${src_file}" \
                -destination="${dst_file}" \
                -package="${PKG_NAME}"
        done < <(grep -H "^type [[:alnum:]]* interface" "$go_file")
    done
done