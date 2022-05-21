package generator

import (
	"fmt"
	"os/exec"
	"strings"
)

// return a string with all dependecies the given cmd returned
func getPackages(cmd *exec.Cmd, targetDir string) string {
	if cmd == nil || targetDir == "" {
		fmt.Println("No executable command or target directory given");
		return ""
	}

	cmd.Dir = targetDir
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return string(stdout)
}

// function to check if an string array contains an object
func contains(array []string, object string) bool {
    for _, a := range array {
        if a == object {
            return true
        }
    }
    return false
}


func GenerateModFile() {
	
	targetDir := "./build/cmd"

	// get all dependencies abd convert them to string array
	allDepsString := getPackages(exec.Command("go", "list", "-f", "'{{ .Deps }}'"), targetDir) 
	allDeps := strings.Split(string(allDepsString), " ")
	
	// remove '[ form first item and ]' from last
	allDeps[0] = allDeps[0][2:]
	allDeps[len(allDeps) - 1] = allDeps[len(allDeps) - 1][:len(allDeps[len(allDeps) - 1]) - 3]

	// for debugging
	//fmt.Println(allDeps[0])
	//fmt.Println(allDeps[len(allDeps) - 1])

	// get standard dependencies
	standardDepsString := getPackages(exec.Command("go", "list", "std"), targetDir)
    standardDeps := strings.Split(string(standardDepsString), "\n")

	// for debugging
	//fmt.Println(standardDeps[0])
	//fmt.Println(standardDeps)
	
	// subtract standard dependencies from all dependencies
	var dependencies []string

	for i := 0; i < len(allDeps); i++ {
		if !contains(standardDeps, allDeps[i]) {
			dependencies = append(dependencies, allDeps[i])
			/* fmt.Printf("%T", allDeps[i])
			fmt.Println(allDeps[i]) */
		}
	}

	// generate mod file
	fmt.Println(dependencies)
}