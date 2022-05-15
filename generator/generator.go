package generator

import (
	"path/filepath"
)

const (
	Build = "build"
	Cmd   = "cmd"
	Pkg   = "pkg"
)

func CreateBuildDirectory() {
	// Removes previously generated folder structure
	deleteFolderRecursively(Build)

	// Generates basic folder structure
	generateFolder(Build)
	generateFolder(filepath.Join(Build, Cmd))
	generateFolder(filepath.Join(Build, Pkg))
}
