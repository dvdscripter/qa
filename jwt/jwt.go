package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	defaultHeader     = joseHeader{Typ: "JWT", Alg: "HS256"}
	DefaultExpiration = time.Now().AddDate(0, 0, 7).Unix()
	defaultLeeway     = 10 * time.Minute
)

type Token struct {
	header    joseHeader
	payload   Payload
	signature []byte
	key       []byte
}

type joseHeader struct {
	Typ string `json:"typ,omitempty"`
	Alg string `json:"alg,omitempty"`
}

type Payload struct {
	Exp   int64  `json:"exp,omitempty"`
	Nbf   int64  `json:"nbf,omitempty"`
	Email string `json:"email,omitempty"`
}

func (p Payload) ValidNbf() bool {
	nbf := time.Unix(p.Nbf, 0)
	return nbf.Before(time.Now())
}

func (p Payload) ValidExp() bool {
	exp := time.Unix(p.Exp, 0)
	return time.Now().Before(exp)
}

func (p Payload) String() (string, error) {
	return b64EncodeMarshal(p)
}

func (h joseHeader) String() (string, error) {
	return b64EncodeMarshal(h)
}

func NewFromFile(claims Payload, keyfile string) *Token {
	claims.Nbf = time.Now().Add(-defaultLeeway).Unix()
	token := &Token{
		header:  defaultHeader,
		payload: claims,
	}
	token.loadKeyFromFile(keyfile)
	return token
}

func (t *Token) loadKeyFromFile(keyfile string) error {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return err
	}
	t.key = make([]byte, len(key))
	copy(t.key, key)
	return nil
}

func (t Token) Encode() (string, error) {
	b64Header, err := t.header.String()
	if err != nil {
		return "", err
	}

	b64Payload, err := t.payload.String()
	if err != nil {
		return "", err
	}

	hp := b64Header + "." + b64Payload
	t.sign(hp)
	return hp + "." + base64.RawURLEncoding.EncodeToString(t.signature), nil
}

func (t *Token) Decode(s string) error {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return fmt.Errorf("Invalid token")
	}

	if err := b64DecodeUnmarshal(parts[0], &t.header); err != nil {
		return fmt.Errorf("Cannot decode header %s", s)
	}
	if err := b64DecodeUnmarshal(parts[1], &t.payload); err != nil {
		return fmt.Errorf("Cannot decode payload %s", s)
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return fmt.Errorf("Cannot decode signature %s", s)
	}
	t.signature = make([]byte, len(signature))
	copy(t.signature, signature)

	return nil
}

func DecodePayload(r *http.Request) (payload Payload) {
	header := r.Header.Get("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return
	}

	rawToken := header[len("Bearer "):]
	parts := strings.Split(rawToken, ".")
	if len(parts) != 3 {
		return
	}

	b64DecodeUnmarshal(parts[1], &payload)
	return
}

func b64EncodeMarshal(src interface{}) (string, error) {
	r, err := json.Marshal(src)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(r), nil
}

func b64DecodeUnmarshal(s string, dst interface{}) error {
	raw, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return err
	}
	return nil
}

func (t Token) Check() error {
	b64Header, err := t.header.String()
	if err != nil {
		return err
	}

	b64Payload, err := t.payload.String()
	if err != nil {
		return err
	}

	if t.payload.Nbf != 0 {
		if !t.payload.ValidNbf() {
			return fmt.Errorf("Invalid nbf claim")
		}
	}
	if t.payload.Exp != 0 {
		if !t.payload.ValidExp() {
			return fmt.Errorf("Invalid exp claim")
		}
	}

	hp := b64Header + "." + b64Payload
	signature := sign(hp, t.key)
	if check := hmac.Equal(t.signature, signature); !check {
		return fmt.Errorf("Invalid signature")
	}
	return nil
}

func (t *Token) sign(s string) {
	signature := sign(s, t.key)
	t.signature = make([]byte, len(signature))
	copy(t.signature, signature)
}

func sign(s string, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(s))
	return mac.Sum(nil)
}
