package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/go-yaaf/yaaf-code-gen"
)

func TestGooXGenerator(t *testing.T) {
	skipCI(t)

	gen := NewCodeGenerator()

	// Get all source folders
	gp := os.Getenv("GOPATH")

	// Model
	gen.WithSourceFolder(fmt.Sprintf("%s/src/github.com/mottyc/goox-api/model", gp), "model")

	// Services
	gen.WithSourceFolder(fmt.Sprintf("%s/src/github.com/mottyc/goox-api/rest", gp), "services")

	// Refer only to src files including the following path
	gen.WithPathFilter("/github.com/mottyc/")

	// Output folder
	outDir := fmt.Sprintf("%s/src/github.com/mottyc/goox-api/client_lib/ng-workspace/projects/ngx-goox-lib2/src/lib", gp)
	err := os.MkdirAll(outDir, os.ModePerm)
	require.Nil(t, err)

	gen.WithTargetFolder(outDir)

	err = gen.Process()
	require.Nil(t, err)
}
