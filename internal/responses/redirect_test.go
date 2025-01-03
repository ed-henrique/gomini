package responses

import (
	"bytes"
	"errors"
	"testing"
)

func TestProcessRedirectResponse(t *testing.T) {
	testTable := []struct {
		name        string
		response    []byte
		expected    RedirectResponse
		expectedErr error
	}{
		{
			name:        "invalid redirect response status code",
			response:    []byte("45 gemini://example.net/search\r\n"),
			expected:    RedirectResponse{},
			expectedErr: ErrRedirectResponseStatusCode,
		},
		{
			name:     "valid redirect response",
			response: []byte("30 gemini://example.net/search\r\n"),
			expected: RedirectResponse{
				StatusCode: 30,
				URI:        "gemini://example.net/search",
			},
			expectedErr: nil,
		},
		{
			name:     "valid redirect response with params and fragment",
			response: []byte("30 gemini://example.net/search?aaa#oi\r\n"),
			expected: RedirectResponse{
				StatusCode: 30,
				URI:        "gemini://example.net/search?aaa#oi",
			},
			expectedErr: nil,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			br := bytes.NewReader(tt.response)
			got, err := processRedirectResponse(br)
			if tt.expectedErr != nil && err == nil {
				t.Errorf("got nil expected %s", tt.expectedErr.Error())
				return
			} else if tt.expectedErr == nil && err != nil {
				t.Errorf("got %s expected nil", err.Error())
				return
			} else if tt.expectedErr != nil && err != nil {
				if !errors.Is(tt.expectedErr, err) {
					t.Errorf("got %s expected %s", err.Error(), tt.expectedErr.Error())
					return
				}
			}

			if got.StatusCode != tt.expected.StatusCode {
				t.Errorf("got %d expected %d", got.StatusCode, tt.expected.StatusCode)
			}

			if got.URI != tt.expected.URI {
				t.Errorf("got %s expected %s", got.URI, tt.expected.URI)
			}
		})
	}
}
