package parser

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type Parser struct {
	sourceFolders map[string]string
	targetFolder  string
}

func NewParser() *Parser {
	return &Parser{
		sourceFolders: make(map[string]string),
	}
}

// AddSourceFolder adds new Go source folder with pkg name
func (p *Parser) AddSourceFolder(path string, pkg string) *Parser {
	p.sourceFolders[path] = pkg
	return p
}

// SetTargetFolder set target folder to put the artifacts
func (p *Parser) SetTargetFolder(path string) *Parser {
	p.targetFolder = path
	return p
}

// Parse run the parser and generate files
func (p *Parser) Parse() error {

	model := NewMetaModel()

	for path, pkg := range p.sourceFolders {
		if err := p.parseFolder(path, pkg, model); err != nil {
			return err
		}
	}

	// For debugging purpose, print the model
	if bytes, err := json.MarshalIndent(model, "", "  "); err == nil {
		fmt.Println(string(bytes))
	}
	return nil
}

// parse single folder
func (p *Parser) parseFolder(path, packageName string, model *MetaModel) error {
	fSet := token.NewFileSet() // positions are relative to file set

	// Parse src folder including comments
	if packages, err := parser.ParseDir(fSet, path, nil, parser.ParseComments); err != nil {
		return err
	} else {

		for _, pkg := range packages {
			p.processPackage(pkg, model)
		}
	}
	return nil
}

// process single package
func (p *Parser) processPackage(astPackage *ast.Package, model *MetaModel) {
	// Get or create package
	pi := &PackageInfo{
		Name:     astPackage.Name,
		Docs:     make([]string, 0),
		Classes:  make([]*ClassInfo, 0),
		Enums:    make([]*EnumInfo, 0),
		Services: make([]*ServiceInfo, 0),
		Sockets:  make([]*WebSocketInfo, 0),
	}
	if val, found := model.Packages[astPackage.Name]; found {
		pi = val
	} else {
		model.Packages[astPackage.Name] = pi
	}

	for _, astFile := range astPackage.Files {
		p.processPackageFile(astFile, pi)
	}
}

// process single package file
func (p *Parser) processPackageFile(astFile *ast.File, pi *PackageInfo) {
	for _, decl := range astFile.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			p.processGenDecl(decl.(*ast.GenDecl), pi)
		case *ast.FuncDecl:
			p.processFuncDecl(decl.(*ast.FuncDecl), pi)
		}
	}
}

// process general declaration
func (p *Parser) processGenDecl(decl *ast.GenDecl, pi *PackageInfo) {
	docs, md := parseComments(decl.Doc)
	if md.IsValid() == false {
		return
	}
	for _, spec := range decl.Specs {
		switch spec.(type) {
		case *ast.TypeSpec:
			p.processTypeSpec(spec.(*ast.TypeSpec), pi, docs, md)
		case *ast.ValueSpec:
			p.processEnumSpec(spec.(*ast.ValueSpec), pi, docs, md)
		}
	}
}

// If the field includes th json tag, return the tag, otherwise, return Json style name
func getJsonFieldName(astField *ast.Field) string {

	fieldName := astField.Names[0].String()
	if fieldName == "<nil>" {
		return fieldName
	}
	json := fmt.Sprintf("%s%s", strings.ToLower(fieldName[0:1]), fieldName[1:])
	if astField.Tag != nil {
		if astField.Tag.Kind == token.STRING {
			tag := astField.Tag.Value
			idx := strings.Index(tag, "json:")
			if idx < 0 {
				return json
			}
			json = tag[idx+6:]
			idx = strings.Index(json, "\"")
			json = json[0:idx]
			return json
		}
	}
	return json
}

// Return field documentation
func getFieldDoc(astField *ast.Field) []string {
	docs, _ := parseComments(astField.Doc, astField.Comment)
	return docs
}

// Extract field type
func setFieldType(fi *FieldInfo, astExpr ast.Expr) {

	switch v := astExpr.(type) {
	case *ast.Ident:
		fi.Type = v.String()
	case *ast.MapType:
		fi.IsMap = true
		fi.MapKey = getFieldType(v.Key)
		fi.Type = getFieldType(v.Value)
	case *ast.ArrayType:
		fi.IsArray = true
		fi.Type = getFieldType(v.Elt)
	case *ast.StarExpr:
		fi.Type = getFieldType(v.X)
	default:
		fi.Type = ""
		return
	}
}

func getFieldType(astExpr ast.Expr) string {
	switch v := astExpr.(type) {
	case *ast.Ident:
		return v.String()
	case *ast.MapType:
		return getFieldType(v.Key)
	case *ast.StarExpr:
		return getFieldType(v.X)
	}
	return ""
}

// Build documentation from comments group and enrich metadata
func parseComments(cgList ...*ast.CommentGroup) (doc []string, md *MetaData) {
	docs := make([]string, 0)
	md = &MetaData{}

	for _, cg := range cgList {
		if cg != nil {
			for _, c := range cg.List {
				text := strings.Trim(strings.ReplaceAll(c.Text, "//", ""), " ")
				if updateMetaData(text, md) == false {
					docs = append(docs, text)
				}
			}
		}
	}
	return docs, md
}

// Analyze comment line and update metadata flags
func updateMetaData(text string, md *MetaData) bool {
	if idx := strings.Index(text, "@Entity:"); idx > -1 {
		md.SetEntity(strings.Trim(text[idx+len("@Entity:"):], " "))
		return true
	} else if idx = strings.Index(text, "@Data"); idx > -1 {
		md.SetData(strings.Trim(text[idx+len("@Data"):], " "))
		return true
	} else if idx = strings.Index(text, "@Enum:"); idx > -1 {
		md.SetEnum(strings.Trim(text[idx+len("@Enum:"):], " "))
		return true
	} else if idx = strings.Index(text, "@Message:"); idx > -1 {
		md.SetMessage(strings.Trim(text[idx+len("@Message:"):], " "))
		return true
	} else if idx = strings.Index(text, "@Context:"); idx > -1 {
		md.SetContext(strings.Trim(text[idx+len("@Context:"):], " "))
		return true
	} else if idx = strings.Index(text, "@Path:"); idx > -1 {
		md.SetPath(strings.Trim(text[idx+len("@Path:"):], " "))
		return true
	} else if idx = strings.Index(text, "@ResourceGroup:"); idx > -1 {
		md.SetGroup(strings.Trim(text[idx+len("@ResourceGroup:"):], " "))
		return true
	} else if idx = strings.Index(text, "@RequestHeader:"); idx > -1 {
		md.AddHeader(strings.Trim(text[idx+len("@RequestHeader:"):], " "))
		return true
	} else {
		return false
	}
}
