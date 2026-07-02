package geoip

import "testing"

func TestExtractHostFromURI_MieruSimple(t *testing.T) {
	got := extractHostFromURI("mierus://user:pass@example.com?profile=default&port=2999&protocol=TCP")
	if got != "example.com" {
		t.Fatalf("extractHostFromURI() = %q, want example.com", got)
	}
}
