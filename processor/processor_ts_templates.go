package processor

// region TypeScript index file template -------------------------------------------------------------------------------

var indexTsTemplate = `
{{range .}}export * from './{{.}}';
{{end}}

`

// endregion
