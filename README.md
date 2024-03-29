# go-open-api-generator

A Command-Line Interface generator for REST APIs using Go's <a href="https://echo.labstack.com/">Echo</a> from a given <a href="https://www.openapis.org/">OpenAPI 3</a> Specification file in either JSON or YAML format.

# Purpose

We aim to make the life of Golang REST API developers (or non technical users) easier by creating a tool which takes an OpenAPI 3 Specification file as input and generates a basic project structure from it so the developers can focus on the business logic. Our currently supported features are:

- Stub handlers for each endpoint and request specified in the input file.
- Go structs for each schema present in the input file.
- Validation of parameters from the OpenAPI rules.
- Config file for environment variables.
- A frontend interface which allows the developer to interact with the generated endpoints without using an external tool like Postman.
- Generation of a simple web client written in react.js to test CRUD operations for the component schemas.
- An option to add middleware logging into a file for the generated REST API.
- An option to integrate boilerplate code for a small SQLite3 database.
- An option to integrate HTTP2 support.
- Supports API keys authorization (globally or per operation).
- Dockerfile and Lifecycle functions generation.
- Generation of BDD Testing from a `.feature` file.

# Contributions

This project was made by 6 students of the TU Berlin as part of the module "Moderne Verteilte Anwendungen Programmierpraktikum" when studying B.Sc Computer Science.

# Prerequisites

Golang (You can find an installation guide for Golang <a href="https://go.dev/">here</a>).

Godog (Only for BDD testing. You can find an installation guide here [godog](https://github.com/cucumber/godog)).

Prerequisite for HTTP/2 is a TLS connection, to generate a quick localhost certificate use either openssl or
`go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost`.

NodeJS and NPM (You can find the installer here <a href="https://nodejs.org/en/">here</a>).

# Setup

After cloning this repository, run the following command inside the repository folder to get all the required dependencies:

`go mod tidy`

# Usage

Generates a REST API template from a given OpenAPI Specification file.
Let's take one of the example files that are already in the project. For the sake of convenience we're going to be using `stores.yaml`.</br>
You can check how the file looks like <a href="https://github.com/MVA-OpenApi/go-open-api-generator/blob/main/examples/stores.yaml">here</a></br>

<b>Step 1: </b>We navigate to the repository folder</br>

<b>Step 2: </b>Run the command `go run main.go generate ./examples/stores.yaml -o ./build -n build -l -d`. A description of the flags can be found [below](https://github.com/MVA-OpenApi/go-open-api-generator/edit/main/README.md#flags).</br>

<b>Step 3: </b> We can now navigate to the output folder (in this case `build`) and run `go run main.go` to launch the REST API.

<b>Step 4: </b> The React Frontend is available at 'server-url'/frontend the server documentation is available under 'server-url'/doc.

## Flags

- `-o [Output path]`. Specifies the output path for the generated REST API.
- `-n [Module name]`. Specifies the go module name.
- `-l [Use logger]`. Enables logging middleware on the generated REST API.
- `-d [Use database]`. Generates boilerplate code for a basic SQLite3 database.
- `-H [Use HTTP2]`. Enables HTTP2 support.
- `-L [Use Lifecycle endpoints]`. Generates Livez and Readyz endpoints.

# Makefile

Available makefile commands:

- `make generate OPEN_API_PATH=path/to/open-api-file`. This command will generate the minimum project structure (no optional flags are set). The parameter `OPEN_API_PATH` is required.
- `make generate-all-flags OPEN_API_PATH=path/to/open-api-file MODULE_NAME=module-name`. This command will generate the maximum project structure (all optional flags are set). The parameter `OPEN_API_PATH` is required.
- `make build OUTPUT_NAME=executable-name`. This command will build an executable which can be used by the developer outside of the project repository.
- `make test`. This command runs the unit tests for the generator.

# BDD Generation

The input feature file has to have the following structure in order for the generator to create the godog test file

- ```Scenario: Test GET Request for url <"regex of the url">
    When I send GET request to <"actual endpoint"> with payload <"payload that needs to be sent">
    Then The response for url <"endpoint again"> with request method <"request method"> should be <status code>
  ```
- So an example for a scenario would be
- ```
    Scenario: Test GET Request for url "/store/{id}"
    When I send GET request to "/store/100" with payload ""
    Then The response for url "/store/100" with request method "GET" should be 404
  ```
- After you have created a feature file that follows this structure, you can generate the godog file by running `go run main.go generate-bdd <path to the file>`

# Examples

You can find a few OpenAPI 3 Specification file examples <a href="https://github.com/MVA-OpenApi/go-open-api-generator/tree/main/examples">here</a>.
