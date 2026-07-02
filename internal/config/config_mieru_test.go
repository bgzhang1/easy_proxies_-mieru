package config

import (
	"net/url"
	"testing"
)

func TestIsProxyURI_Mieru(t *testing.T) {
	if !IsProxyURI("mieru://base64payload") {
		t.Fatalf("expected mieru:// URI to be recognized")
	}
	if !IsProxyURI("mierus://user:pass@example.com?profile=default&port=2999&protocol=TCP") {
		t.Fatalf("expected mierus:// URI to be recognized")
	}
}

func TestParseClashYAML_Mieru(t *testing.T) {
	content := `proxies:
  - name: "test-mieru"
    type: "mieru"
    server: example.com
    username: user
    password: pass
    port: 2999
    transport: TCP
    multiplexing: MULTIPLEXING_LOW
    traffic-pattern: GgQIARAK
`

	nodes, err := parseClashYAML(content)
	if err != nil {
		t.Fatalf("parse clash yaml failed: %v", err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}
	u, err := url.Parse(nodes[0].URI)
	if err != nil {
		t.Fatalf("parse generated uri failed: %v", err)
	}
	if u.Scheme != "mierus" {
		t.Fatalf("expected scheme mierus, got %q", u.Scheme)
	}
	if u.Hostname() != "example.com" {
		t.Fatalf("expected host example.com, got %q", u.Hostname())
	}
	if u.User.Username() != "user" {
		t.Fatalf("expected user, got %q", u.User.Username())
	}
	password, _ := u.User.Password()
	if password != "pass" {
		t.Fatalf("expected password, got %q", password)
	}
	query := u.Query()
	if query.Get("port") != "2999" {
		t.Fatalf("expected port=2999, got %q", query.Get("port"))
	}
	if query.Get("protocol") != "TCP" {
		t.Fatalf("expected protocol=TCP, got %q", query.Get("protocol"))
	}
	if query.Get("multiplexing") != "MULTIPLEXING_LOW" {
		t.Fatalf("expected multiplexing, got %q", query.Get("multiplexing"))
	}
	if query.Get("traffic-pattern") != "GgQIARAK" {
		t.Fatalf("expected traffic-pattern, got %q", query.Get("traffic-pattern"))
	}
}
