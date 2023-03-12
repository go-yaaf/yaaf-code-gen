package parser

import (
	"fmt"
	"go/ast"
	"strings"
)

// Process public function
func (p *Parser) processFuncDecl(decl *ast.FuncDecl, pi *PackageInfo) {

	name := decl.Name.String()
	Name := fmt.Sprintf("%s%s", strings.ToLower(name[0:1]), name[1:])

	mi := &MethodInfo{
		Name:              Name,
		Method:            "",
		Path:              "",
		Docs:              nil,
		Headers:           nil,
		PathParams:        nil,
		QueryParams:       nil,
		BodyParam:         nil,
		FileParam:         nil,
		StreamsRequest:    false,
		Return:            nil,
		Context:           "",
		IsSocketMessage:   false,
		SocketMessageType: "",
	}
	parseFunctionComments(decl.Doc, mi)
	pi.AddServiceMethod(mi)
}

// Build documentation from comments group and enrich metadata
func parseFunctionComments(cg *ast.CommentGroup, mi *MethodInfo) {

	if cg == nil {
		return
	}

	for _, c := range cg.List {
		text := strings.Trim(strings.ReplaceAll(c.Text, "//", ""), " ")
		updateMethodInfo(text, mi)
	}
}

// Analyze comment line and update method info flags
func updateMethodInfo(text string, mi *MethodInfo) {
	if parseHttpTag(text, mi) {
		return
	} else if parseContextTag(text, mi) {
		return
	} else if parseParamTag("@QueryParam:", text, mi) {
		return
	} else if parseParamTag("@PathParam:", text, mi) {
		return
	} else if parseParamTag("@BodyParam:", text, mi) {
		return
	} else if parseReturnTag(text, mi) {
		return
	} else {
		mi.Docs = append(mi.Docs, text)
	}
}

// parse line with @Http prefix, format is: @Http <name> [<type>]: <description>
func parseHttpTag(text string, mi *MethodInfo) bool {
	if idx := strings.Index(text, "@Http:"); idx > -1 {
		http := strings.Trim(text[idx+len("@Http:"):], " ")
		parts := strings.Split(http, " ")
		mi.Method = strings.Trim(parts[0], " ")
		if len(parts) > 1 {
			mi.Path = strings.Trim(parts[1], " ")
		}
		return true
	} else {
		return false
	}
}

// parse line with @Context prefix, format is: @Context <name>
func parseContextTag(text string, mi *MethodInfo) bool {
	if idx := strings.Index(text, "@Context:"); idx > -1 {
		mi.Context = strings.Trim(text[idx+len("@Context:"):], " ")
		return true
	} else {
		return false
	}
}

// parse line with @QueryParam / @PathParam / @BodyParam tad, format is: @tagType: <name> | <type> | <description>
func parseParamTag(tagType, text string, mi *MethodInfo) bool {
	idx := strings.Index(text, tagType)
	if idx < 0 {
		return false
	}
	query := strings.Trim(text[idx+len(tagType):], " ")
	params := strings.Split(query, "|")
	paramName := strings.Trim(params[0], " ")
	paramType := "UNKNOWN"
	if len(params) > 1 {
		paramType = strings.Trim(params[1], " ")
	}
	paramDesc := "UNKNOWN"
	if len(params) > 2 {
		paramDesc = strings.Trim(params[2], " ")
	}

	pi := parseParamType(paramName, paramType)
	if pi != nil {
		pi.Sequence = len(mi.QueryParams)
		pi.Docs = []string{paramDesc}
		pi.ParamType = tagType[1 : len(tagType)-1]
		mi.QueryParams = append(mi.QueryParams, pi)
	}
	return true
}

// parse param type line and create ParamInfo
func parseParamType(paramName, paramType string) *ParamInfo {

	pi := &ParamInfo{
		Name:    paramName,
		Json:    paramName,
		Type:    paramType,
		IsArray: false,
		IsMap:   false,
		MapKey:  "",
	}
	// handle array
	if idx := strings.Index(paramType, "[]"); idx > -1 {
		pi.ParamType = paramType[idx+len("[]"):]
		pi.IsArray = true
	}

	// handle map
	if idx := strings.Index(paramType, "map["); idx > -1 {
		if idz := strings.Index(paramType, "]"); idz > -1 {
			pi.ParamType = paramType[idz+1:]
			pi.MapKey = paramType[idx+len("map[") : idz]
			pi.IsMap = true
		}
	}
	return pi
}

// parse line with @Return prefix, format is: @Return <type>
func parseReturnTag(text string, mi *MethodInfo) bool {
	idx := strings.Index(text, "@Return:")
	if idx < 0 {
		return false
	}
	ret := strings.Trim(text[idx+len("@Return:"):], " ")

	ret = strings.ReplaceAll(ret, "<", "[")
	ret = strings.ReplaceAll(ret, ">", "]")

	typ := ret
	gen := ""
	if idy := strings.Index(ret, "["); idy > -1 {
		if idz := strings.Index(ret, "]"); idz > -1 {
			typ = ret[idy+1 : idz]
			gen = ret[:idy]
		}
	}

	mi.Return = &ReturnInfo{
		Name:        ret,
		Type:        typ,
		GenericType: gen,
		IsArray:     false,
		IsMap:       false,
		MapKey:      "",
		Docs:        nil,
	}
	return true
}
