package parser

// create field
func newField(seq int, name, json, typ, doc string) *FieldInfo {
	return &FieldInfo{
		Name:     name,
		Json:     json,
		Type:     typ,
		Sequence: seq,
		Docs:     []string{doc},
	}
}

// create array field
func newArrField(seq int, name, json, typ, doc string) *FieldInfo {
	return &FieldInfo{
		Name:     name,
		Json:     json,
		Type:     typ,
		Sequence: seq,
		IsArray:  true,
		Docs:     []string{doc},
	}
}

// NewBaseEntityModel create BaseEntity metamodel
func NewBaseEntityModel() *ClassInfo {
	ci := &ClassInfo{
		ID:      "BaseEntity",
		Name:    "BaseEntity",
		Package: "yaaf-common",
		Docs:    []string{"Base class for all persistent entities"},
	}
	ci.Fields = append(ci.Fields, newField(0, "Id", "id", "string", "Unique object Id"))
	ci.Fields = append(ci.Fields, newField(1, "Key", "key", "string", "Shard (tenant) key"))
	ci.Fields = append(ci.Fields, newField(2, "CreatedOn", "createdOn", "Timestamp", "When the object was created [Epoch milliseconds Timestamp]"))
	ci.Fields = append(ci.Fields, newField(3, "UpdatedOn", "updatedOn", "Timestamp", "When the object was last updated [Epoch milliseconds Timestamp]"))

	return ci
}

// NewActionResponseModel create ActionResponse metamodel
func NewActionResponseModel() *ClassInfo {
	ci := &ClassInfo{
		ID:      "ActionResponse",
		Name:    "ActionResponse",
		Package: "messages",
		Docs:    []string{"Return message for any action on entity with no return data (e.d. delete)"},
	}
	ci.Fields = append(ci.Fields, newField(0, "Code", "code", "int", "Error code (0 for success)"))
	ci.Fields = append(ci.Fields, newField(1, "Error", "error", "string", "Error message"))
	ci.Fields = append(ci.Fields, newField(2, "Key", "key", "string", "The entity key (Id)"))
	ci.Fields = append(ci.Fields, newField(3, "Data", "data", "string", "Additional data"))

	return ci
}

// NewEntityResponseModel create NewEntityResponseModel metamodel
func NewEntityResponseModel() *ClassInfo {
	ci := &ClassInfo{
		ID:      "EntityResponse",
		Name:    "EntityResponse",
		Package: "messages",
		Docs:    []string{"Return message for any create/update action on entity"},
	}
	ci.Fields = append(ci.Fields, newField(0, "Code", "code", "int", "Error code (0 for success)"))
	ci.Fields = append(ci.Fields, newField(1, "Error", "error", "string", "Error message"))
	ci.Fields = append(ci.Fields, newField(2, "Entity", "entity", "Entity", "The entity"))

	return ci
}

// NewEntitiesResponseModel create NewEntitiesResponseModel metamodel
func NewEntitiesResponseModel() *ClassInfo {
	ci := &ClassInfo{
		ID:      "EntitiesResponse",
		Name:    "EntitiesResponse",
		Package: "messages",
		Docs:    []string{"Return message for any action returning multiple entities"},
	}
	ci.Fields = append(ci.Fields, newField(0, "Code", "code", "int", "Error code (0 for success)"))
	ci.Fields = append(ci.Fields, newField(1, "Error", "error", "string", "Error message"))
	ci.Fields = append(ci.Fields, newField(2, "Page", "page", "int", "Current page (Bulk) number"))
	ci.Fields = append(ci.Fields, newField(3, "Size", "size", "int", "Size of page (items in bulk)"))
	ci.Fields = append(ci.Fields, newField(4, "Pages", "pages", "int", "Total number of pages"))
	ci.Fields = append(ci.Fields, newField(5, "Total", "total", "int", "Total number of items in the query"))
	ci.Fields = append(ci.Fields, newArrField(6, "List", "list", "Entity", "List of objects in the current result set"))

	return ci
}
