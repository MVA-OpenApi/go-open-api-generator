package generator

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func check(e error) {
	if e != nil {
		log.Error().Err(e).Msg("")
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

func CopyFile(sourcePath string, destinationPath string, fileName string) {
	// check if file at sourcePath exists
	if !CheckIfFileExists(sourcePath)  {
		log.Error().Msg("Failed to copy file from " + sourcePath + " to " + destinationPath + "because file doesn't exists.")
	}

	//Read all the contents of the  original file
    bytesRead, err := ioutil.ReadFile(sourcePath)
    if err != nil {
        log.Error().Err(err).Msg("")
    }

    //Copy all the contents to the desitination file
    err = ioutil.WriteFile(filepath.Join(destinationPath, fileName), bytesRead, 0755)
    if err != nil {
        log.Error().Err(err).Msg("")
    }
}