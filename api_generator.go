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
	"github.com/go-yaaf/yaaf-code-gen/processor_html"
)

// ApiGenerator is the main tool to parse source folder
type ApiGenerator struct {
	sourceFolders map[string]string // Map of source folders to namespaces
	targetFolder  string            // Root target folder for the artifacts
	pathFilters   []string          // Filter to process only files that their path includes the filter
	Model         *model.MetaModel  // The generated abstract model
	ApiName       string            // API Documentation name
	ApiVersion    string            // API version
}

func NewApiGenerator() *ApiGenerator {
	return &ApiGenerator{
		Model:         model.NewMetaModel(),
		sourceFolders: make(map[string]string),
		pathFilters:   make([]string, 0),
	}
}

// WithSourceFolder adds new Go source folder with pkg name
func (cg *ApiGenerator) WithSourceFolder(path string, pkg string) *ApiGenerator {
	cg.sourceFolders[path] = pkg
	return cg
}

// WithTargetFolder sets the target artifacts folders
func (cg *ApiGenerator) WithTargetFolder(path string) *ApiGenerator {
	cg.targetFolder = path
	return cg
}

// WithPathFilter add a filter to process only files that their path includes the filter
func (cg *ApiGenerator) WithPathFilter(filter string) *ApiGenerator {
	cg.pathFilters = append(cg.pathFilters, filter)
	return cg
}

// WithApiName sets the API documentation name
func (cg *ApiGenerator) WithApiName(name string) *ApiGenerator {
	cg.ApiName = name
	return cg
}

// WithApiVersion sets the API version
func (cg *ApiGenerator) WithApiVersion(version string) *ApiGenerator {
	cg.ApiVersion = version
	return cg
}

// WithEnumTemplate sets the enum template and map of functions
func (cg *ApiGenerator) WithEnumTemplate(template string, funcMap template.FuncMap) *ApiGenerator {
	processor.AddExternalTemplate("enum", template, funcMap)
	return cg
}

// WithClassTemplate sets the class template and map of functions
func (cg *ApiGenerator) WithClassTemplate(template string, funcMap template.FuncMap) *ApiGenerator {
	processor.AddExternalTemplate("class", template, funcMap)
	return cg
}

// WithServiceTemplate sets the service template and map of functions
func (cg *ApiGenerator) WithServiceTemplate(template string, funcMap template.FuncMap) *ApiGenerator {
	processor.AddExternalTemplate("service", template, funcMap)
	return cg
}

// Process the source folders and generate artifacts
func (cg *ApiGenerator) Process() error {

	// run the file parser to fill the metamodel
	if err := cg.parseSourceFiles(); err != nil {
		return fmt.Errorf("failed to parse source files: %s", err.Error())
	}

	// replace all aliases
	cg.Model.ReplaceAliases()

	// fill the dependencies
	cg.Model.FillDependencies()

	// generate the artifacts
	return cg.createHTMLFiles()
}

// Parse all files in the list of folders and fill the metamodel
func (cg *ApiGenerator) parseSourceFiles() error {
	fileParser := parser.NewFileParser(cg.Model, cg.pathFilters)
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
func (cg *ApiGenerator) parseFile(fileParser *parser.FileParser, filePath string, info os.FileInfo, err error) error {
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
func (cg *ApiGenerator) checkFilter(filePath string) bool {
	if len(cg.pathFilters) == 0 {
		return true
	}

	// Check the filters
	for _, filter := range cg.pathFilters {
		if strings.Contains(filePath, filter) {
			return true
		}
	}

	// Check yaaf-common
	if strings.Contains(filePath, "/yaaf-common/") {
		return true
	}

	return false
}

// Create HTML files
func (cg *ApiGenerator) createHTMLFiles() error {
	p := processor_html.NewHtmlProcessor(cg.Model, cg.targetFolder, cg.ApiName, cg.ApiVersion)
	return p.Start()
}
