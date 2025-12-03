package parser

import (
	"fmt"
	"github.com/go-yaaf/yaaf-code-gen/model"
	"go/ast"
	"strconv"
	"strings"
)

// process enum type
func (p *FileParser) processEnumType(ti *model.TypeInfo, decl *ast.GenDecl) error {
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
	ei := model.NewEnumInfo(ti.Name)
	ei.TsName = ti.TsName
	ei.PackageFullName = ti.PackageFullName
	ei.PackageShortName = ti.PackageShortName
	ei.Docs = ti.Docs
	ei.TableName = ti.TableName

	spec, _ := decl.Specs[0].(*ast.TypeSpec)

	if ident, ok := spec.Type.(*ast.Ident); ok {
		ei.Type = ident.Name
	}

	// Add enum to package
	p.Model.AddEnumInfo(ei)
	return nil
}

// process enum values
func (p *FileParser) processEnumValues(ti *model.TypeInfo, decl *ast.GenDecl) error {

	// Get the enum (table name)
	enm := p.Model.GetEnum(ti.TableName)
	if enm == nil {
		return fmt.Errorf("enum %s not found", ti.TableName)
	}
	if len(decl.Specs) < 1 {
		return fmt.Errorf("no specs found")
	}
	spec, ok := decl.Specs[0].(*ast.TypeSpec)
	if !ok {
		return fmt.Errorf("enum spec type %T not supported", spec)
	}
	for _, fld := range spec.Type.(*ast.StructType).Fields.List {
		ev := model.NewEnumValueInfo(fld.Names[0].Name)
		if err := p.processEnumValueComments(ev, fld.Doc, fld.Comment); err != nil {
			continue
		}
		if fld.Tag == nil {
			continue
		}
		if err := p.processEnumValueTags(ev, fld.Tag.Value); err != nil {
			continue
		}
		enm.AddValue(ev)
	}
	return nil
}

// Extract tag value from comment with tag
func (p *FileParser) processEnumValueComments(evi *model.EnumValueInfo, groups ...*ast.CommentGroup) error {
	// collect all comments
	list := make([]*ast.Comment, 0)
	for _, group := range groups {
		if group != nil {
			list = append(list, group.List...)
		}
	}

	// process comments
	for _, comment := range list {
		doc := strings.Replace(comment.Text, "//", "", -1)
		evi.Docs = append(evi.Docs, strings.TrimSpace(doc))
	}
	return nil
}

// Extract tag value from comment with tag
func (p *FileParser) processEnumValueTags(evi *model.EnumValueInfo, tag string) error {
	value := strings.Replace(tag, "`", "", -1)
	value = strings.TrimSpace(value)
	idx := strings.Index(value, "value:")
	if idx == -1 {
		return fmt.Errorf("enum value tags has no value")
	}
	str := strings.Replace(value[idx+6:], "\"", "", -1)
	if num, err := strconv.Atoi(str); err == nil {
		evi.Value = num
		return nil
	} else {
		return fmt.Errorf("enum value tags has no valid value")
	}
}
