package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleNonce(t *testing.T) {
	t.Parallel()

	// init router
	r := New("secret-for-test-purposes")

	tests := []struct {
		name     string
		request  []byte
		wantCode int
	}{
		{
			name:     "valid request",
			request:  []byte(`{"address":"0x1234567890abcdef"}`),
			wantCode: http.StatusOK,
		},
		{
			name:     "invalid request",
			request:  []byte(`{"address":""}`),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid json",
			request:  []byte(`{`),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "empty json",
			request:  []byte(`{}`),
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest("POST", "/nonce", bytes.NewReader(tt.request))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("got %d; want %d", w.Code, tt.wantCode)
			}

			if w.Code == http.StatusOK {
				var response struct {
					Nonce string `json:"nonce"`
				}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}

				if response.Nonce == "" {
					t.Errorf("got empty nonce")
				}
			}
		})
	}
}
