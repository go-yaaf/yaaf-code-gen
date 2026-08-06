package processor_html

import (
	"html/template"
	"log"
	"os"
	"path"

	. "github.com/go-yaaf/yaaf-code-gen/model"
	. "github.com/go-yaaf/yaaf-code-gen/processor"
	. "github.com/go-yaaf/yaaf-code-gen/templates/docs"
)

var htmlProcessorInst *HtmlProcessor = nil

// HtmlProcessor - HTML processor converts proto files to HTML files for documentation
type HtmlProcessor struct {
	BaseProcessor
}

// NewHtmlProcessor - Factory method
func NewHtmlProcessor(model *MetaModel, output, apiName, apiVersion string) Processor {
	htmlProcessorInst = &HtmlProcessor{BaseProcessor{
		Output:     output,
		Model:      model,
		ApiName:    apiName,
		ApiVersion: apiVersion,
	}}
	return htmlProcessorInst
}

func GetHtmlProcessor() *HtmlProcessor {
	return htmlProcessorInst
}

// var classPackageMap = make(map[string]string)

// Start the processor
func (p *HtmlProcessor) Start() error {

	// Build and get root template
	folder := path.Join(p.Output, "")
	p.makeDir(folder)

	root := p.getRootTemplate()

	// Generate styles
	p.generateStylesheet(root)

	// Generate root Index
	p.generateIndexHtml(root)

	// Generate all enums
	p.generateEnumsHtml(root)

	// Generate all classes
	p.generateTypesHtml(root)

	// Generate all services
	p.generateServicesHtml(root)

	return nil
}

// create directory
func (p *HtmlProcessor) makeDir(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatal("Error creating folder: "+path, err)
	}
}

// Get the root template
func (p *HtmlProcessor) getRootTemplate() *template.Template {
	// 1. Create the root template container
	tmpl := template.New("base")

	// 2. Define a map of template names to their string content
	templatesToParse := map[string]string{
		"INDEX_TMPL":        INDEX_TMPL,
		"STYLES_TMPL":       STYLES_TMPL,
		"STYLES_FULL_TMPL":  STYLES_FULL_TMPL,
		"HEAD_TMPL":         HEAD_TMPL,
		"HEADER_TMPL":       HEADER_TMPL,
		"ENUMS_TMPL":        ENUMS_TMPL,
		"TYPES_TMPL":        TYPES_TMPL,
		"TYPE_TMPL":         TYPE_TMPL,
		"FIELDS_TMPL":       FIELDS_TMPL,
		"TYPE_FIELDS_TMPL":  TYPE_FIELDS_TMPL,
		"TYPE_EXAMPLE_TMPL": TYPE_EXAMPLE_TMPL,
		"SERVICE_TMPL":      SERVICE_TMPL,
		"METHOD_TMPL":       METHOD_TMPL,
		"EXAMPLE_TMPL":      EXAMPLE_TMPL,
		"PARAM_TMPL":        PARAM_TMPL,
		"PARAMS_TMPL":       PARAMS_TMPL,
	}

	// 3. Loop through and parse each variable into the template tree
	for name, content := range templatesToParse {
		// New() associates the new template name with the root tmpl object
		_, err := tmpl.New(name).Parse(content)
		if err != nil {
			panic(err) // Handle your error properly in production
		}
	}

	return tmpl
}

// Generate main index file
func (p *HtmlProcessor) generateIndexHtml(root *template.Template) {

	fileName := path.Join(p.Output, "index.html")

	data := NewTemplateData(p.ApiName, p.ApiVersion)
	data.ServiceGroups = p.Model.ListServiceGroups("")

	if f, err := os.Create(fileName); err != nil {
		log.Fatal("Error creating file: ", fileName, err)
	} else if err = root.ExecuteTemplate(f, "INDEX_TMPL", data); err != nil {
		panic(err)
	} else {
		_ = f.Close()
	}
}

// Generate main stylesheet file
func (p *HtmlProcessor) generateStylesheet(root *template.Template) {

	fileName := path.Join(p.Output, "styles.css")
	data := NewTemplateData(p.ApiName, p.ApiVersion)

	if f, err := os.Create(fileName); err != nil {
		log.Fatal("Error creating file: ", fileName, err)
	} else if err = root.ExecuteTemplate(f, "STYLES_FULL_TMPL", data); err != nil {
		panic(err)
	} else {
		_ = f.Close()
	}
}
