package processor

import (
	"fmt"
	"github.com/emicklei/proto"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Base processor parses proto files and generates abstract meta Model
type BaseProcessor struct {
	Output   string
	Model    *MetaModel
	ClassMap map[string]*ClassInfo
}

// region General Methods ----------------------------------------------------------------------------------------------

// Start the processor
func (p *BaseProcessor) Parse(root string) error {

	// First, create meta Model
	p.Model = &MetaModel{}
	p.Model.Packages = make(map[string]*PackageInfo)
	p.ClassMap = make(map[string]*ClassInfo)

	// Walk through a folder hierarchy and process proto files types
	err := filepath.Walk(root,
		func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path.Ext(filePath) == ".proto" {
				_ = p.processFileTypes(filePath)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	// Build the class inheritance
	p.buildClassesInheritance()

	// Walk through a folder hierarchy and process proto files services
	err = filepath.Walk(root,
		func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path.Ext(filePath) == ".proto" {
				_ = p.processFileServices(filePath)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return nil
}

// Process single proto file for messages and enums
func (p *BaseProcessor) processFileTypes(filePath string) error {
	reader, _ := os.Open(filePath)
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, _ := parser.Parse()

	proto.Walk(definition, proto.WithMessage(p.handleMessage), proto.WithEnum(p.handleEnum))

	return nil
}

// Process single proto file for services
func (p *BaseProcessor) processFileServices(filePath string) error {
	reader, _ := os.Open(filePath)
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, _ := parser.Parse()

	proto.Walk(definition, proto.WithService(p.handleService))

	return nil
}

// Build classes inheritance
func (p *BaseProcessor) buildClassesInheritance() {
	// Go through the classes
	for _, cInfo := range p.ClassMap {
		p.buildClassInheritance(cInfo, cInfo.BaseClass)
	}
}

// Build class inheritance
func (p *BaseProcessor) buildClassInheritance(ci *ClassInfo, base string) {
	if len(base) == 0 {
		return
	}

	// Add all fields of base class
	bi, ok := p.ClassMap[base]
	if !ok {
		return
	}

	for _, fi := range bi.Fields {
		ci.Fields = append(ci.Fields, fi)
	}

	// Call this class again
	p.buildClassInheritance(ci, bi.BaseClass)
	return
}

// endregion

// region Message Processing -------------------------------------------------------------------------------------------

/**
 * Process message element in the Model
 */
func (p *BaseProcessor) handleMessage(m *proto.Message) {

	// Get package
	pkg := p.getPackage(m.Parent)
	shortName := p.getSimpleName(pkg.Name)

	ci := &ClassInfo{Name: m.Name, IsVisible: true, PackageFullName: pkg.Name, PackageShortName: shortName}

	// Add class documentation
	p.processClassComments(ci, m.Comment)

	// Check if class is inherited
	ci.IsExtend = m.IsExtend

	// Add field elements to class
	for _, elem := range m.Elements {
		field, ok := elem.(*proto.NormalField)
		if ok {
			if fi := p.processClassField(field, ci); fi != nil {
				ci.Fields = append(ci.Fields, fi)
			}
		}
	}

	// Add class to package
	pkg.Classes = append(pkg.Classes, ci)

	// Add class to class map
	p.ClassMap[ci.Name] = ci
}

/**
 * Process class comments
 */
func (p *BaseProcessor) processClassComments(ci *ClassInfo, c *proto.Comment) {
	if c == nil {
		return
	}

	for _, line := range c.Lines {
		p.processClassComment(line, ci)
	}
}

/**
 * Process class comment and extract tags to enrich class meta data. The following tags are expected:
 * @Table: - the table in the database
 */
func (p *BaseProcessor) processClassComment(line string, ci *ClassInfo) {

	line = p.trimComment(line)
	if len(line) == 0 {
		return
	}

	if strings.HasPrefix(line, "@Table:") {
		ci.TableName = p.getTagValue(line, "@Table:")
	} else if strings.HasPrefix(line, "@Request") {
		ci.IsVisible = false
	} else if strings.HasPrefix(line, "@Response") {
		ci.IsVisible = true
	} else {
		ci.Docs = append(ci.Docs, line)
	}
	return
}

/**
 * Process class field
 */
func (p *BaseProcessor) processClassField(field *proto.NormalField, ci *ClassInfo) *FieldInfo {

	ignoreField := false

	fieldType := p.getSimpleName(field.Type)
	fieldJson := p.smallCaps(field.Name)
	fi := &FieldInfo{
		Name:       field.Name,
		TsName:     p.smallCaps(field.Name),
		Type:       fieldType,
		Sequence:   field.Sequence,
		Json:       fieldJson,
		IsArray:    field.Repeated,
		IsRequired: field.Required,
	}

	// Process upper comments
	if field.Comment != nil {
		for _, line := range field.Comment.Lines {
			if !p.processFieldComments(fi, ci, line) {
				ignoreField = true
			}
		}
	}

	// Process inline comments
	if field.InlineComment != nil {
		for _, line := range field.InlineComment.Lines {
			if !p.processFieldComments(fi, ci, line) {
				ignoreField = true
			}
		}
	}

	if ignoreField {
		return nil
	} else {
		return fi
	}
}

/**
 * Process field comments and extract tags to enrich class and field meta data. The following tags are expected:
 * @InheritFrom: - the field type is the parent class
 * @Json: - the json name of the field
 */
func (p *BaseProcessor) processFieldComments(fi *FieldInfo, ci *ClassInfo, line string) bool {

	line = p.trimComment(line)
	if len(line) == 0 {
		return true
	}

	if strings.HasPrefix(line, "@InheritFrom") {
		ci.BaseClass = fi.Type
		ci.IsExtend = true
		return false
	} else if strings.HasPrefix(line, "@Json:") {
		fi.Json = p.getTagValue(line, "@Json:")
	} else if strings.HasPrefix(line, "@Alias:") {
		fi.Alias = p.getTagValue(line, "@Alias:")
	} else if strings.HasPrefix(line, "@PathParam") {
		fi.ParamType = "path"
	} else if strings.HasPrefix(line, "@QueryParam") {
		fi.ParamType = "query"
	} else if strings.HasPrefix(line, "@BodyParam") {
		fi.ParamType = "body"
	} else if strings.HasPrefix(line, "@FileParam") {
		fi.ParamType = "file"
	} else {
		fi.Docs = append(fi.Docs, line)
	}

	return true
}

// endregion

// region Enum Processing -------------------------------------------------------------------------------------------

// Process Enum element in the Model
func (p *BaseProcessor) handleEnum(m *proto.Enum) {

	// Get package
	pkg := p.getPackage(m.Parent)

	ei := &EnumInfo{Name: m.Name, IsFlags: false}

	// Add Enum documentation
	p.addEnumComments(ei, m.Comment)

	// Add Enum values to the enum
	for _, elem := range m.Elements {
		field, ok := elem.(*proto.EnumField)
		if ok {
			ev := p.addEnumValue(field)
			ei.Values = append(ei.Values, ev)
		}
	}

	// Add enum to package
	pkg.Enums = append(pkg.Enums, ei)
}

// Add enum comments
func (p *BaseProcessor) addEnumComments(ei *EnumInfo, c *proto.Comment) {
	if c == nil {
		return
	}

	for _, line := range c.Lines {
		trimmed := p.trimComment(line)
		if len(trimmed) > 0 {
			if strings.HasPrefix(trimmed, "@Flags") {
				ei.IsFlags = true
			} else {
				ei.Docs = append(ei.Docs, trimmed)
			}
		}
	}
}

// Add enum value
func (p *BaseProcessor) addEnumValue(field *proto.EnumField) *EnumValueInfo {
	evi := &EnumValueInfo{Name: field.Name, Value: field.Integer}
	// Add comments and inline comments
	p.addEnumValueComments(evi, field.Comment)
	p.addEnumValueComments(evi, field.InlineComment)

	return evi
}

// Add field comments
func (p *BaseProcessor) addEnumValueComments(e *EnumValueInfo, c *proto.Comment) {
	if c == nil {
		return
	}

	for _, line := range c.Lines {
		trimmed := p.trimComment(line)
		if len(trimmed) > 0 {
			e.Docs = append(e.Docs, trimmed)
		}
	}
}

// endregion

// region Service Processing -------------------------------------------------------------------------------------------

// Process service element in the Model
func (p *BaseProcessor) handleService(s *proto.Service) {

	// Check if the service is marked as @WebSocket
	if p.isWebSocketService(s) {
		p.handleWebSocketService(s)
		return
	}

	// Get package
	pkg := p.getPackage(s.Parent)

	si := &ServiceInfo{
		Name:   s.Name,
		TsName: p.smallCaps(s.Name),
	}

	// Add service documentation
	p.processServiceComments(si, s.Comment)

	// Add methods to service
	for _, elem := range s.Elements {
		rpc, ok := elem.(*proto.RPC)
		if ok {

			if mi := p.processServiceMethod(rpc, si); mi != nil {
				si.Methods = append(si.Methods, mi)
			}
		}
	}

	// Add service to package
	pkg.Services = append(pkg.Services, si)
}

/**
 * Process service comments
 */
func (p *BaseProcessor) processServiceComments(si *ServiceInfo, c *proto.Comment) {
	if c == nil {
		return
	}

	for _, line := range c.Lines {
		p.processServiceComment(line, si)
	}
}

/**
 * Process service comment and extract tags to enrich class meta data. The following tags are expected:
 * @Path: - the service base path
 * @RequestHeader: X-API-KEY The key to identify the application (portal)
 * @RequestHeader: X-ACCESS-TOKEN The token to identify the logged-in user
 * @ResourceGroup: User Actions
 */
func (p *BaseProcessor) processServiceComment(line string, si *ServiceInfo) {

	line = p.trimComment(line)
	if len(line) == 0 {
		return
	}

	if strings.HasPrefix(line, "@Path:") {
		si.Path = p.getTagValue(line, "@Path:")
	} else if strings.HasPrefix(line, "@RequestHeader:") {
		if header := p.getTagValue(line, "@RequestHeader:"); len(header) > 0 {
			si.Headers = append(si.Headers, header)
		}
	} else if strings.HasPrefix(line, "@ResourceGroup:") {
		si.Group = p.getTagValue(line, "@ResourceGroup:")
	} else if strings.HasPrefix(line, "@Context:") {
		si.Context = p.getTagValue(line, "@Context:")
	} else {
		si.Docs = append(si.Docs, line)
	}
	return
}

/**
 * Process rpc service method
 */
func (p *BaseProcessor) processServiceMethod(rpc *proto.RPC, si *ServiceInfo) *MethodInfo {
	mi := &MethodInfo{
		Name:           rpc.Name,
		Context:        si.Context,
		StreamsRequest: rpc.StreamsRequest,
	}

	// Process upper comments
	if rpc.Comment != nil {
		for _, line := range rpc.Comment.Lines {
			p.processMethodComments(mi, si, line)
		}
	}

	// Process inline comments
	if rpc.InlineComment != nil {
		for _, line := range rpc.InlineComment.Lines {
			p.processMethodComments(mi, si, line)
		}
	}

	// Process input parameters
	requestType := p.getSimpleName(rpc.RequestType)
	if ci, ok := p.ClassMap[requestType]; ok {
		for _, fi := range ci.Fields {
			if pi := p.buildInputParams(fi); pi != nil {
				if pi.ParamType == "path" {
					mi.PathParams = append(mi.PathParams, pi)
				} else if pi.ParamType == "query" {
					mi.QueryParams = append(mi.QueryParams, pi)
				} else if pi.ParamType == "body" {
					mi.BodyParam = pi
				} else if pi.ParamType == "file" {
					mi.FileParam = pi
				}
			}
		}
	}

	// Process output parameters
	returnType := p.getSimpleName(rpc.ReturnsType)
	if ci, ok := p.ClassMap[returnType]; ok {
		ci.IsStream = rpc.StreamsReturns
		mi.Return = ci
	}

	return mi
}

/**
 * Process method comments and extract tags to enrich service and method meta data. The following tags are expected:
 * @Http: - the HTTP method: GET | POST | PUT | DELETE
 * @Path: - the method relative path
 */
func (p *BaseProcessor) processMethodComments(mi *MethodInfo, si *ServiceInfo, line string) {

	line = p.trimComment(line)
	if len(line) == 0 {
		return
	}

	if strings.HasPrefix(line, "@Http:") {
		mi.Method = p.getTagValue(line, "@Http:")
	} else if strings.HasPrefix(line, "@Path:") {
		mi.Path = p.getTagValue(line, "@Path:")
	} else if strings.HasPrefix(line, "@Context:") {
		mi.Context = p.getTagValue(line, "@Context:")
	} else {
		mi.Docs = append(mi.Docs, line)
	}

	return
}

/**
 * Build input parameters from the class info
 */
func (p *BaseProcessor) buildInputParams(fi *FieldInfo) *ParamInfo {

	pi := &ParamInfo{
		Name:      fi.Name,
		TsName:    fi.TsName,
		Json:      fi.Json,
		Type:      fi.Type,
		IsArray:   fi.IsArray,
		Docs:      fi.Docs,
		ParamType: fi.ParamType,
	}
	return pi
}

// endregion

// region Web Socket Processing ----------------------------------------------------------------------------------------

/**
 * Analyze service comments to determine if this service describes web socket
 */
func (p *BaseProcessor) isWebSocketService(s *proto.Service) bool {

	for _, line := range s.Doc().Lines {
		line = p.trimComment(line)
		if strings.HasPrefix(line, "@WebSocket:") {
			return true
		}
	}
	return false
}

/**
 * Process service element in the Model
 */
func (p *BaseProcessor) handleWebSocketService(s *proto.Service) {
	// Get package
	pkg := p.getPackage(s.Parent)

	si := &WebSocketInfo{
		Name:   s.Name,
		TsName: p.smallCaps(s.Name),
	}

	// Add service documentation
	p.processWebSocketComments(si, s.Comment)

	// Add methods to service
	for _, elem := range s.Elements {
		rpc, ok := elem.(*proto.RPC)
		if ok {

			if mi := p.processWebSocketMethod(rpc, si); mi != nil {
				si.Methods = append(si.Methods, mi)
			}
		}
	}

	// Add socket service to package
	pkg.Sockets = append(pkg.Sockets, si)
}

/**
 * Process web socket service comments
 */
func (p *BaseProcessor) processWebSocketComments(si *WebSocketInfo, c *proto.Comment) {
	if c == nil {
		return
	}

	for _, line := range c.Lines {
		p.processWebSocketComment(line, si)
	}
}

/**
 * Process web socket service comment and extract tags to enrich class meta data. The following tags are expected:
 * @Path: - the service base path
 * @RequestHeader: X-API-KEY The key to identify the application (portal)
 * @RequestHeader: X-ACCESS-TOKEN The token to identify the logged-in user
 * @ResourceGroup: User Actions
 */
func (p *BaseProcessor) processWebSocketComment(line string, si *WebSocketInfo) {

	line = p.trimComment(line)
	if len(line) == 0 {
		return
	}

	if strings.HasPrefix(line, "@Path:") {
		si.Path = p.getTagValue(line, "@Path:")
	} else if strings.HasPrefix(line, "@ResourceGroup:") {
		si.Group = p.getTagValue(line, "@ResourceGroup:")
	} else if strings.HasPrefix(line, "@Usage:") {
		si.Usage = p.getTagValue(line, "@Usage:")
	} else if strings.HasPrefix(line, "@WebSocket:") {
		si.Name = p.getTagValue(line, "@WebSocket:")
	} else {
		si.Docs = append(si.Docs, line)
	}
	return
}

/**
 * Process rpc service method
 */
func (p *BaseProcessor) processWebSocketMethod(rpc *proto.RPC, si *WebSocketInfo) *MethodInfo {
	mi := &MethodInfo{
		Name: rpc.Name,
	}

	// Process upper comments
	if rpc.Comment != nil {
		for _, line := range rpc.Comment.Lines {
			p.processWebSocketMethodComments(mi, si, line)
		}
	}

	// Process inline comments
	if rpc.InlineComment != nil {
		for _, line := range rpc.InlineComment.Lines {
			p.processWebSocketMethodComments(mi, si, line)
		}
	}

	// Process input parameters
	requestType := p.getSimpleName(rpc.RequestType)
	if ci, ok := p.ClassMap[requestType]; ok {
		for _, fi := range ci.Fields {
			if pi := p.buildInputParams(fi); pi != nil {
				if pi.ParamType == "path" {
					mi.PathParams = append(mi.PathParams, pi)
				} else if pi.ParamType == "query" {
					mi.QueryParams = append(mi.QueryParams, pi)
				} else if pi.ParamType == "body" {
					mi.BodyParam = pi
				}
			}
		}
	}

	// Process output parameters
	returnType := p.getSimpleName(rpc.ReturnsType)
	if ci, ok := p.ClassMap[returnType]; ok {
		mi.Return = ci
	}
	return mi
}

/**
 * Process method comments and extract tags to enrich service and method meta data. The following tags are expected:
 * @Http: - the HTTP method: GET | POST | PUT | DELETE
 * @Path: - the method relative path
 */
func (p *BaseProcessor) processWebSocketMethodComments(mi *MethodInfo, si *WebSocketInfo, line string) {

	line = p.trimComment(line)
	if len(line) == 0 {
		return
	}

	if strings.HasPrefix(line, "@Http:") {
		mi.Method = p.getTagValue(line, "@Http:")
	} else if strings.HasPrefix(line, "@Path:") {
		mi.Path = p.getTagValue(line, "@Path:")
	} else if strings.HasPrefix(line, "@SocketMessage:") {
		mi.SocketMessageType = p.getTagValue(line, "@SocketMessage:")
		mi.IsSocketMessage = true
	} else {
		mi.Docs = append(mi.Docs, line)
	}

	return
}

// endregion

// region Internal helpers for proto processing ------------------------------------------------------------------------

// Get package
func (p *BaseProcessor) getPackage(v proto.Visitee) *PackageInfo {
	// Get package name
	pkgName := p.getPackageName(v)
	if pkg, ok := p.Model.Packages[pkgName]; !ok {
		pkg = &PackageInfo{Name: pkgName}
		p.Model.Packages[pkgName] = pkg
		return pkg
	} else {
		return pkg
	}
}

// Extract package name
func (p *BaseProcessor) getPackageName(v proto.Visitee) string {

	root, ok := v.(*proto.Proto)
	if ok {
		for _, elem := range root.Elements {
			pkg, ok := elem.(*proto.Package)
			if ok {
				return pkg.Name
			}
		}
	}
	return "default"
}

// Trim comments
func (p *BaseProcessor) trimComment(line string) string {
	trimmed := strings.TrimSpace(line)

	if strings.HasPrefix(trimmed, "/*") {
		trimmed = strings.Replace(trimmed, "/*", "", 1)
		trimmed = strings.TrimSpace(trimmed)
	}

	if strings.HasPrefix(trimmed, "*") {
		trimmed = strings.Replace(trimmed, "*", "", 1)
		trimmed = strings.TrimSpace(trimmed)
	}

	if strings.HasPrefix(trimmed, "//") {
		trimmed = strings.Replace(trimmed, "//", "", 1)
		trimmed = strings.TrimSpace(trimmed)
	}
	return trimmed
}

// Get simple name (not canonical name)
func (p *BaseProcessor) getSimpleName(name string) string {

	idx := strings.LastIndex(name, ".")
	if idx > 0 {
		return name[idx+1:]
	} else {
		return name
	}
}

// Convert big Caps to small Caps
func (p *BaseProcessor) smallCaps(name string) string {
	return fmt.Sprintf("%s%s", strings.ToLower(name[0:1]), name[1:])
}

// Extract tag value from comment with tag
func (p *BaseProcessor) getTagValue(line string, tag string) string {
	value := strings.Replace(line, tag, "", 1)
	value = strings.TrimSpace(value)
	return value
}

// endregion
