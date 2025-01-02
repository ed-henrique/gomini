package main

import (
	"errors"
	"strings"
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

func TestProcessInputResponse(t *testing.T) {
	testTable := []struct {
		name        string
		response    string
		expected    InputResponse
		expectedErr error
	}{
		{
			name:        "invalid response status code",
			response:    "25 Digite seu nome\r\n",
			expected:    InputResponse{},
			expectedErr: ErrInputResponseStatusCode,
		},
		{
			name:     "valid input response",
			response: "10 O que você deseja buscar?\r\n",
			expected: InputResponse{
				StatusCode: 10,
				Prompt:     "O que você deseja buscar?",
			},
			expectedErr: nil,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sr := strings.NewReader(tt.response)
			got, err := processInputResponse(sr)
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

			if got.Prompt != tt.expected.Prompt {
				t.Errorf("got %s expected %s", got.Prompt, tt.expected.Prompt)
			}
		})
	}
}
