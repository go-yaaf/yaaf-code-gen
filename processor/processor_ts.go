package processor

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/go-yaaf/yaaf-code-gen/model"
)

var tsTypes = map[string]string{
	"double":    "number",
	"float":     "number",
	"float32":   "number",
	"float64":   "number",
	"int":       "number",
	"int32":     "number",
	"int64":     "number",
	"uint32":    "number",
	"uint64":    "number",
	"sint32":    "number",
	"sint64":    "number",
	"fixed32":   "number",
	"fixed64":   "number",
	"sfixed32":  "number",
	"sfixed64":  "number",
	"bool":      "boolean",
	"string":    "string",
	"bytes":     "File",
	"any":       "any",
	"Timestamp": "number",
	"Json":      "Map<string,object>",
}

// TsProcessor - TS processor converts proto files to TypeScript files
type TsProcessor struct {
	BaseProcessor
}

// NewTsProcessor - Factory method
func NewTsProcessor(model *model.MetaModel, output string) Processor {
	return &TsProcessor{BaseProcessor{
		Output: output,
		Model:  model,
	}}
}

// var classPackageMap = make(map[string]string)

// Start the processor
func (p *TsProcessor) Start() error {

	// Build package map
	//p.buildClassPackageMap()

	// Generate all enums
	p.handleTsEnums()

	// Generate all enums mapping
	p.handleTsEnumsMapping()

	// Generate all classes
	p.handleTsClasses()

	// Generate all services
	p.handleTsServices()

	// Generate service exports
	p.generateServicesExports()

	// Generate all index.ts files (barrels)
	p.generateIndexes()
	return nil
}

// Generate service exports
func (p *TsProcessor) generateServicesExports() {
	var content []string
	for _, pkg := range p.Model.Packages {
		for sn := range pkg.Services {
			content = append(content, sn)
		}
	}
	if len(content) == 0 {
		return
	}
	tmpl, _ := template.New("services.index.ts.tpl").Parse(servicesIndexTsTemplate)
	fileName := path.Join(p.Output, "services.export.ts")

	if f, err := os.Create(fileName); err != nil {
		log.Fatal("Error creating file: ", fileName, err)
	} else if er := tmpl.Execute(f, content); er != nil {
		log.Fatal("Error executing template [services.index.ts.tpl]: ", er)
	} else {
		_ = f.Close()
	}
}

// build imports
//func (p *TsProcessor) buildClassPackageMap() {
//	for _, v := range p.Model.Packages {
//		for _, class := range v.Classes {
//			// Do not create import for parameter messages
//			if !class.IsParam {
//				classPackageMap[class.Name] = class.PackageShortName
//			}
//		}
//		for _, enum := range v.Enums {
//			classPackageMap[enum.Name] = "enums"
//		}
//	}
//}

// create directory
func (p *TsProcessor) makeDir(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatal("Error creating folder: "+path, err)
	}
}

func toCamelCase(s string) string {
	return fmt.Sprintf("%s%s", strings.ToLower(s[0:1]), s[1:])
}

func toCapitalCase(s string) string {
	return fmt.Sprintf("%s%s", strings.ToUpper(s[0:1]), strings.ToLower(s[1:]))
}

func toDisplayName(s string) string {
	parts := strings.Split(s, "_")
	caps := make([]string, 0)
	for _, p := range parts {
		caps = append(caps, toCapitalCase(p))
	}
	return strings.Join(caps, " ")
}

