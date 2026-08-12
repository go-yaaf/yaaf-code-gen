package model

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/go-yaaf/yaaf-common/entity"
)

// CreateClassInfoSample creates a sample JSON representation of a class
// instance based on its fields, using a random value per field type.
// Nested classes and enums are resolved via the meta-model; arrays get 2
// sample entries and maps/Json get 2 sample key-value pairs.
func CreateClassInfoSample(m *MetaModel, c *ClassInfo) string {
	return renderJSON(buildClass(m, c, map[string]bool{}), 0)
}

// jsonObject is an order-preserving JSON object (map[string]any sorts keys).
type jsonObject struct {
	keys []string
	vals []any
}

func (o *jsonObject) add(k string, v any) {
	o.keys = append(o.keys, k)
	o.vals = append(o.vals, v)
}

// buildClass builds a value tree for a class. seen guards against reference cycles.
func buildClass(m *MetaModel, c *ClassInfo, seen map[string]bool) any {
	if c == nil || len(c.Fields) == 0 || seen[c.Name] {
		return &jsonObject{}
	}
	seen[c.Name] = true
	defer delete(seen, c.Name)

	obj := &jsonObject{}
	fields := c.Fields
	if m != nil && c.BaseClass != "" {
		fields = m.GetAllClassFields(c.Name)
	}
	for _, fi := range fields {
		key := fi.Json
		if key == "" {
			key = fi.TsName
		}
		obj.add(key, buildField(m, fi, seen))
	}
	return obj
}

// buildField resolves a field's value, honoring array (2 entries) and map wrappers.
func buildField(m *MetaModel, fi *FieldInfo, seen map[string]bool) any {
	if fi.IsMap {
		o := &jsonObject{}
		o.add(fmt.Sprintf("%s_0", randomWord(fi.Name)), buildElem(m, fi, seen))
		o.add(fmt.Sprintf("%s_1", randomWord(fi.Name)), buildElem(m, fi, seen))
		return o
	}
	if fi.IsArray {
		return []any{buildElem(m, fi, seen), buildElem(m, fi, seen)}
	}
	return buildElem(m, fi, seen)
}

// buildElem builds a single value for the field's (element) type.
func buildElem(m *MetaModel, fi *FieldInfo, seen map[string]bool) any {
	base := strings.TrimPrefix(fi.Type, "[]")

	// Timestamp serializes as epoch-millis number.
	if base == "Timestamp" {
		return 1700000000000 + rand.Int63n(1_000_000_000)
	}

	if m != nil {
		if e := m.GetEnum(base); e != nil {
			return enumSample(e)
		}
		if ci := m.FindClass(base); ci != nil {
			return buildClass(m, ci, seen)
		}
	}

	ts := fi.TsType
	if ts == "" {
		ts = GetTsType(base)
	}
	switch ts {
	case "string":
		return randomWord(fi.Name)
	case "number":
		return rand.Intn(1000)
	case "boolean":
		return rand.Intn(2) == 0
	case "File":
		return ""
	case "any", "Json", "Record<string,any>":
		// Untyped map / Json object — a couple of arbitrary pairs.
		o := &jsonObject{}
		o.add("key_0", randomWord(fi.Name))
		o.add("key_1", rand.Intn(1000))
		return o
	default:
		// Unknown complex type with no meta-model to resolve — empty object.
		return &jsonObject{}
	}
}

// buildEntityResponse builds EntityResponse class
func buildEntityResponse(m *MetaModel, c *ClassInfo) any {
	obj := &jsonObject{}
	obj.add("code", 0)
	obj.add("entity", buildClass(m, c, map[string]bool{}))
	return obj
}

// buildEntityResponse builds EntityResponse class
func buildEntitiesResponse(m *MetaModel, c *ClassInfo) any {
	obj := &jsonObject{}
	obj.add("code", 0)
	obj.add("page", 1)
	obj.add("size", 100)
	obj.add("pages", 1)
	obj.add("total", 2)

	list := make([]any, 0)
	list = append(list, buildClass(m, c, map[string]bool{}))
	list = append(list, buildClass(m, c, map[string]bool{}))

	obj.add("list", list)
	return obj
}

// enumSample returns a random enum value (its numeric value, as serialized).
func enumSample(e *EnumInfo) any {
	if len(e.Values) == 0 {
		return 0
	}
	return e.Values[rand.Intn(len(e.Values))].Value
}

var sampleWords = []string{"alpha", "bravo", "charlie", "delta", "echo", "foxtrot", "golf", "hotel", "india", "juliett", "kilo", "lima", "mike", "november", "oscar", "papa", "quebec", "romeo", "sierra", "tango", "uniform", "victor", "whiskey", "xray", "yankee", "zulu"}

