package parser

// MetaModel represents the info about the domain model
type MetaModel struct {
	Packages    map[string]*PackageInfo
	BaseClasses map[string]*ClassInfo
}

// PackageInfo package metamodel
type PackageInfo struct {
	// Name of the package as documentation
	Name     string           // Package name
	Docs     []string         // Package documentation
	Classes  []*ClassInfo     // List of classes in package
	Enums    []*EnumInfo      // List of enums in package
	Services []*ServiceInfo   // List of services in package
	Sockets  []*WebSocketInfo // List of web sockets
}

func (pi *PackageInfo) AddServiceMethod(mi *MethodInfo) {
	// Find related service by context
	for _, srv := range pi.Services {
		if srv.Context == mi.Context {
			srv.Methods = append(srv.Methods, mi)
			return
		}
	}
}

// ClassInfo  domain model class metamodel
type ClassInfo struct {
	ID          string       // Class full name (including package)
	Name        string       // Name of class
	Package     string       // Package name
	Docs        []string     // Class documentation
	IsExtend    bool         // Is extended
	IsStream    bool         // Is this class represented as stream
	BaseClasses []string     // List of base classes (empty if class is not extended)
	TableName   string       // Name of table (in case it is a persistent entity in database
	Fields      []*FieldInfo // List of class fields
}

// FieldInfo class field metamodel
type FieldInfo struct {
	Name     string   // Field name
	Json     string   // Json name (small capital)
	Type     string   // Field type
	Sequence int      // Field ordinal number
	IsArray  bool     // If type is a list, the type in the list: []Type
	IsMap    bool     // If type is a map
	MapKey   string   // If type is a map, the map key: map[MapKey]Type
	Docs     []string // Field documentation
}

// ReturnInfo method return value metamodel
type ReturnInfo struct {
	Name        string   // Return full type
	Type        string   // Return type
	GenericType string   // Return generic type
	IsArray     bool     // If type is a list, the type in the list: []Type
	IsMap       bool     // If type is a map
	MapKey      string   // If type is a map, the map key: map[MapKey]Type
	Docs        []string // Field documentation
}

// EnumInfo enum metamodel
type EnumInfo struct {
	Name    string           // Enum name
	Docs    []string         // Enum documentation
	Values  []*EnumValueInfo // Enum values
	IsFlags bool             // Is this enum is marked as flags (allow using bitwise operators)
}

// EnumValueInfo enum value metamodel
type EnumValueInfo struct {
	Name  string   // Name of value
	Docs  []string // Documentation
	Value int      // Numeric value
}

// ServiceInfo REST service metamodel
type ServiceInfo struct {
	Name    string        // Name of the service
	Group   string        // Name of the service group
	Path    string        // Service URI path
	Docs    []string      // Documentation
	Methods []*MethodInfo // List of class fields
	Headers []string      // List of Http headers (common to all service methods)
	Context string        // Context (objects)
}

// MethodInfo REST service method metamodel
type MethodInfo struct {
	Name              string       // Name of the service method
	Method            string       // HTTP method: GET | POST | PUT | DELETE
	Path              string       // Method URI path
	Docs              []string     // Documentation
	Headers           []string     // List of Http headers for this method
	PathParams        []*ParamInfo // List of service path parameters
	QueryParams       []*ParamInfo // List of service query parameters
	BodyParam         *ParamInfo   // Body
	FileParam         *ParamInfo   // File param (for upload)
	StreamsRequest    bool         // Flag to indicate the method should return stream
	Return            *ReturnInfo  // Return class info
	Context           string       // Context (objects)
	IsSocketMessage   bool         // Is this method represents socket message
	SocketMessageType string       // Is method is socket message of type Request | Response
}

// ParamInfo REST service method parameter metamodel
type ParamInfo struct {
	Name      string   // Field name
	Json      string   // Json name (small capital)
	Type      string   // Field type
	Sequence  int      // Field ordinal number
	IsArray   bool     // If type is a list, the type in the list: []Type
	IsMap     bool     // If type is a map
	MapKey    string   // If type is a map, the map key: map[MapKey]Type
	Docs      []string // Field documentation
	ParamType string   // How parameter is passed: QueryParam | PathParam | BodyParam
}

// MessageInfo Web Socket Message metamodel
type MessageInfo struct {
	Name      string     // Name of message
	Docs      []string   // Class documentation
	IsRequest bool       // Is request message (client-server) or response(server-client)
	Message   *ClassInfo // List of class fields
}

// WebSocketInfo service endpoint metamodel
type WebSocketInfo struct {
	Name    string        // Name of the socket
	Group   string        // Name of the socket group
	Path    string        // Web socket URI path
	Usage   string        // Web socket Usage sample
	Docs    []string      // Web socket Documentation
	Methods []*MethodInfo // List of class fields
}

// FindClass search for class info by class name in the model
func (m *MetaModel) FindClass(name string) *ClassInfo {
	for _, p := range m.Packages {
		for _, c := range p.Classes {
			if c.ID == name {
				return c
			}
			if c.Name == name {
				return c
			}
		}
	}
	return nil
}

// NewMetaModel initialize metamodel with all base classes in the yaaf-common
func NewMetaModel() *MetaModel {

	model := &MetaModel{
		Packages:    make(map[string]*PackageInfo),
		BaseClasses: make(map[string]*ClassInfo),
	}
	model.BaseClasses["BaseEntity"] = NewBaseEntityModel()
	model.BaseClasses["ActionResponse"] = NewActionResponseModel()
	model.BaseClasses["EntityResponse"] = NewEntityResponseModel()
	model.BaseClasses["EntityResponse"] = NewEntityResponseModel()
	model.BaseClasses["EntitiesResponse"] = NewEntitiesResponseModel()

	return model
}