// Build method content - invoke rest utils http call
func methodContent(methodInfo model.MethodInfo) string {

	fmt.Println(methodInfo.Path, "method", methodInfo.Name)
	if methodInfo.Name == "SetAccount" {
		fmt.Println("stop here")
	}

	url := methodInfo.Path
	if url == "/" {
		url = ""
	}
	queryParamArg := ""
	content := ""
	bodyParam := ", ''"
	queryParam := ""

	if methodInfo.BodyParam != nil {
		bodyParam = fmt.Sprintf(", typeof %s === 'object' ? JSON.stringify(%s) : %s",
			methodInfo.BodyParam.Json,
			methodInfo.BodyParam.Json,
			methodInfo.BodyParam.Json)
	}

	for _, param := range methodInfo.PathParams {
		url = strings.Replace(
			url,
			fmt.Sprintf("{%s}", param.Json),
			"${"+param.Json+"}",
			-1)
	}

	for _, param := range methodInfo.QueryParams {
		queryParam += fmt.Sprintf(
			"    if (%s != null) { params.push(`%s=${%s}`); }\n",
			param.Json,
			param.Json,
			param.Json)
	}

	if len(queryParam) > 0 {
		content += fmt.Sprintf("const params = [];\t\t\n%s\n\t\t", queryParam)
		queryParamArg = ", ...params"
	}

	if methodInfo.Method == "GET" || methodInfo.Method == "DELETE" {
		bodyParam = ""
	}

	urlSuffix := url
	if urlSuffix == "/" {
		urlSuffix = ""
	}

	// Motty - create getUploadURL method
	if methodInfo.IsFileUpload {
		functionLine := fmt.Sprintf("return `${this.baseUrl}%s`;", urlSuffix)
		return content + functionLine
	}

	returnType := convertToTypeScript(methodInfo.ReturnClass)
	functionLine := fmt.Sprintf(
		"return this.rest.%s<%s>(`${this.baseUrl}%s`%s%s);",
		strings.ToLower(methodInfo.Method),
		returnType,
		urlSuffix,
		bodyParam,
		queryParamArg,
	)

	// If the response is a stream, apply http.download
	if methodInfo.Return.IsStream {
		fileName := methodInfo.Context
		if len(fileName) == 0 {
			fileName = "export"
		}
		functionLine = fmt.Sprintf(
			"return this.rest.download(`%s`,`${this.baseUrl}%s`%s%s);",
			fileName,
			url,
			bodyParam,
			queryParamArg,
		)
	}

	// If the request is a stream, apply http.upload
	if methodInfo.StreamsRequest {
		functionLine = fmt.Sprintf(
			"return this.rest.upload(%s,`${this.baseUrl}%s`%s);",
			methodInfo.FileParam.Json,
			url,
			queryParamArg,
		)
	}

	return content + functionLine
}

// Build method input parameters list
func handleMethodParams(methodInfo model.MethodInfo) string {
	p := ""

	//if methodInfo.ReturnClass == "EntitiesResponse<Syllabus>" {
	//	fmt.Println("stop here")
	//}
	if methodInfo.FileParam != nil {
		p += fmt.Sprintf("%s: %s, ", methodInfo.FileParam.Json, getTsType(methodInfo.FileParam.Type))
	}
	for _, param := range methodInfo.PathParams {
		if param.IsArray {
			p += fmt.Sprintf("%s?: %s[], ", param.Json, getTsType(param.Type))
		} else {
			p += fmt.Sprintf("%s?: %s, ", param.Json, getTsType(param.Type))
		}
	}
	for _, param := range methodInfo.QueryParams {
		if param.IsArray {
			p += fmt.Sprintf("%s?: %s[], ", param.Json, getTsType(param.Type))
		} else {
			p += fmt.Sprintf("%s?: %s, ", param.Json, getTsType(param.Type))
		}
	}
	if methodInfo.BodyParam != nil {
		if methodInfo.BodyParam.IsArray {
			p += fmt.Sprintf("%s?: %s[], ", methodInfo.BodyParam.Json, getTsType(methodInfo.BodyParam.Type))
		} else {
			p += fmt.Sprintf("%s?: %s, ", methodInfo.BodyParam.Json, getTsType(methodInfo.BodyParam.Type))
		}
	}

	if len(p) > 0 {
		p = p[0 : len(p)-2]
	}
	return p
}

func addServiceImports(methods []*model.MethodInfo) string {
	output := ""
	imports := make(map[string]string)
	for _, method := range methods {
		if method.BodyParam != nil && isObject(method.BodyParam.Type) {
			if _, ok := imports[method.BodyParam.Type]; !ok {
				imports[method.BodyParam.Type] = method.BodyParam.Type
			}
		}
		for _, param := range method.QueryParams {
			if isObject(param.Type) {
				if _, ok := imports[param.Type]; !ok {
					imports[param.Type] = param.Type
				}
			}
		}
		for _, param := range method.PathParams {
			if isObject(param.Type) {
				if _, ok := imports[param.Type]; !ok {
					imports[param.Type] = param.Type
				}
			}
		}
		returnTypes := getReturnTypes(method.ReturnClass)
		for _, rt := range returnTypes {
			imports[rt] = rt
		}
	}

	if len(imports) == 0 {
		return ""
	}

	for _, class := range imports {
		//if _, ok := classPackageMap[class]; ok {
		output += fmt.Sprintf("import { %s } from '.';\n", class)
		//}
	}

	return output
}

