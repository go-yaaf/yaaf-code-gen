package processor_html

import (
	"html/template"
	"log"
	"os"
	"path"

	. "github.com/go-yaaf/yaaf-code-gen/model"
)

// Generate enums index file
func (p *HtmlProcessor) generateEnumsHtml(root *template.Template) {

	fileName := path.Join(p.Output, "enums.html")

	data := NewTemplateData(p.ApiName, p.ApiVersion)
	data.EnumGroups = p.Model.ListEnumGroups()

	if f, err := os.Create(fileName); err != nil {
		log.Fatal("Error creating file: ", fileName, err)
	} else if err = root.ExecuteTemplate(f, "ENUMS_TMPL", data); err != nil {
		panic(err)
	} else {
		_ = f.Close()
	}
}
