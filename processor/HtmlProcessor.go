package processor

import (
	"fmt"
	"github.com/go-yaaf/yaaf-code-gen/model"
	"os"
	"strings"
	"text/template"
)

// HtmlProcessor - Html processor converts proto files to documentation files
type HtmlProcessor struct {
	outputFolder string
	Model        *model.MetaModel
}

// NewHtmlProcessor - Factory method
func NewHtmlProcessor(outputFolder string) Processor {
	return &HtmlProcessor{outputFolder: outputFolder}
}

// Name returns the processor name
func (p *HtmlProcessor) Name() string {
	return "HTML Processor"
}

// Process starts the processor
func (p *HtmlProcessor) Process(metaModel *model.MetaModel) error {

	p.Model = metaModel

	// First, ensure output directory
	if err := os.MkdirAll(p.outputFolder, os.ModePerm); err != nil {
		return fmt.Errorf("error creating folders: %s: %s", p.outputFolder, err.Error())
	}

	// Iterate over the packages and merge all classes
	for _, v := range p.Model.Packages {

		// Generate class pages and append to all classes
		for _, class := range v.Classes {
			if err := p.generateClassPage(class); err != nil {
				fmt.Println(err.Error())
			}
		}

		// Generate enum pages and append to all enums
		for _, enum := range v.Enums {
			if err := p.generateEnumPage(enum); err != nil {
				fmt.Println(err.Error())
			}
		}

		// Generate service pages and append to all services
		for _, service := range v.Services {
			if err := p.generateServicePage(service); err != nil {
				fmt.Println(err.Error())
			}
		}
	}

	p.generateClassesTable()
	p.generateEnumsTable()
	p.generateResources()
	p.generateCSS()

	return nil
}

func (p *HtmlProcessor) generateCSS() {
	tmpl, _ := template.New("style.css").ParseFiles("templates/html/style.css")
	f, err := os.Create("./output/html/style.css")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}
	tmpl.Execute(f, "")
}

// Create HTML class page
func (p *HtmlProcessor) generateClassPage(class *model.ClassInfo) error {
	funcMap := template.FuncMap{
		"getType": getType,
	}
	tmpl, terr := template.New("base.html").Funcs(funcMap).ParseFiles("templates/html/footer.html", "templates/html/json_data_class.html", "templates/html/base.html")
	if terr != nil {
		return fmt.Errorf("error parsing template: %s: %s", "base.html", terr)
	}

	filePath := fmt.Sprintf("%s/json_%s.html", p.outputFolder, class.Name)
	if f, err := os.Create(filePath); err != nil {
		return fmt.Errorf("error creating file: %s: %s", filePath, err)
	} else {
		return tmpl.Execute(f, class)
	}
}

// Create HTML enum page
func (p *HtmlProcessor) generateEnumPage(enum *model.EnumInfo) error {

	tmpl, terr := template.New("enum").ParseFiles("templates/html/footer.html", "templates/html/json_data_enum.html", "templates/html/base.html")
	if terr != nil {
		return fmt.Errorf("error parsing template: %s: %s", "enum", terr)
	}

	filePath := fmt.Sprintf("%s/json_%s.html", p.outputFolder, enum.Name)
	if f, err := os.Create(filePath); err != nil {
		return fmt.Errorf("error creating file: %s: %s", filePath, err)
	} else {
		return tmpl.Execute(f, enum)
	}
}

// Create HTML service page
func (p *HtmlProcessor) generateServicePage(service *model.ServiceInfo) error {
	funcMap := template.FuncMap{
		"getType":      getType,
		"addBodyParam": addBodyParam,
		"addParams":    addParams,
	}

	tm := template.New("service")
	tm.Funcs(funcMap)
	tmpl, terr := tm.ParseFiles("templates/html/footer.html", "templates/html/json_data_service.html", "templates/html/base.html")

	//tmpl, terr := template.New("base.html").Funcs(funcMap).ParseFiles("templates/html/footer.html", "templates/html/json_data_service.html", "templates/html/base.html")
	if terr != nil {
		return fmt.Errorf("error parsing template: %s: %s", "service", terr)
	}

	filePath := fmt.Sprintf("%s/resource_%s.html", p.outputFolder, service.Name)
	if f, err := os.Create(filePath); err != nil {
		return fmt.Errorf("error creating file: %s: %s", filePath, err.Error())
	} else {
		return tmpl.Execute(f, *service)
	}
}

