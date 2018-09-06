package sql

import (
	"html"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"securecodewarrior.com/ddias/heapoverflow/model"
	"securecodewarrior.com/ddias/heapoverflow/model/storage"
)

func (db *DB) CreateComment(c model.Comment) (model.Comment, error) {
	question, err := db.FindQuestion(c.QuestionID)
	if err != nil {
		return model.Comment{}, storage.ErrQuestionNotFound
	}

	if _, err := db.FindUser(c.UserID); err != nil {
		return model.Comment{}, storage.ErrUserNotFound
	}

	c.When = time.Now()
	c.LastEdit = time.Now()
	c.QuestionID = question.ID
	c.Content = html.EscapeString(c.Content)
	c.Votes = 0

	if err := c.Valid(); err != nil {
		return model.Comment{}, model.ErrInvalidComment
	}

	if err := db.Create(&c).Error; err != nil {
		return model.Comment{}, err
	}

	return c, nil
}
func (db *DB) UpdateComment(c model.Comment) (model.Comment, error) {
	if err := c.Valid(); err != nil {
		return model.Comment{}, model.ErrInvalidComment
	}
	if _, err := db.FindQuestion(c.QuestionID); err != nil {
		return model.Comment{}, storage.ErrQuestionNotFound
	}
	if _, err := db.FindUser(c.UserID); err != nil {
		return model.Comment{}, storage.ErrUserNotFound
	}

	comment, err := db.FindComment(c.ID)
	if err != nil {
		return model.Comment{}, storage.ErrCommentNotFound
	}

	comment.Content = html.EscapeString(c.Content)
	comment.LastEdit = time.Now()

	if err := db.Save(&comment).Error; err != nil {
		return model.Comment{}, err
	}

	return comment, nil
}
func (db *DB) FindComment(id int) (model.Comment, error) {
	var comment model.Comment

	if err := db.First(&comment, id).Error; err != nil {
		return model.Comment{}, storage.ErrCommentNotFound
	}

	return comment, nil
}
func (db *DB) FindCommentByAuthor(id int) ([]model.Comment, error) {
	var comment []model.Comment

	if err := db.Where("UserID = ?").Find(&comment).Error; err != nil {
		return nil, storage.ErrCommentNotFound
	}

	return comment, nil
}

func (db *DB) FindCommentByQuestion(id int) ([]model.Comment, error) {
	var comment []model.Comment

	if err := db.Where("QuestionID = ?").Find(&comment).Error; err != nil {
		return nil, storage.ErrCommentNotFound
	}

	return comment, nil

}
func (db *DB) UpComment(id int) error {
	comment, err := db.FindComment(id)
	if err != nil {
		return storage.ErrQuestionNotFound
	}
	if err := db.Model(&comment).UpdateColumn("votes",
		gorm.Expr("votes + ?", 1)).Error; err != nil {
		return storage.ErrCannotVote
	}
	return nil
}
func (db *DB) DownComment(id int) error {
	comment, err := db.FindComment(id)
	if err != nil {
		return storage.ErrQuestionNotFound
	}
	if err := db.Model(&comment).UpdateColumn("votes",
		gorm.Expr("votes + ?", 1)).Error; err != nil {
		return storage.ErrCannotVote
	}
	return nil
}