// Generate all services
func (p *TsProcessor) handleTsServices() {
	var serviceList []model.ServiceInfo
	for _, pkg := range p.Model.Packages {
		for _, service := range pkg.Services {
			serviceList = append(serviceList, *service)
		}
	}

	funcMap := template.FuncMap{
		"toLowerCase":        strings.ToLower,
		"toCamelCase":        toCamelCase,
		"methodContent":      methodContent,
		"handleMethodParams": handleMethodParams,
		"addServiceImports":  addServiceImports,
	}

	tmpl, _ := template.New("base_service.ts.tpl").Funcs(funcMap).Parse(serviceTsTemplate)
	for _, service := range serviceList {

		if service.TsName == "UsrDevicesService" {
			fmt.Println("stop here")
		}
		fName := service.Name
		if len(service.TsName) > 0 {
			fName = service.TsName
		}
		fileName := path.Join(p.Output, fmt.Sprintf("%s.ts", fName))
		if f, err := os.Create(fileName); err != nil {
			log.Fatal("Error creating file: ", fileName, err)
		} else if er := tmpl.Execute(f, service); er != nil {
			log.Fatal("Error executing template [base_service.ts.tpl]: ", er)
		} else {
			_ = f.Close()
		}
	}
}

// Add constructor method
func addConstructor(class model.ClassInfo) string {
	output := "    constructor("
	for _, field := range class.Fields {
		output += field.TsName + "?: " + getTsType(field.Type)
		if field.IsArray {
			output += "[]"
		}
		output += ", "
	}
	if len(class.Fields) > 0 {
		output = output[0 : len(output)-2]
	}
	output += ") { \n"
	for _, field := range class.Fields {
		line := fmt.Sprintf(`if (%s !== undefined) { this.%s = %s; }`, field.TsName, field.Json, field.TsName)
		output += "        " + line + "\n"
	}
	output += "    }\n"
	return output
}

// Generate classes
func (p *TsProcessor) handleTsClasses() {
	funcMap := template.FuncMap{
		"getTsType":      getTsType,
		"addImports":     addImports,
		"addConstructor": addConstructor,
		"join":           strings.Join,
		"genericsParam":  genericsParam,
	}

	var classList []model.ClassInfo
	for _, v := range p.Model.Packages {
		for _, class := range v.Classes {
			classList = append(classList, *class)
		}
	}

	tmpl, _ := template.New("base_class.ts.tpl").Funcs(funcMap).Parse(classTsTemplate)
	for _, class := range classList {

		// For parameter classes, do not create TS file
		if !class.IsParam {
			fileName := path.Join(p.Output, fmt.Sprintf("%s.ts", class.Name))

			if f, err := os.Create(fileName); err != nil {
				log.Fatal("Error creating file: ", fileName, err)
			} else if er := tmpl.Execute(f, class); er != nil {
				log.Fatal("Error executing template [base_class.ts.tpl]: ", er)
			} else {
				_ = f.Close()
			}
		}
	}
}

// Generate indexes
func (p *TsProcessor) generateIndexes() {
	var content []string

	// Order is important, start with system classes first, than enums, classes and services
	//basePkg := p.Model.Packages["base"]
	//for en := range basePkg.Enums {
	//	content = append(content, fmt.Sprintf("%s", en))
	//}
	//for cn := range basePkg.Classes {
	//	content = append(content, fmt.Sprintf("%s", cn))
	//}
	//for sn := range basePkg.Services {
	//	content = append(content, fmt.Sprintf("%s", sn))
	//}

	//mainPkg := p.Model.Packages["model"]
	//for en := range mainPkg.Enums {
	//	content = append(content, fmt.Sprintf("%s", en))
	//}
	//for cn := range mainPkg.Classes {
	//	content = append(content, fmt.Sprintf("%s", cn))
	//}
	//for sn := range mainPkg.Services {
	//	content = append(content, fmt.Sprintf("%s", sn))
	//}

	for _, pkg := range p.Model.Packages {
		for en := range pkg.Enums {
			content = append(content, fmt.Sprintf("%s", en))
		}
		for cn := range pkg.Classes {
			content = append(content, fmt.Sprintf("%s", cn))
		}
		for sn := range pkg.Services {
			content = append(content, fmt.Sprintf("%s", sn))
		}
	}

	p.generateIndexTs(content)
}

