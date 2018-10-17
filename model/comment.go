package model

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Comment struct {
	ID         int       `json:"id" bson:"_id"`
	QuestionID int       `json:"question" bson:"question_id"`
	UserID     int       `json:"author" bson:"user_id"`
	Content    string    `json:"content,omitempty;size:2000"`
	Votes      int       `json:"votes"`
	When       time.Time `json:"when,omitempty"`
	LastEdit   time.Time `json:"last_edit,omitempty" bson:"last_edit"`
}

func (c Comment) validContent() error {
	if len(c.Content) > defaultContentMaxSize {
		return errors.Errorf("Invalid content length must be below %d characters",
			defaultContentMaxSize)
	}
	return nil
}

func (c Comment) Valid() error {
	validation := [](func() error){
		c.validContent,
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

	return errors.Errorf("%s", strings.Join(errFound, "\n"))
}
