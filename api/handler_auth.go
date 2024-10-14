package api

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v5"

	"github.com/fmiskovic/eth-auth/logging"
)

func (h handler) handleAuth(w http.ResponseWriter, r *http.Request) error {
	logger := logging.Logger()

	var request struct {
		Address   string `json:"address"`
		Signature string `json:"signature"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.ErrorContext(r.Context(), "failed to decode request", "error", err)
		return newError(http.StatusBadRequest, "invalid request", err)
	}

	// validate the request
	if request.Address == "" {
		logger.ErrorContext(r.Context(), "missing address")
		return newError(http.StatusBadRequest, "missing address", nil)
	}
	if request.Signature == "" {
		logger.ErrorContext(r.Context(), "missing signature")
		return newError(http.StatusBadRequest, "missing signature", nil)
	}

	address := strings.ToLower(request.Address)
	signature := request.Signature

	// get the nonce from the store
	nonce, exist := h.store.Get(address)
	if !exist {
		logger.ErrorContext(r.Context(), "failed to find nonce", "address", address)
		return newError(http.StatusBadRequest, "failed to find nonce", errors.New("nonce not found"))
	}

	// recover the address from the signature
	recoveredAddress, err := h.recoverAddressFromSignature(nonce.(string), signature)
	if err != nil {
		logger.ErrorContext(r.Context(), "failed to recover address from signature", "error", err)
		return newError(http.StatusBadRequest, "failed to recover address from signature", err)
	}

	// verify the signature
	if strings.ToLower(recoveredAddress) != address {
		logger.ErrorContext(r.Context(), "signature verification failed", "address", address)
		return newError(
			http.StatusBadRequest,
			"signature verification failed",
			errors.New("signature verification failed"),
		)
	}

	// delete the nonce from the store to prevent replay attacks
	h.store.Delete(address)

	token, err := h.generateJwt(address)
	if err != nil {
		logger.ErrorContext(r.Context(), "failed to generate JWT", "error", err)
		return newError(http.StatusInternalServerError, "failed to generate JWT token", err)
	}

	var response struct {
		Token string `json:"access_token"`
	}
	response.Token = token

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		logger.ErrorContext(r.Context(), "failed to encode response", "error", err)
		return newError(http.StatusInternalServerError, "failed to encode response", err)
	}

	return nil
}

// recoverAddressFromSignature verifies the signature and returns the signer's address
func (h handler) recoverAddressFromSignature(nonce, signature string) (string, error) {
	// add the ethereum personal message prefix
	messageHash := crypto.Keccak256Hash([]byte("\x19Ethereum Signed Message:\n" + fmt.Sprint(len(nonce)) + nonce))

	// decode the signature hex string into bytes
	sigBytes, err := hex.DecodeString(signature[2:]) // remove '0x' prefix
	if err != nil {
		return "", err
	}

	// handle the 'v' value in the signature
	if sigBytes[64] == 27 || sigBytes[64] == 28 {
		sigBytes[64] -= 27
	}

	// recover the public key from the signature and the hashed message
	pubKey, err := crypto.SigToPub(messageHash.Bytes(), sigBytes)
	if err != nil {
		return "", err
	}

	// get the Ethereum address from the public key
	recoveredAddress := crypto.PubkeyToAddress(*pubKey).Hex()

	return recoveredAddress, nil
}

// generateJwt generates a JWT token for the authenticated user
func (h handler) generateJwt(address string) (string, error) {
	claims := jwt.MapClaims{
		"address": address,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.secret))
}
