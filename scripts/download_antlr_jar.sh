#!/bin/bash

# Check if at least 3 arguments are provided
if [ $# -lt 3 ]; then
    echo "Usage: $0 <version> <antlr-jar> <output-dir>"
    exit 1
fi

ANTLR_VER=$1
ANTLR_JAR=$2
ANTLR_DIR=$3

# Download ANTLR if not already downloaded
if [ ! -f $ANTLR_JAR ]; then
    mkdir -p $ANTLR_DIR

    echo "Downloading ANTLR Jar"
    curl -o $ANTLR_JAR https://www.antlr.org/download/antlr-$ANTLR_VER-complete.jar
    echo "Done downloading ANTLR Jar in $ANTLR_JAR"
fi