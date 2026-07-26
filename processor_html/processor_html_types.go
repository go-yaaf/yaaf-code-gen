package processor_html

import (
	"html/template"
	"log"
	"os"
	"path"

	. "github.com/go-yaaf/yaaf-code-gen/model"
)

// Generate types index file
func (p *HtmlProcessor) generateTypesHtml(root *template.Template) {

	fileName := path.Join(p.Output, "types.html")

	data := NewTemplateData(p.ApiName, p.ApiVersion)
	data.ClassGroups = p.Model.ListClassGroups()

	if f, err := os.Create(fileName); err != nil {
		log.Fatal("Error creating file: ", fileName, err)
	} else if err = root.ExecuteTemplate(f, "TYPES_TMPL", data); err != nil {
		panic(err)
	} else {
		_ = f.Close()
	}
}
