package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fmiskovic/eth-auth/store"
)

func TestHandleAuth(t *testing.T) {
	nonce := "f6eccf7aee12fb29793a4147410820e7b22fea422162e83214e798c626f187f2"
	address := "0x52Aa0e471F234CfC4997b6f7A526DB4eDee152E4"
	signature := "0x5e32205e0b731461223cda8b6b68a1abd1a8e25ff903f677ed5196eb5a86a2e86f7a4ba10e4c55bc5c28b8fdffc79537593b177ca5f2b45eaad5b11e073a3f841b"

	tests := []struct {
		name     string
		request  []byte
		wantCode int
	}{
		{
			name:     "valid request",
			request:  []byte(fmt.Sprintf(`{"address":"%s","signature":"%s"}`, address, signature)),
			wantCode: http.StatusOK,
		},
		{
			name:     "invalid address",
			request:  []byte(fmt.Sprintf(`{"address":"0x1234567890abcdef","signature":"%s"}`, signature)),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid signature",
			request:  []byte(fmt.Sprintf(`{"address":"%s","signature":"0x1234567890abcdef"}`, address)),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "missing address",
			request:  []byte(fmt.Sprintf(`{"signature":"%s"}`, signature)),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "missing signature",
			request:  []byte(fmt.Sprintf(`{"address":"%s"}`, address)),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "empty json",
			request:  []byte(`{}`),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid json",
			request:  []byte(`{`),
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storer, err := store.New(10)
			if err != nil {
				t.Fatalf("failed to create store: %v", err)
			}

			// add nonce to the store
			storer.Add(strings.ToLower(address), nonce)

			h := newHandler(storer, "secret-for-test-purposes")

			req := httptest.NewRequest("POST", "/auth", bytes.NewReader(tt.request))
			w := httptest.NewRecorder()

			if err = h.handleAuth(w, req); err != nil {
				var apiErr *Error
				if errors.As(err, &apiErr) {
					if apiErr.Code != tt.wantCode {
						t.Errorf("got %d; want %d", apiErr.Code, tt.wantCode)
					}
				} else {
					t.Fatalf("failed to handle auth: %v", err)
				}
				return
			}

			if w.Code != tt.wantCode {
				t.Errorf("got %d; want %d", w.Code, tt.wantCode)
			}

			if w.Code == http.StatusOK {
				var response struct {
					Token string `json:"access_token"`
				}

				if err = json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}

				if response.Token == "" {
					t.Errorf("got empty token")
				}
			}
		})
	}
}
