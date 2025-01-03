package main

import (
	"testing"
)

func TestGetAbsoluteURI(t *testing.T) {
	testTable := []struct {
		name      string
		authority string
		expected  string
	}{
		{
			name:      "authority without / suffix",
			authority: "a",
			expected:  "gemini://a:1965/\r\n",
		},
		{
			name:      "authority with / suffix",
			authority: "a/",
			expected:  "gemini://a:1965/\r\n",
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := getAbsoluteURI(tt.authority)
			if got != tt.expected {
				t.Errorf("got %s expected %s", got, tt.expected)
			}
		})
	}
}
