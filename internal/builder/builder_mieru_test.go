package builder

import (
	"net/url"
	"strings"
	"testing"

	mieruappctl "github.com/enfein/mieru/v3/pkg/appctl"
	mierupb "github.com/enfein/mieru/v3/pkg/appctl/appctlpb"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"google.golang.org/protobuf/proto"
)

func TestBuildNodeOutbound_MieruSimpleURL(t *testing.T) {
	raw := "mierus://user:pass@example.com?profile=default&port=2999&protocol=TCP&port=3000-3002&protocol=TCP#Mieru"
	outbound, err := buildNodeOutbound("test-mieru", raw, false)
	if err != nil {
		t.Fatalf("build node outbound failed: %v", err)
	}

	if outbound.Type != C.TypeMieru {
		t.Fatalf("expected type %q, got %q", C.TypeMieru, outbound.Type)
	}
	opts, ok := outbound.Options.(*option.MieruOutboundOptions)
	if !ok {
		t.Fatalf("expected *option.MieruOutboundOptions, got %T", outbound.Options)
	}
	if opts.Server != "example.com" {
		t.Fatalf("expected server example.com, got %q", opts.Server)
	}
	if opts.ServerPort != 2999 {
		t.Fatalf("expected port 2999, got %d", opts.ServerPort)
	}
	if len(opts.ServerPortRanges) != 1 || opts.ServerPortRanges[0] != "3000-3002" {
		t.Fatalf("expected server port range [3000-3002], got %v", opts.ServerPortRanges)
	}
	if opts.Transport != "TCP" {
		t.Fatalf("expected transport TCP, got %q", opts.Transport)
	}
	if opts.UserName != "user" || opts.Password != "pass" {
		t.Fatalf("unexpected credentials %q/%q", opts.UserName, opts.Password)
	}
}

func TestBuildNodeOutbound_MieruFullConfigURL(t *testing.T) {
	link, err := mieruappctl.ClientConfigToURL(&mierupb.ClientConfig{
		ActiveProfile: proto.String("default"),
		Profiles: []*mierupb.ClientProfile{{
			ProfileName: proto.String("default"),
			User: &mierupb.User{
				Name:     proto.String("user"),
				Password: proto.String("pass"),
			},
			Servers: []*mierupb.ServerEndpoint{{
				DomainName: proto.String("example.org"),
				PortBindings: []*mierupb.PortBinding{{
					Port:     proto.Int32(443),
					Protocol: mierupb.TransportProtocol_TCP.Enum(),
				}},
			}},
			Multiplexing: &mierupb.MultiplexingConfig{
				Level: mierupb.MultiplexingLevel_MULTIPLEXING_LOW.Enum(),
			},
		}},
	})
	if err != nil {
		t.Fatalf("build mieru config URL failed: %v", err)
	}

	parsed, err := url.Parse(link)
	if err != nil {
		t.Fatalf("parse generated link failed: %v", err)
	}
	opts, err := buildMieruOptions(link, parsed)
	if err != nil {
		t.Fatalf("build mieru options failed: %v", err)
	}
	if opts.Server != "example.org" {
		t.Fatalf("expected server example.org, got %q", opts.Server)
	}
	if opts.ServerPort != 443 {
		t.Fatalf("expected port 443, got %d", opts.ServerPort)
	}
	if opts.Transport != "TCP" {
		t.Fatalf("expected transport TCP, got %q", opts.Transport)
	}
	if opts.Multiplexing != "MULTIPLEXING_LOW" {
		t.Fatalf("expected multiplexing MULTIPLEXING_LOW, got %q", opts.Multiplexing)
	}
}

func TestBuildNodeOutbound_MieruMixedTransportRejected(t *testing.T) {
	raw := "mierus://user:pass@example.com?profile=default&port=2999&protocol=TCP&port=3000&protocol=UDP"
	_, err := buildNodeOutbound("bad-mieru", raw, false)
	if err == nil {
		t.Fatalf("expected mixed transport error, got nil")
	}
	if !strings.Contains(err.Error(), "mixed TCP/UDP") {
		t.Fatalf("expected mixed transport error, got %v", err)
	}
}
