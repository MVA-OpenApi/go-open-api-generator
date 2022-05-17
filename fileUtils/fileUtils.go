package generator

import (
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func GenerateFolder(path string) {
	err := os.MkdirAll(path, os.ModePerm)

	check(err)
}

func DeleteFolderRecursively(path string) {
	err := os.RemoveAll(path)

	check(err)
}

// Creates a file in a specific path.
func GenerateFile(path string) {
	file, err := os.Create(path)

	check(err)

	defer file.Close()
}

func CheckIfFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