func randomWord(name string) string {
	n := strings.ToLower(name)
	if strings.HasSuffix(n, "id") {
		return entity.GUID()
	}
	if strings.HasSuffix(n, "key") {
		return entity.GUID()
	}
	if strings.HasSuffix(n, "code") {
		return "0"
	}
	if strings.HasSuffix(n, "error") {
		return ""
	}
	if strings.HasSuffix(n, "email") {
		return fmt.Sprintf("%s@%s.com", sampleWords[rand.Intn(len(sampleWords))], sampleWords[rand.Intn(len(sampleWords))])
	}
	return sampleWords[rand.Intn(len(sampleWords))]
}

// renderJSON pretty-prints the value tree, preserving object key order.
func renderJSON(v any, level int) string {
	switch t := v.(type) {
	case *jsonObject:
		if len(t.keys) == 0 {
			return "{}"
		}
		pad := strings.Repeat("  ", level+1)
		end := strings.Repeat("  ", level)
		lines := make([]string, 0, len(t.keys))
		for i, k := range t.keys {
			kb, _ := json.Marshal(k)
			lines = append(lines, fmt.Sprintf("%s%s: %s", pad, kb, renderJSON(t.vals[i], level+1)))
		}
		return "{\n" + strings.Join(lines, ",\n") + "\n" + end + "}"
	case []any:
		if len(t) == 0 {
			return "[]"
		}
		pad := strings.Repeat("  ", level+1)
		end := strings.Repeat("  ", level)
		lines := make([]string, 0, len(t))
		for _, e := range t {
			lines = append(lines, pad+renderJSON(e, level+1))
		}
		return "[\n" + strings.Join(lines, ",\n") + "\n" + end + "]"
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

// CreateMethodResponseSample creates a sample JSON representation of a method
// instance based on its fields, using a random value per field type.
// Nested classes and enums are resolved via the meta-model; arrays get 2
// sample entries and maps/Json get 2 sample key-value pairs.
func CreateMethodResponseSample(m *MetaModel, c *MethodInfo) string {

	if c.ReturnClass == "StreamContent" {
		return "Binary stream of byte array (file download or image)"
	}

	if resp := m.FindClass(c.ReturnClass); resp != nil {
		return renderJSON(buildClass(m, resp, map[string]bool{}), 0)
	}

	getGenericTypeClass := func(node *TypeNode) *ClassInfo {
		if len(node.Args) == 0 {
			return nil
		} else {
			return m.FindClass(node.Args[0].Name)
		}
	}

	if strings.HasPrefix(c.ReturnClass, "EntityResponse") {
		if ci := getGenericTypeClass(c.ReturnType); ci != nil {
			return renderJSON(buildEntityResponse(m, ci), 0)
		} else {
			return c.ReturnClass
		}
	}
	if strings.HasPrefix(c.ReturnClass, "EntitiesResponse") {
		if ci := getGenericTypeClass(c.ReturnType); ci != nil {
			return renderJSON(buildEntitiesResponse(m, ci), 0)
		} else {
			return c.ReturnClass
		}
	}
	return c.ReturnClass
}

// CreateMethodRequestSample creates a sample JSON representation of a method
// instance based on its fields, using a random value per field type.
// Nested classes and enums are resolved via the meta-model; arrays get 2
// sample entries and maps/Json get 2 sample key-value pairs.
func CreateMethodRequestSample(m *MetaModel, c *MethodInfo) string {

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("[%s] https://path/to/api/v1%s\n", c.Method, c.FullPath()))
	sb.WriteString("\t-H \"Content-Type: application/json\"\n")
	sb.WriteString("\t-H \"X-API-KEY: {{your-api-key}}\"\n")
	sb.WriteString("\t-H \"X-ACCESS-TOKEN: {{your-access-token}}\"\n")

	if c.BodyParam != nil {
		bodyType := c.BodyParam.Type

		if ci := m.FindClass(bodyType); ci != nil {
			bodyType = renderJSON(buildClass(m, ci, map[string]bool{}), 0)
		}

		if c.BodyParam.IsArray {
			sb.WriteString(fmt.Sprintf("[\n%s\n]\n", bodyType))
		} else {
			sb.WriteString(fmt.Sprintf("%s\n", bodyType))
		}
	}

	/*
		{{.Method}} https://path/to/api/v1{{.Path}} \
		  -H "X-API-KEY: your_api_key"
		  -H "X-ACCESS-TOKEN: your_access_token"
	*/

	/*
		GET http://{{host}}/v1/accounts
		Content-Type: application/json
		X-API-KEY: {{goox-dashboard-key}}
		X-Access-Token: {{goox-access-token}}


	*/

	return sb.String()
}
