package payment

import "testing"

func TestNormalizeJSONBool(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		want   bool
		wantOK bool
	}{
		{name: "bool true", input: true, want: true, wantOK: true},
		{name: "bool false", input: false, want: false, wantOK: true},
		{name: "string true", input: "true", want: true, wantOK: true},
		{name: "string false", input: "false", want: false, wantOK: true},
		{name: "string one", input: "1", want: true, wantOK: true},
		{name: "string zero", input: "0", want: false, wantOK: true},
		{name: "number one", input: float64(1), want: true, wantOK: true},
		{name: "number zero", input: float64(0), want: false, wantOK: true},
		{name: "invalid string", input: "sandbox", want: false, wantOK: false},
		{name: "empty string", input: "", want: false, wantOK: false},
		{name: "nil", input: nil, want: false, wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := normalizeJSONBool(tt.input)
			if got != tt.want || ok != tt.wantOK {
				t.Fatalf("normalizeJSONBool(%v) = (%v, %v), want (%v, %v)", tt.input, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}
