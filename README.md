# ConfigMate - Your Configs' Best Friend 🤝 [![Lint Test and Build](https://github.com/ConfigMate/configmate/actions/workflows/lint_test_build.yml/badge.svg)](https://github.com/ConfigMate/configmate/actions/workflows/lint_test_build.yml)

ConfigMate is a tool engineered to scrutinize the contents of configuration files against a user-defined specification. The goal is to reduce errors and streamline the validation process, thereby alleviating the mental burden on developers and improving operational efficiency for businesses.

## Features
- **Custom Validation**: Define a specification for your own configuration files and use checks and structural semantics to validate them.
- **Multi-file & Multi-field Checks**: Import other specifications and compare configuration fields across files.
- **CLI Interface**: Easily run configuration checks directly from your command line and receive detailed error descriptions with visual location information.
- **VS Code Extension**: Get in-editor validation and error highlighting within Visual Studio Code. Write specification with syntax highlighting.

## Build Instructions

### Prerequisites
- [Go](https://golang.org/doc/install) >= 1.19
- [Java](https://www.java.com/en/download/help/download_options.html) >= 8

### Build
1. Clone the repository
2. Run `make configm` to build. The binary will be located in the `bin` directory.

## Usage

You can use the `help` command to get information about the available commands and flags.

## License
[MIT](https://github.com/ConfigMate/configmate/blob/master/LICENSE)
