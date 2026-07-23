package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/go-yaaf/yaaf-code-gen"
)

func TestGooXApiDoc(t *testing.T) {
	skipCI(t)

	gen := NewApiGenerator()

	// Get all source folders
	gp := os.Getenv("GOPATH")

	// Model
	gen.WithSourceFolder(fmt.Sprintf("%s/src/github.com/mottyc/goox-api/model", gp), "model")

	// Services
	gen.WithSourceFolder(fmt.Sprintf("%s/src/github.com/mottyc/goox-api/rest", gp), "services")

	// Refer only to src files including the following path
	gen.WithPathFilter("/github.com/mottyc/")

	// Output folder
	outDir := fmt.Sprintf("%s/src/github.com/go-yaaf/yaaf-code-gen/test_out/goox/api", gp)
	err := os.MkdirAll(outDir, os.ModePerm)
	require.Nil(t, err)

	gen.WithTargetFolder(outDir)

	gen.WithApiName("GooX API").WithApiVersion("v1.0.5")

	err = gen.Process()
	require.Nil(t, err)
}
