package parser

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/go-yaaf/yaaf-code-gen/model"
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
		for i, field := range structType.Fields.List {
			if fi := p.processClassField(i, field, ci); fi != nil {
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
		for i, field := range structType.Fields.List {
			if fi := p.processClassField(i, field, ci); fi != nil {
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

// Process class field
func (p *FileParser) processClassField(idx int, field *ast.Field, ci *model.ClassInfo) *model.FieldInfo {

	ignoreField := false

	// Field has no name, it is inherited class
	if len(field.Names) < 1 {
		p.processInheritedClass(field, ci)
		return nil
	}

	fi := &model.FieldInfo{
		Name:       field.Names[0].Name,
		TsName:     model.SmallCaps(field.Names[0].Name),
		Json:       model.SmallCaps(field.Names[0].Name),
		Sequence:   idx,
		IsArray:    false,
		IsRequired: true,
	}

	// DEBUG point
	if ci.Name == "TimeSeriesConsumption" && fi.Name == "Values" {
		fmt.Println("stop here")
	}

	switch ft := field.Type.(type) {
	case *ast.Ident:
		p.processFieldTypeIdent(fi, ft)
	case *ast.SelectorExpr:
		p.processFieldTypeIdent(fi, ft.Sel)
	case *ast.ArrayType:
		p.processFieldTypeArray(fi, ft)
	case *ast.MapType:
		p.processFieldTypeMap(fi, ft)
	case *ast.IndexExpr:
		p.processFieldTypeGeneric(fi, ft)
	default:
		err := fmt.Errorf("error: [%s.%s]: field type is not of type *ast.Ident or *ast.ArrayType", ci.Name, fi.Name)
		panic(err)
	}

	if field.Tag != nil {
		p.processFieldTag(fi, ci, field.Tag.Value)
	}

	// Process inline comments
	if field.Comment != nil {
		if !p.processFieldComments(fi, ci, field.Comment.List) {
			ignoreField = true
		}
	}

	// Process upper documentation
	if field.Doc != nil {
		if !p.processFieldComments(fi, ci, field.Doc.List) {
			ignoreField = true
		}
	}

	if ignoreField {
		return nil
	} else {
		return fi
	}
}

// If class has field with no name, it is inherited class
func (p *FileParser) processInheritedClass(field *ast.Field, ci *model.ClassInfo) {

	ci.IsExtend = true

	// safe type cast
	switch ft := field.Type.(type) {
	case *ast.Ident:
		ci.BaseClass = ft.Name
	case *ast.SelectorExpr:
		ci.BaseClass = ft.Sel.Name
	default:
		//fmt.Println("error: field type is not of type *ast.ArrayType")
	}
}

// Process field comments and extract tags to enrich class and field metadata. The following tags are expected:
// @InheritFrom: - the field type is the parent class
// @Json: - the json name of the field
func (p *FileParser) processFieldComments(fi *model.FieldInfo, ci *model.ClassInfo, comments []*ast.Comment) bool {

	for _, comment := range comments {
		line := p.trimComment(comment.Text)
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "@InheritFrom") {
			ci.BaseClass = fi.Type
			ci.IsExtend = true
			return false
		} else if strings.HasPrefix(line, "@Json:") {
			fi.Json = p.getTagValue(line, "@Json:")
		} else if strings.HasPrefix(line, "@Alias:") {
			fi.Alias = p.getTagValue(line, "@Alias:")
		} else if strings.HasPrefix(line, "@Format:") {
			fi.Format = p.getTagValue(line, "@Format:")
		} else if strings.HasPrefix(line, "@PathParam") {
			fi.ParamType = "path"
		} else if strings.HasPrefix(line, "@QueryParam") {
			fi.ParamType = "query"
		} else if strings.HasPrefix(line, "@BodyParam") {
			fi.ParamType = "body"
		} else if strings.HasPrefix(line, "@FileParam") {
			fi.ParamType = "file"
		} else {
			fi.Docs = append(fi.Docs, line)
		}
	}

	return true
}

// Process field comments and extract tags to enrich class and field metadata. The following tags are expected:
// @InheritFrom: - the field type is the parent class
// @Json: - the json name of the field
func (p *FileParser) processFieldTag(fi *model.FieldInfo, ci *model.ClassInfo, tag string) {
	items := strings.Split(tag, ":")
	if len(items) != 2 {
		return
	}
	key := strings.ReplaceAll(items[0], "`", "")
	key = strings.ReplaceAll(key, "\"", "")

	val := strings.ReplaceAll(items[1], "`", "")
	val = strings.ReplaceAll(val, "\"", "")
	val = strings.ReplaceAll(val, "omitempty", "")
	val = strings.ReplaceAll(val, ",", "")
	val = strings.TrimSpace(val)

	if key == "json" {
		if val != "-" {
			fi.Json = val
		}
	}
}

// process simple type
func (p *FileParser) processFieldTypeIdent(fi *model.FieldInfo, ft *ast.Ident) {
	fi.Type = ft.Name
	fi.TsType = model.GetTsType(ft.Name)

	// Set default format
	switch fi.Type {
	case "double":
		fi.Format = "decimal"
		break
	case "float":
		fi.Format = "decimal"
		break
	case "float32":
		fi.Format = "decimal"
		break
	case "float64":
		fi.Format = "decimal"
		break
	case "Timestamp":
		fi.Format = "datetime"
		break
	}
	if fi.Type == "number" || fi.Type == "string" || fi.Type == "boolean" {
		fi.IsComplex = false
	} else {
		fi.IsComplex = true
	}
}

// process array type
func (p *FileParser) processFieldTypeArray(fi *model.FieldInfo, arrType *ast.ArrayType) {
	fi.IsArray = true
	switch fieldType := arrType.Elt.(type) {
	case *ast.Ident:
		p.processFieldTypeIdent(fi, fieldType)
	case *ast.SelectorExpr:
		p.processFieldTypeIdent(fi, fieldType.Sel)
	case *ast.IndexExpr:
		if tmplType, ok := fieldType.X.(*ast.Ident); ok {
			p.processFieldTypeIdent(fi, tmplType)
			p.processFieldTypeGeneric(fi, fieldType)
		} else {
			//fmt.Println("processFieldTypeArray: error processing type", fi.Name)
		}
	case *ast.StarExpr:
		if tmplType, ok := fieldType.X.(*ast.Ident); ok {
			p.processFieldTypeIdent(fi, tmplType)
		}
	case *ast.IndexListExpr:
		fi.TsType = "any"
		fi.Type = "any"
		fi.IsArray = true
	default:
		//fmt.Println("processFieldTypeArray error: field type of", fi.Name)
	}
}

// process array type
func (p *FileParser) processFieldTypeMap(fi *model.FieldInfo, mapType *ast.MapType) {
	keyName := ""
	valName := ""

	if keyType, ok := mapType.Key.(*ast.Ident); ok {
		keyName = keyType.Name
	}

	if valType, ok := mapType.Value.(*ast.Ident); ok {
		valName = valType.Name
	} else {
		valName = "any"
	}

	if keyName == "string" && valName == "any" {
		fi.Type = "Json"
		fi.TsType = model.GetTsType(fi.Type)
	} else if keyName == "string" && valName == "intefrace{}" {
		fi.Type = "Json"
		fi.TsType = model.GetTsType(fi.Type)
	} else {
		fi.Type = fmt.Sprintf("map[%s]%s", keyName, valName)
		tsKey := model.GetTsType(keyName)
		tsVal := model.GetTsType(valName)
		fi.TsType = fmt.Sprintf("Map<%s,%s>", tsKey, tsVal)
	}
}

// process generic type in the form ox X[ind1, ind2, ...]
func (p *FileParser) processFieldTypeGeneric(fi *model.FieldInfo, genType *ast.IndexExpr) {

	xName := ""
	idxName := ""

	// Extract x value
	switch xType := genType.X.(type) {
	case *ast.Ident:
		xName = xType.Name
	case *ast.SelectorExpr:
		xName = xType.Sel.Name
	default:
		err := fmt.Errorf("processFieldTypeGeneric: the X value of the type %s is not *ast.Ident | *ast.SelectorExpr", fi.Name)
		panic(err)
	}

	// Extract index value
	switch idxType := genType.Index.(type) {
	case *ast.Ident:
		idxName = idxType.Name
	case *ast.SelectorExpr:
		idxName = idxType.Sel.Name
	default:
		err := fmt.Errorf("processFieldTypeGeneric: the index value of the type %s is not *ast.Ident | *ast.SelectorExpr", fi.Name)
		panic(err)
	}

	tsName := model.GetTsType(xName)
	tsIndex := model.GetTsType(idxName)

	fi.Type = fmt.Sprintf("%s[%s]", xName, idxName)
	fi.TsType = fmt.Sprintf("%s<%s>", tsName, tsIndex)
	fi.IsGeneric = true
	fi.GenericTypes = append(fi.GenericTypes, idxName)
}
