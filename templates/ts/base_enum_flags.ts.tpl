/* {{range .Docs}}
   {{.}}{{end}} 
*/
export enum {{.Name}} {
 {{range .Values}}
    // {{range .Docs}}{{.}} {{end}}
    {{.Name}} = {{.Shifter}} << {{.Value}},
 {{end}}
}