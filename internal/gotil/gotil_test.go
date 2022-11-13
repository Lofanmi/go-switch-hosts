package gotil

import (
	"testing"
)

func TestStringCut(t *testing.T) {
	type args struct {
		s         string
		begin     string
		end       string
		withBegin bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test", args{"", "a", "d", false}, ""},
		{"test", args{"abcd", "", "d", false}, "abc"},
		{"test", args{"abcd", "a", "", false}, ""},
		{"test", args{"abcd", "e", "d", false}, ""},
		{"test", args{"abcd", "a", "f", false}, ""},
		{"test", args{"abcd", "a", "d", false}, "bc"},
		{"test", args{"abcd", "a", "d", true}, "abc"},
		{"test", args{"abcd", "abcd", "", true}, "abcd"},
		{"test", args{"abcd", "", "abcd", true}, ""},
		{"test", args{"abcd", "ab", "cd", false}, ""},
		{"test", args{"abcd", "abcd", "e", false}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringCut(tt.args.s, tt.args.begin, tt.args.end, tt.args.withBegin); got != tt.want {
				t.Errorf("StringCut() = %v, want %v", got, tt.want)
			}
		})
	}
}
