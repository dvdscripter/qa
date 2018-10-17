package model

import (
	"bytes"
	"image"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/vincent-petithory/dataurl"
	"securecodewarrior.com/ddias/heapoverflow/crypto/argon2"
)

const (
	defaultNickMaxSize     = 16
	defaultAvatarMaxSize   = 4096
	defaultAvatarDim       = 40
	defaultMinPasswordSize = 7
	defaultMaxPasswordSize = 128
)

type User struct {
	ID       int       `json:"id" bson:"_id"`
	Since    time.Time `json:"since,omitempty"`
	Email    string    `json:"email,omitempty" gorm:"unique_index;size:320"`
	Nick     string    `json:"nick,omitempty" gorm:"unique_index;size:16"`
	Avatar   string    `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Password string    `json:"password,omitempty" gorm:"not null"`
}

func GenPass(password string) (string, error) {
	// pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	pass, err := argon2.GenerateFromPassword([]byte(password), nil, nil)
	if err != nil {
		return "", errors.Wrap(err, "cannot generate hash")
	}
	return string(pass), nil
}

// need to use copy because slice receiver/param can have elements changed even
// without pointer receiver
func OmitPass(users []User) []User {
	usersPassOmit := make([]User, len(users))
	copy(usersPassOmit, users)
	for i := range usersPassOmit {
		usersPassOmit[i].Password = ""
	}
	return usersPassOmit
}

func (u User) validEmail() error {
	digits := regexp.MustCompile(`^\d+$`)
	emailRE := regexp.MustCompile(`^[\w` +
		regexp.QuoteMeta("!#$%&‘*+–/=?^_`.{|}~") + `]+@[\w\.]+$`)
	if emailRE.MatchString(u.Email) == false {
		return ErrInvalidEmail
	}
	email := strings.Split(u.Email, "@")
	if len(email) != 2 {
		return ErrInvalidEmail
	}
	local, domain := email[0], email[1]
	if len(local) > 64 || len(domain) > 255 {
		return ErrInvalidEmail
	}

	for _, dnsLabel := range strings.Split(domain, ".") {
		if len(dnsLabel) > 63 || digits.MatchString(dnsLabel) ||
			dnsLabel[0] == '-' || dnsLabel[len(dnsLabel)-1] == '-' {
			return ErrInvalidEmail
		}
	}

	if local[0] == '.' || local[len(local)-1] == '.' ||
		strings.Contains(local, "..") {
		return ErrInvalidEmail
	}

	return nil
}

func (u User) ValidNick() error {
	reNick := regexp.MustCompile(`^\w+$`)
	if len(u.Nick) > defaultNickMaxSize || !reNick.MatchString(u.Nick) {
		return errors.Errorf("Invalid nick format, only letters, digits and hyphens")
	}
	return nil
}

func (u User) ValidAvatar() error {
	if u.Avatar == "" {
		return nil
	}
	avatar, err := dataurl.DecodeString(u.Avatar)
	if err != nil {
		return errors.Wrap(err, "cannot decode avatar")
	}
	if len(avatar.Data) > defaultAvatarMaxSize {
		return errors.Errorf("Max avatar size: %d", defaultAvatarMaxSize)
	}
	avatarBuffer := bytes.NewBuffer(avatar.Data)
	avatarImg, _, err := image.Decode(avatarBuffer)
	if err != nil {
		return errors.Errorf("Cannot decode avatar image")
	}
	if avatarImg.Bounds().Max.X > defaultAvatarDim ||
		avatarImg.Bounds().Max.Y > defaultAvatarDim {
		return errors.Errorf("Avatar exceeds %d dimensions", defaultAvatarDim)
	}
	return nil
}

func (u User) ValidPassword() error {

	rules := []passRules{
		{oneDigit, "at least 1 digit (0-9)"},
		{oneSpecialCase, "at least 1 special character"},
		{oneLowerCase, "at least 1 lowercase character (a-z)"},
		{oneUpperCase, "at least 1 uppercase character (A-Z)"},
	}

	rulesFailing := []string{}

	for _, rule := range rules {
		if !rule.fn(u.Password) {
			rulesFailing = append(rulesFailing, rule.missing)
		}
	}

	if len(rulesFailing) > 1 {
		return errors.Errorf("Invalid password %s", strings.Join(rulesFailing, "\n"))
	}

	if len(u.Password) < 10 || len(u.Password) > 128 {
		return errors.Errorf("Invalid password: length must be between %d and %d characters",
			defaultMinPasswordSize, defaultMaxPasswordSize)
	}

	if seqOf(u.Password) {
		return errors.Errorf("Invalid password: not more than 2 identical characters in a row (e.g., 111 not allowed)")
	}

	return nil
}

func (u User) Valid() error {
	validation := [](func() error){
		u.validEmail,
		u.ValidAvatar,
		u.ValidNick,
		u.ValidPassword,
	}

	var errFound []string
	for _, fn := range validation {
		err := fn()
		if err != nil {
			errFound = append(errFound, err.Error())
		}
	}
	if errFound == nil {
		return nil
	}

	return errors.Errorf("Invalid user model: %s", strings.Join(errFound, "\n"))
}
