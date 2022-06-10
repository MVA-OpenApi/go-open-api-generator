package generator

import (
	fs "go-open-api-generator/fileUtils"
	"os"
	"strings"
	"testing"
)

func TestGenerationOfProjectJson(t *testing.T) {

	config := GeneratorConfig{OpenAPIPath: "../examples/stores.yaml", OutputPath: "../build", ModuleName: "build"}
	err := GenerateServer(config)
	if err != nil {
		return
	}
	if _, err := os.Stat("../build"); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Project not generated")
		}
	}
	fs.DeleteFolderRecursively("../build")
}

func TestGenerationOfProjectYaml(t *testing.T) {
	config := GeneratorConfig{OpenAPIPath: "../examples/stores.yaml", OutputPath: "../build", ModuleName: "build"}
	err := GenerateServer(config)
	if err != nil {
		return
	}
	if _, err := os.Stat("../build"); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Project not generated")
		}
	}
	fs.DeleteFolderRecursively("../build")
}

func TestHandlers(t *testing.T) {
	config := GeneratorConfig{OpenAPIPath: "../examples/stores.yaml", OutputPath: "./build", ModuleName: "build"}
	err := GenerateServer(config)
	if err != nil {
		return
	}
	handlers := []string{"createStore.go", "deleteStoreByID.go", "getAllStores.go", "getStoreByID.go", "updateStoreByID.go",
		"handler.go"}
	for _, name := range handlers {
		if _, err := os.Stat("./build/pkg/handler/" + name); err != nil {
			if os.IsNotExist(err) {
				t.Errorf("Handler " + name + " not generated")
			}
		}
	}
	fs.DeleteFolderRecursively("./build")
}

func TestDeletionOfProjects(t *testing.T) {
	config := GeneratorConfig{OpenAPIPath: "../examples/stores.yaml", OutputPath: "./build", ModuleName: "build"}
	err := GenerateServer(config)
	if err != nil {
		return
	}
	fs.DeleteFolderRecursively("./build")
	config = GeneratorConfig{OpenAPIPath: "stores.yaml", OutputPath: "./mock", ModuleName: "mock"}
	err1 := GenerateServer(config)
	if err1 != nil {
		return
	}
	_, err2 := os.Stat("./build")
	if !os.IsNotExist(err2) {
		t.Errorf("Projects name should be mock not build")
	}
	fs.DeleteFolderRecursively("./mock")
}

func TestCreationOfFileFromTemplate(t *testing.T) {

	config := GeneratorConfig{OpenAPIPath: "../examples/stores.yaml", OutputPath: "../build", ModuleName: "build"}
	err := GenerateServer(config)
	if err != nil {
		return
	}

	_, err1 := os.Stat("../build/cmd/main.go")
	if os.IsNotExist(err1) {
		t.Errorf("Main file not created")
	}
	content, errorFile := os.ReadFile("../build/cmd/main.go")
	if errorFile != nil {
		t.Errorf("Error reading file")
	}
	if !strings.Contains(string(content), "e.Logger.Fatal(e.Start(\":8000\"))") {
		t.Errorf("False port number")
	}
	fs.DeleteFolderRecursively("../build")
}

func TestGenerationOfProjectWithFalseName(t *testing.T) {
	config := GeneratorConfig{OpenAPIPath: "stres.yaml", OutputPath: "./build", ModuleName: "build"}
	err := GenerateServer(config)
	if err != nil {
		return
	}
	_, err1 := os.Stat("./build")
	if !os.IsNotExist(err1) {
		t.Errorf("Project generated but yaml file does not exist")
	}
}
