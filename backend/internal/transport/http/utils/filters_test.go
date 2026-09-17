package utils

import (
	"testing"
	"time"
)

func TestNormalizeFilterValue(t *testing.T) {
	tests := []struct {
		name  string
		field string
		value string
		want  string
	}{
		{
			name:  "nextVerificationDate epoch seconds converted",
			field: "nextVerificationDate",
			value: "1774983600",
			want:  "2026-04-01 00:00:00",
		},
		{
			name:  "ISO string left as is",
			field: "nextVerificationDate",
			value: "2026-04-01T00:00:00.000Z",
			want:  "2026-04-01T00:00:00.000Z",
		},
		{
			name:  "NullDate left as is",
			field: "nextVerificationDate",
			value: "0001-01-01T00:00:00Z",
			want:  "0001-01-01T00:00:00Z",
		},
		{
			name:  "non date field left as is",
			field: "place",
			value: "1774983600",
			want:  "1774983600",
		},
		{
			name:  "text date value left as is",
			field: "nextVerificationDate",
			value: "nextMonth",
			want:  "nextMonth",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := time.Local
			defer func() { time.Local = old }()
			time.Local = time.FixedZone("UTC+5", 5*60*60)

			got := NormalizeFilterValue(tt.field, tt.value)
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}