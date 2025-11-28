package processor

import (
	"fmt"
	"github.com/go-yaaf/yaaf-code-gen/model"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var tsTypes = map[string]string{
	"double":   "number",
	"float":    "number",
	"int32":    "number",
	"int64":    "number",
	"uint32":   "number",
	"uint64":   "number",
	"sint32":   "number",
	"sint64":   "number",
	"fixed32":  "number",
	"fixed64":  "number",
	"sfixed32": "number",
	"sfixed64": "number",
	"bool":     "boolean",
	"string":   "string",
	"bytes":    "File",
	"Any":      "any",
}

var fileImports = make(map[string]string)

// TsProcessor - TS processor converts meta model to TypeScript files
type TsProcessor struct {
	Model *model.MetaModel
}

// NewTsProcessor - Factory method
func NewTsProcessor() Processor {
	return &TsProcessor{}
}

// Process starts the processor
func (p *TsProcessor) Process(metaModel *model.MetaModel) error {

	p.Model = metaModel

	// build list of file imports
	buildFileImports(p)

	// First, ensure output directories exist
	makeAllDirs(p)

	// Generate all enums
	handleTsEnums(p)

	// Generate all classes
	handleTsClasses(p)

	// Generate all services
	handleTsServices(p)

	// Generate all services exports
	p.generateServicesExports()

	// Generate all index.ts files (barrels)
	generateIndexes(p)

	return nil
}

func (p *TsProcessor) generateServicesExports() {
	var services []*ServiceInfo
	for _, v := range p.Model.Packages {
		for _, service := range v.Services {
			services = append(services, service)
		}
	}
	tmpl, _ := template.New("services.index.ts.tpl").ParseFiles("templates/ts/services.index.ts.tpl")
	f, err := os.Create("./output/ts/services/services.export.ts")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}
	tmpl.Execute(f, services)
}

func (p *TsProcessor) buildFileImports() {
	for _, v := range p.Model.Packages {
		for _, class := range v.Classes {
			fileImports[class.Name] = class.Package
		}
		for _, enum := range v.Enums {
			fileImports[enum.Name] = "enums"
		}
	}
}

func (p *TsProcessor) makeAllDirs() {
	mkDir("./output/ts/enums")
	mkDir("./output/ts/services")
	for _, v := range p.Model.Packages {
		for _, class := range v.Classes {
			mkDir("./output/ts/" + class.Package)
		}
	}
}

func mkDir(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatal("Error creating folder: "+path, err)
	}
}

func toCamelCase(s string) string {
	firstLetter := s[0:1]
	return fmt.Sprintf("%s%s", strings.ToLower(firstLetter), s[1:])
}

func methodContent(methodInfo model.MethodInfo) string {
	path := methodInfo.Path
	if path == "/" {
		path = ""
	}
	queryParamArg := ""
	content := ""
	bodyParam := ", null"
	queryParam := ""

	if methodInfo.BodyParam != nil {
		bodyParam = fmt.Sprintf(", typeof %s === 'object' ? JSON.stringify(%s) : %s",
			methodInfo.BodyParam.Json,
			methodInfo.BodyParam.Json,
			methodInfo.BodyParam.Json)
	}

	for _, param := range methodInfo.PathParams {
		path = strings.Replace(
			path,
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
		content += fmt.Sprintf("const params = new Array();\t\t\n%s\n\t\t", queryParam)
		queryParamArg = ", ...params"
	}

	if methodInfo.Method == "GET" || methodInfo.Method == "DELETE" {
		bodyParam = ""
	}

	functionLine := fmt.Sprintf(
		"return this.rest.%s(`${this.baseUrl}%s`%s%s);",
		strings.ToLower(methodInfo.Method),
		path,
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
			path,
			bodyParam,
			queryParamArg,
		)
	}

	// If the request is a stream, apply http.upload
	if methodInfo.StreamsRequest {
		functionLine = fmt.Sprintf(
			"return this.rest.upload(%s,`${this.baseUrl}%s`%s);",
			methodInfo.FileParam.Json,
			path,
			queryParamArg,
		)
	}

	return content + functionLine
}

func handleMethodParams(methodInfo MethodInfo) string {
	p := ""
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

	if methodInfo.FileParam != nil {
		p += fmt.Sprintf("%s?: %s, ", methodInfo.FileParam.Json, getTsType(methodInfo.FileParam.Type))
	}

	if len(p) > 0 {
		p = p[0 : len(p)-2]
	}
	return p
}

func addServiceImports(methods []*MethodInfo) string {
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
	}

	if len(imports) == 0 {
		return ""
	}

	for _, class := range imports {
		if _, ok := fileImports[class]; ok {
			output += fmt.Sprintf("import { %s } from '../%s/%s';\n", class, fileImports[class], class)
		}
	}

	return output
}

