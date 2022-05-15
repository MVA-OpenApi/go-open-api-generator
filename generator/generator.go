package generator

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"
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
	deleteFolderRecursively(Build)

	// Generates basic folder structure
	generateFolder(Build)
	generateFolder(filepath.Join(Build, Cmd))
	generateFolder(filepath.Join(Build, Pkg))
}

func GenerateServerTemplate(port int16) {
	vars := PortConfig{port}
	fileName := "main.go"
	templateFile := "templates/server.go.tmpl"
	templateName := path.Base(templateFile)

	// Create main.go and open it
	generateFile(filepath.Join(Build, Cmd, fileName))
	file, fErr := os.Open(fileName)
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
