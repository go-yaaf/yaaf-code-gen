package processor_html

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path"

	. "github.com/go-yaaf/yaaf-code-gen/model"
)

// Generate services documentation files
func (p *HtmlProcessor) generateServicesHtml(root *template.Template) {

	// serviceGroups := p.Model.ListServiceGroups("")
	list := p.Model.ListServices()
	for _, serviceInfo := range list {

		fileName := path.Join(p.Output, fmt.Sprintf("service_%s.html", serviceInfo.Name))

		data := NewTemplateData(p.ApiName, p.ApiVersion)
		data.ServiceInfo = serviceInfo
		data.ServiceGroups = p.Model.ListServiceGroups(serviceInfo.Name)

		if f, err := os.Create(fileName); err != nil {
			log.Fatal("Error creating file: ", fileName, err)
		} else if err = root.ExecuteTemplate(f, "SERVICE_TMPL", data); err != nil {
			panic(err)
		} else {
			_ = f.Close()
		}
	}
}
