package tests

import (
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
	config := generator.GeneratorConfig{OpenAPIPath: "stres.yaml", OutputPath: "./build", ModuleName: "build"}
	generator.GenerateServer(config)
	if _, err := os.Stat("./build"); err != nil {
		if !os.IsNotExist(err) {
			t.Errorf("Project generated but yaml file does not exist")
		}
	}
}
