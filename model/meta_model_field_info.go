package model

import (
	"fmt"
	"strings"
)

// region Field Info structure -----------------------------------------------------------------------------------------

// FieldInfo field information
type FieldInfo struct {
	Name         string           // Field name
	FullName     string           // Field full canonical name (including class)
	TsName       string           // TypeScript field name (small caps)
	Json         string           // Json name (small capital)
	Type         string           // Field original type
	TsType       string           // Field typescript type
	Alias        string           // Type alias
	Format       string           // Display format hint
	IsArray      bool             // Is it array
	IsMap        bool             // Is it map field
	IsComplex    bool             // Is complex type (NOT number | string | boolean)
	IsGeneric    bool             // Is this is generic type
	GenericTypes []StringKeyValue // List of generics name to type
	Docs         []string         // Field documentation
	ParamType    string           // How parameter is passed: Query | Path | Body
}

func NewFieldInfo(name string, doc ...string) *FieldInfo {
	fi := &FieldInfo{
		Name:         name,
		TsName:       SmallCaps(name),
		IsComplex:    true,
		IsGeneric:    false,
		GenericTypes: make([]StringKeyValue, 0),
	}
	fi.Docs = append(fi.Docs, doc...)
	return fi
}

// GetGenericsIndicesList returns the generics indices list
func (f *FieldInfo) GetGenericsIndicesList() string {
	list := make([]string, 0)
	for _, svk := range f.GenericTypes {
		ind := fmt.Sprintf("%s %s", svk.Key, svk.Value)
		list = append(list, ind)
	}
	return strings.Join(list, ",")
}

// endregion
