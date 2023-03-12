package parser

import (
	"fmt"
	"go/ast"
)

// Process structure type
func (p *Parser) processTypeSpec(spec *ast.TypeSpec, pi *PackageInfo, docs []string, md *MetaData) {

	if md.IsData {
		p.processClassSpec(spec, pi, docs, md)
	}

	if md.IsService {
		p.processServiceSpec(spec, pi, docs, md)
	}
}

// Process structure type for class
func (p *Parser) processClassSpec(spec *ast.TypeSpec, pi *PackageInfo, docs []string, md *MetaData) {

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
		Name:     astField.Names[0].String(),
		Json:     getJsonFieldName(astField),
		Sequence: idx,
		IsArray:  false,
		Docs:     getFieldDoc(astField),
	}

	if field.Name == "PathParams" {
		fmt.Println("stop here, let's check the type")
	}
	setFieldType(field, astField.Type)
	ci.Fields = append(ci.Fields, field)
}

// Process structure type for REST service
func (p *Parser) processServiceSpec(spec *ast.TypeSpec, pi *PackageInfo, docs []string, md *MetaData) {

	// Create service info
	si := &ServiceInfo{
		Name:    spec.Name.String(),
		Group:   md.SrvGroup,
		Path:    md.SrvPath,
		Docs:    docs,
		Methods: nil,
		Headers: md.SrvHeaders,
		Context: md.SrvContext,
	}
	pi.Services = append(pi.Services, si)
}
