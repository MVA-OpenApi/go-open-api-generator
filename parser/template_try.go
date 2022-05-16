package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

// func createFile(name string) (fp *os.File){
//     nameOfFile := parseFileName(name)
//     nameOfFile = nameOfFile + ".go"
//     f, err := os.Create(nameOfFile)
//     check(err)
//     defer f.Close()
//     return f
// }

func check(e error){
    if e != nil {
        panic(e)    
    }
}

func parseFileName(filename string) string {
    s := filename[:len(filename)-5]
    return s
}

func Parse(jsonFile string, templateFile string) {
    file, _ := ioutil.ReadFile(jsonFile)  
    m := map[string]interface{}{}
    err := json.Unmarshal(file, &m)
    if err != nil {
        fmt.Println(err.Error)
    }
    filename := parseFileName(jsonFile)
    fileToWrite, err := os.Create(filename+".go")
    // filename := parseFileName("example.json")
    var restArray []string
    var strArray []string
    for _, record := range m {
        if rec, ok := record.(map[string]interface{}); ok {
            for key, val := range rec {
                val1 := fmt.Sprint(val)
                restArray = append(restArray, key)
                strArray = append(strArray, val1)
            }
        } 
    }
    for i := 0; i < len(restArray); i++{
        temp, err := template.ParseFiles(templateFile)
        if err != nil {
            fmt.Println(err)
        }
        var tpl bytes.Buffer
        input := temp.Execute(&tpl, strArray[i])
        if(input != nil){
            fmt.Println("Error")
        }
        result := tpl.String()
            // f.WriteString(string(input))
        fmt.Fprintf(fileToWrite, "%+v", result)
    }
    //We want to iterate inside type map[coordinate:map[maximum:4 minimum:1 type:integer]] struct {}
}





