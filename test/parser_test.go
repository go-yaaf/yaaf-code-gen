package test

import (
	"fmt"
	p "github.com/go-yaaf/yaaf-code-gen/parser"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	skipCI(t)

	gp := os.Getenv("GOPATH")
	//f1 := fmt.Sprintf("%s/src/bitbucket.org/shieldiot/pulse/pulse-model/common", gp)
	//f2 := fmt.Sprintf("%s/src/bitbucket.org/shieldiot/pulse/pulse-model/entities", gp)
	// f3 := fmt.Sprintf("%s/src/bitbucket.org/shieldiot/pulse/pulse-model/enums", gp)
	f4 := fmt.Sprintf("%s/src/bitbucket.org/shieldiot/pulse/pulse-api/rest/system", gp)

	parser := p.NewParser().AddSourceFolder(f4, "")

	if err := parser.Parse(); err != nil {
		t.Fail()
	} else {
		return
	}
}
