package processor_html

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path"

	. "github.com/go-yaaf/yaaf-code-gen/model"
)

// Generate enums index file
func (p *HtmlProcessor) generateEnumsIndexHtml(root *template.Template) {

	fileName := path.Join(p.Output, "enums.html")

	data := NewTemplateData(p.ApiName, p.ApiVersion)
	data.EnumList = p.Model.ListEnums()

	if f, err := os.Create(fileName); err != nil {
		log.Fatal("Error creating file: ", fileName, err)
	} else if err = root.ExecuteTemplate(f, "ENUMS_TMPL", data); err != nil {
		panic(err)
	} else {
		_ = f.Close()
	}
}

// Generate services documentation files
func (p *HtmlProcessor) generateEnumsHtml(root *template.Template) {

	// Create the index
	p.generateEnumsIndexHtml(root)

	// create each type
	list := p.Model.ListEnums()
	for _, enumInfo := range list {

		fileName := path.Join(p.Output, fmt.Sprintf("type_%s.html", enumInfo.Name))

		data := NewTemplateData(p.ApiName, p.ApiVersion)
		data.EnumList = p.Model.ListEnums()
		data.EnumInfo = enumInfo

		if f, err := os.Create(fileName); err != nil {
			log.Fatal("Error creating file: ", fileName, err)
		} else if err = root.ExecuteTemplate(f, "ENUM_TMPL", data); err != nil {
			panic(err)
		} else {
			_ = f.Close()
		}
	}
}
