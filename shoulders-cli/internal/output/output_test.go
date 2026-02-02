package output

import "testing"

func TestParseFormat(t *testing.T) {
	format, err := ParseFormat("json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if format != JSON {
		t.Fatalf("expected json format, got %s", format)
	}
}

func TestRenderJSON(t *testing.T) {
	payload, err := Render(map[string]string{"hello": "world"}, JSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payload) == 0 {
		t.Fatalf("expected json payload")
	}
}

func TestRenderYAML(t *testing.T) {
	payload, err := Render(map[string]string{"hello": "world"}, YAML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payload) == 0 {
		t.Fatalf("expected yaml payload")
	}
}
