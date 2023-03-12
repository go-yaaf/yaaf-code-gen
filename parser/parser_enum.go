package parser

import (
	"go/ast"
	"strconv"
	"strings"
)

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
