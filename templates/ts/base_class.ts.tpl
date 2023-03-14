{{. | addImports}}
/* {{range .Docs}}
   {{.}}{{end}} 
*/
export class {{.Name}}{{template "extend" .}} {
 {{range .Fields}}
    // {{range .Docs}}{{.}} {{end}}
    public {{.Json}}: {{.Type | getTsType }}{{ if .IsArray }}[]{{ end }};
 {{end}}
{{ if not .IsExtend }}{{. | addConstructor }}{{end}}
}

{{define "extend"}}{{ if .IsExtend }} extends {{.BaseClass}}{{ end }}{{end}}
