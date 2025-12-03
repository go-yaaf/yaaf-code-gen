package model

// region Class Info structure -----------------------------------------------------------------------------------------

// ClassInfo class information
type ClassInfo struct {
	TypeInfo
	IsExtend     bool              // Is extended
	IsGeneric    bool              // Is this is generic class
	GenericTypes []StringKeyValue  // List of generics name to type
	IsVisible    bool              // Is this class is visible for documentation
	IsStream     bool              // Is this class represented as stream
	BaseClass    string            // Base class (empty if class is not extended)
	IsParam      bool              // IS this message is a method input / output param
	Fields       []*FieldInfo      // List of class fields
	Dependencies map[string]string // List of dependencies (class->model)
}

func NewClassInfo(name string, doc ...string) *ClassInfo {
	ci := &ClassInfo{
		TypeInfo: TypeInfo{
			Name:   name,
			TsName: SmallCaps(name),
			Docs:   make([]string, 0),
		},
		IsVisible:    true,
		GenericTypes: make([]StringKeyValue, 0),
		Fields:       make([]*FieldInfo, 0),
		Dependencies: make(map[string]string),
	}

	ci.Docs = append(ci.Docs, doc...)
	return ci
}

// AddField adds field to class
func (ci *ClassInfo) AddField(fName, fType string, doc ...string) {
	fi := NewFieldInfo(fName, doc...)
	fi.Json = SmallCaps(fName)
	fi.Type = fType
	fi.Sequence = len(ci.Fields) + 1
	ci.Fields = append(ci.Fields, fi)
}

// GetField gets field by name
func (ci *ClassInfo) GetField(fName string) *FieldInfo {
	for _, fi := range ci.Fields {
		if fi.Name == fName {
			return fi
		}
	}
	return nil
}

// Fill the dependencies map
func (ci *ClassInfo) fillDependencies(mm *MetaModel) {
	// Add dependencies for complex fields
	for _, fi := range ci.Fields {
		isNative, arr := isNativeType(fi.Type)
		if !isNative && !ci.isGenericFieldType(fi.Type) {
			ci.Dependencies[fi.TsType] = arr
		}
	}

	if len(ci.BaseClass) > 0 {
		ci.Dependencies[ci.BaseClass] = ""
	}
}

// Check if the field type is not part of the generic type list
func (ci *ClassInfo) isGenericFieldType(fieldType string) bool {
	for _, g := range ci.GenericTypes {
		if g.Key == fieldType {
			return true
		}
	}
	return false
}

// endregion
