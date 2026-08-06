package model

// region Enum Info structure ------------------------------------------------------------------------------------------

// EnumInfo enum information
type EnumInfo struct {
	TypeInfo
	Values  []*EnumValueInfo // Enum values
	IsFlags bool             // Is this enum is marked as flags (allow using bitwise operators)
}

func NewEnumInfo(name string, doc ...string) *EnumInfo {
	ei := &EnumInfo{
		TypeInfo: TypeInfo{
			Name: name,
			Docs: make([]string, 0),
		},
		Values: make([]*EnumValueInfo, 0),
	}
	ei.Docs = append(ei.Docs, doc...)
	return ei
}

// AddValue adds enum value to the enum
func (e *EnumInfo) AddValue(val *EnumValueInfo) {
	e.Values = append(e.Values, val)
}

// Clone creates deep cloned object
func (e *EnumInfo) Clone() *EnumInfo {
	if e == nil {
		return nil
	}

	// 1. Shallow copy all scalar fields at once
	cloned := *e
	for _, val := range e.Values {
		cloned.Values = append(cloned.Values, val.Clone())
	}
	return &cloned
}

// EnumValueInfo enum value information
type EnumValueInfo struct {
	Name  string   // Name of value
	Docs  []string // Documentation
	Value int      // Numeric value
}

func NewEnumValueInfo(name string, doc ...string) *EnumValueInfo {
	evi := &EnumValueInfo{
		Name: name,
		Docs: make([]string, 0),
	}
	evi.Docs = append(evi.Docs, doc...)
	return evi
}

// Clone creates deep cloned object
func (ev *EnumValueInfo) Clone() *EnumValueInfo {
	if ev == nil {
		return nil
	}

	// 1. Shallow copy all scalar fields at once
	cloned := *ev
	return &cloned
}

// endregion
