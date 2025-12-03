package model

// region Field Info structure -----------------------------------------------------------------------------------------

// FieldInfo field information
type FieldInfo struct {
	Name       string   // Field name
	TsName     string   // TypeScript field name (small caps)
	Json       string   // Json name (small capital)
	Type       string   // Field original type
	TsType     string   // Field typescript type
	Alias      string   // Type alias
	Format     string   // Display format hint
	Sequence   int      // Field ordinal number
	IsArray    bool     // Is it array
	IsComplex  bool     // Is complex type (NOT number | string | boolean)
	IsRequired bool     // Is field required
	Docs       []string // Field documentation
	ParamType  string   // How parameter is passed: Query | Path | Body
}

func NewFieldInfo(name string, doc ...string) *FieldInfo {
	fi := &FieldInfo{
		Name:      name,
		TsName:    SmallCaps(name),
		IsComplex: true,
	}
	fi.Docs = append(fi.Docs, doc...)
	return fi
}

// endregion
