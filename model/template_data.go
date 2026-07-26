package model

// TemplateData is a unified structure to inject to templates
type TemplateData struct {
	ApiName       string          // API documentation name
	ApiVersion    string          // API version
	EnumGroups    []*EnumGroup    // Enums groups list
	ClassGroups   []*ClassGroup   // Classes groups list
	ServiceGroups []*ServiceGroup // Services groups
}

func NewTemplateData(name string, version string) TemplateData {
	return TemplateData{
		ApiName:       name,
		ApiVersion:    version,
		EnumGroups:    make([]*EnumGroup, 0),
		ClassGroups:   make([]*ClassGroup, 0),
		ServiceGroups: make([]*ServiceGroup, 0),
	}
}
