#!/bin/bash

# Check if at least 3 arguments are provided
if [ $# -lt 2 ]; then
    echo "Usage: $0 <version> <output-dir>"
    exit 1
fi

ANTLR_VER=$1
OUTPUT_DIR=$2

# Download ANTLR if not already downloaded
if [ ! -f $OUTPUT_DIR/antlr-$ANTLR_VER-complete.jar ]; then
    mkdir -p $OUTPUT_DIR

    echo "Downloading ANTLR Jar"
    curl -o $OUTPUT_DIR/antlr-$ANTLR_VER-complete.jar https://www.antlr.org/download/antlr-$ANTLR_VER-complete.jar
    echo "Done downloading ANTLR Jar"
fi