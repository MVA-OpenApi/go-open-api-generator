package generator

import (
	fs "go-open-api-generator/fileUtils"
	"os"
	"strings"
	"testing"
)

func TestGenerationOfProjectJson(t *testing.T) {

	config := GeneratorConfig{OpenAPIPath: "../examples/stores.yaml", OutputPath: "../build", ModuleName: "build"}
	GenerateServer(config)
	if _, err := os.Stat("../build"); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Project not generated")
		}
	}
}

func TestGenerationOfProjectYaml(t *testing.T) {
	fs.DeleteFolderRecursively("../build")
	config := GeneratorConfig{OpenAPIPath: "../examples/stores.yaml", OutputPath: "../build", ModuleName: "build"}
	GenerateServer(config)
	if _, err := os.Stat("../build"); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Project not generated")
		}
	}
}

func TestGenerationOfProjectWithFalseName(t *testing.T) {
	fs.DeleteFolderRecursively("./build")
	config := GeneratorConfig{OpenAPIPath: "stres.yaml", OutputPath: "./build", ModuleName: "build"}
	GenerateServer(config)
	_, err := os.Stat("./build")
	if !os.IsNotExist(err) {
		t.Errorf("Project generated but yaml file does not exist")
	}
}

func TestHandlers(t *testing.T) {
	fs.DeleteFolderRecursively("./build")
	config := GeneratorConfig{OpenAPIPath: "../examples/stores.yaml", OutputPath: "./build", ModuleName: "build"}
	GenerateServer(config)
	handlers := []string{"createStore.go", "deleteStoreByID.go", "getAllStores.go", "getStoreByID.go", "updateStoreByID.go",
		"handler.go"}
	for _, name := range handlers {
		if _, err := os.Stat("./build/pkg/handler/" + name); err != nil {
			if os.IsNotExist(err) {
				t.Errorf("Handler " + name + " not generated")
			}
		}
	}
}

func TestDeletionOfProjects(t *testing.T) {
	config := GeneratorConfig{OpenAPIPath: "../examples/stores.yaml", OutputPath: "./build", ModuleName: "build"}
	GenerateServer(config)
	fs.DeleteFolderRecursively("./build")
	config = GeneratorConfig{OpenAPIPath: "stores.yaml", OutputPath: "./mock", ModuleName: "mock"}
	GenerateServer(config)
	_, err := os.Stat("./build")
	if !os.IsNotExist(err) {
		t.Errorf("Projects name should be mock not build")
	}
}

func TestCreationOfFileFromTemplate(t *testing.T) {
	fs.DeleteFolderRecursively("../build")
	fs.DeleteFolderRecursively("../mock")

	config := GeneratorConfig{OpenAPIPath: "../examples/stores.yaml", OutputPath: "../build", ModuleName: "build"}
	GenerateServer(config)

	_, err := os.Stat("../build/cmd/main.go")
	if os.IsNotExist(err) {
		t.Errorf("Main file not created")
	}
	content, errorFile := os.ReadFile("../build/cmd/main.go")
	if errorFile != nil {
		t.Errorf("Error reading file")
	}
	if !strings.Contains(string(content), "\te.Logger.Fatal(e.Start(\":8000\"))\n") {
		t.Errorf("False port number")
	}
}
