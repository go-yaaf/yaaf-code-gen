package processor_ts

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/go-yaaf/yaaf-code-gen/model"
	"github.com/go-yaaf/yaaf-common/utils/collections"

	. "github.com/go-yaaf/yaaf-code-gen/processor"
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
	//types := make([]string, 0)
	//for className, _ := range class.Dependencies {
	//	types = append(types, className)
	//}
	//if len(types) == 0 {
	//	return ""
	//}
	//return fmt.Sprintf("import { %s } from '.';\n", strings.Join(types, ", "))

	output := ""
	for className, _ := range class.Dependencies {
		output += fmt.Sprintf("import { %s } from './%s';\n", className, className)
	}
	return output
}

// Add class factory functions imports based on the class dependencies
func addFactoriesImports(class model.ClassInfo) string {
	factories := make([]string, 0)

	// Get all fields and filter arrays
	p := GetTsProcessor()
	fields := p.Model.GetAllClassFields(class.Name)
	for _, field := range fields {
		if field.IsArray || field.IsMap || field.IsGeneric || !field.IsComplex {
			continue
		}

		if collections.Include[string]([]string{"string", "number", "boolean", "any", "json", "Record<string,any>"}, field.TsType) {
			continue
		}

		if ci := p.Model.FindClass(field.Type); ci != nil {
			if ci.IsGeneric == false {
				factories = append(factories, fmt.Sprintf("New%s", field.Type))
			}
		}

	}
	if len(factories) == 0 {
		return ""
	} else {
		reduced := collections.Distinct[string](factories)
		return fmt.Sprintf("import { %s } from '.';\n", strings.Join(reduced, ", "))
	}
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

func addNewInstance(class model.ClassInfo) string {
	if class.IsGeneric {
		return ""
	}

	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("export function New%s() : %s {\n", class.Name, class.Name))
	builder.WriteString(fmt.Sprintf("\tlet result : %s = new %s();\n", class.Name, class.Name))

	// Add fields
	p := GetTsProcessor()
	fields := p.Model.GetAllClassFields(class.Name)
	for _, field := range fields {
		if field.IsArray {
			builder.WriteString(fmt.Sprintf("\tresult.%s = [];\n", field.TsName))
		} else if len(field.DefaultValue) > 0 {
			builder.WriteString(fmt.Sprintf("\tresult.%s = %s;\n", field.Json, field.DefaultValue))
		}
	}

	builder.WriteString("\treturn result;\n")
	builder.WriteString("}")
	return builder.String()
}

// Generate classes
func (p *TsProcessor) handleTsClasses() {
	funcMap := template.FuncMap{
		"getTsType":      getTsType,
		"addImports":     addClassImports,
		"addImports2":    addFactoriesImports,
		"addConstructor": addClassConstructor,
		"join":           strings.Join,
		"genericsParam":  genericsParam,
		"addNewInstance": addNewInstance,
	}

	var classList []model.ClassInfo
	for _, v := range p.Model.Packages {
		for _, class := range v.Classes {
			classList = append(classList, *class)
		}
	}

	folder := path.Join(p.Output, "model")
	p.makeDir(folder)

	tp := GetExternalTemplate("class", classTsTemplate, funcMap)
	tmpl, _ := template.New("base_class.ts.tpl").Funcs(tp.FuncMap).Parse(tp.Template)
	for _, class := range classList {

		// For parameter classes, do not create TS file
		if !class.IsParam {

			var tpl bytes.Buffer
			if err := tmpl.Execute(&tpl, class); err != nil {
				log.Fatal("Error executing template [base_class.ts.tpl]: ", err)
			}
			// Remove newlines
			processedContent := p.TrimNewLines(tpl.String())

			safeName := SanitizeName(class.Name)
			if safeName == "" {
				log.Printf("skipping class with unsafe name: %q", class.Name)
				continue
			}
			fileName, err := ConfinedJoin(folder, fmt.Sprintf("%s.ts", safeName))
			if err != nil {
				log.Printf("skipping class %q: %v", class.Name, err)
				continue
			}
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

// Add this line only if it is required to add columns definitions
// {{ if .IsExtend }}{{template "getColumnDef" .}}{{end}}
var classTsTemplate = `
{{. | addImports}}
{{. | addImports2}}

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



// New empty instance
{{. | addNewInstance}}


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
