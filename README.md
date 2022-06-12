# go-open-api-generator

# Contributions

This project was made by 6 students of the TU Berlin as part of the module "Programmierpraktikum" when studying B.Sc Computer Science.

# Examples

You can find a few examples <a href="https://github.com/MVA-OpenApi/go-cookbook">here</a> as they are part of the cookbook. 

# Usage
Generates a REST API template from a given OpenAPI Specification file.
Let's take one of the example files that are already in the project. For the sake of convenience we're going to be using ```stores.yaml```.</br>
You can check how the file looks like <a href="https://github.com/MVA-OpenApi/go-open-api-generator/blob/main/examples/stores.yaml">here</a></br>

<b>Step 1: </b>We navigate to the folder and `go-open-api-generator`</br>

<b>Step 2: </b>Run the command `go run main.go generate ./examples/stores.yaml -o ./build -n build`. A description of the flags can be found [below](https://github.com/MVA-OpenApi/go-open-api-generator/edit/main/README.md#flags).</br>

<b>Step 3: </b> We can now navigate to the folder `build` and we can see that there are 2 folders present, `cmd` and `pkg`. The `cmd` folder contains the `main.go` file. This file is used to initialise the server. The other package(`pkg`) is used for the database and for the handlers. If we open the handler files we'll see that the functions are there but they are empty as it is part of the developer's job to fill them with functionality. 
## Flags
```
-o [Output path]
-n [Module name]
-l [Use logger]
```

# Build
`go build -o [Executable name]`
