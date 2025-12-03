package model

// region Type Info structure ------------------------------------------------------------------------------------------

// TypeInfo type information
type TypeInfo struct {
	Name             string   // Name of class
	TsName           string   // TypeScript name (small caps)
	PackageFullName  string   // Full name of the package
	PackageShortName string   // Short package name (suffix only)
	Docs             []string // Class documentation
	TableName        string   // Name of table (in case it is a persistent entity in database
	Type             string   // Meta type
	Headers          []string // List of Http headers (common to all service methods)
	Group            string   // Name of the service group
	Context          string   // Context (objects)
	Path             string   // Context (objects)
}

func NewTypeInfo(name string) *TypeInfo {
	return &TypeInfo{
		Name:    name,
		TsName:  SmallCaps(name),
		Docs:    make([]string, 0),
		Headers: make([]string, 0),
	}
}

// AddHeader add new header
func (t *TypeInfo) AddHeader(header string) {
	if len(header) > 0 {
		t.Headers = append(t.Headers, header)
	}
}

// endregion
