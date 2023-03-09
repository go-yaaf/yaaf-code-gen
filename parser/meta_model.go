package parser

// MetaData list of metadata flags and values
type MetaData struct {
	Data    string // Mark structure as data model, the value is the struct name
	Entity  string // Mark structure as Entity type, value is the table / index name pattern
	Enum    string // Mark variable as Enum type, the value is the enum name
	Message string // Mark structure as Message type, the value is the topic name pattern
	Valid   bool   // Mark if metadata flags were detected
}

func (md *MetaData) IsEmpty() bool {
	return (len(md.Data) + len(md.Entity) + len(md.Enum) + len(md.Message)) == 0
}

// MetaModel represents the info about the domain model
type MetaModel struct {
	Packages map[string]*PackageInfo
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
	Name      string   // Field name
	TsName    string   // TypeScript field name (small caps)
	Json      string   // Json name (small capital)
	Type      string   // Field type
	Sequence  int      // Field ordinal number
	IsArray   bool     // If type is a list, the type in the list: []Type
	IsMap     bool     // If type is a map
	MapKey    string   // If type is a map, the map key: map[MapKey]Type
	Docs      []string // Field documentation
	ParamType string   // How parameter is passed: Query | Path | Body
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
	TsName  string        // TypeScript service name (small caps)
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
	TsName            string       // TypeScript method name (small caps)
	Method            string       // HTTP method: GET | POST | PUT | DELETE
	Path              string       // Method URI path
	Docs              []string     // Documentation
	Headers           []string     // List of Http headers for this method
	PathParams        []*ParamInfo // List of service path parameters
	QueryParams       []*ParamInfo // List of service query parameters
	BodyParam         *ParamInfo   // Body
	FileParam         *ParamInfo   // File param (for upload)
	StreamsRequest    bool
	Return            *ClassInfo // Return class info
	Context           string     // Context (objects)
	IsSocketMessage   bool       // Is this method represents socket message
	SocketMessageType string     // Is method is socket message of type Request | Response
}

// ParamInfo REST service method parameter metamodel
type ParamInfo struct {
	Name      string   // Parameter name
	TsName    string   // TypeScript field name (small caps)
	Json      string   // Json name (small capital)
	Type      string   // Parameter value type
	IsArray   bool     // Is it array
	ParamType string   // How parameter is passed: Query | Path | Body
	Docs      []string // Field documentation
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
	TsName  string        // TypeScript socket name (small caps)
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
