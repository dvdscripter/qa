package sql

import (
	"html"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"securecodewarrior.com/ddias/heapoverflow/model"
	"securecodewarrior.com/ddias/heapoverflow/model/storage"
)

func (db *DB) FindAllQuestion() ([]model.Question, error) {
	var questions []model.Question
	if err := db.Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func (db *DB) CreateQuestion(q model.Question) (model.Question, error) {
	if _, err := db.FindQuestionByTitle(q.Title); err == nil {
		return model.Question{}, storage.ErrQuestionAlreadyExist
	}

	if _, err := db.FindUser(q.UserID); err != nil {
		return model.Question{}, storage.ErrUserNotFound
	}

	q.When = time.Now()
	q.LastEdit = time.Now()
	q.Votes = 0
	q.Title = html.EscapeString(q.Title)
	q.Content = html.EscapeString(q.Content)

	if err := q.Valid(); err != nil {
		return model.Question{}, err
	}

	if err := db.Create(&q).Error; err != nil {
		return model.Question{}, err
	}

	return q, nil
}

func (db *DB) UpdateQuestion(q model.Question) (model.Question, error) {
	if err := q.Valid(); err != nil {
		return model.Question{}, err
	}

	question, err := db.FindQuestion(q.ID)
	if err != nil {
		return model.Question{}, storage.ErrQuestionNotFound
	}

	question.Title = html.EscapeString(q.Title)
	question.Content = html.EscapeString(q.Content)
	question.LastEdit = time.Now()
	if err := db.Save(&question).Error; err != nil {
		return model.Question{}, err
	}

	return question, nil
}

func (db *DB) FindQuestion(id int) (model.Question, error) {
	var question model.Question

	if err := db.First(&question, id).Error; err != nil {
		return model.Question{}, storage.ErrQuestionNotFound
	}

	return question, nil
}

func (db *DB) FindQuestionByTitle(title string) (model.Question, error) {
	var question model.Question

	if err := db.Where("title = ?", title).First(&question).Error; err != nil {
		return model.Question{}, storage.ErrQuestionNotFound
	}

	return question, nil
}

func (db *DB) FindQuestionByAuthor(author int) ([]model.Question, error) {
	var question []model.Question

	if err := db.Where("user_id = ?", author).Find(&question).Error; err != nil {
		return nil, storage.ErrQuestionNotFound
	}

	return question, nil
}

func (db *DB) UpQuestion(id int) error {
	question, err := db.FindQuestion(id)
	if err != nil {
		return storage.ErrQuestionNotFound
	}
	if err := db.Model(&question).UpdateColumn("votes",
		gorm.Expr("votes + ?", 1)).Error; err != nil {
		return storage.ErrCannotVote
	}
	return nil
}

func (db *DB) DownQuestion(id int) error {
	question, err := db.FindQuestion(id)
	if err != nil {
		return storage.ErrQuestionNotFound
	}
	if err := db.Model(&question).UpdateColumn("votes",
		gorm.Expr("votes - ?", 1)).Error; err != nil {
		return storage.ErrCannotVote
	}
	return nil
}
