package generator

import (
	"fmt"
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

	file := generateFile(filepath.Join(Build, Cmd, "main.go"))
	defer file.Close()

	templateFile := "templates/server.go.tmpl"
	templateName := path.Base(templateFile)
	tmpl := template.Must(template.New(templateName).ParseFiles(templateFile))

	err := tmpl.Execute(file, vars)
	if err != nil {
		fmt.Println(err.Error())
	}
}
