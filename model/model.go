package model

import (
	"regexp"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/pkg/errors"
)

const (
	defaultContentMaxSize = 2000
)

var (
	ErrInvalidEmail    = errors.New("Invalid e-mail")
	ErrInvalidUser     = errors.New("Invalid User structure")
	ErrInvalidQuestion = errors.New("Invalid Question structure")
	ErrInvalidComment  = errors.New("Invalid Comment structure")
)

func oneUpperCase(s string) bool {
	re := regexp.MustCompile(`[A-Z]`)
	return re.MatchString(s)
}

func oneLowerCase(s string) bool {
	re := regexp.MustCompile(`[a-z]`)
	return re.MatchString(s)
}

func oneDigit(s string) bool {
	re := regexp.MustCompile(`[0-9]`)
	return re.MatchString(s)
}

func oneSpecialCase(s string) bool {
	re := regexp.MustCompile("[" + regexp.QuoteMeta("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~") + "]")
	return re.MatchString(s)
}

type passRules struct {
	fn      func(string) bool
	missing string
}

func seqOf(s string) bool {
	for i := 0; i < (len(s) - 1); i++ {
		if s[i] == s[i+1] {
			return true
		}
	}
	return false
}
