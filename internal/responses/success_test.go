package responses

import (
	"bytes"
	"errors"
	"testing"
	"testing/iotest"
)

func TestProcessSuccessResponse(t *testing.T) {
	testTable := []struct {
		name        string
		response    []byte
		expected    SuccessResponse
		expectedErr error
	}{
		{
			name:        "invalid success response status code",
			response:    []byte("32 text/plain\r\noi"),
			expected:    SuccessResponse{},
			expectedErr: ErrSuccessResponseStatusCode,
		},
		{
			name:        "invalid success response mime type",
			response:    []byte("20 edcba/abcde\r\nOi"),
			expected:    SuccessResponse{},
			expectedErr: ErrSuccessResponseMIMEType,
		},
		{
			name:     "valid success response",
			response: []byte("20 text/plain\r\nOi"),
			expected: SuccessResponse{
				StatusCode:     20,
				rawBody:        []byte("Oi"),
				MIMEType:       "text/plain",
				MIMETypeParams: map[string]string{},
			},
			expectedErr: nil,
		},
		{
			name:     "valid success response with mime type params",
			response: []byte("20 text/plain;charset=UTF-8\r\nOi"),
			expected: SuccessResponse{
				StatusCode: 20,
				rawBody:    []byte("Oi"),
				MIMEType:   "text/plain",
				MIMETypeParams: map[string]string{
					"charset": "UTF-8",
				},
			},
			expectedErr: nil,
		},
		{
			name:     "valid success response with UTF-8 characters",
			response: []byte("20 text/plain;charset=UTF-8\r\n💲 🐖 📨 ⤴️ 🏡 🐣 ☣ ☝️ ⏱ 1️⃣"),
			expected: SuccessResponse{
				StatusCode: 20,
				rawBody:    []byte("💲 🐖 📨 ⤴️ 🏡 🐣 ☣ ☝️ ⏱ 1️⃣"),
				MIMEType:   "text/plain",
				MIMETypeParams: map[string]string{
					"charset": "UTF-8",
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			br := bytes.NewReader(tt.response)
			got, err := processSuccessResponse(br)
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

			if got.MIMEType != tt.expected.MIMEType {
				t.Errorf("got %s expected %s", got.MIMEType, tt.expected.MIMEType)
			}

			if tt.expected.MIMETypeParams != nil && got.MIMETypeParams == nil {
				t.Errorf("got nil expected %+v", tt.expected.MIMETypeParams)
			} else if tt.expected.MIMETypeParams == nil && got.MIMETypeParams != nil {
				t.Errorf("got %+v expected nil", got.MIMETypeParams)
			} else if got.MIMETypeParams != nil && tt.expected.MIMETypeParams != nil {
				for param, value := range tt.expected.MIMETypeParams {
					v, ok := got.MIMETypeParams[param]

					if !ok {
						t.Errorf("got nothing expected %s", value)
					}

					if v != value {
						t.Errorf("got %s expected %s", v, value)
					}
				}
			}

			if !bytes.Equal(got.rawBody, tt.expected.rawBody) {
				t.Errorf("got %s expected %s", got.rawBody, tt.expected.rawBody)
			}

			if tt.expected.rawBody != nil {
				if err := iotest.TestReader(got.Body, tt.expected.rawBody); err != nil {
					t.Errorf("got %s expected nil", err.Error())
				}
			}
		})
	}
}
