package processor_html

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path"

	. "github.com/go-yaaf/yaaf-code-gen/model"
)

// Generate types index file
func (p *HtmlProcessor) generateTypesIndexHtml(root *template.Template) {

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

// Generate services documentation files
func (p *HtmlProcessor) generateTypesHtml(root *template.Template) {

	// Create the index
	p.generateTypesIndexHtml(root)

	// create each type
	list := p.Model.ListClasses()
	for _, classInfo := range list {

		fileName := path.Join(p.Output, fmt.Sprintf("type_%s.html", classInfo.Name))

		data := NewTemplateData(p.ApiName, p.ApiVersion)
		data.ClassInfo = classInfo
		data.ClassGroups = p.Model.ListClassGroups()
		//data.ClassGroups = p.Model.ListClassGroups(classInfo.Name)

		if f, err := os.Create(fileName); err != nil {
			log.Fatal("Error creating file: ", fileName, err)
		} else if err = root.ExecuteTemplate(f, "TYPE_TMPL", data); err != nil {
			panic(err)
		} else {
			_ = f.Close()
		}
	}
}
