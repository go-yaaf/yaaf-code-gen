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

// region TS template Classes Processor --------------------------------------------------------------------------------

// Add class constructor method
func addClassConstructor(class model.ClassInfo) string {
	output := "    constructor("
	for _, field := range class.Fields {
		output += field.TsName + "?: " + field.TsType
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

// Add class imports based on the class dependencies
func addClassImports(class model.ClassInfo) string {
	output := ""
	for className, _ := range class.Dependencies {
		output += fmt.Sprintf("import { %s } from './%s';\n", className, className)
	}
	return output
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

// Generate classes
func (p *TsProcessor) handleTsClasses() {
	funcMap := template.FuncMap{
		"getTsType":      getTsType,
		"addImports":     addClassImports,
		"addConstructor": addClassConstructor,
		"join":           strings.Join,
		"genericsParam":  genericsParam,
	}

	var classList []model.ClassInfo
	for _, v := range p.Model.Packages {
		for _, class := range v.Classes {
			classList = append(classList, *class)
		}
	}

	folder := path.Join(p.Output, "model")
	p.makeDir(folder)

	tmpl, _ := template.New("base_class.ts.tpl").Funcs(funcMap).Parse(classTsTemplate)
	for _, class := range classList {

		// For parameter classes, do not create TS file
		if !class.IsParam {

			var tpl bytes.Buffer
			if err := tmpl.Execute(&tpl, class); err != nil {
				log.Fatal("Error executing template [base_class.ts.tpl]: ", err)
			}
			// Remove newlines
			processedContent := p.trimNewLines(tpl.String())

			fileName := path.Join(folder, fmt.Sprintf("%s.ts", class.Name))
			if f, err := os.Create(fileName); err != nil {
				log.Fatal("Error creating file: ", fileName, err)
			} else {
				if _, err = f.WriteString(processedContent); err != nil {
					log.Fatal("Error writing to file: ", fileName, err)
				}
				_ = f.Close()
			}
		}
	}

	// Create the enums and classes index file
	var list []string

	for _, v := range p.Model.Packages {
		for _, enm := range v.Enums {
			list = append(list, enm.Name)
		}
		for _, class := range v.Classes {
			list = append(list, class.Name)
		}

	}
	p.generateIndexTs(list, folder)
}

// endregion

// region TypeScript class file template -------------------------------------------------------------------------------

var classTsTemplate = `
{{. | addImports}}

{{range .Docs}}
// {{.}}{{end}}
export class {{.Name}}{{. | genericsParam }}{{template "extend" .}} {
{{range .Fields}}
	// {{range .Docs}}{{.}} {{end}}
	public {{.Json}}: {{.TsType }}{{ if .IsArray }}[]{{ end }};
{{end}}
{{ if not .IsExtend }}{{. | addConstructor }}{{end}}

{{ if eq .Name "BaseEntityEx" }}

{{end}}
}

{{ if .IsExtend }}{{template "getColumnDef" .}}{{end}}

{{define "extend"}}{{ if .IsExtend }} extends {{.BaseClass}}{{ end }}{{end}}


{{define "getColumnDef"}}
export function Get{{.Name}}ColumnsDef() : ColumnDef[] {
    let result : ColumnDef[] = [];
	result.push(new ColumnDef("", "id", "string", ""));
	result.push(new ColumnDef("", "createdOn", "number", "datetime"));
	result.push(new ColumnDef("", "updatedOn", "number", "datetime"));
	{{range .Fields}}result.push(new ColumnDef("", "{{.Json}}", "{{.TsType}}", "{{.Format}}"));
	{{end}}

	return result;
}
{{end}}
`

// endregion
