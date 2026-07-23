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

	list := p.Model.ListEnums()
	data := struct {
		Name    string
		Version string
		Enums   []*EnumInfo
	}{
		Name:    p.ApiName,
		Version: p.ApiVersion,
		Enums:   list,
	}
	if f, err := os.Create(fileName); err != nil {
		log.Fatal("Error creating file: ", fileName, err)
	} else if err = root.ExecuteTemplate(f, "ENUMS_TMPL", data); err != nil {
		panic(err)
	} else {
		_ = f.Close()
	}
}
