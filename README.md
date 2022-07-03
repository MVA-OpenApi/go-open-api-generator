# go-open-api-generator
A Command-Line Interface generator for REST APIs using Go's <a href="https://echo.labstack.com/">Echo</a> from a given <a href="https://www.openapis.org/">OpenAPI 3</a> Specification file in either JSON or YAML format.

# Purpose
We aim to make the life of Golang REST API developers easier by creating a tool which takes an OpenAPI 3 Specification file as input and generates a basic project structure from it so the developers can focus on the business logic. Our currently supported features are:
- Stub handlers for each endpoint and request specified in the input file.
- Go structs for each schema present in the input file.
- Config file for environment variables.
- A frontend interface which allows the developer to interact with the generated endpoints without using an external tool like Postman.
- An option to add middleware logging into a file for the generated REST API.
- An option to integrate boilerplate code for a small SQLite3 database.
- An option to integrate HTTP2 support.
- Supports API keys authorization (globally or per operation).
# Contributions

This project was made by 6 students of the TU Berlin as part of the module "Moderne Verteilte Anwendungen Programmierpraktikum" when studying B.Sc Computer Science.

# Prerequisites
Golang (You can find an installation guide for Golang <a href="https://go.dev/">here</a>).

Prerequisite for HTTP/2 is a TLS connection, to generate a quick localhost certificate use either openssl or
`go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost`

# Setup
After cloning this repository, run the following command inside the repository folder to get all the required dependencies:

```go mod tidy```

# Usage
Generates a REST API template from a given OpenAPI Specification file.
Let's take one of the example files that are already in the project. For the sake of convenience we're going to be using ```stores.yaml```.</br>
You can check how the file looks like <a href="https://github.com/MVA-OpenApi/go-open-api-generator/blob/main/examples/stores.yaml">here</a></br>

<b>Step 1: </b>We navigate to the repository folder</br>

<b>Step 2: </b>Run the command `go run main.go generate ./examples/stores.yaml -o ./build -n build -l -d`. A description of the flags can be found [below](https://github.com/MVA-OpenApi/go-open-api-generator/edit/main/README.md#flags).</br>

<b>Step 3: </b> We can now navigate to the output folder (in this case `build`) and run `go run main.go` to launch the REST API.
## Flags
- `-o [Output path]`. Specifies the output path for the generated REST API.
- `-n [Module name]`. Specifies the go module name.
- `-l [Use logger]`. Enables logging middleware on the generated REST API.
- `-d [Use database]`. Generates boilerplate code for a basic SQLite3 database.
- `-H [Use HTTP2]`. Enables HTTP2 support.

# Makefile
Available makefile commands:
- `make generate OPEN_API_PATH=path/to/open-api-file`. This command will generate the minimum project structure (no optional flags are set). The parameter `OPEN_API_PATH` is required.
- `make generate-all-flags OPEN_API_PATH=path/to/open-api-file MODULE_NAME=module-name`. This command will generate the maximum project structure (all optional flags are set). The parameter `OPEN_API_PATH` is required.
- `make build OUTPUT_NAME=executable-name`. This command will build an executable which can be used by the developer outside of the project repository.
- `make test`. ~~This command runs the unit tests for the generator.~~ (WIP)


# Examples

You can find a few OpenAPI 3 Specification file examples <a href="https://github.com/MVA-OpenApi/go-open-api-generator/tree/main/examples">here</a>. 
