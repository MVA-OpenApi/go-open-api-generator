package tests

import (
	generator "go-open-api-generator/generator"
	"os"
	"testing"
)

func TestGenerationOfProject(t *testing.T) {
	generator.GenerateServer("stores.yaml", "./build", "build")
	if _, err := os.Stat("./build"); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Project not generated")
		}
	}
}

func TestGenerationOfProjectWithFalseName(t *testing.T) {
	generator.GenerateServer("stres.yaml", "./build", "build")
	if _, err := os.Stat("./build"); err != nil {
		if !os.IsNotExist(err) {
			t.Errorf("Project generated but yaml file does not exist")
		}
	}
}
