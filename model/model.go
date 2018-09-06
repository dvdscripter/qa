package model

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"regexp"
	"strings"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/crypto/bcrypt"
)

const (
	defaultNickMaxSize     = 16
	defaultAvatarMaxSize   = 4096
	defaultMinPasswordSize = 7
	defaultMaxPasswordSize = 128

	defaultTitleMaxSize   = 140
	defaultTitleMinSize   = 30
	defaultContentMaxSize = 2000
)

var (
	ErrInvalidEmail    = errors.New("Invalid e-mail")
	ErrInvalidUser     = errors.New("Invalid User structure")
	ErrInvalidQuestion = errors.New("Invalid Question structure")
	ErrInvalidComment  = errors.New("Invalid Comment structure")
)

type User struct {
	ID       int       `json:"id,omitempty"`
	Since    time.Time `json:"since,omitempty"`
	Email    string    `json:"email,omitempty" gorm:"unique_index"`
	Nick     string    `json:"nick,omitempty" gorm:"unique_index"`
	Avatar   string    `json:"avatar,omitempty"`
	Password string    `json:"password,omitempty" gorm:"not null"`
}

type Question struct {
	ID       int       `json:"id,omitempty"`
	Title    string    `json:"title,omitempty" gorm:"unique_index"`
	Content  string    `json:"content,omitempty" gorm:"index"`
	Votes    int       `json:"votes,omitempty"`
	UserID   int       `json:"author,omitempty"`
	When     time.Time `json:"when,omitempty"`
	LastEdit time.Time `json:"last_edit,omitempty"`
}

type Comment struct {
	ID         int       `json:"id,omitempty"`
	QuestionID int       `json:"question,omitempty"`
	UserID     int       `json:"author,omitempty"`
	Content    string    `json:"content,omitempty"`
	Votes      int       `json:"votes,omitempty"`
	When       time.Time `json:"when,omitempty"`
	LastEdit   time.Time `json:"last_edit,omitempty"`
}

func GenPass(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
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

func (c Comment) validContent() error {
	if len(c.Content) > defaultContentMaxSize {
		return fmt.Errorf("Invalid content length must be below %d characters",
			defaultContentMaxSize)
	}
	return nil
}

func (c Comment) Valid() []error {
	validation := [](func() error){
		c.validContent,
	}

	var errFound []error
	for _, fn := range validation {
		err := fn()
		if err != nil {
			errFound = append(errFound, err)
		}
	}

	return errFound
}

func (u Question) validTitle() error {
	if len(u.Title) < defaultTitleMinSize && len(u.Title) > defaultTitleMaxSize {
		return fmt.Errorf("Invalid title length must be between %d and %d characters",
			defaultTitleMinSize, defaultTitleMaxSize)
	}
	return nil
}

func (u Question) validContent() error {
	if len(u.Content) > defaultContentMaxSize {
		return fmt.Errorf("Invalid content length must be below %d characters",
			defaultContentMaxSize)
	}
	return nil
}

func (q Question) Valid() []error {
	validation := [](func() error){
		q.validContent,
		q.validTitle,
	}

	var errFound []error
	for _, fn := range validation {
		err := fn()
		if err != nil {
			errFound = append(errFound, err)
		}
	}

	return errFound
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
		return fmt.Errorf("Invalid nick format, only letters, digits and hyphens")
	}
	return nil
}

func (u User) ValidAvatar() error {
	if u.Avatar == "" {
		return nil
	}
	avatar, err := base64.StdEncoding.DecodeString(u.Avatar)
	if err != nil {
		return fmt.Errorf("Cannot decode avatar")
	}
	if len(avatar) > defaultAvatarMaxSize {
		return fmt.Errorf("Max avatar size: %d", defaultAvatarMaxSize)
	}
	avatarBuffer := bytes.NewBuffer(avatar)
	avatarImg, _, err := image.Decode(avatarBuffer)
	if err != nil {
		return fmt.Errorf("Cannot decode avatar image")
	}
	dim := defaultAvatarMaxSize / 4
	if avatarImg.Bounds().Max.X > dim || avatarImg.Bounds().Max.Y > dim {
		return fmt.Errorf("Avatar exceeds %d dimensions", dim)
	}
	return nil
}

func (u User) ValidPassword() error {
	if len(u.Password) < 7 || len(u.Password) > 128 {
		return fmt.Errorf("Invalid password length must be between %d and %d characters",
			defaultMinPasswordSize, defaultMaxPasswordSize)
	}
	return nil
}

func (u User) Valid() []error {
	validation := [](func() error){
		u.validEmail,
		u.ValidAvatar,
		u.ValidNick,
		u.ValidPassword,
	}

	var errFound []error
	for _, fn := range validation {
		err := fn()
		if err != nil {
			errFound = append(errFound, err)
		}
	}

	return errFound
}
