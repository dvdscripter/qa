package model

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	defaultTitleMaxSize = 140
	defaultTitleMinSize = 30
)

type Question struct {
	ID       int       `json:"id" bson:"_id"`
	Title    string    `json:"title,omitempty" gorm:"unique_index;size:140"`
	Content  string    `json:"content,omitempty" gorm:"index;size:2000"`
	Votes    int       `json:"votes"`
	UserID   int       `json:"author" bson:"user_id"`
	When     time.Time `json:"when,omitempty"`
	LastEdit time.Time `json:"last_edit,omitempty"`
}

func (u Question) validTitle() error {
	if len(u.Title) < defaultTitleMinSize || len(u.Title) > defaultTitleMaxSize {
		return errors.Errorf("Invalid title length must be between %d and %d characters",
			defaultTitleMinSize, defaultTitleMaxSize)
	}
	return nil
}

func (u Question) validContent() error {
	if len(u.Content) > defaultContentMaxSize {
		return errors.Errorf("Invalid content length must be below %d characters",
			defaultContentMaxSize)
	}
	return nil
}

func (q Question) Valid() error {
	validation := [](func() error){
		q.validContent,
		q.validTitle,
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

	return errors.Errorf("Invalid question model: %s", strings.Join(errFound, "\n"))
}
