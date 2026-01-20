package generator

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/go-yaaf/yaaf-code-gen/model"
	"github.com/go-yaaf/yaaf-code-gen/parser"
	"github.com/go-yaaf/yaaf-code-gen/processor"
)

// CodeGenerator is the main tool to parse source folder
type CodeGenerator struct {
	sourceFolders map[string]string // Map of source folders to namespaces
	targetFolder  string            // Root target folder for the artifacts
	pathFilter    string            // Filter to process only files that their path includes the filter
	Model         *model.MetaModel  // The generated abstract model
}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		Model:         model.NewMetaModel(),
		sourceFolders: make(map[string]string),
	}
}

// WithSourceFolder adds new Go source folder with pkg name
func (cg *CodeGenerator) WithSourceFolder(path string, pkg string) *CodeGenerator {
	cg.sourceFolders[path] = pkg
	return cg
}

// WithTargetFolder sets the target artifacts folders
func (cg *CodeGenerator) WithTargetFolder(path string) *CodeGenerator {
	cg.targetFolder = path
	return cg
}

// WithPathFilter sets the filter to process only files that their path includes the filter
func (cg *CodeGenerator) WithPathFilter(filter string) *CodeGenerator {
	cg.pathFilter = filter
	return cg
}

// WithEnumTemplate sets the enum template and map of functions
func (cg *CodeGenerator) WithEnumTemplate(template string, funcMap template.FuncMap) *CodeGenerator {
	processor.AddExternalTemplate("enum", template, funcMap)
	return cg
}

// WithClassTemplate sets the class template and map of functions
func (cg *CodeGenerator) WithClassTemplate(template string, funcMap template.FuncMap) *CodeGenerator {
	processor.AddExternalTemplate("class", template, funcMap)
	return cg
}

// WithServiceTemplate sets the service template and map of functions
func (cg *CodeGenerator) WithServiceTemplate(template string, funcMap template.FuncMap) *CodeGenerator {
	processor.AddExternalTemplate("service", template, funcMap)
	return cg
}

// Process the source folders and generate artifacts
func (cg *CodeGenerator) Process() error {

	// run the file parser to fill the metamodel
	if err := cg.parseSourceFiles(); err != nil {
		return fmt.Errorf("failed to parse source files: %s", err.Error())
	}

	// replace all aliases
	cg.Model.ReplaceAliases()

	// fill the dependencies
	cg.Model.FillDependencies()

	// generate the artifacts
	return cg.createTSFiles()
}

// Parse all files in the list of folders and fill the metamodel
func (cg *CodeGenerator) parseSourceFiles() error {
	fileParser := parser.NewFileParser(cg.Model, cg.pathFilter)
	for folder, _ := range cg.sourceFolders {
		if err := filepath.Walk(folder, func(filePath string, info os.FileInfo, err error) error {
			return cg.parseFile(fileParser, filePath, info, err)
		}); err != nil {
			return err
		}
	}
	return nil
}

// Parse specific file
func (cg *CodeGenerator) parseFile(fileParser *parser.FileParser, filePath string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if path.Ext(filePath) == ".go" {
		if cg.checkFilter(filePath) {
			if e := fileParser.ParseFile(filePath); e != nil {
				fmt.Println("error", e.Error())
			}
		}
	}
	return nil
}

// Check file path filter
func (cg *CodeGenerator) checkFilter(filePath string) bool {
	if len(cg.pathFilter) == 0 {
		return true
	}

	// Check the filter
	if strings.Contains(filePath, cg.pathFilter) {
		return true
	}

	// Check yaaf-common
	if strings.Contains(filePath, "/yaaf-common/") {
		return true
	}

	return false
}

// Create Typescript files
func (cg *CodeGenerator) createTSFiles() error {
	p := processor.NewTsProcessor(cg.Model, cg.targetFolder)
	return p.Start()
}
