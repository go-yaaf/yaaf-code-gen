package model

// TemplateData is a unified structure to inject to templates
type TemplateData struct {
	ApiName       string          // API documentation name
	ApiVersion    string          // API version
	EnumGroups    []*EnumGroup    // Enums groups list
	ClassGroups   []*ClassGroup   // Classes groups list
	ServiceGroups []*ServiceGroup // Services groups
	EnumList      []*EnumInfo     // Enums list
	ClassList     []*ClassInfo    // Classes list
	ServiceList   []*ServiceInfo  // Services list
	EnumInfo      *EnumInfo       // Enum info
	ClassInfo     *ClassInfo      // Class info
	ServiceInfo   *ServiceInfo    // Service info
}

func NewTemplateData(name string, version string) TemplateData {
	return TemplateData{
		ApiName:       name,
		ApiVersion:    version,
		EnumGroups:    make([]*EnumGroup, 0),
		ClassGroups:   make([]*ClassGroup, 0),
		ServiceGroups: make([]*ServiceGroup, 0),
		EnumList:      make([]*EnumInfo, 0),
		ClassList:     make([]*ClassInfo, 0),
		ServiceList:   make([]*ServiceInfo, 0),
	}
}
