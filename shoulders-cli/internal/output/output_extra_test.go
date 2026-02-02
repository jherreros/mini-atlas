package output

import "testing"

func TestParseFormatInvalid(t *testing.T) {
	_, err := ParseFormat("xml")
	if err == nil {
		t.Fatalf("expected error for invalid format")
	}
}

func TestRenderTableReturnsNil(t *testing.T) {
	payload, err := Render(map[string]string{"hello": "world"}, Table)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payload != nil {
		t.Fatalf("expected nil payload for table format")
	}
}
