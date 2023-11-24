#!/bin/bash

# Setup the ANTLR parser
ANTLR_VERSION="4.13.1"
ANTLR4_JAR_OUTPUT_DIR="./lib"
ANTLR4_JAR="$ANTLR4_JAR_OUTPUT_DIR/antlr-$ANTLR_VERSION-complete.jar"
GRAMMAR_DIR="./parsers/grammars"
ANTLR4_OUTPUT_DIR="./parsers/gen"
ANTLR4_CMD="java -jar $ANTLR4_JAR"
ANTLR_FLAGS="-Dlanguage=Go -Xexact-output-dir"

echo "ANTLR Version: $ANTLR_VERSION"
echo "ANTLR JAR Directory: $ANTLR4_JAR_OUTPUT_DIR"
echo "ANTLR JAR Path: $ANTLR4_JAR"
echo "Grammar Directory: $GRAMMAR_DIR"
echo "Output Directory: $ANTLR4_OUTPUT_DIR"

# Download ANTLR jar if not present
download_antlr_jar() {
    if [ ! -f "$ANTLR4_JAR" ]; then
        echo "ANTLR JAR not found. Downloading..."
        chmod +x ./scripts/download_antlr_jar.sh
        ./scripts/download_antlr_jar.sh $ANTLR_VERSION $ANTLR4_JAR_OUTPUT_DIR
    else
        echo "ANTLR JAR is already present."
    fi
}

# Generate parsers
generate_parsers() {
    download_antlr_jar
    export CLASSPATH=".:$ANTLR4_JAR:$$CLASSPATH"

    echo "Starting parser generation for combined grammars..."

    # First loop: Find and generate parsers for combined grammar files (those not ending in Parser or Lexer)
    for grammar_file in $(find "$GRAMMAR_DIR" -name '*.g4' ! -name '*Lexer.g4' ! -name '*Parser.g4'); do
        echo "Processing combined grammar file: $grammar_file"
        base_name=$(basename "$grammar_file" ".g4")
        package_name="parser_${base_name,,}"
        output_dir="$ANTLR4_OUTPUT_DIR/$package_name"

        echo "Compiling combined grammar for: $base_name"
        $ANTLR4_CMD $ANTLR_FLAGS -package "$package_name" -o "$output_dir" "$grammar_file"
        echo "Compilation completed for combined grammar: $base_name"
    done

    echo "Starting parser generation for separated grammars..."

    # Second loop: Find Parser files and compile them together with the corresponding Lexer files
    for parser_file in $(find "$GRAMMAR_DIR" -name '*Parser.g4'); do
        echo "Processing parser file: $parser_file"
        base_name=$(basename "$parser_file" "Parser.g4")
        lexer_file="$GRAMMAR_DIR/${base_name}Lexer.g4"
        package_name="parser_${base_name,,}"
        output_dir="$ANTLR4_OUTPUT_DIR/$package_name"

        if [ -f "$lexer_file" ]; then
            echo "Compiling Lexer and Parser together for: $base_name"
            $ANTLR4_CMD $ANTLR_FLAGS -package "$package_name" -o "$output_dir" "$lexer_file" "$parser_file"
            echo "Compilation completed for Lexer and Parser: $base_name"
        else
            echo "Lexer file not found for $base_name"
        fi
    done

    echo "Parser generation process completed."
}

generate_parsers