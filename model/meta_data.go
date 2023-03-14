package model

// MetaData list of metadata flags and values
type MetaData struct {
	Data       string   // Mark structure as data model, the value is the struct name
	Entity     string   // Mark structure as Entity type, value is the table / index name pattern
	Enum       string   // Mark variable as Enum type, the value is the enum name
	Message    string   // Mark structure as Message type, the value is the topic name pattern
	SrvContext string   // REST service context
	SrvPath    string   // REST service path
	SrvGroup   string   // REST service resource group
	SrvHeaders []string // REST service headers
	IsData     bool     // Mark if metadata flags were detected for domain data model
	IsMessage  bool     // Mark if metadata flags were detected for domain message model
	IsService  bool     // Mark if metadata flags were detected for domain service
}

func (md *MetaData) IsValid() bool {
	return md.IsData || md.IsMessage || md.IsService
}

func (md *MetaData) SetData(value string) {
	md.Data = value
	md.IsData = true
}

func (md *MetaData) SetEntity(value string) {
	md.Entity = value
	md.IsData = true
}

func (md *MetaData) SetEnum(value string) {
	md.Enum = value
	md.IsData = true
}

func (md *MetaData) SetMessage(value string) {
	md.Message = value
	md.IsMessage = true
}

func (md *MetaData) SetContext(value string) {
	md.SrvContext = value
	md.IsService = true
}

func (md *MetaData) SetPath(value string) {
	md.SrvPath = value
	md.IsService = true
}

func (md *MetaData) SetGroup(value string) {
	md.SrvGroup = value
	md.IsService = true
}

func (md *MetaData) AddHeader(value string) {
	md.SrvHeaders = append(md.SrvHeaders, value)
	md.IsService = true
}
