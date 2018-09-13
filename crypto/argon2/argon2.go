package argon2

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

const (
	defaultTime    = 1
	defaultMemory  = 64 * 1024
	defaultThreads = 4
	defaultLength  = 32
	defaultVersion = 0x13
)

type ARGON2 struct {
	time, memory uint32
	threads      uint8
	keyLen       uint32
}

type ARGON2Option func(*ARGON2)

func New(options ...ARGON2Option) *ARGON2 {
	argon := ARGON2{defaultTime, defaultMemory, defaultThreads, defaultLength}

	for _, option := range options {
		option(&argon)
	}

	return &argon
}

func newSalt() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		errors.Wrap(err, "unable to generate salt")
	}
	return base64.RawStdEncoding.EncodeToString(b), nil
}

// GenerateFromPassword hash password using argon2id and return phc format string
func GenerateFromPassword(password, salt []byte, parameters *ARGON2) (string, error) {

	if parameters == nil {
		parameters = New()
	}

	if salt == nil {
		rsalt, err := newSalt()
		if err != nil {
			return "", err
		}
		salt = []byte(rsalt)
	}

	pass := argon2.IDKey(password, salt, parameters.time,
		parameters.memory, parameters.threads, parameters.keyLen)

	return fmt.Sprintf("$argon2id$v=%d$m=%d$t=%d$p=%d$%s$%s", defaultVersion,
		parameters.memory, parameters.time, parameters.threads, string(salt),
		base64.RawStdEncoding.EncodeToString(pass)), nil
}

func newFromHash(hash string) (*ARGON2, []byte, error) {

	var argontype string
	var time, memory, threads uint64

	parts := strings.Split(hash[1:], "$")
	if len(parts) != 7 {
		return nil, nil, errors.New("incorrect number of segments on hash")
	}

	argontype = parts[0]
	if argontype != "argon2id" {
		return nil, nil, errors.New("incorrect argon2 type")
	}

	if parts[1] != ("v=" + strconv.Itoa(defaultVersion)) {
		return nil, nil, errors.New("incorrect version")
	}

	rmemory, err := splitAndGet(parts[2], "=", 1)
	if err != nil {
		return nil, nil, err
	}
	memory, err = strconv.ParseUint(rmemory, 10, 32)
	if err != nil {
		return nil, nil, errors.New("cannot convert memory value")
	}

	rtime, err := splitAndGet(parts[3], "=", 1)
	if err != nil {
		return nil, nil, err
	}
	time, err = strconv.ParseUint(rtime, 10, 32)
	if err != nil {
		return nil, nil, errors.New("cannot convert time value")
	}

	rthreads, err := splitAndGet(parts[4], "=", 1)
	if err != nil {
		return nil, nil, err
	}
	threads, err = strconv.ParseUint(rthreads, 10, 32)
	if err != nil {
		return nil, nil, errors.New("cannot convert threads value")
	}

	salt := parts[5]

	argon := New(
		func(a *ARGON2) {
			a.time = uint32(time)
			a.memory = uint32(memory)
			a.threads = uint8(threads)
		},
	)

	return argon, []byte(salt), nil
}

func splitAndGet(s, sep string, index int) (string, error) {
	split := strings.Split(s, sep)
	if len(split) <= index {
		return "", errors.New("not enough elements to segment")
	}
	return split[index], nil
}

func CompareHashAndPassword(hashed, pass []byte) error {
	argon, salt, err := newFromHash(string(hashed))
	if err != nil {
		return err
	}

	expected, err := GenerateFromPassword(pass, salt, argon)
	if err != nil {
		return err
	}

	if subtle.ConstantTimeCompare([]byte(expected), hashed) == 1 {
		return nil
	}

	return errors.New("not equal hash and password")
}
