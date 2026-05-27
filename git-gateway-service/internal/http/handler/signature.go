package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func readRequestBody(r *http.Request) ([]byte, error) {
	return io.ReadAll(r.Body)
}

func verifyHMACSHA256(payload, secret []byte, headerValue, prefix string) error {
	if headerValue == "" {
		return fmt.Errorf("missing signature header")
	}
	raw := strings.TrimPrefix(headerValue, prefix)
	raw = strings.TrimSpace(raw)
	expected, err := hex.DecodeString(raw)
	if err != nil {
		return fmt.Errorf("invalid signature encoding")
	}
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write(payload)
	if !hmac.Equal(expected, mac.Sum(nil)) {
		return fmt.Errorf("signature mismatch")
	}
	return nil
}
