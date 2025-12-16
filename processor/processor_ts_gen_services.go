package processor

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/go-yaaf/yaaf-code-gen/model"
)

// region TS template Service Processor --------------------------------------------------------------------------------

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

	folder := path.Join(p.Output, "services")
	p.makeDir(folder)

	var list []string

	tmpl, _ := template.New("base_service.ts.tpl").Funcs(funcMap).Parse(serviceTsTemplate)
	for _, service := range serviceList {

		fName := service.Name
		if len(service.TsName) > 0 {
			fName = service.TsName
		}

		var tpl bytes.Buffer
		if err := tmpl.Execute(&tpl, service); err != nil {
			log.Fatal("Error executing template [base_service.ts.tpl]: ", err)
		}
		// Remove newlines
		processedContent := p.trimNewLines(tpl.String())

		list = append(list, fName)

		fileName := path.Join(folder, fmt.Sprintf("%s.ts", fName))
		if f, err := os.Create(fileName); err != nil {
			log.Fatal("Error creating file: ", fileName, err)
		} else {
			if _, err = f.WriteString(processedContent); err != nil {
				log.Fatal("Error writing to file: ", fileName, err)
			}
			_ = f.Close()
		}
	}

	// Create the enums index file
	p.generateIndexTs(list, folder)
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

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, content); err != nil {
		log.Fatal("Error executing template [services.index.ts.tpl]: ", err)
	}

	// Remove newlines
	processedContent := p.trimNewLines(tpl.String())

	if f, err := os.Create(fileName); err != nil {
		log.Fatal("Error creating file: ", fileName, err)
	} else {
		if _, err = f.WriteString(processedContent); err != nil {
			log.Fatal("Error writing to file: ", fileName, err)
		}
		_ = f.Close()
	}
}

// Build method content - invoke rest utils http call
func methodContent(methodInfo model.MethodInfo) string {

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
		content += fmt.Sprintf("const params = [];\n%s\n\t\t", queryParam)
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

func addServiceImports(service model.ServiceInfo) string {
	output := ""
	for className, _ := range service.Dependencies {
		output += fmt.Sprintf("import { %s } from '../model';\n", className)
	}
	return output
}

// endregion

// region TypeScript service file template -----------------------------------------------------------------------------

var serviceTsTemplate = `
import { Injectable, Inject } from '@angular/core';
import { RestUtils } from '../../rest-utils';
import { APP_CONFIG, AppConfig } from '../../config';

{{. | addServiceImports}}

{{range .Docs}}
// {{.}} {{end}}
@Injectable({
  providedIn: 'root'
})
export class {{.TsName}} {

  // URL to web api
  private baseUrl = '{{.Path}}';

  // Class constructor
  constructor(@Inject(APP_CONFIG) private config: AppConfig, private rest: RestUtils) {
    this.baseUrl = this.config.api + this.baseUrl;
  }

{{range .Methods}}
  /**{{range .Docs}}
   * {{.}}{{end}}
   */
  {{.Name | toCamelCase}}({{. | handleMethodParams}}) {
    {{. | methodContent}}
  }
{{end}}
}
`

// endregion

// region TypeScript index file template -------------------------------------------------------------------------------

var servicesIndexTsTemplate = `
{{range .}}import { {{.}} } from '.';
{{end}}
export const Services = [
    {{range .}}{{.}},
    {{end}}
]
`

// endregion
