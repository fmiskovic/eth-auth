package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/fmiskovic/eth-auth/logging"
)

// handleNonce generates a random nonce and returns it to the client
func (h handler) handleNonce(w http.ResponseWriter, r *http.Request) error {
	logger := logging.Logger()

	var request struct {
		Address string `json:"address"`
	}

	// decode request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.ErrorContext(r.Context(), "failed to decode request", "error", err)
		return newError(http.StatusBadRequest, "failed to decode request", err)
	}

	// validate address
	if request.Address == "" {
		err := errors.New("invalid address")
		return newError(http.StatusBadRequest, "invalid address", err)
	}

	// generate nonce
	nonce, err := h.generateRandomNonce()
	if err != nil {
		logger.ErrorContext(r.Context(), "failed to generate nonce", "error", err)
		return newError(http.StatusInternalServerError, "failed to generate nonce", err)
	}

	// store nonce
	address := strings.ToLower(request.Address)
	logger.InfoContext(r.Context(), "storing nonce", "address", address, "nonce", nonce)
	evicted := h.store.Add(address, nonce)
	if evicted {
		logger.WarnContext(r.Context(), "old nonce evicted", "address", address)
	}

	var response struct {
		Nonce string `json:"nonce"`
	}
	response.Nonce = nonce

	w.WriteHeader(http.StatusOK)
	// encode response
	if err = json.NewEncoder(w).Encode(response); err != nil {
		logger.ErrorContext(r.Context(), "Failed to encode response", "error", err)
		return newError(http.StatusInternalServerError, "failed to encode response", err)
	}

	return nil
}

func (h handler) generateRandomNonce() (string, error) {
	nonce := make([]byte, 32)
	_, err := rand.Read(nonce)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(nonce), nil
}
