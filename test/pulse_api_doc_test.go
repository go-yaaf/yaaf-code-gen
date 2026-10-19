package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/go-yaaf/yaaf-code-gen"
)

func TestPulseApiDoc(t *testing.T) {
	skipCI(t)

	gen := NewApiGenerator()

	// Get all source folders
	gp := os.Getenv("GOPATH")

	// Model
	gen.WithSourceFolder(fmt.Sprintf("%s/src/bitbucket.org/shieldiot/pulse/pulse-model", gp), "model")

	// Services
	gen.WithSourceFolder(fmt.Sprintf("%s/src/bitbucket.org/shieldiot/pulse/pulse-api/rest", gp), "services")

	// Refer only to src files including the following path
	gen.WithPathFilter("/bitbucket.org/shieldiot/")

	// Output folder
	outDir := fmt.Sprintf("%s/src/github.com/go-yaaf/yaaf-code-gen/test_out/pulse/api", gp)
	err := os.MkdirAll(outDir, os.ModePerm)
	require.Nil(t, err)

	gen.WithTargetFolder(outDir)

	gen.WithApiName("ShieldEdge API").WithApiVersion("v2.0.x")

	err = gen.Process()
	require.Nil(t, err)
}
