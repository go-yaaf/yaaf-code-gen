package parser

import (
	"fmt"
	"github.com/go-yaaf/yaaf-code-gen/model"
	"go/ast"
)

// process entity type
func (p *FileParser) processEntityType(ti *model.TypeInfo, decl *ast.GenDecl) error {
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
	ci := model.NewClassInfo(ti.Name)
	ci.TsName = ti.TsName
	ci.PackageFullName = ti.PackageFullName
	ci.PackageShortName = ti.PackageShortName
	ci.Docs = ti.Docs
	ci.TableName = ti.TableName

	spec, _ := decl.Specs[0].(*ast.TypeSpec)
	if spec.TypeParams != nil {
		p.processClassGenericsParams(ci, spec.TypeParams)
	}

	// Add field elements to class
	if structType, ok := spec.Type.(*ast.StructType); ok {
		for _, field := range structType.Fields.List {
			if fi := p.processClassField(field, ci); fi != nil {
				ci.Fields = append(ci.Fields, fi)
			}
		}
	} else {
		//fmt.Println("error: spec.Type is not of type *ast.StructType")
	}

	// Add class to model
	p.Model.AddClassInfo(ci)
	return nil
}

// process data type
func (p *FileParser) processDataType(ti *model.TypeInfo, decl *ast.GenDecl) error {
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
	ci := model.NewClassInfo(ti.Name)
	ci.TsName = ti.TsName
	ci.PackageFullName = ti.PackageFullName
	ci.PackageShortName = ti.PackageShortName
	ci.Docs = ti.Docs
	ci.TableName = ti.TableName

	spec, _ := decl.Specs[0].(*ast.TypeSpec)
	if spec.TypeParams != nil {
		p.processClassGenericsParams(ci, spec.TypeParams)
	}

	// Add field elements to class
	if structType, ok := spec.Type.(*ast.StructType); ok {
		for _, field := range structType.Fields.List {
			if fi := p.processClassField(field, ci); fi != nil {
				ci.Fields = append(ci.Fields, fi)
			}
		}
	} else {
		//fmt.Println("processDataType: spec.Type is not of type *ast.StructType")
	}

	// Add class to model
	p.Model.AddClassInfo(ci)
	return nil
}

// Process class generics type arguments
func (p *FileParser) processClassGenericsParams(ci *model.ClassInfo, params *ast.FieldList) {

	if params.List == nil {
		return
	}
	ci.IsGeneric = true
	for _, field := range params.List {
		for _, fn := range field.Names {
			key := fn.Name

			if val, ok := field.Type.(*ast.Ident); ok {
				ci.GenericTypes = append(ci.GenericTypes, model.StringKeyValue{Key: key, Value: val.Name})
			} else {
				//fmt.Println("processClassGenericsParams: field type is not of type *ast.Ident")
			}
		}
	}
}
