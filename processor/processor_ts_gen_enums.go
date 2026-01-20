package processor

import (
	"fmt"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/go-yaaf/yaaf-code-gen/model"
)

// region TS template Enums Processor ----------------------------------------------------------------------------------

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

	folder := path.Join(p.Output, "model")
	p.makeDir(folder)

	var list []string

	tp := GetExternalTemplate("enum", enumTsTemplate, funcMap)
	tmpl, _ := template.New("base_enum.ts.tpl").Funcs(tp.FuncMap).Parse(tp.Template)
	for _, enum := range enumList {
		list = append(list, enum.Name)
		fileName := path.Join(folder, fmt.Sprintf("%s.ts", enum.Name))
		if f, err := os.Create(fileName); err != nil {
			log.Fatal("Error creating file: ", fileName, err)
		} else if er := tmpl.Execute(f, enum); er != nil {
			log.Fatal("Error executing template [base_enum.ts.tpl]: ", er)
		} else {
			_ = f.Close()
		}
	}

	// Create the enums index file
	//p.generateIndexTs(list, folder)
}

// Generate enums mapping
/*
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
	fileName := path.Join(p.Output, "model", "enums.ts")

	if f, err := os.Create(fileName); err != nil {
		log.Fatal("Error creating file: ", fileName, err)
	} else if er := tmpl.Execute(f, enumList); er != nil {
		log.Fatal("Error executing template [enum_map.ts.tpl]: ", er)
	} else {
		_ = f.Close()
	}
}
*/

// endregion

// region TypeScript enum file template -------------------------------------------------------------------------------

var enumTsTemplate = `
import { Tuple } from '.';

{{range .Docs}}
// {{.}}{{end}}
export enum {{.Name}} {
{{range .Values}}
    // {{range .Docs}}{{.}} {{end}}
    {{.Name}} = {{.Value}},
{{end}}
}

// Return list of {{.Name}} values and their display names
export function Get{{.Name}}s() : Tuple<{{.Name}}, string>[] {
	let result : Tuple<{{.Name}}, string>[] = [];
    {{range .Values}}{{ if ne .Name "UNDEFINED" }}result.push(new Tuple<{{$.Name}}, string>({{$.Name}}.{{.Name}}, '{{$.Name}}.{{.Name}}')){{ end }}
	{{end}}
    return result;
}

// Return map of {{.Name}} values and their display names
export function Map{{.Name}}s() : Map<{{.Name}}, string> {
    let result = new Map<{{.Name}}, string>();
    {{range .Values}}result.set({{$.Name}}.{{.Name}}, '{{.Name | toDisplayName}}');
	{{end}}
    return result;
}
`

var enumFlagsTsTemplate = `
// {{range .Docs}} {{.}} {{end}}
export enum {{.Name}} {
 {{range .Values}}
    // {{range .Docs}}{{.}} {{end}}
    {{.Name}} = {{.Shifter}} << {{.Value}},
 {{end}}
}
`
var enumMapTsTemplate = `
// Return a map of Enum values to enum strings
export function MapEnumValues() : Map<string, string> {

    let result: Map<string, string> = new Map<string, string>();
    {{range .}}{{ $dv := .Name }}{{range .Values}}
	result.set("{{$dv}}.{{.Value}}", "{{.Name}}");{{end}}{{end}}
    
	return result;
}
`

// endregion
