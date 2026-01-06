package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"

	"github.com/go-yaaf/yaaf-code-gen/model"
)

var goPrimitiveTypes = map[string]string{
	"double":   "number",
	"float":    "number",
	"float32":  "number",
	"float64":  "number",
	"int":      "number",
	"int32":    "number",
	"int64":    "number",
	"uint":     "number",
	"uint32":   "number",
	"uint64":   "number",
	"sint":     "number",
	"sint32":   "number",
	"sint64":   "number",
	"fixed32":  "number",
	"fixed64":  "number",
	"sfixed32": "number",
	"sfixed64": "number",
	"bool":     "boolean",
	"string":   "string",
	"bytes":    "File",
	"any":      "any",
}

// FileParser is used to parse go file and extract the meta model
type FileParser struct {
	Model       *model.MetaModel
	ClassMap    map[string]*model.ClassInfo
	pathFilter  string // Filter to process only files that their path includes the filter
	parsedFiles map[string]bool
}

func NewFileParser(model *model.MetaModel, filter string) *FileParser {
	return &FileParser{
		Model:       model,
		pathFilter:  filter,
		parsedFiles: make(map[string]bool),
	}
}

// region General Methods ----------------------------------------------------------------------------------------------

// ParseFile parse go file
func (p *FileParser) ParseFile(path string) error {

	// Check filters
	// If file was already parsed, skip it
	if parsed, ok := p.parsedFiles[path]; ok {
		if parsed {
			return nil
		}
	}

	fSet := token.NewFileSet()
	result, err := parser.ParseFile(fSet, path, nil, parser.ParseComments)
	if err != nil {
		return err
	} else {
		p.parsedFiles[path] = true
	}

	for _, imp := range result.Imports {
		imPath := p.getAbsolutePath(imp.Path.Value)
		p.ParseFolder(imPath)
	}
	for _, dcl := range result.Decls {
		switch spec := dcl.(type) {
		case *ast.GenDecl:
			_ = p.processType(result.Name.Name, spec)
		case *ast.FuncDecl:
			_ = p.processServiceMethod(spec)
		default:
			//fmt.Println("ParseFile: spec is not *ast.GenDecl or *ast.FuncDecl")
		}
	}
	return nil
}

func (p *FileParser) ParseFolder(folderPath string) {

	if files, err := os.ReadDir(folderPath); err == nil {
		for _, fe := range files {
			filePath := path.Join(folderPath, fe.Name())
			if strings.HasSuffix(filePath, ".go") {
				if p.checkFilter(filePath) {
					_ = p.ParseFile(filePath)
				}
			}
		}
	}
}

// Get absolute path from relative path
func (p *FileParser) getAbsolutePath(relPath string) string {
	relPath = strings.ReplaceAll(relPath, "\"", "")
	gp := os.Getenv("GOPATH")
	srcPath := fmt.Sprintf("%s/src", gp)
	return path.Join(srcPath, relPath)
}

func (p *FileParser) checkFilter(filePath string) bool {
	if len(p.pathFilter) == 0 {
		return true
	}

	// Check the filter
	if strings.Contains(filePath, p.pathFilter) {
		return true
	}

	// Check yaaf-common
	if strings.Contains(filePath, "/yaaf-common/") {
		return true
	}

	// Check yaaf-common-net
	if strings.Contains(filePath, "/yaaf-common-net/") {
		return true
	}

	return false
}

// process file
func (p *FileParser) processType(pkgName string, decl *ast.GenDecl) error {
	if len(decl.Specs) < 1 {
		return fmt.Errorf("no specs found")
	}

	switch spec := decl.Specs[0].(type) {
	case *ast.TypeSpec:
		break
	case *ast.ImportSpec:
		return nil
	default:
		return fmt.Errorf("unknown spec type %T", spec)
	}

	// At this point, it is known that spec is of type ast.TypeSpec
	spec, _ := decl.Specs[0].(*ast.TypeSpec)

	ti := model.NewTypeInfo(spec.Name.Name)
	ti.PackageFullName = "model"  // pkgName
	ti.PackageShortName = "model" // pkgName

	p.processTypeComments(ti, decl.Doc, spec.Comment)
	switch ti.Type {
	case "@Entity":
		return p.processEntityType(ti, decl)
	case "@Data":
		return p.processDataType(ti, decl)
	case "@Enum":
		return p.processEnumType(ti, decl)
	case "@EnumValues":
		return p.processEnumValues(ti, decl)
	case "@Service":
		return p.processServiceType(ti, decl)
	default:
		return fmt.Errorf("unknown type %s", ti.Type)
	}
}

// Process class comments
func (p *FileParser) processTypeComments(ti *model.TypeInfo, groups ...*ast.CommentGroup) {

	// collect all comments
	list := make([]*ast.Comment, 0)
	for _, group := range groups {
		if group != nil {
			list = append(list, group.List...)
		}
	}

	// process comments
	for _, comment := range list {
		if line := p.trimComment(comment.Text); len(line) == 0 {
			continue
		} else {
			if strings.HasPrefix(line, "@Entity:") {
				ti.TableName = p.getTagValue(line, "@Entity:")
				ti.Type = "@Entity"
			} else if strings.HasPrefix(line, "@Data") {
				ti.Type = "@Data"
			} else if strings.HasPrefix(line, "@EnumValuesFor") {
				ti.TableName = p.getTagValue(line, "@EnumValuesFor:")
				ti.Type = "@EnumValues"
			} else if strings.HasPrefix(line, "@Enum") {
				ti.Type = "@Enum"
			} else if strings.HasPrefix(line, "@Service") {
				altName := p.getTagValue(line, "@Service:")
				if altName != "@Service" {
					ti.TsName = altName
				} else {
					ti.TsName = model.Title(ti.TsName)
				}
				ti.Type = "@Service"
			} else if strings.HasPrefix(line, "@Path") {
				ti.Path = p.getTagValue(line, "@Path:")
			} else if strings.HasPrefix(line, "@RequestHeader") {
				ti.AddHeader(p.getTagValue(line, "@RequestHeader:"))
			} else if strings.HasPrefix(line, "@Context") {
				ti.Context = p.getTagValue(line, "@Context:")
			} else if strings.HasPrefix(line, "@ResourceGroup") {
				ti.Group = p.getTagValue(line, "@ResourceGroup:")
			} else {
				ti.Docs = append(ti.Docs, line)
			}
		}
	}
}

// region Internal helpers for proto processing ------------------------------------------------------------------------

// Trim comments
func (p *FileParser) trimComment(line string) string {
	trimmed := strings.TrimSpace(line)

	if strings.HasPrefix(trimmed, "/*") {
		trimmed = strings.Replace(trimmed, "/*", "", 1)
		trimmed = strings.TrimSpace(trimmed)
	}

	if strings.HasPrefix(trimmed, "*") {
		trimmed = strings.Replace(trimmed, "*", "", 1)
		trimmed = strings.TrimSpace(trimmed)
	}

	if strings.HasPrefix(trimmed, "//") {
		trimmed = strings.Replace(trimmed, "//", "", 1)
		trimmed = strings.TrimSpace(trimmed)
	}
	return trimmed
}

// Get simple name (not canonical name)
func (p *FileParser) getSimpleName(name string) string {

	idx := strings.LastIndex(name, ".")
	if idx > 0 {
		return name[idx+1:]
	} else {
		return name
	}
}

// Extract tag value from comment with tag
func (p *FileParser) getTagValue(line string, tag string) string {
	value := strings.Replace(line, tag, "", 1)
	value = strings.TrimSpace(value)
	return value
}

// endregion
