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
	correct := CheckIfFileExists("./fileUtils.go")
	wrong := CheckIfFileExists("./fileUtils.txt")
	if !correct {
		t.Errorf("File exists but it could not be found")
	}
	if wrong {
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
}

func TestDeleteFolderRecursively(t *testing.T) {
	DeleteFolderRecursively("../test_folder")
	if _, err := os.Stat("../test_folder"); err != nil {
		if !os.IsNotExist(err) {
			t.Errorf("Folder was not deleted")
		}
	}
}
