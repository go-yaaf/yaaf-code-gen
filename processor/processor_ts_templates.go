package processor

// region TypeScript index file template -------------------------------------------------------------------------------

var indexTsTemplate = `
{{range .}}export * from './{{.}}';
{{end}}
export * from './Enums';
`

// endregion

// region TypeScript class file template -------------------------------------------------------------------------------

var classTsTemplate = `
{{. | addImports}}

{{range .Docs}}
// {{.}}{{end}}
export class {{.Name}}{{. | genericsParam }}{{template "extend" .}} {
 {{range .Fields}}
    // {{range .Docs}}{{.}} {{end}}
    public {{.Json}}: {{.Type | getTsType }}{{ if .IsArray }}[]{{ end }};
 {{end}}
 {{ if not .IsExtend }}{{. | addConstructor }}{{end}}

 {{ if eq .BaseClass "BaseEntityEx" }}
	get(field: string) : any {
		if (!this.props) {
			return "";
		}
		let val = this.props.get(field);
		return String(val);
	}
 {{end}}
}

{{define "extend"}}{{ if .IsExtend }} extends {{.BaseClass}}{{ end }}{{end}}
`

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

    {{range .Values}}{{ if ne .Name "UNDEFINED" }}
	result.push(new Tuple<{{$.Name}}, string>({{$.Name}}.{{.Name}}, '{{$.Name}}.{{.Name}}'))
	{{ end }}{{end}}

    return result;
}

// Return map of {{.Name}} values and their display names
export function Map{{.Name}}s() : Map<{{.Name}}, string> {
    let result = new Map<{{.Name}}, string>();

    {{range .Values}}
	result.set({{$.Name}}.{{.Name}}, '{{.Name | toDisplayName}}');
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

// region TypeScript service file template -------------------------------------------------------------------------------

var serviceTsTemplate = `
import { Injectable, Inject } from '@angular/core';
import { RestUtil } from '../utils';
import { GooxConfig } from '../config';

{{.Methods | addServiceImports}}

{{range .Docs}}
// {{.}} {{end}}
@Injectable()
export class {{.TsName}} {

  // URL to web api
  private baseUrl = '{{.Path}}';

  // Class constructor
  constructor(@Inject('config') private config: GooxConfig, private rest: RestUtil) {
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
