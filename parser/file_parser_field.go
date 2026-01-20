package parser

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/go-yaaf/yaaf-code-gen/model"
)

// Process class field
func (p *FileParser) processClassField(field *ast.Field, ci *model.ClassInfo) *model.FieldInfo {

	ignoreField := false

	// Field has no name, it is inherited class
	if len(field.Names) < 1 {
		p.processInheritedClass(field, ci)
		return nil
	}

	fi := &model.FieldInfo{
		Name:     field.Names[0].Name,
		FullName: fmt.Sprintf("%s.%s", ci.Name, field.Names[0].Name),
		TsName:   model.SmallCaps(field.Names[0].Name),
		Json:     model.SmallCaps(field.Names[0].Name),
		IsArray:  false,
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
		} else if strings.HasPrefix(line, "@Type:") {
			fi.TsType = p.getTagValue(line, "@Type:")
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
		if tmplType, ok := fieldType.X.(*ast.Ident); ok {
			p.processFieldTypeIdent(fi, tmplType)
			p.processFieldTypeGenerics(fi, fieldType)
		}
		if tmplType, ok := fieldType.X.(*ast.SelectorExpr); ok {
			p.processFieldTypeIdent(fi, tmplType.Sel)
			p.processFieldTypeGenerics(fi, fieldType)
		}
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
		//fi.Type = fmt.Sprintf("map[%s]%s", keyName, valName)
		//tsKey := model.GetTsType(keyName)
		//tsVal := model.GetTsType(valName)
		//fi.TsType = fmt.Sprintf("Map<%s,%s>", tsKey, tsVal)
		fi.Type = "Json"
		fi.TsType = model.GetTsType(fi.Type)
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
	fi.GenericTypes = append(fi.GenericTypes, model.StringKeyValue{Key: "", Value: idxName})
}

// process generic type in the form ox X[ind1, ind2, ...]
func (p *FileParser) processFieldTypeGenericNew(fi *model.FieldInfo, genType *ast.IndexExpr) {

	xName := ""

	// Extract X name
	switch xType := genType.X.(type) {
	case *ast.Ident:
		xName = xType.Name
	case *ast.SelectorExpr:
		xName = xType.Sel.Name
	default:
		err := fmt.Errorf("processFieldTypeGeneric: the X value of the type %s is not *ast.Ident | *ast.SelectorExpr", fi.Name)
		panic(err)
	}

	// Extract ident value
	var ident *ast.Ident = nil

	switch xType := genType.Index.(type) {
	case *ast.Ident:
		ident = xType
	case *ast.SelectorExpr:
		ident = xType.Sel
	default:
		err := fmt.Errorf("processFieldTypeGeneric: the X value of the type %s is not *ast.Ident | *ast.SelectorExpr", fi.Name)
		panic(err)
	}

	tsName := model.GetTsType(ident.Name)
	switch decl := ident.Obj.Decl.(type) {
	case *ast.Field:
		fi.GenericTypes = append(fi.GenericTypes, model.StringKeyValue{Key: decl.Type.(*ast.Ident).Name, Value: xName})
	case *ast.TypeSpec:
		for _, fl := range decl.TypeParams.List {
			fi.GenericTypes = append(fi.GenericTypes, model.StringKeyValue{Key: fl.Names[0].Name, Value: fl.Type.(*ast.Ident).Name})
		}
	default:
		err := fmt.Errorf("processFieldTypeGeneric: the ident.Obj.Decl value of the type %s is not *ast.Field | *ast.TypeSpec", fi.Name)
		panic(err)
	}
	fi.Type = fmt.Sprintf("%s[%s]", ident.Name, fi.GetGenericsIndicesList())
	fi.TsType = fmt.Sprintf("%s<%s>", tsName, p.getTsGenericsIndicesList(fi))
	fi.IsGeneric = true
}

// Get TS generics list
func (p *FileParser) getTsGenericsIndicesList(f *model.FieldInfo) string {
	list := make([]string, 0)
	for _, svk := range f.GenericTypes {

		tsType := model.GetTsType(svk.Value)
		ind := fmt.Sprintf("%s %s", svk.Key, tsType)
		list = append(list, ind)
	}
	return strings.Join(list, ",")
}

// process generic type in the form ox X[ind1, ind2, ...]
func (p *FileParser) processFieldTypeGenerics(fi *model.FieldInfo, genTypes *ast.IndexListExpr) {

	xName := ""

	// Extract x value
	switch xType := genTypes.X.(type) {
	case *ast.Ident:
		xName = xType.Name
	case *ast.SelectorExpr:
		xName = xType.Sel.Name
	default:
		err := fmt.Errorf("processFieldTypeGeneric: the X value of the type %s is not *ast.Ident | *ast.SelectorExpr", fi.Name)
		panic(err)
	}

	idxNames := make([]string, 0)
	tsIndexes := make([]string, 0)

	// Extract index values
	for _, ind := range genTypes.Indices {
		switch idxType := ind.(type) {
		case *ast.Ident:
			idxNames = append(idxNames, idxType.Name)
			tsIndexes = append(tsIndexes, model.GetTsType(idxType.Name))
			p.addFieldGenericTypes(fi, "", idxType.Name)
		case *ast.SelectorExpr:
			idxNames = append(idxNames, idxType.Sel.Name)
			tsIndexes = append(tsIndexes, model.GetTsType(idxType.Sel.Name))
			p.addFieldGenericTypes(fi, "", idxType.Sel.Name)
		case *ast.IndexListExpr:
			xname := idxType.X.(*ast.Ident).Name
			xList := make([]string, 0)
			tsList := make([]string, 0)
			for _, xind := range idxType.Indices {
				switch xidxType := xind.(type) {
				case *ast.Ident:
					xList = append(xList, xidxType.Name)
					tsList = append(tsList, model.GetTsType(xidxType.Name))
				case *ast.SelectorExpr:
					xList = append(xList, xidxType.Sel.Name)
					tsList = append(tsList, model.GetTsType(xidxType.Sel.Name))
				}
			}

			canon := fmt.Sprintf("%s[%s]", xname, strings.Join(xList, ","))
			tsCanon := fmt.Sprintf("%s<%s>", xname, strings.Join(tsList, ","))

			idxNames = append(idxNames, canon)
			tsIndexes = append(tsIndexes, tsCanon)
			p.addFieldGenericTypes(fi, xname, xname)
		default:
			err := fmt.Errorf("processFieldTypeGeneric: the index value of the type %s is not *ast.Ident | *ast.SelectorExpr", fi.Name)
			panic(err)
		}
	}

	tsName := model.GetTsType(xName)
	idxNameList := strings.Join(idxNames, ", ")
	tsIndexList := strings.Join(tsIndexes, ", ")

	fi.Type = fmt.Sprintf("%s[%s]", xName, idxNameList)
	fi.TsType = fmt.Sprintf("%s<%s>", tsName, tsIndexList)
	fi.IsGeneric = true
}

// Add only complex type to generics list, ignore primitive types
func (p *FileParser) addFieldGenericTypes(fi *model.FieldInfo, name, idxType string) {
	if _, ok := goPrimitiveTypes[idxType]; ok {
		return
	} else {
		fi.GenericTypes = append(fi.GenericTypes, model.StringKeyValue{Key: name, Value: idxType})
	}
}
