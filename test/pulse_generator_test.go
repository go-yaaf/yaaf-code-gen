package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/go-yaaf/yaaf-code-gen"
)

func TestPulseGenerator(t *testing.T) {
	skipCI(t)

	gen := NewCodeGenerator()

	// Get all source folders
	gp := os.Getenv("GOPATH")
	f1 := fmt.Sprintf("%s/src/bitbucket.org/shieldiot/pulse/pulse-model", gp)
	gen.WithSourceFolder(f1, "model")

	f2 := fmt.Sprintf("%s/src/bitbucket.org/shieldiot/pulse/pulse-api/rest", gp)
	gen.WithSourceFolder(f2, "services")

	gen.WithPathFilter("/bitbucket.org/shieldiot/")

	outDir := fmt.Sprintf("%s/src/bitbucket.org/shieldiot/pulse/pulse-api/client_lib/axios/src/lib", gp)

	err := os.MkdirAll(outDir, os.ModePerm)
	require.Nil(t, err)

	gen.WithTargetFolder(outDir)
	gen.WithTargetPlatform("axios")
	err = gen.Process()
	require.Nil(t, err)
}
