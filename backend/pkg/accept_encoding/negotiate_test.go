package accept_encoding

import (
	"testing"
)

func TestNegotiate(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
	}{
		{"empty header", "", ""},
		{"gzip only", "gzip", "gzip"},
		{"br only", "br", "br"},
		{"unsupported encoding", "deflate", ""},
		{"identity only", "identity", ""},
		{"gzip preferred over br", "gzip, br", "br"},
		{"br first in header", "br, gzip", "br"},
		{"br excluded", "br;q=0, gzip", "gzip"},
		{"gzip excluded", "gzip;q=0, br", "br"},
		{"wildcard", "*", "gzip"},
		{"wildcard excluded", "*;q=0", ""},
		{"wildcard with gzip excluded", "br;q=0, *", "gzip"},
		{"all excluded", "gzip;q=0, br;q=0, *", ""},
		{"quality ordering", "gzip;q=0.5, br;q=0.8", "br"},
		{"quality ordering inverted", "br;q=0.2, gzip;q=0.9", "gzip"},
		{"malformed quality", "gzip;q=abc", "gzip"},
		{"case sensitive ignored", "GZip", ""},
		{"gzip excluded with identity", "gzip;q=0, identity", ""},
		{"wildcard lower quality than gzip", "gzip;q=0.5, *;q=0.2", "gzip"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Negotiate(tt.header); got != tt.want {
				t.Errorf("Negotiate(%q) = %q, want %q", tt.header, got, tt.want)
			}
		})
	}
}

func TestParseField(t *testing.T) {
	tests := []struct {
		field     string
		wantName  string
		wantQual  float64
	}{
		{"gzip", "gzip", 1.0},
		{" br ", "br", 1.0},
		{"", "", 0},
		{"gzip;q=0.5", "gzip", 0.5},
		{"gzip;q=0", "gzip", 0},
		{"br; q=0.8; foo=1", "br", 0.8},
		{"gzip;q=abc", "gzip", 1.0},
	}
	for _, tt := range tests {
		name, qual := parseField(tt.field)
		if name != tt.wantName || qual != tt.wantQual {
			t.Errorf("parseField(%q) = (%q, %v), want (%q, %v)", tt.field, name, qual, tt.wantName, tt.wantQual)
		}
	}
}
