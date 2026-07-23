package processor

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-yaaf/yaaf-code-gen/model"
)

// SanitizeName reduces an arbitrary (possibly annotation-derived) identifier to
// a single, path-safe file-name component. Model names such as service TsName,
// @Context and @Path values originate from free-text source comments; using them
// verbatim in an output path would allow path traversal (e.g. "../../../etc/x").
// This keeps only the final path component and strips separators / parent refs.
func SanitizeName(name string) string {
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

// ConfinedJoin joins base with the given (already file-name-shaped) segment and
// verifies the result stays within base. It is a defense-in-depth check on top
// of sanitizeName: it returns an error if the resulting path would escape base.
func ConfinedJoin(base, segment string) (string, error) {
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

// Processor interface
type Processor interface {
	Start() error
}

// BaseProcessor parses proto files and generates abstract meta Model
type BaseProcessor struct {
	Output     string
	Model      *model.MetaModel
	ApiName    string // API Documentation name
	ApiVersion string // API version
}

// FileCopy copies a single file from src to dst
func (p *BaseProcessor) FileCopy(src, dst string) error {
	var err error
	var srcFd *os.File
	var dstFd *os.File
	var srcInfo os.FileInfo

	if srcFd, err = os.Open(src); err != nil {
		return err
	}
	defer func() {
		_ = srcFd.Close()
	}()

	if dstFd, err = os.Create(dst); err != nil {
		return err
	}
	defer func() {
		_ = dstFd.Close()
	}()

	if _, err = io.Copy(dstFd, srcFd); err != nil {
		return err
	}
	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// DirCopy copies a whole directory recursively
func (p *BaseProcessor) DirCopy(src string, dst string) error {
	var err error
	var fds []os.DirEntry
	var fileInfo os.FileInfo

	if fileInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, fileInfo.Mode()); err != nil {
		return err
	}

	if fds, err = os.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcPath := path.Join(src, fd.Name())
		dstPath := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = p.DirCopy(srcPath, dstPath); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = p.FileCopy(srcPath, dstPath); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

// TrimNewLines remove multiple new lines for better readability
func (p *BaseProcessor) TrimNewLines(source string) string {
	// Remove newlines
	result := strings.ReplaceAll(source, "\n\n\n\n", "\n\n")
	result = strings.ReplaceAll(result, "\n\n\n", "\n\n")
	return result
}
