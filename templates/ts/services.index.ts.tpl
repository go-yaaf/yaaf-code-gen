{{range .}}import { {{.Name}} } from './{{.TsName}}';
{{end}}
export const Services = [
    {{range .}}{{.Name}},
    {{end}}
]