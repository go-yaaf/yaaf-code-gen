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
