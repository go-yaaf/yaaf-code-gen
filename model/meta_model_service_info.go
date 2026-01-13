package model

import (
	"strings"
)

// region Service Info structure ---------------------------------------------------------------------------------------

// ServiceInfo service information
type ServiceInfo struct {
	TypeInfo
	Path         string            // Service URI path
	Methods      []*MethodInfo     // List of class fields
	Dependencies map[string]string // List of dependencies (class->model)
}

func NewServiceInfo(name string, doc ...string) *ServiceInfo {
	si := &ServiceInfo{
		TypeInfo: TypeInfo{
			Name:    name,
			TsName:  SmallCaps(name),
			Docs:    make([]string, 0),
			Headers: make([]string, 0),
		},
		Methods:      make([]*MethodInfo, 0),
		Dependencies: make(map[string]string),
	}
	si.Docs = append(si.Docs, doc...)
	return si
}

// Fill the dependencies map
func (s *ServiceInfo) fillDependencies(mm *MetaModel) {

	// Add dependencies for complex fields
	for _, mi := range s.Methods {
		// Check Path parameters
		for _, pp := range mi.PathParams {
			s.addDependency(pp.Type)
		}

		// Check Query parameters
		for _, qp := range mi.QueryParams {
			s.addDependency(qp.Type)
		}

		// Check Body parameter
		if mi.BodyParam != nil {
			tn := NewTypeNode(mi.BodyParam.Type)
			s.addNodeDependencies(tn)
		}

		// Check Return parameter
		if mi.ReturnType == nil {
			return
		} else {
			s.addNodeDependencies(mi.ReturnType)
		}
	}
}

func (s *ServiceInfo) addNodeDependencies(node *TypeNode) {
	s.addDependency(node.Name)
	for _, arg := range node.Args {
		s.addNodeDependencies(arg)
	}
}

func (s *ServiceInfo) addDependency(name string) {
	if len(name) == 0 {
		return
	}

	if isNative, arr := isNativeType(name); isNative == false {
		s.Dependencies[name] = arr
	}
}

// Replace all aliases
func (s *ServiceInfo) replaceAliases(pi *PackageInfo) {

	// for every method
	for _, mi := range s.Methods {
		// TODO: replace Path parameters
		//for _, pp := range mi.PathParams {
		//	s.addDependency(pp.Type)
		//}

		// TODO: replace Query parameters
		//for _, qp := range mi.QueryParams {
		//	s.addDependency(qp.Type)
		//}

		// TODO: replace Body parameter
		//if mi.BodyParam != nil {
		//	tn := NewTypeNode(mi.BodyParam.Type)
		//	s.addNodeDependencies(tn)
		//}

		// Check Return parameter
		if returnClass, ok := pi.Aliases[mi.ReturnClass]; ok {
			mi.SetReturnType(returnClass)
		}
	}
}

// MethodInfo service method information
type MethodInfo struct {
	Name              string       // Name of the service method
	TsName            string       // Type Script method name (small caps)
	Method            string       // HTTP method: GET | POST | PUT | DELETE | PATCH
	Path              string       // Method URI path
	Docs              []string     // Documentation
	Headers           []string     // List of Http headers for this method
	PathParams        []*ParamInfo // List of service path parameters
	QueryParams       []*ParamInfo // List of service query parameters
	BodyParam         *ParamInfo   // Body
	FileParam         *ParamInfo   // File param (for upload)
	StreamsRequest    bool         // Is stream
	Return            *ClassInfo   // Return class info
	ReturnType        *TypeNode    // Return type node
	ReturnClass       string       // Return class name
	Context           string       // Context (objects)
	IsSocketMessage   bool         // Is this method represents socket message
	IsFileUpload      bool         // Is this method represents file upload handler
	SocketMessageType string       // Is method is socket message of type Request | Response
}

func NewMethodInfo(name string) *MethodInfo {
	return &MethodInfo{
		Name:        name,
		TsName:      SmallCaps(name),
		Docs:        make([]string, 0),
		Headers:     make([]string, 0),
		PathParams:  make([]*ParamInfo, 0),
		QueryParams: make([]*ParamInfo, 0),
	}
}

