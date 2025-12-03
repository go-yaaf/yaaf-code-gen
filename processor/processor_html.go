package processor

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/go-yaaf/yaaf-code-gen/model"
)

var classes []model.ClassInfo
var enums []model.EnumInfo
var services []model.ServiceInfo
var sockets []model.WebSocketInfo

// HtmlProcessor - Html processor converts proto files to documentation files
type HtmlProcessor struct {
	BaseProcessor
}

// NewHtmlProcessor - Factory method
func NewHtmlProcessor(model *model.MetaModel, output string) Processor {
	return &HtmlProcessor{BaseProcessor{
		Output: output,
		Model:  model,
	}}
}

// Start the processor
func (p *HtmlProcessor) Start() error {

	// First, ensure output directory
	if err := os.MkdirAll("./output/html/img", os.ModePerm); err != nil {
		log.Fatal("Error creating folder: ./output/html", err)
	}

	if err := p.dirCopy("./templates/html/img", "./output/html/img"); err != nil {
		log.Fatal("Error copy folder: ./templates/html/img: ", err)
	}

	// Iterate over the packages and merge all classes
	for _, v := range p.Model.Packages {

		// Generate class pages and append to all classes
		for _, class := range v.Classes {
			if class.IsVisible {
				classes = append(classes, *class)
				p.generateClassPage(*class)
			}
		}

		// Generate enum pages and append to all enums
		for _, enum := range v.Enums {
			enums = append(enums, *enum)
			p.generateEnumPage(*enum)
		}

		// Generate service pages and append to all services
		for _, service := range v.Services {
			services = append(services, *service)
			p.generateServicePage(*service)
		}

		// Generate Web Socket service pages and append to all socket services
		for _, socket := range v.Sockets {
			sockets = append(sockets, *socket)
			p.generateWebSocketPage(*socket)
		}
	}

	p.generateClassesTable(classes)
	p.generateEnumsTable(enums)
	p.generateResources(services)
	p.generateWebSocketsPage(sockets)
	p.generateCSS()

	return nil
}

// Generate CSS files
func (p *HtmlProcessor) generateCSS() {
	tmpl, _ := template.New("style.css").ParseFiles("templates/html/style.css")
	f, err := os.Create("./output/html/style.css")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}
	if er := tmpl.Execute(f, ""); er != nil {
		log.Fatal("Error executing template", tmpl.Name(), er)
	}
}

// Generate class page
func (p *HtmlProcessor) generateClassPage(class model.ClassInfo) {
	funcMap := template.FuncMap{
		"getType": getType,
	}
	tmpl, _ := template.New("base.html").Funcs(funcMap).ParseFiles("templates/html/footer.html", "templates/html/json_data_class.html", "templates/html/base.html")
	f, err := os.Create("./output/html/json_" + class.Name + ".html")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}
	if er := tmpl.Execute(f, class); er != nil {
		log.Fatal("Error executing template", tmpl.Name(), er)
	}
}

// Generate enum page
func (p *HtmlProcessor) generateEnumPage(enum model.EnumInfo) {
	tmpl, _ := template.New("base.html").ParseFiles("templates/html/footer.html", "templates/html/json_data_enum.html", "templates/html/base.html")
	f, err := os.Create("./output/html/json_" + enum.Name + ".html")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}
	if er := tmpl.Execute(f, enum); er != nil {
		log.Fatal("Error executing template", tmpl.Name(), er)
	}
}

// Generate service page
func (p *HtmlProcessor) generateServicePage(service model.ServiceInfo) {
	funcMap := template.FuncMap{
		"getType":      getType,
		"addBodyParam": addBodyParam,
		"addParams":    addParams,
	}
	tmpl, _ := template.New("base.html").Funcs(funcMap).ParseFiles("templates/html/footer.html", "templates/html/json_data_service.html", "templates/html/base.html")
	f, err := os.Create("./output/html/resource_" + service.Name + ".html")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}
	if er := tmpl.Execute(f, service); er != nil {
		log.Fatal("Error executing template", tmpl.Name(), er)
	}
}

// Generate web-socket page
func (p *HtmlProcessor) generateWebSocketPage(socket model.WebSocketInfo) {
	funcMap := template.FuncMap{
		"getType":      getType,
		"addBodyParam": addBodyParam,
		"addParams":    addParams,
	}
	tmpl, _ := template.New("base.html").Funcs(funcMap).ParseFiles("templates/html/footer.html", "templates/html/json_web_socket.html", "templates/html/base.html")
	f, err := os.Create("./output/html/web_socket_" + socket.Name + ".html")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}
	if er := tmpl.Execute(f, socket); er != nil {
		log.Fatal("Error executing template", tmpl.Name(), er)
	}
}

func listServicesGroups(services []model.ServiceInfo) map[string][]model.ServiceInfo {
	groups := make(map[string][]model.ServiceInfo)
	for _, service := range services {
		if len(service.Group) > 0 {
			groups[service.Group] = append(groups[service.Group], service)
		} else {
			groups[""] = append(groups[""], service)
		}
	}
	return groups
}

func removeSpaces(str string) string {
	return strings.Replace(str, " ", "_", -1)
}

