package model

import (
	"fmt"
	"path"
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

	// remove array mark
	name = strings.Replace(name, "[]", "", -1)

	if isNative, arr := isNativeType(name); isNative == false {
		s.Dependencies[name] = arr
	}
}

// Replace all aliases
func (s *ServiceInfo) replaceAliases(pi *PackageInfo) {
	// for every method
	for _, mi := range s.Methods {
		// Replace Path parameters
		//for _, pp := range mi.PathParams {
		//	pp.Type = replaceClassNode(pp.Type, pi)
		//}

		// Replace Query parameters
		//for _, qp := range mi.QueryParams {
		//	qp.Type = replaceClassNode(qp.Type, pi)
		//}

		// Replace Body parameter
		//if mi.BodyParam != nil {
		//	mi.BodyParam.Type = replaceClassNode(mi.BodyParam.Type, pi)
		//}

		// Replace Return parameter
		replaceTypeNode(mi.ReturnType, pi)
		replaceClassNode(mi.ReturnClass, pi)
	}
}

func (s *ServiceInfo) NewMethodInfo(name string) *MethodInfo {
	return &MethodInfo{
		Name:        name,
		TsName:      SmallCaps(name),
		Docs:        make([]string, 0),
		Headers:     make([]string, 0),
		PathParams:  make([]*ParamInfo, 0),
		QueryParams: make([]*ParamInfo, 0),
		ServicePath: s.Path,
	}
}

// Clone create deep clone of this object
func (s *ServiceInfo) Clone() *ServiceInfo {
	if s == nil {
		return nil
	}

	// 1. Shallow copy all scalar fields at once
	cloned := *s

	for _, mi := range s.Methods {
		cloned.Methods = append(cloned.Methods, mi.Clone())
	}
	return &cloned
}

// replaceClassNode will replace all aliases in a generic class string like EntityResponse<StringIntValue<int>>
func replaceClassNode(class string, pi *PackageInfo) string {
	// Split the generic class string into parts
	parts := strings.FieldsFunc(class, func(r rune) bool {
		return r == '<' || r == '>' || r == ','
	})

	// Replace each part with its alias if it exists
	for i, part := range parts {
		if alias, ok := pi.Aliases[strings.TrimSpace(part)]; ok {
			parts[i] = alias
		}
	}

	// Rejoin the parts back into a generic class string
	result := parts[0]
	if len(parts) > 1 {
		result += "<"
		result += strings.Join(parts[1:], ", ")
		result += ">"
	}
	return result
}

func replaceTypeNode(node *TypeNode, pi *PackageInfo) {
	if node == nil {
		return
	}
	for _, arg := range node.Args {
		replaceTypeNode(arg, pi)
	}
	if returnClass, ok := pi.Aliases[node.Name]; ok {
		node.Name = returnClass
	}
}

