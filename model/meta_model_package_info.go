package model

// PackageInfo package information
type PackageInfo struct {
	Name     string                    // Package name
	Docs     []string                  // Package documentation
	Classes  map[string]*ClassInfo     // Map of classes in package
	Enums    map[string]*EnumInfo      // Map of enums in package
	Services map[string]*ServiceInfo   // Map of services in package
	Sockets  map[string]*WebSocketInfo // Map of web sockets
	Aliases  map[string]string         // Map of type aliases
}

func NewPackageInfo(name string) *PackageInfo {
	return &PackageInfo{
		Name:     name,
		Docs:     make([]string, 0),
		Classes:  make(map[string]*ClassInfo, 0),
		Enums:    make(map[string]*EnumInfo),
		Services: make(map[string]*ServiceInfo),
		Sockets:  make(map[string]*WebSocketInfo),
		Aliases:  make(map[string]string),
	}
}

// fill dependencies
func (p *PackageInfo) fillDependencies(mm *MetaModel) {
	for _, pkg := range mm.Packages {
		for _, ci := range pkg.Classes {
			ci.fillDependencies(mm)
		}
		for _, si := range pkg.Services {
			si.fillDependencies(mm)
		}
	}
}

// AddAlias add entry to aliases
func (p *PackageInfo) AddAlias(alias, name string) {
	p.Aliases[alias] = name
}

// ReplaceAliases replace all aliases
func (p *PackageInfo) replaceAliases(mm *MetaModel) {
	for _, pkg := range mm.Packages {
		//for _, ci := range pkg.Classes {
		//	ci.replaceAliases(pkg)
		//}
		for _, si := range pkg.Services {
			si.replaceAliases(pkg)
		}
	}
}
