package main

import (
	gen "go-open-api-generator/generator"
)

func main() {
	gen.CreateBuildDirectory()
	gen.GenerateServerTemplate(3000)
}
