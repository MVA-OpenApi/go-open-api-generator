package tests

import (
	"go-open-api-generator/fileUtils"
	_ "go-open-api-generator/tests"
	"log"
	"os"
	"testing"
)

func TestFileCreation(t *testing.T) {
	generator.GenerateFile("test_file.txt")
	if _, err := os.Stat("test_file.txt"); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("File does not exist")
		}
	}
	e := os.Remove("test_file.txt")
	if e != nil {
		log.Fatal(e)
	}
}

func TestCheckIfFileExiscdts(t *testing.T) {
	if !generator.CheckIfFileExists("fileUtils/fileUtils.go") {
		t.Errorf("File exists but it could not be found")
	}
	if generator.CheckIfFileExists("fileUtils/fileUtils.txt") {
		t.Errorf("File does not exist but CheckIfFileExists returned true")
	}
}

func TestGenerateFolder(t *testing.T) {
	generator.GenerateFolder("test_folder")
	if _, err := os.Stat("test_folder"); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Folder does not exist")
		}
	}
	generator.DeleteFolderRecursively("test_folder")
}

func TestDeleteFolderRecursively(t *testing.T) {
	generator.GenerateFolder("test_folder")
	generator.DeleteFolderRecursively("test_folder")
	if _, err := os.Stat("test_folder"); err != nil {
		if !os.IsNotExist(err) {
			t.Errorf("Folder was not deleted")
		}
	}
}
