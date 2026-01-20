package processor

import (
	"text/template"
)

// TemplateParams is a template parameters
type TemplateParams struct {
	Template string
	FuncMap  template.FuncMap
}

var ExtTemplates = make(map[string]*TemplateParams)

func AddExternalTemplate(name, template string, funcMap template.FuncMap) {
	tp := &TemplateParams{Template: template, FuncMap: funcMap}
	ExtTemplates[name] = tp
}

func GetExternalTemplate(name string, defaultTemplate string, defaultMap template.FuncMap) *TemplateParams {
	if tmpl, ok := ExtTemplates[name]; ok {
		return tmpl
	} else {
		return &TemplateParams{Template: defaultTemplate, FuncMap: defaultMap}
	}
}
