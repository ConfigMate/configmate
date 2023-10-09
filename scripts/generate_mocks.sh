#!/bin/bash

# list of package paths to generate mocks for
PACKAGES=(
    "analyzer"
    "parsers"
    "server"
)

for pkg in "${PACKAGES[@]}"; do
    SRC_DIR="${pkg}"
    DST_DIR="${SRC_DIR}"

    # Create the destination directory if it doesn't exist
    mkdir -p "${DST_DIR}"

    # Find all interfaces defined in the package and their source files
    while IFS=: read -r src_file line; do
        # Skip if inside generated directory (/gen)
        if [[ "${src_file}" == *"gen/"* ]]; then
            continue
        fi

        # Get the base name of the source file without the extension
        src_base=$(basename "${src_file}" .go)

        # Set the destination filename to {originalfile}_mocks.go
        dst_file="${DST_DIR}/${src_base}_mocks.go"

        # Generate a mock for all interfaces in the source file
        mockgen \
            -source="${src_file}" \
            -destination="${dst_file}" \
            -package="${pkg}" \
            "${pkg}"
    done < <(grep -r "^type .* interface" "${SRC_DIR}")
done