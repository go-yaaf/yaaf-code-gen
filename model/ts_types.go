package model

import (
	"fmt"
	"strings"
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
	//"Json":          "Map<string,object>",
	"Json":          "any",
	"StreamContent": "File",
}

// GetTsType - convert variables types to known TypeScript types
func GetTsType(pType string) string {
	// Handle generics
	if strings.Contains(pType, "[") {
		return getGenericTsType(pType)
	}
	if _, ok := tsTypes[pType]; ok {
		return tsTypes[pType]
	}
	return pType
}

// GetGenericTsType - convert variables generics types to known TypeScript types
func getGenericTsType(pType string) string {
	// Extract type and index
	start := strings.Index(pType, "[")
	end := strings.Index(pType, "]")
	x := pType[0:start]
	y := pType[start+1 : end]

	xt := GetTsType(x)
	yt := GetTsType(y)

	return fmt.Sprintf("%s<%s>", xt, yt)
}

// Check if the provided tsType is in the list of primitive types
func isNativeType(tsType string) (isNative bool, arr string) {

	isNative = true
	arr = ""

	if strings.HasPrefix(tsType, "[]") {
		tsType = tsType[2:]
		arr = "[]"
	}

	for k, v := range tsTypes {
		if tsType == k || tsType == v {
			isNative = true
			return
		}
	}

	if strings.HasPrefix(tsType, "Map<") {
		isNative = true
		return
	}

	if strings.ToLower(tsType) == "streamcontent" {
		isNative = true
		return
	}

	isNative = false
	return
}
