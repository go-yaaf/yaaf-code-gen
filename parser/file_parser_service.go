package parser

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/go-yaaf/yaaf-code-gen/model"
)

// process entity type
func (p *FileParser) processServiceType(ti *model.TypeInfo, decl *ast.GenDecl) error {
	if len(decl.Specs) < 1 {
		return fmt.Errorf("no specs found")
	}

	switch spec := decl.Specs[0].(type) {
	case *ast.TypeSpec:
		break
	case *ast.ImportSpec:
		return nil
	default:
		return fmt.Errorf("unknown spec type %T", spec)
	}

	// At this point, it is known that spec is of type ast.TypeSpec
	si := model.NewServiceInfo(ti.Name, ti.Docs...)
	si.PackageFullName = ti.PackageFullName
	si.PackageShortName = ti.PackageShortName
	si.TsName = ti.TsName
	si.Headers = ti.Headers
	si.Context = ti.Context
	si.Group = ti.Group
	si.Path = ti.Path

	// Add class to model
	p.Model.AddServiceInfo(si)
	return nil
}

// process service method
func (p *FileParser) processServiceMethod(decl *ast.FuncDecl) error {
	if decl.Recv == nil {
		return fmt.Errorf("no reciever fount")
	}

	starExp, ok := decl.Recv.List[0].Type.(*ast.StarExpr)
	if !ok {
		return nil
	}

	serviceName := ""
	switch ident := starExp.X.(type) {
	case *ast.Ident:
		serviceName = ident.Name
	case *ast.IndexExpr:
		// do nothing
	case *ast.IndexListExpr:
		//fmt.Println("processServiceMethod: starExp.X is *ast.IndexListExpr, what to do?")
	case *ast.SelectorExpr:
		//fmt.Println("processServiceMethod: starExp.X is *ast.SelectorExpr, what to do?")
	default:
		//fmt.Println("processServiceMethod: starExp.X is not a type of the above, what to do?")
	}

	if len(serviceName) == 0 {
		return nil
	}

	si := p.Model.GetService(serviceName)
	if si == nil {
		return fmt.Errorf("service %s not found", serviceName)
	}

	if decl.Doc == nil {
		return nil
	}

	p.processServiceMethodComments(si, decl.Name.Name, decl.Doc.List)

	return nil
}

// Process service endpoint comments and extract tags to enrich service class. The following tags are expected:
// @InheritFrom: - the field type is the parent class
// @Json: - the json name of the field
func (p *FileParser) processServiceMethodComments(si *model.ServiceInfo, name string, comments []*ast.Comment) {

	mi := model.NewMethodInfo(model.Title(name))

	for _, comment := range comments {
		line := p.trimComment(comment.Text)
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "@Http") {
			action := p.getTagValue(line, "@Http:")
			mi.SetAction(action)
		} else if strings.HasPrefix(line, "@Context") {
			mi.Context = p.getTagValue(line, "@Context:")
		} else if strings.HasPrefix(line, "@Return:") {
			returnClass := p.getTagValue(line, "@Return:")
			mi.Return = model.NewClassInfo(returnClass)
			mi.SetReturnType(returnClass)
		} else if strings.HasPrefix(line, "@PathParam") {
			mi.AddPathParam(p.getTagValue(line, "@PathParam:"))
		} else if strings.HasPrefix(line, "@QueryParam") {
			mi.AddQueryParam(p.getTagValue(line, "@QueryParam:"))
		} else if strings.HasPrefix(line, "@BodyParam") {
			mi.AddBodyParam(p.getTagValue(line, "@BodyParam:"))
		} else if strings.HasPrefix(line, "@FileParam") {
			mi.AddFileParam(p.getTagValue(line, "@FileParam:"))
		} else if strings.HasPrefix(line, "@Upload") {
			functionName := p.getTagValue(line, "@Upload:")
			mi.SetUploadFunction(functionName)
		} else {
			mi.Docs = append(mi.Docs, line)
		}
	}

	// Add only REST service methods (that Method name is not empty
	if mi != nil {
		if len(mi.Method) > 0 {
			si.Methods = append(si.Methods, mi)
		}
	}
}
