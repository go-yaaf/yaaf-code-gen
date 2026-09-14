package processor_ng

// region TypeScript index file template -------------------------------------------------------------------------------

var indexTsTemplate = `
{{range .}}export * from './{{.}}';
{{end}}

`

// endregion
