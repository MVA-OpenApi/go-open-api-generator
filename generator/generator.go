package generator

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"

	fs "go-open-api-generator/fileUtils"
)

const (
	Build = "build"
	Cmd   = "cmd"
	Pkg   = "pkg"
)

type PortConfig struct {
	Port int16
}

func CreateBuildDirectory() {
	// Removes previously generated folder structure
	fs.DeleteFolderRecursively(Build)

	// Generates basic folder structure
	fs.GenerateFolder(Build)
	fs.GenerateFolder(filepath.Join(Build, Cmd))
	fs.GenerateFolder(filepath.Join(Build, Pkg))
}

func GenerateServerTemplate(port int16) {
	vars := PortConfig{port}
	fileName := "main.go"
	templateFile := "templates/server.go.tmpl"
	templateName := path.Base(templateFile)

	// Create main.go and open it
	mainPath := filepath.Join(Build, Cmd, fileName)
	fs.GenerateFile(mainPath)
	file, fErr := os.OpenFile(mainPath, os.O_WRONLY, os.ModeAppend)
	if fErr != nil {
		fmt.Println(fErr.Error())
	}
	defer file.Close()

	// Parse the tempalte and write into main.go
	tmpl := template.Must(template.New(templateName).ParseFiles(templateFile))
	tmplErr := tmpl.Execute(file, vars)
	if tmplErr != nil {
		fmt.Println(tmplErr.Error())
	}
}