// Generate resources page
func (p *HtmlProcessor) generateResources(services []model.ServiceInfo) {
	if len(services) == 0 {
		return
	}

	// sort services
	sort.Slice(services, func(i, j int) bool {
		return services[i].Name < services[j].Name
	})

	funcMap := template.FuncMap{
		"listServiceMethods":  listServiceMethods,
		"listPathMethodTypes": listPathMethodTypes,
		"listServicesGroups":  listServicesGroups,
		"removeSpaces":        removeSpaces,
	}
	tmpl, _ := template.New("base.html").Funcs(funcMap).ParseFiles("templates/html/base.html", "templates/html/index.html", "templates/html/footer.html")
	f, err := os.Create("./output/html/index.html")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}
	if err := tmpl.Execute(f, services); err != nil {
		fmt.Println("Error generating: ", err)
		return
	}
}

// Generate web sockets page
func (p *HtmlProcessor) generateWebSocketsPage(sockets []model.WebSocketInfo) {
	if len(sockets) == 0 {
		return
	}

	tmpl, _ := template.New("base.html").ParseFiles("templates/html/footer.html", "templates/html/web_sockets.html", "templates/html/base.html")
	f, err := os.Create("./output/html/webSockets.html")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}
	if err := tmpl.Execute(f, sockets); err != nil {
		fmt.Println("Error generateWebSocketsPage: ", err)
		return
	}
}

func contains(classes []model.ClassInfo, className string) bool {
	for _, c := range classes {
		if c.Name == className {
			return true
		}
	}
	return false
}

func addBodyParam(bodyParam *model.ParamInfo) string {
	rows := ""
	dataTypeRef := ""
	arrayPrefix := ""
	if bodyParam != nil {
		if bodyParam.IsArray {
			arrayPrefix = "array of "
		}
		if contains(classes, bodyParam.Name) {
			dataTypeRef = fmt.Sprintf(`%s<a href="json_%s.html">%s</a> (JSON)`, arrayPrefix, bodyParam.Name, bodyParam.Name)
		} else {
			dataTypeRef = fmt.Sprintf(`%s<a href="json_%s.html">%s</a> (JSON)`, arrayPrefix, bodyParam.Type, bodyParam.Type)
		}
		rows += fmt.Sprintf(
			`
			<tr>
				<td><abbr data-toggle="tooltip" data-placement="top"
				title="Use the &quot;Content-Type: application/json&quot; HTTP header to specify this media type to the server."><span
				class="request-type">application/json</span></abbr></td>
				<td><span class="datatype-reference">%s</span></td>
				<td><span class="request-description">%s</span></td>
			</tr>
			 `,
			dataTypeRef,
			strings.Join(bodyParam.Docs, "<br>"),
		)
	}
	return "<tr>" + rows + "</tr>"
}

func listServiceMethods(service model.ServiceInfo) string {
	methods := make(map[string]string)
	output := ""
	for _, method := range service.Methods {
		methods[method.Path] = method.Path
	}
	for _, method := range methods {
		output += fmt.Sprintf(`
		<li>
			<samp>
				<span class="resource-path">%s%s</span>
			</samp>
		</li>
		`,
			service.Path,
			method,
		)
	}
	return output
}

func listPathMethodTypes(service model.ServiceInfo) string {
	methods := make(map[string][]string)
	output := ""
	for _, method := range service.Methods {
		methods[method.Path] = append(methods[method.Path],
			fmt.Sprintf(`<span class="label label-default resource-method">%s</span>`,
				method.Method),
		)
	}

	for _, method := range methods {
		output += fmt.Sprintf(`
		<li>
			<samp>
				%s
			</samp>
		</li>
		`,
			strings.Join(method, "&nbsp;"),
		)
	}
	return output
}

func addParams(params []*model.ParamInfo) string {

	rows := ""
	for _, param := range params {
		rows += fmt.Sprintf(
			`
			<tr>
				<td><span class="parameter-name"> %s </span></td>
				<td>%s</td>
				<td><span class="parameter-description">%s</span></td>
			</tr>`,
			param.TsName,
			param.ParamType,
			strings.Join(param.Docs, "<br>"))
	}
	return rows
}

// Generate enums table
func (p *HtmlProcessor) generateEnumsTable(enums []model.EnumInfo) {
	if len(enums) == 0 {
		return
	}

	// sort enums
	sort.Slice(enums, func(i, j int) bool {
		return enums[i].Name < enums[j].Name
	})

	tmpl, _ := template.New("base.html").ParseFiles("templates/html/footer.html", "templates/html/enums.html", "templates/html/base.html")
	f, err := os.Create("./output/html/enums.html")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}
	if err := tmpl.Execute(f, enums); err != nil {
		fmt.Println("Error generateEnumsTable: ", err)
		return
	}
}

// Generate classes table
func (p *HtmlProcessor) generateClassesTable(classes []model.ClassInfo) {
	if len(classes) == 0 {
		return
	}

	// sort classes
	sort.Slice(classes, func(i, j int) bool {
		return classes[i].Name < classes[j].Name
	})

	tmpl, _ := template.New("base.html").ParseFiles("templates/html/footer.html", "templates/html/data_types.html", "templates/html/base.html")
	f, err := os.Create("./output/html/dataTypes.html")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}
	if err := tmpl.Execute(f, classes); err != nil {
		fmt.Println("Error generating: ", err)
		return
	}
}

func getType(pType string) string {
	types := map[string]string{
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
		"bytes":    "string",
	}

	if _, ok := types[pType]; ok {
		return types[pType]
	}
	return "<a href='json_" + pType + ".html'>" + pType + "</a>"
}
