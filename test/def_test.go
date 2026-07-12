package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
}

// sanitizeName reduces an arbitrary (possibly annotation-derived) identifier to
// a single, path-safe file-name component. Model names such as service TsName,
// @Context and @Path values originate from free-text source comments; using them
// verbatim in an output path would allow path traversal (e.g. "../../../etc/x").
// This keeps only the final path component and strips separators / parent refs.
func sanitizeName(name string) string {
	// Normalize Windows-style separators so traversal is collapsed regardless of
	// the host OS, then keep only the final path component.
	normalized := strings.ReplaceAll(strings.TrimSpace(name), "\\", "/")
	base := filepath.Base(filepath.Clean(normalized))
	base = strings.ReplaceAll(base, "/", "")
	base = strings.ReplaceAll(base, "\\", "")
	if base == "." || base == ".." {
		return ""
	}
	return base
}

// confinedJoin joins base with the given (already file-name-shaped) segment and
// verifies the result stays within base. It is a defense-in-depth check on top
// of sanitizeName: it returns an error if the resulting path would escape base.
func confinedJoin(base, segment string) (string, error) {
	joined := filepath.Join(base, segment)
	absBase, err := filepath.Abs(base)
	if err != nil {
		return "", err
	}
	absJoined, err := filepath.Abs(joined)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(absBase, absJoined)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("refusing to write outside target folder: %q", joined)
	}
	return joined, nil
}