func handleTsServices(p *TsProcessor) {
	var services []ServiceInfo
	for _, v := range p.Model.Packages {
		for _, service := range v.Services {
			services = append(services, *service)
		}
	}

	funcMap := template.FuncMap{
		"toLowerCase":        strings.ToLower,
		"toCamelCase":        toCamelCase,
		"methodContent":      methodContent,
		"handleMethodParams": handleMethodParams,
		"addServiceImports":  addServiceImports,
	}

	tmpl, _ := template.New("base_service.ts.tpl").Funcs(funcMap).ParseFiles("templates/ts/base_service.ts.tpl")
	for _, service := range services {
		f, err := os.Create("./output/ts/services/" + service.TsName + ".ts")
		if err != nil {
			fmt.Println("create file: ", err)
			return
		}
		tmpl.Execute(f, service)
	}
}

func addConstructor(class ClassInfo) string {
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
		output += "        this." + field.Json + " = " + field.TsName + ";\n"
	}
	output += "    }\n"
	return output
}

// handleTsClasses - Generate classes
func handleTsClasses(p *TsProcessor) {
	funcMap := template.FuncMap{
		"getTsType":      getTsType,
		"addImports":     addImports,
		"addConstructor": addConstructor,
	}

	var classes []ClassInfo
	for _, v := range p.Model.Packages {
		for _, class := range v.Classes {
			classes = append(classes, *class)
		}
	}

	tmpl, _ := template.New("base_class.ts.tpl").Funcs(funcMap).ParseFiles("templates/ts/base_class.ts.tpl")
	for _, class := range classes {
		f, err := os.Create("./output/ts/" + class.PackageShortName + "/" + class.Name + ".ts")
		if err != nil {
			fmt.Println("create file: ", err)
			return
		}
		tmpl.Execute(f, class)
	}
}

func generateIndexes(p *TsProcessor) {
	var path = "./output/ts/"
	folders := make(map[string]string)
	content := []string{}
	folders["enums"] = "enums"
	folders["services"] = "services"
	for _, v := range p.Model.Packages {
		for _, class := range v.Classes {
			folders[class.PackageShortName] = class.PackageShortName
		}
	}

	// generate a index.ts file for each folder
	for _, folder := range folders {
		files, err := ioutil.ReadDir(path + folder)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			var filename = f.Name()
			var extension = filepath.Ext(filename)
			var name = filename[0 : len(filename)-len(extension)]
			content = append(content, name)
		}
		generateIndexTs(path+folder+"/", content)
		content = []string{} // reset content
	}

	// Then generate index.ts for all folders
	content = []string{} // reset content
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		content = append(content, f.Name())
	}
	generateIndexTs(path+"/", content)

}

// getTsType - convert variables types to known TypeScript types
func getTsType(pType string) string {
	if _, ok := tsTypes[pType]; ok {
		return tsTypes[pType]
	}
	return pType
}

// handleTsEnums - Generate enums
func handleTsEnums(p *TsProcessor) {
	var enums []EnumInfo
	for _, v := range p.Model.Packages {
		for _, enum := range v.Enums {
			enums = append(enums, *enum)

			// If this is a flag enum, put the native value in the class field
			if enum.IsFlags {
				tsTypes[enum.Name] = "number"
			}
		}
	}
	tmpl, _ := template.New("base_enum.ts.tpl").ParseFiles("templates/ts/base_enum.ts.tpl")
	for _, enum := range enums {
		f, err := os.Create("./output/ts/enums/" + enum.Name + ".ts")
		if err != nil {
			fmt.Println("create file: ", err)
			return
		}
		tmpl.Execute(f, enum)
	}
}

func generateIndexTs(path string, data []string) {
	tmpl, _ := template.New("index.ts.tpl").ParseFiles("templates/ts/index.ts.tpl")
	f, err := os.Create(path + "index.ts")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}
	tmpl.Execute(f, data)
}

func addImports(class ClassInfo) string {
	output := ""
	imports := make(map[string]string)
	for _, field := range class.Fields {
		if isObject(field.Type) && field.Type != class.Name {
			if _, ok := imports[field.Type]; !ok {
				imports[field.Type] = field.Type
			}
		}
	}

	if class.IsExtend {
		imports[class.BaseClass] = class.BaseClass
	}

	if len(imports) == 0 {
		return ""
	}

	for _, class := range imports {
		if _, ok := fileImports[class]; ok {
			output += fmt.Sprintf("import { %s } from '../%s/%s';\n", class, fileImports[class], class)
		}
	}

	return output
}

func isObject(pType string) bool {
	if _, ok := tsTypes[pType]; ok {
		return false
	}
	return true
}
