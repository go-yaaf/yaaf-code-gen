package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-yaaf/yaaf-common/entity"
)

// StringKeyValue is a key-value pair of strings
type StringKeyValue entity.Tuple[string, string]

// region Meta Model structure -----------------------------------------------------------------------------------------

// MetaModel is the meta description of all types and services
type MetaModel struct {
	Packages map[string]*PackageInfo
}

func NewMetaModel() *MetaModel {
	return &MetaModel{
		Packages: make(map[string]*PackageInfo),
	}
}

// GetPackage get package by name or create one if not exists
func (m *MetaModel) GetPackage(name string) *PackageInfo {
	if len(name) == 0 {
		name = "model"
	}

	// Get package name
	if pkg, ok := m.Packages[name]; !ok {
		pkg = NewPackageInfo(name)
		m.Packages[name] = pkg
		return pkg
	} else {
		return pkg
	}
}

// AddClassInfo add new class to the model
func (m *MetaModel) AddClassInfo(ci *ClassInfo) {
	pkg := m.GetPackage(ci.PackageFullName)
	pkg.Classes[ci.Name] = ci
}

// AddAlias add new type alias
func (m *MetaModel) AddAlias(packageName string, alias, name string) {
	pkg := m.GetPackage(packageName)
	pkg.AddAlias(alias, name)
}

// AddEnumInfo add new class to the model
func (m *MetaModel) AddEnumInfo(ei *EnumInfo) {
	pkg := m.GetPackage(ei.PackageFullName)
	pkg.Enums[ei.Name] = ei
}

// AddServiceInfo add new service to the model
func (m *MetaModel) AddServiceInfo(si *ServiceInfo) {
	pkg := m.GetPackage(si.PackageFullName)
	pkg.Services[si.Name] = si
}

// GetEnum look for the enum by name in all the packages
func (m *MetaModel) GetEnum(name string) *EnumInfo {
	for _, pkg := range m.Packages {
		for key, val := range pkg.Enums {
			if key == name {
				return val
			}
		}
	}
	return nil
}

// GetService look for the service by name in all the packages
func (m *MetaModel) GetService(name string) *ServiceInfo {
	for _, pkg := range m.Packages {
		for key, val := range pkg.Services {
			if key == name {
				return val
			}
		}
	}
	return nil
}

func (m *MetaModel) String() string {
	if bytes, err := json.MarshalIndent(m, "", "    "); err != nil {
		return err.Error()
	} else {
		return string(bytes)
	}
}

// FillDependencies fill class dependencies
func (m *MetaModel) FillDependencies() {
	for _, pkg := range m.Packages {
		pkg.fillDependencies(m)
	}
}

func (m *MetaModel) ReplaceAliases() {
	for _, pkg := range m.Packages {
		pkg.replaceAliases(m)
	}
}

// GetAllClassFields return an array of field info for all class fields
func (m *MetaModel) GetAllClassFields(className string) []*FieldInfo {
	result := make([]*FieldInfo, 0)
	m.getAllClassFields(className, &result)

	// fill default value
	for _, fi := range result {
		fi.DefaultValue = m.getDefaultValue(fi.TsType)
	}
	return result
}

// getAllClassFields internal recursive function to return an array of field info for all class fields
func (m *MetaModel) getAllClassFields(className string, arr *[]*FieldInfo) {
	ci := m.FindClass(className)
	if ci == nil {
		return
	}
	// Get parent class fields
	if len(ci.BaseClass) > 0 {
		m.getAllClassFields(ci.BaseClass, arr)
	}
	for _, field := range ci.Fields {
		*arr = append(*arr, field)
	}
}

// GetFactoryMethods accepts type list and return factory methods for all types that are complex classes
func (m *MetaModel) GetFactoryMethods(list []string) []string {
	result := make([]string, 0)

	for _, className := range list {
		if ci := m.FindClass(className); ci != nil {
			result = append(result, fmt.Sprintf("New%s", className))
		}
	}

	return result
}

// FindClass find class info in all the packages
func (m *MetaModel) FindClass(className string) *ClassInfo {
	for _, pkg := range m.Packages {
		for _, ci := range pkg.Classes {
			if ci.Name == className {
				return ci
			}
		}
	}
	return nil
}

func (m *MetaModel) getDefaultValue(tsType string) string {
	// First find primitive types
	if tsType == "string" {
		return "''"
	}
	if tsType == "number" {
		return "0"
	}
	if tsType == "boolean" {
		return "false"
	}
	if tsType == "any" || tsType == "Json" {
		return "{}"
	}
	if ei := m.GetEnum(tsType); ei != nil {
		return "0"
	}
	if ci := m.FindClass(tsType); ci != nil {
		return fmt.Sprintf("New%s()", tsType)
	}
	return ""
}

// endregion

// region Internal helper functions ------------------------------------------------------------------------------------

// SmallCaps Convert big Caps to small Caps
func SmallCaps(name string) string {
	if len(name) > 0 {
		return fmt.Sprintf("%s%s", strings.ToLower(name[0:1]), name[1:])
	} else {
		return ""
	}
}

// Title converts name to Title (first letter upper Caps)
func Title(name string) string {
	if len(name) > 0 {
		return fmt.Sprintf("%s%s", strings.ToUpper(name[0:1]), name[1:])
	} else {
		return ""
	}
}

// endregion