// MethodInfo service method information
type MethodInfo struct {
	Name              string       // Name of the service method
	TsName            string       // Type Script method name (small caps)
	Method            string       // HTTP method: GET | POST | PUT | DELETE | PATCH
	Path              string       // Method URI path
	ServicePath       string       // Parent service path
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

// FullPath returns the path of the method combined with the service path
func (m *MethodInfo) FullPath() string {
	return path.Join(m.ServicePath, m.Path)
}

// ReturnLink creates link to a complex object
func (m *MethodInfo) ReturnLink() string {
	if m.ReturnType == nil {
		return fmt.Sprintf("#")
	}

	if len(m.ReturnType.Args) == 0 {
		return fmt.Sprintf("type_%s.html", m.ReturnType.Name)
	} else {
		return fmt.Sprintf("type_%s.html", m.ReturnType.Args[0].Name)
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

// GetTsReturnType build the TypeScript representation of the return type
func (m *MethodInfo) GetTsReturnType() string {
	if m.ReturnType == nil {
		return "void"
	}
	return buildTsType(m.ReturnType)
}

// Clone creates deep cloned object
func (m *MethodInfo) Clone() *MethodInfo {
	if m == nil {
		return nil
	}

	// 1. Shallow copy all scalar fields at once
	cloned := *m

	for _, pi := range m.PathParams {
		cloned.PathParams = append(cloned.PathParams, pi.Clone())
	}
	for _, pi := range m.QueryParams {
		cloned.QueryParams = append(cloned.QueryParams, pi.Clone())
	}
	cloned.BodyParam = m.BodyParam.Clone()
	cloned.FileParam = m.FileParam.Clone()
	cloned.Return = m.Return.Clone()
	cloned.ReturnType = m.ReturnType.Clone()
	return &cloned
}

// SampleRequest create sample HTTP request
func (m *MethodInfo) SampleRequest() string {
	return CreateMethodRequestSample(GetMetaModel(), m)
}

// SampleResponse create sample HTTP response
func (m *MethodInfo) SampleResponse() string {
	return CreateMethodResponseSample(GetMetaModel(), m)
}

func buildTsType(node *TypeNode) string {
	if node == nil {
		return ""
	}

	var builder strings.Builder
	builder.WriteString(GetTsType(node.Name))

	if len(node.Args) > 0 {
		builder.WriteString("<")
		for i, arg := range node.Args {
			builder.WriteString(buildTsType(arg))
			if i < len(node.Args)-1 {
				builder.WriteString(", ")
			}
		}
		builder.WriteString(">")
	}

	if node.IsArray {
		builder.WriteString("[]")
	}
	return builder.String()
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

// Clone creates deep cloned object
func (p *ParamInfo) Clone() *ParamInfo {
	if p == nil {
		return nil
	}

	// Shallow copy all scalar fields at once
	cloned := *p
	return &cloned
}

// Link creates link to a complex object
func (p *ParamInfo) Link() string {
	if ci := GetMetaModel().FindClass(p.Type); ci != nil {
		return fmt.Sprintf("type_%s.html", ci.Name)
	} else if ei := GetMetaModel().FindEnum(p.Type); ei != nil {
		return fmt.Sprintf("type_%s.html", ei.Name)
	} else {
		return fmt.Sprintf("#%s", p.Type)
	}
}

type TypeNode struct {
	Name    string      `json:"name"`
	IsArray bool        `json:"isArray,omitempty"`
	Args    []*TypeNode `json:"args,omitempty"`
}

func NewTypeNode(input string) *TypeNode {
	p := newGenericsParser(input)
	if node, err := p.parseType(); err != nil {
		return nil
	} else {
		return node
	}
}

// Clone creates deep cloned object
func (t *TypeNode) Clone() *TypeNode {
	if t == nil {
		return nil
	}

	// Shallow copy all scalar fields at once
	cloned := *t

	for _, arg := range t.Args {
		cloned.Args = append(cloned.Args, arg.Clone())
	}
	return &cloned
}

// endregion

// region Service Group structure --------------------------------------------------------------------------------------

// ServiceGroup service information
type ServiceGroup struct {
	Name string         // Name of resource group
	List []*ServiceInfo // List of services
}

func NewServiceGroup(name string, services ...*ServiceInfo) *ServiceGroup {
	sg := &ServiceGroup{
		Name: name,
		List: make([]*ServiceInfo, 0),
	}
	sg.List = append(sg.List, services...)
	return sg
}

func (sg *ServiceGroup) AddService(info *ServiceInfo) {
	sg.List = append(sg.List, info)
}

// endregion

// region Type Group structure -----------------------------------------------------------------------------------------

// ClassGroup service information
type ClassGroup struct {
	Name string       // Name of resource group
	List []*ClassInfo // List of types
}

func NewClassGroup(name string, types ...*ClassInfo) *ClassGroup {
	sg := &ClassGroup{
		Name: name,
		List: make([]*ClassInfo, 0),
	}
	sg.List = append(sg.List, types...)
	return sg
}

func (sg *ClassGroup) AddClass(info *ClassInfo) {
	sg.List = append(sg.List, info)
}

// endregion

// region Enum Group structure -----------------------------------------------------------------------------------------

// EnumGroup groups of related enums
type EnumGroup struct {
	Name string      // Name of resource group
	List []*EnumInfo // List of types
}

func NewEnumsGroup(name string, enums ...*EnumInfo) *EnumGroup {
	sg := &EnumGroup{
		Name: name,
		List: make([]*EnumInfo, 0),
	}
	sg.List = append(sg.List, enums...)
	return sg
}

func (sg *EnumGroup) AddEnum(info *EnumInfo) {
	sg.List = append(sg.List, info)
}

// endregion