// getTsType - convert variables types to known TypeScript types
func getTsType(pType string) string {
	if _, ok := tsTypes[pType]; ok {
		return tsTypes[pType]
	}
	return pType
}

// Generate enums
func (p *TsProcessor) handleTsEnums() {

	funcMap := template.FuncMap{
		"toDisplayName": toDisplayName,
	}

	var enumList []model.EnumInfo
	for _, v := range p.Model.Packages {
		for _, enum := range v.Enums {
			enumList = append(enumList, *enum)

			// If this is a flag enum, put the native value in the class field
			if enum.IsFlags {
				tsTypes[enum.Name] = "number"
			}
		}
	}
	tmpl, _ := template.New("base_enum.ts.tpl").Funcs(funcMap).Parse(enumTsTemplate)
	for _, enum := range enumList {
		fileName := path.Join(p.Output, fmt.Sprintf("%s.ts", enum.Name))
		if f, err := os.Create(fileName); err != nil {
			log.Fatal("Error creating file: ", fileName, err)
		} else if er := tmpl.Execute(f, enum); er != nil {
			log.Fatal("Error executing template [base_enum.ts.tpl]: ", er)
		} else {
			_ = f.Close()
		}
	}
}

// Generate enums mapping
func (p *TsProcessor) handleTsEnumsMapping() {

	funcMap := template.FuncMap{
		"toDisplayName": toDisplayName,
	}

	var enumList []model.EnumInfo
	for _, v := range p.Model.Packages {
		for _, enum := range v.Enums {
			enumList = append(enumList, *enum)

			// If this is a flag enum, put the native value in the class field
			if enum.IsFlags {
				tsTypes[enum.Name] = "number"
			}
		}
	}
	tmpl, _ := template.New("enum_map.ts.tpl").Funcs(funcMap).Parse(enumMapTsTemplate)
	fileName := path.Join(p.Output, "Enums.ts")

	if f, err := os.Create(fileName); err != nil {
		log.Fatal("Error creating file: ", fileName, err)
	} else if er := tmpl.Execute(f, enumList); er != nil {
		log.Fatal("Error executing template [enum_map.ts.tpl]: ", er)
	} else {
		_ = f.Close()
	}
}

// Generate TypeScript index
func (p *TsProcessor) generateIndexTs(data []string) {
	tmpl, _ := template.New("index.ts.tpl").Parse(indexTsTemplate)
	fileName := path.Join(p.Output, "index.ts")

	if f, err := os.Create(fileName); err != nil {
		log.Fatal("Error creating file: ", fileName, err)
	} else if er := tmpl.Execute(f, data); er != nil {
		log.Fatal("Error executing template [index.ts.tpl]: ", er)
	} else {
		_ = f.Close()
	}

}

// Add imports based on the class dependencies
func addImports(class model.ClassInfo) string {
	output := ""
	for className, _ := range class.Dependencies {
		output += fmt.Sprintf("import { %s } from '.';\n", className)
	}
	return output
}

func stripGeneric(name string) string {
	idx := strings.Index(name, "<")
	if idx > 0 {
		return name[:idx]
	} else {
		return name
	}
}

func isObject(pType string) bool {
	if _, ok := tsTypes[pType]; ok {
		return false
	}
	return true
}

func genericsParam(class model.ClassInfo) string {
	if class.IsGeneric {
		list := make([]string, 0)
		for _, kv := range class.GenericTypes {
			list = append(list, kv.Key)
		}
		return fmt.Sprintf("<%s>", strings.Join(list, ","))
	} else {
		return ""
	}
}

func convertToTypeScript(name string) string {
	tokens := strings.Split(name, "<")
	types := make([]string, 0)
	for _, token := range tokens {
		types = append(types, strings.ReplaceAll(token, ">", ""))
	}

	for _, t := range types {
		tsType := getTsType(t)
		if tsType != t {
			name = strings.ReplaceAll(name, t, tsType)
		}
	}

	return name
}

func getReturnTypes(name string) []string {
	result := make([]string, 0)
	tokens := strings.Split(name, "<")
	types := make([]string, 0)
	for _, token := range tokens {
		types = append(types, strings.ReplaceAll(token, ">", ""))
	}

	for _, t := range types {
		if isObject(t) {
			result = append(result, t)
		}
	}

	return result
}
