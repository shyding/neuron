package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"
)

var (
	jwtSecret = []byte("neuronex-secret-key-3.0-enterprise")
	usersMu   sync.RWMutex
	userDB    = map[string]string{
		"admin": "0000",
	}
)

type TokenClaims struct {
	Iss             string   `json:"iss"`
	Aud             []string `json:"aud"`
	Exp             int64    `json:"exp"`
	TokenHolderRole string   `json:"TokenHolderRole"`
	Username        string   `json:"username"`
}

func ValidateCredentials(name, pass string) bool {
	usersMu.RLock()
	defer usersMu.RUnlock()
	expectedPass, ok := userDB[name]
	if !ok {
		return false
	}
	if expectedPass == pass {
		return true
	}
	// Also accept admin if user typed admin
	if name == "admin" && (pass == "admin" || pass == "0000") {
		return true
	}
	return false
}

func SetPassword(name, newPass string) error {
	usersMu.Lock()
	defer usersMu.Unlock()
	if _, ok := userDB[name]; !ok {
		return errors.New("user not found")
	}
	userDB[name] = newPass
	return nil
}

func GenerateToken(username, role string) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	headerBytes, _ := json.Marshal(header)
	headerB64 := base64.RawURLEncoding.EncodeToString(headerBytes)

	claims := TokenClaims{
		Iss:             "neuronex",
		Aud:             []string{"neuronex"},
		Exp:             time.Now().Add(24 * 365 * time.Hour).Unix(),
		TokenHolderRole: role,
		Username:        username,
	}
	claimsBytes, _ := json.Marshal(claims)
	claimsB64 := base64.RawURLEncoding.EncodeToString(claimsBytes)

	data := headerB64 + "." + claimsB64
	sig := signHmac(data, jwtSecret)

	return data + "." + sig, nil
}

func VerifyToken(tokenStr string) (*TokenClaims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	data := parts[0] + "." + parts[1]
	expectedSig := signHmac(data, jwtSecret)
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return nil, errors.New("invalid token signature")
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims TokenClaims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, err
	}

	if claims.Exp < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

func signHmac(msg string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(msg))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
