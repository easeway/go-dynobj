package dynobj

import "testing"

const (
	MAP_DATA = `{
	"key-int": 10,
	"key-str": "string",
	"key-arr": [
		1, "two",
		{ "nested_key": "nested" }
	],
	"key-map": {
		"subkey_true": true,
		"subkey_false": false,
		"subkey_yes": "yes",
		"subkey_no": "no"
	}
}`
)

func TestParsePathSimple(t *testing.T) {
	parsed := ParsePath("a/b/c")
	expected := []string{"a", "b", "c"}
	if len(parsed) != len(expected) {
		t.Errorf("Incorrect length: %v", len(parsed))
	}
	for i, v := range parsed {
		if v != expected[i] {
			t.Errorf("Mismatch at %v: expect %v, got %v", i, expected[i], v)
		}
	}
}

func TestParsePathDup(t *testing.T) {
	parsed := ParsePath("/a///b/c")
	expected := []string{"a", "b", "c"}
	if len(parsed) != len(expected) {
		t.Errorf("Incorrect length: %v", len(parsed))
	}
	for i, v := range parsed {
		if v != expected[i] {
			t.Errorf("Mismatch at %v: expect %v, got %v", i, expected[i], v)
		}
	}
}

func TestAsStr(t *testing.T) {
	if obj, err := NewJsonStringObj(MAP_DATA); err != nil {
		t.Error(err)
	} else {
		v := obj.AsStrD("/non-existed", "defaultValue")
		if v != "defaultValue" {
			t.Errorf("Unexpected default value: %v", v)
		}
		v = obj.AsStrD("/key-str", "defaultValue")
		if v != "string" {
			t.Errorf("Expect string but got %v", v)
		}
		v = obj.AsStr("/non-existed")
		if v != "" {
			t.Errorf("Unexpected %v", v)
		}
		v = obj.AsStr("/key-arr/2/nested_key")
		if v != "nested" {
			t.Errorf("Expect nested but got %v", v)
		}
	}
}

func TestAsInt(t *testing.T) {
	if obj, err := NewJsonStringObj(MAP_DATA); err != nil {
		t.Error(err)
	} else {
		v := obj.AsIntD("key-str", 100)
		if v != 100 {
			t.Errorf("Expect 100 but got %v", v)
		}
		v = obj.AsIntD("key-int", 100)
		if v != 10 {
			t.Errorf("Expect 10 but got %v", v)
		}
		v = obj.AsInt("/key-arr/0")
		if v != 1 {
			t.Errorf("Expect 1 but got %v", v)
		}
		v = obj.AsInt("/key-arr/2")
		if v != 0 {
			t.Errorf("Expect 0 but got %v", v)
		}
	}
}

func TestAsBool(t *testing.T) {
	if obj, err := NewJsonStringObj(MAP_DATA); err != nil {
		t.Error(err)
	} else {
		v := obj.AsBoolD("key-str", true)
		if !v {
			t.Errorf("Expect true but got %v", v)
		}
		v = obj.AsBoolD("key-map/subkey_yes", false)
		if !v {
			t.Errorf("Expect true but got %v", v)
		}
		v = obj.AsBool("/key-arr/0")
		if !v {
			t.Errorf("Expect true but got %v", v)
		}
		v = obj.AsBool("/key-arr/2")
		if v {
			t.Errorf("Expect false but got %v", v)
		}
		v = obj.AsBool("key-map/subkey_true")
		if !v {
			t.Errorf("Expect true but got %v", v)
		}
		v = obj.AsBool("key-map/subkey_false")
		if v {
			t.Errorf("Expect false but got %v", v)
		}
		v = obj.AsBool("key-map/subkey_no")
		if v {
			t.Errorf("Expect false but got %v", v)
		}
	}
}
