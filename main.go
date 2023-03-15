package main

import (
	_ "flag"
	"fmt"
	p "github.com/go-yaaf/yaaf-code-gen/parser"
	"os"
	_ "strings"

	. "github.com/go-yaaf/yaaf-code-gen/processor"
)

// main entry point
func main() {

	// Set command line arguments
	//var format string
	//flag.StringVar(&format, "f", "", "Output format: html | ts")
	//flag.Parse()
	//fmt.Println(format)
	//
	//if strings.ToLower(format) == "html" {
	//	runHtmlProcessor()
	//} else if strings.ToLower(format) == "ts" {
	//	runTsProcessor()
	//} else {
	//	fmt.Println("Unrecognized format. Must be: html | ts")
	//}

	gp := os.Getenv("GOPATH")
	f1 := fmt.Sprintf("%s/src/bitbucket.org/shieldiot/pulse/pulse-model/common", gp)
	f2 := fmt.Sprintf("%s/src/bitbucket.org/shieldiot/pulse/pulse-model/entities", gp)
	f3 := fmt.Sprintf("%s/src/bitbucket.org/shieldiot/pulse/pulse-model/enums", gp)
	f4 := fmt.Sprintf("%s/src/bitbucket.org/shieldiot/pulse/pulse-api/rest/system", gp)
	f5 := fmt.Sprintf("%s/src/bitbucket.org/shieldiot/pulse/pulse-api/rest/user", gp)

	parser := p.NewParser().
		AddSourceFolder(f4, "").
		AddSourceFolder(f5, "").
		AddSourceFolder(f1, "").
		AddSourceFolder(f2, "").
		AddSourceFolder(f3, "").
		AddProcessor(NewLogProcessor("./output")).
		AddProcessor(NewHtmlProcessor("./output/html"))

	if err := parser.Parse(); err != nil {
		fmt.Println(err.Error())
	} else {
		return
	}
}

func runHtmlProcessor() {
	//processor := NewHtmlProcessor()
	//processor.Start("./proto")
	//fmt.Println("Done")
}

func runGoProcessor() {
	fmt.Println("Go processor not implemented")
}

func runTsProcessor() {

	fmt.Println("Ts processor not implemented")

	//processor := NewTsProcessor()
	//processor.Start("./proto")
	//fmt.Println("Done")
}
