package processor

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/go-yaaf/yaaf-code-gen/model"
)

var tsTypes = map[string]string{
	"double":    "number",
	"float":     "number",
	"float32":   "number",
	"float64":   "number",
	"int":       "number",
	"int32":     "number",
	"int64":     "number",
	"uint":      "number",
	"uint32":    "number",
	"uint64":    "number",
	"sint":      "number",
	"sint32":    "number",
	"sint64":    "number",
	"fixed32":   "number",
	"fixed64":   "number",
	"sfixed32":  "number",
	"sfixed64":  "number",
	"bool":      "boolean",
	"string":    "string",
	"bytes":     "File",
	"any":       "any",
	"Timestamp": "number",
	"Json":      "Record<string,object>",
}

// TsProcessor - TS processor converts proto files to TypeScript files
type TsProcessor struct {
	BaseProcessor
}

// NewTsProcessor - Factory method
func NewTsProcessor(model *model.MetaModel, output string) Processor {
	return &TsProcessor{BaseProcessor{
		Output: output,
		Model:  model,
	}}
}

// var classPackageMap = make(map[string]string)

// Start the processor
func (p *TsProcessor) Start() error {

	// Generate all enums
	p.handleTsEnums()

	// Generate all classes
	p.handleTsClasses()

	// Generate all services
	p.handleTsServices()

	// Generate service exports
	//p.generateServicesExports()

	// Generate all index.ts files (barrels)
	//p.generateIndexes()
	return nil
}

// create directory
func (p *TsProcessor) makeDir(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatal("Error creating folder: "+path, err)
	}
}

func toCamelCase(s string) string {
	return fmt.Sprintf("%s%s", strings.ToLower(s[0:1]), s[1:])
}

func toCapitalCase(s string) string {
	return fmt.Sprintf("%s%s", strings.ToUpper(s[0:1]), strings.ToLower(s[1:]))
}

func toDisplayName(s string) string {
	parts := strings.Split(s, "_")
	caps := make([]string, 0)
	for _, p := range parts {
		caps = append(caps, toCapitalCase(p))
	}
	return strings.Join(caps, " ")
}

// Generate indexes
/*
func (p *TsProcessor) generateIndexes() {
	var content []string

	for _, pkg := range p.Model.Packages {
		for en := range pkg.Enums {
			content = append(content, fmt.Sprintf("%s", en))
		}
		for cn := range pkg.Classes {
			content = append(content, fmt.Sprintf("%s", cn))
		}
		for sn := range pkg.Services {
			content = append(content, fmt.Sprintf("%s", sn))
		}
	}

	p.generateIndexTs(content)
}
*/

// getTsType - convert variables types to known TypeScript types
func getTsType(pType string) string {

	if !strings.Contains(pType, ",") {
		return getTsTypeInternal(pType)
	}

	// In case of generic type with comma-separated list of arguments, (e.g. EntityResponse<Distribution<FlightClassCode,float64>>)
	// Split arguments, convert each one and join them back
	typeList := strings.Split(pType, ",")
	tsTypeList := make([]string, 0)
	for _, t := range typeList {
		tsTypeList = append(tsTypeList, getTsTypeInternal(t))
	}
	return strings.Join(tsTypeList, ",")
}

// getTsTypeInternal - convert variables types to known TypeScript types
func getTsTypeInternal(pType string) string {

	pType = strings.TrimSpace(pType)

	// Handle maps
	if strings.Contains(pType, "map[") {
		return getGenericTsMap(pType)
	}

	// Handle generics
	if strings.Contains(pType, "[") && strings.Contains(pType, "]") {
		return getGenericTsType(pType)
	}

	if _, ok := tsTypes[pType]; ok {
		return tsTypes[pType]
	}
	return pType
}

// convert variables generics types to known TypeScript types
func getGenericTsType(pType string) string {
	// Extract type and index
	start := strings.Index(pType, "[")
	end := strings.Index(pType, "]")
	x := pType[0:start]
	xt := getTsType(x)

	y := pType[start+1 : end]
	yt := getTsType(y)

	// Handle multiple generics arguments
	if strings.Contains(y, ",") {
		args := strings.Split(y, ",")
		tsArgs := make([]string, 0)
		for _, a := range args {
			tsArgs = append(tsArgs, getTsType(a))
		}
		yt = strings.Join(tsArgs, ", ")
	}

	return fmt.Sprintf("%s<%s>", xt, yt)
}

// convert Map generics types to known TypeScript types
func getGenericTsMap(pType string) string {

	// Extract type and index
	start := strings.Index(pType, "[")
	end := strings.Index(pType, "]")
	x := pType[start+1 : end]
	xt := getTsType(x)

	y := pType[end+1:]
	yt := getTsType(y)

	return fmt.Sprintf("Record<%s,%s>", xt, yt)
}

// Generate TypeScript index
func (p *TsProcessor) generateIndexTs(data []string, folder string) {
	tmpl, _ := template.New("index.ts.tpl").Parse(indexTsTemplate)
	fileName := path.Join(folder, "index.ts")

	if f, err := os.Create(fileName); err != nil {
		log.Fatal("Error creating file: ", fileName, err)
	} else if er := tmpl.Execute(f, data); er != nil {
		log.Fatal("Error executing template [index.ts.tpl]: ", er)
	} else {
		_ = f.Close()
	}
}

//func convertToTypeScript(name string) string {
//	tokens := strings.Split(name, "<")
//	types := make([]string, 0)
//	for _, token := range tokens {
//		types = append(types, strings.ReplaceAll(token, ">", ""))
//	}
//
//	for _, t := range types {
//		tsType := getTsType(t)
//		if tsType != t {
//			name = strings.ReplaceAll(name, t, tsType)
//		}
//	}
//
//	return name
//}
