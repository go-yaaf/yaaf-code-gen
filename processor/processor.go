package processor

import (
	"fmt"
	"github.com/go-yaaf/yaaf-code-gen/model"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Processor interface
type Processor interface {
	Start() error
}

// BaseProcessor parses proto files and generates abstract meta Model
type BaseProcessor struct {
	Output string
	Model  *model.MetaModel
}

// File copies a single file from src to dst
func (p *BaseProcessor) fileCopy(src, dst string) error {
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

// Dir copies a whole directory recursively
func (p *BaseProcessor) dirCopy(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var fileInfo os.FileInfo

	if fileInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, fileInfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcPath := path.Join(src, fd.Name())
		dstPath := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = p.dirCopy(srcPath, dstPath); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = p.fileCopy(srcPath, dstPath); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

// remove multiple new lines for better readability
func (p *BaseProcessor) trimNewLines(source string) string {
	// Remove newlines
	result := strings.ReplaceAll(source, "\n\n\n\n", "\n\n")
	result = strings.ReplaceAll(result, "\n\n\n", "\n\n")
	//result = strings.ReplaceAll(result, "\n\n", "\n")
	return result
}
