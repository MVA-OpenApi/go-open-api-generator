package generator

import (
	"os"
	"testing"
)

func TestFileCreation(t *testing.T) {
	GenerateFile("../test_file.txt")
	if _, err := os.Stat("../test_file.txt"); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("File does not exist")
		}
	}
}

func TestCheckIfFileExists(t *testing.T) {

	if CheckIfFileExists("./fileUtils.go") {
		t.Errorf("File exists but it could not be found")
	}
	if CheckIfFileExists("./fileUtils.txt") {
		t.Errorf("File does not exist but CheckIfFileExists returned true")
	}
}

func TestGenerateFolder(t *testing.T) {
	GenerateFolder("../test_folder")
	if _, err := os.Stat("../test_folder"); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Folder does not exist")
		}
	}
	DeleteFolderRecursively("../test_folder")
}

func TestDeleteFolderRecursively(t *testing.T) {
	GenerateFolder("../test_folder")
	DeleteFolderRecursively("../test_folder")
	if _, err := os.Stat("../test_folder"); err != nil {
		if !os.IsNotExist(err) {
			t.Errorf("Folder was not deleted")
		}
	}
}
