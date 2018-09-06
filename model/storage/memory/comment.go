package memory

import (
	"html"
	"time"

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

	c.ID = len(db.comments) + 1
	c.When = time.Now()
	c.LastEdit = time.Now()
	c.QuestionID = question.ID
	c.Content = html.EscapeString(c.Content)
	c.Votes = 0

	if err := c.Valid(); err != nil {
		return model.Comment{}, model.ErrInvalidComment
	}

	db.comments = append(db.comments, c)

	return c, nil
}

func (db *DB) UpdateComment(c model.Comment) (model.Comment, error) {
	if err := c.Valid(); err != nil {
		return model.Comment{}, model.ErrInvalidComment
	}
	if _, err := db.FindQuestion(c.QuestionID); err != nil {
		return model.Comment{}, storage.ErrQuestionNotFound
	}

	for i, comment := range db.comments {
		if c.ID == comment.ID {
			db.comments[i].Content = html.EscapeString(c.Content)
			db.comments[i].LastEdit = time.Now()
			return db.comments[i], nil
		}
	}
	return model.Comment{}, storage.ErrCommentNotFound
}

func (db *DB) FindComment(id int) (model.Comment, error) {
	for _, comment := range db.comments {
		if id == comment.ID {
			return comment, nil
		}
	}
	return model.Comment{}, storage.ErrCommentNotFound
}

func (db *DB) FindCommentByAuthor(author int) ([]model.Comment, error) {
	found := []model.Comment{}
	for _, comment := range db.comments {
		if comment.UserID == author {
			found = append(found, comment)
		}
	}
	if len(found) == 0 {
		return nil, storage.ErrCommentNotFound
	}
	return found, nil
}

func (db *DB) FindCommentByQuestion(question int) ([]model.Comment, error) {
	if _, err := db.FindQuestion(question); err != nil {
		return nil, storage.ErrQuestionNotFound
	}

	found := []model.Comment{}
	for _, comment := range db.comments {
		if comment.QuestionID == question {
			found = append(found, comment)
		}
	}
	if len(found) == 0 {
		return nil, storage.ErrCommentNotFound
	}
	return found, nil
}

func (db *DB) UpComment(id int) error {
	comment, err := db.FindComment(id)
	if err != nil {
		return err
	}

	comment.Votes++
	return nil
}

func (db *DB) DownComment(id int) error {
	comment, err := db.FindComment(id)
	if err != nil {
		return err
	}

	comment.Votes--
	return nil
}
