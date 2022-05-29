package tests

import (
	fs "go-open-api-generator/fileUtils"
	generator "go-open-api-generator/generator"
	"os"
	"testing"
)

func TestGenerationOfProject(t *testing.T) {
	config := generator.GeneratorConfig{OpenAPIPath: "stores.yaml", OutputPath: "./build", ModuleName: "build"}
	generator.GenerateServer(config)
	if _, err := os.Stat("./build"); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Project not generated")
		}
	}
}

func TestGenerationOfProjectWithFalseName(t *testing.T) {
	fs.DeleteFolderRecursively("./build")
	config := generator.GeneratorConfig{OpenAPIPath: "stres.yaml", OutputPath: "./build", ModuleName: "build"}
	generator.GenerateServer(config)
	_, err := os.Stat("./build")
	if !os.IsNotExist(err) {
		t.Errorf("Project generated but yaml file does not exist")
	}
}

func TestHandlers(t *testing.T) {
	fs.DeleteFolderRecursively("./build")
	config := generator.GeneratorConfig{OpenAPIPath: "stores.yaml", OutputPath: "./build", ModuleName: "build"}
	generator.GenerateServer(config)
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
