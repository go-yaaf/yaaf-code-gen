package test

import (
	"encoding/json"
	"testing"

	. "github.com/go-yaaf/yaaf-code-gen/model"
)

func TestCreateClassInfoSample(t *testing.T) {
	m := NewMetaModel()

	// Enum: Status { Active=1, Blocked=2 }
	status := NewEnumInfo("Status")
	av := NewEnumValueInfo("Active")
	av.Value = 1
	bv := NewEnumValueInfo("Blocked")
	bv.Value = 2
	status.AddValue(av)
	status.AddValue(bv)
	m.AddEnumInfo(status)

	// Nested class: Address { City string }
	addr := NewClassInfo("Address")
	addr.AddField("City", "string")
	m.AddClassInfo(addr)

	// Root class: Account
	acc := NewClassInfo("Account")
	acc.AddField("Name", "string")
	acc.AddField("Age", "int")
	acc.GetField("Age").TsType = "number"
	acc.AddField("CreatedOn", "Timestamp")
	acc.AddField("Status", "Status") // enum
	acc.AddField("Home", "Address")  // nested class
	acc.AddField("Tags", "string")   // array of primitives
	acc.GetField("Tags").IsArray = true
	acc.AddField("Contacts", "Address") // array of classes
	acc.GetField("Contacts").IsArray = true

	out := CreateClassInfoSample(m, acc)

	var res map[string]any
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, out)
	}

	// Enum -> one of the numeric values
	if v, _ := res["status"].(float64); v != 1 && v != 2 {
		t.Errorf("status should be an enum value (1 or 2), got %v", res["status"])
	}
	// Nested class -> object with city
	if home, ok := res["home"].(map[string]any); !ok {
		t.Errorf("home should be a nested object, got %T", res["home"])
	} else if _, ok := home["city"]; !ok {
		t.Errorf("nested Address missing city: %v", home)
	}
	// Array of primitives -> 2 entries
	if tags, ok := res["tags"].([]any); !ok || len(tags) != 2 {
		t.Errorf("tags should be an array of 2, got %v", res["tags"])
	}
	// Array of classes -> 2 nested objects
	if c, ok := res["contacts"].([]any); !ok || len(c) != 2 {
		t.Errorf("contacts should be an array of 2, got %v", res["contacts"])
	}
}
