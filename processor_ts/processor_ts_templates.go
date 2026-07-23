package processor_ts

// region TypeScript index file template -------------------------------------------------------------------------------

var indexTsTemplate = `
{{range .}}export * from './{{.}}';
{{end}}

`

// endregion
