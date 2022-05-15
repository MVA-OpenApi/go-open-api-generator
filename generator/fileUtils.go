package generator

import (
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func generateFolder(path string) {
	err := os.MkdirAll(path, os.ModePerm)

	check(err)
}

func deleteFolderRecursively(path string) {
	err := os.RemoveAll(path)

	check(err)
}

// Creates a file in a specific path.
// Returns a pointer of the created file.
// Responsibility for closing is on the caller.
func generateFile(path string) *os.File {
	file, err := os.Create(path)

	check(err)

	return file
}