// SetAction decompose action parameters (http verb + http path)
func (m *MethodInfo) SetAction(action string) {
	items := strings.Split(action, " ")
	if len(items) > 1 {
		m.Method = strings.TrimSpace(items[0])
		m.Path = strings.TrimSpace(items[1])
	} else if len(items) == 1 {
		m.Method = "GET"
		m.Path = strings.TrimSpace(items[0])
	}
}

// AddPathParam decompose path parameters (name | type |  description)
func (m *MethodInfo) AddPathParam(params string) {
	items := strings.Split(params, "|")

	if len(items) == 0 {
		return
	}

	pi := NewParamInfo(strings.TrimSpace(items[0]))
	pi.ParamType = "path"

	if len(items) > 1 {
		pi.Type = strings.TrimSpace(items[1])
	}
	if len(items) > 2 {
		pi.Docs = append(pi.Docs, strings.TrimSpace(items[2]))
	}
	m.PathParams = append(m.PathParams, pi)
}

// AddQueryParam decompose query parameters (name | type |  description)
func (m *MethodInfo) AddQueryParam(params string) {
	items := strings.Split(params, "|")

	if len(items) == 0 {
		return
	}

	pi := NewParamInfo(strings.TrimSpace(items[0]))
	pi.ParamType = "query"

	if len(items) > 1 {
		pt := strings.TrimSpace(items[1])
		if strings.HasPrefix(pt, "[]") {
			pi.Type = pt[2:]
			pi.IsArray = true
		} else {
			pi.Type = pt
		}
	}
	if len(items) > 2 {
		pi.Docs = append(pi.Docs, strings.TrimSpace(items[2]))
	}
	m.QueryParams = append(m.QueryParams, pi)
}

// AddBodyParam decompose body parameter (name | type |  description)
func (m *MethodInfo) AddBodyParam(params string) {
	items := strings.Split(params, "|")

	if len(items) == 0 {
		return
	}

	pi := NewParamInfo(strings.TrimSpace(items[0]))
	pi.ParamType = "body"

	if len(items) > 1 {
		pi.Type = strings.TrimSpace(items[1])
	}
	if len(items) > 2 {
		pi.Docs = append(pi.Docs, strings.TrimSpace(items[2]))
	}
	m.BodyParam = pi
}

// AddFileParam decompose file parameter (name | type |  description)
func (m *MethodInfo) AddFileParam(params string) {
	items := strings.Split(params, "|")

	if len(items) == 0 {
		return
	}

	pi := NewParamInfo(strings.TrimSpace(items[0]))
	pi.ParamType = "file"

	if len(items) > 1 {
		pi.Type = strings.TrimSpace(items[1])
	}
	if len(items) > 2 {
		pi.Docs = append(pi.Docs, strings.TrimSpace(items[2]))
	}
	m.FileParam = pi
}

// SetUploadFunction decompose upload parameter
func (m *MethodInfo) SetUploadFunction(name string) {
	m.Name = name
	m.IsFileUpload = true
}

func (m *MethodInfo) SetReturnType(returnClass string) {
	m.ReturnClass = returnClass
	if returnClass == "StreamContent" {
		m.Return.IsStream = true
	}

	// Set Generics
	m.ReturnType = NewTypeNode(returnClass)
}

// ParamInfo method parameter information
type ParamInfo struct {
	Name      string   // Parameter name
	TsName    string   // TypeScript field name (small caps)
	Json      string   // Json name (small capital)
	Type      string   // Parameter value type
	IsArray   bool     // Is it array
	ParamType string   // How parameter is passed: path | query | body | file
	Docs      []string // Field documentation
}

func NewParamInfo(name string) *ParamInfo {
	return &ParamInfo{
		Name:   name,
		TsName: SmallCaps(name),
		Json:   SmallCaps(name),
		Docs:   make([]string, 0),
	}
}

type TypeNode struct {
	Name string      `json:"name"`
	Args []*TypeNode `json:"args,omitempty"`
}

func NewTypeNode(input string) *TypeNode {
	p := newGenericsParser(input)
	if node, err := p.parseType(); err != nil {
		return nil
	} else {
		return node
	}
}

// endregion