func (p *HtmlProcessor) listServicesGroups() map[string][]*model.ServiceInfo {
	groups := make(map[string][]*model.ServiceInfo)
	for _, v := range p.Model.Packages {
		for _, service := range v.Services {
			if len(service.Group) > 0 {
				groups[service.Group] = append(groups[service.Group], service)
			} else {
				groups[""] = append(groups[""], service)
			}
		}
	}
	return groups
}

func removeSpaces(str string) string {
	return strings.Replace(str, " ", "_", -1)
}

//func contains(classes []ClassInfo, className string) bool {
//	for _, c := range classes {
//		if c.Name == className {
//			return true
//		}
//	}
//	return false
//}

func addBodyParam(bodyParam *model.ParamInfo) string {
	rows := ""
	dataTypeRef := ""
	arrayPrefix := ""
	if bodyParam != nil {
		if bodyParam.IsArray {
			arrayPrefix = "array of "
		}
		//if contains(classes, bodyParam.Name) {
		dataTypeRef = fmt.Sprintf(`%s<a href="json_%s.html">%s</a> (JSON)`, arrayPrefix, bodyParam.Name, bodyParam.Name)
		//} else {
		//	dataTypeRef = fmt.Sprintf(`%s<a href="json_%s.html">%s</a> (JSON)`, arrayPrefix, bodyParam.Type, bodyParam.Type)
		//}
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

func (p *HtmlProcessor) listServiceMethods(service *model.ServiceInfo) string {
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

func (p *HtmlProcessor) listPathMethodTypes(service model.ServiceInfo) string {
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
				<td><span class="parameter-default-value">%s</span></td>
				<td><span class="parameter-constraints">%s</span></td>
			</tr>`,
			param.TsName(),
			param.ParamType,
			strings.Join(param.Docs, "<br>"),
			"",
			"")
	}
	return rows
}

// Create enums tables
func (p *HtmlProcessor) generateEnumsTable() {
	tmpl, _ := template.New("base.html").ParseFiles("templates/html/footer.html", "templates/html/enums.html", "templates/html/base.html")
	f, err := os.Create("./output/html/enums.html")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}

	// create enum list
	enums := make([]*model.EnumInfo, 0)
	for _, v := range p.Model.Packages {
		enums = append(enums, v.Enums...)
	}

	if err := tmpl.Execute(f, enums); err != nil {
		fmt.Println("Error generateEnumsTable: ", err)
		return
	}
}

// Create classes tables
func (p *HtmlProcessor) generateClassesTable() {
	tmpl, _ := template.New("base.html").ParseFiles("templates/html/footer.html", "templates/html/data_types.html", "templates/html/base.html")
	f, err := os.Create("./output/html/dataTypes.html")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}

	// create classes list
	classes := make([]*model.ClassInfo, 0)
	for _, v := range p.Model.Packages {
		classes = append(classes, v.Classes...)
	}
	if err := tmpl.Execute(f, classes); err != nil {
		fmt.Println("Error generating: ", err)
		return
	}
}

// Create resources tables
func (p *HtmlProcessor) generateResources() {
	funcMap := template.FuncMap{
		"listServiceMethods":  p.listServiceMethods,
		"listPathMethodTypes": p.listPathMethodTypes,
		"listServicesGroups":  p.listServicesGroups,
		"removeSpaces":        removeSpaces,
	}
	tmpl, _ := template.New("base.html").Funcs(funcMap).ParseFiles("templates/html/base.html", "templates/html/resources.html", "templates/html/footer.html")
	f, err := os.Create("./output/html/index.html")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}

	// create service list
	services := make([]*model.ServiceInfo, 0)
	for _, v := range p.Model.Packages {
		services = append(services, v.Services...)
	}
	if err := tmpl.Execute(f, services); err != nil {
		fmt.Println("Error generating: ", err)
		return
	}
}

// region Template Functions -------------------------------------------------------------------------------------------

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

// endregion
