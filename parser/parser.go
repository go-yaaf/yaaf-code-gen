package parser

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
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

	model := &MetaModel{
		Packages: make(map[string]*PackageInfo),
	}

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
		}
	}
}

func (p *Parser) processGenDecl(decl *ast.GenDecl, pi *PackageInfo) {
	docs, md := getComments(decl.Doc)
	if md.Valid == false {
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

// Process structure type
func (p *Parser) processTypeSpec(spec *ast.TypeSpec, pi *PackageInfo, docs []string, md *MetaData) {

	// Create class info
	ci := &ClassInfo{
		ID:          fmt.Sprintf("%s.%s", pi.Name, spec.Name.String()),
		Name:        spec.Name.String(),
		Package:     pi.Name,
		Docs:        docs,
		IsExtend:    false,
		IsStream:    false,
		BaseClasses: make([]string, 0),
		TableName:   md.Entity,
		Fields:      make([]*FieldInfo, 0),
	}

	objType := spec.Type.(*ast.StructType)

	// process fields
	if objType.Fields != nil {
		for i, field := range objType.Fields.List {
			p.processTypeField(i, field, ci)
		}
	}
	pi.Classes = append(pi.Classes, ci)
}

// Process structure type field
func (p *Parser) processTypeField(idx int, astField *ast.Field, ci *ClassInfo) {

	// Check if it is a nameless field, in this case it is base type
	if astField.Names == nil {
		baseType := astField.Type.(*ast.Ident).String()
		ci.BaseClasses = append(ci.BaseClasses, baseType)
		ci.IsExtend = true
		return
	}

	// Do not handle private fields
	if astField.Names[0].IsExported() == false {
		return
	}

	if astField.Names[0].Obj == nil {
		return
	}

	field := &FieldInfo{
		Name:      astField.Names[0].String(),
		TsName:    astField.Names[0].String(),
		Json:      getJsonFieldName(astField),
		Sequence:  idx,
		IsArray:   false,
		Docs:      getFieldDoc(astField),
		ParamType: "",
	}

	if field.Name == "PathParams" {
		fmt.Println("stop here, let's check the type")
	}
	setFieldType(field, astField.Type)
	ci.Fields = append(ci.Fields, field)
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
	docs, _ := getComments(astField.Doc, astField.Comment)
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
func getComments(cgList ...*ast.CommentGroup) (doc []string, md *MetaData) {
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
		md.Entity = strings.Trim(text[idx+8:], " ")
		md.Valid = true
		return true
	} else if idx = strings.Index(text, "@Data"); idx > -1 {
		md.Data = strings.Trim(text[idx+5:], " ")
		md.Valid = true
		return true
	} else if idx = strings.Index(text, "@Enum:"); idx > -1 {
		md.Enum = strings.Trim(text[idx+6:], " ")
		md.Valid = true
		return true
	} else if idx = strings.Index(text, "@Message:"); idx > -1 {
		md.Message = strings.Trim(text[idx+9:], " ")
		md.Valid = true
		return true
	} else {
		return false
	}
}

// Process public function
func (p *Parser) processFunc(astObj *ast.Object, pi *PackageInfo) {
	fmt.Println("processFunc", astObj.Name, astObj.Type)
}

// Process public enum
func (p *Parser) processEnumSpec(spec *ast.ValueSpec, pi *PackageInfo, docs []string, md *MetaData) {

	name := md.Enum
	if len(name) == 0 {
		name = spec.Names[0].String()
	}

	// Create enum info
	ei := &EnumInfo{
		Name:    name,
		Docs:    docs,
		Values:  make([]*EnumValueInfo, 0),
		IsFlags: false,
	}

	// process fields
	for _, val := range spec.Values {
		if exp, ok := val.(*ast.UnaryExpr); ok {
			if x, ok := exp.X.(*ast.CompositeLit); ok {

				// build comments map
				cm := getEnumCommentMap(x.Type.(*ast.Ident))

				for _, elt := range x.Elts {
					if evi := p.getEnumValue(elt.(*ast.KeyValueExpr)); evi != nil {
						evi.Docs = append(evi.Docs, cm[evi.Name])
						ei.Values = append(ei.Values, evi)
					}
				}
			}
		}
	}

	pi.Enums = append(pi.Enums, ei)
}

// get comment
func getEnumCommentMap(expr *ast.Ident) map[string]string {
	result := make(map[string]string)
	for _, f := range expr.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List {
		key := f.Names[0].String()
		value := ""
		if f.Tag != nil {
			value = strings.ReplaceAll(f.Tag.Value, "`", "")
		}
		result[key] = value
	}
	return result
}

// extract enum values
func (p *Parser) getEnumValue(kv *ast.KeyValueExpr) *EnumValueInfo {

	evi := &EnumValueInfo{
		Docs: nil,
	}

	if key, ok := kv.Key.(*ast.Ident); ok {
		evi.Name = key.Name
	} else {
		return nil
	}

	if val, ok := kv.Value.(*ast.BasicLit); ok {
		if i, err := strconv.Atoi(val.Value); err == nil {
			evi.Value = i
		}
	} else {
		return nil
	}
	return evi
}
