package test

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestSanitizeName verifies that annotation-derived names cannot introduce
// path separators or parent-directory references into an output file name.
func TestSanitizeName(t *testing.T) {
	cases := map[string]string{
		"UserService":            "UserService",
		"../../../etc/passwd":    "passwd",
		"..":                     "",
		".":                      "",
		"a/b/c":                  "c",
		`..\..\windows\system32`: "system32",
		"  Padded  ":             "Padded",
		"":                       "",
	}
	for in, want := range cases {
		if got := sanitizeName(in); got != want {
			t.Errorf("sanitizeName(%q) = %q, want %q", in, got, want)
		}
		// The result must never contain a path separator.
		if got := sanitizeName(in); strings.ContainsAny(got, `/\`) {
			t.Errorf("sanitizeName(%q) leaked a separator: %q", in, got)
		}
	}
}

// TestConfinedJoin verifies that a path escaping the target folder is rejected.
func TestConfinedJoin(t *testing.T) {
	base := t.TempDir()

	if got, err := confinedJoin(base, "model.ts"); err != nil {
		t.Errorf("confinedJoin rejected a valid name: %v", err)
	} else if !strings.HasPrefix(got, base) {
		t.Errorf("confinedJoin(%q, %q) = %q, expected inside base", base, "model.ts", got)
	}

	// A traversal segment must be refused.
	if _, err := confinedJoin(base, filepath.Join("..", "..", "escape.ts")); err == nil {
		t.Error("confinedJoin accepted a traversal path that escapes the target folder")
	}
}
