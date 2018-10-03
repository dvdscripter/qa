package mongodb

import (
	"html"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"securecodewarrior.com/ddias/heapoverflow/model"
	"securecodewarrior.com/ddias/heapoverflow/model/storage"
)

func (db *DB) CreateComment(c model.Comment) (model.Comment, error) {
	conn := db.Copy()
	defer conn.Close()

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

	c.ID = db.getID(db.GetCommentC())

	if err := conn.DB(db.GetDatabase()).C(db.GetCommentC()).Insert(&c); err != nil {
		return model.Comment{}, errors.Wrap(err, "cannot create new comment")
	}

	return c, nil
}

func (db *DB) UpdateComment(c model.Comment) (model.Comment, error) {
	conn := db.Copy()
	defer conn.Close()

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

	if err := conn.DB(db.GetDatabase()).C(db.GetCommentC()).UpdateId(comment.ID, &comment); err != nil {
		return model.Comment{}, errors.Wrap(err, "cannot update comment")
	}

	return comment, nil
}

func (db *DB) FindComment(id int) (model.Comment, error) {
	conn := db.Copy()
	defer conn.Close()

	var comment model.Comment

	if err := conn.DB(db.GetDatabase()).C(db.GetCommentC()).FindId(id).One(&comment); err != nil {
		return model.Comment{}, storage.ErrCommentNotFound
	}

	return comment, nil
}

func (db *DB) FindCommentByAuthor(id int) ([]model.Comment, error) {
	conn := db.Copy()
	defer conn.Close()

	var comments []model.Comment

	if err := conn.DB(db.GetDatabase()).C(db.GetCommentC()).Find(bson.M{"user_id": id}).All(&comments); err != nil {
		return nil, storage.ErrCommentNotFound
	}

	return comments, nil
}

func (db *DB) FindCommentByQuestion(id int) ([]model.Comment, error) {
	conn := db.Copy()
	defer conn.Close()

	var comments []model.Comment

	if err := conn.DB(db.GetDatabase()).C(db.GetCommentC()).Find(bson.M{"question_id": id}).All(&comments); err != nil {
		return nil, storage.ErrCommentNotFound
	}

	return comments, nil
}

func (db *DB) UpComment(id int) error {
	conn := db.Copy()
	defer conn.Close()

	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"votes": 1}},
		ReturnNew: false,
	}

	if _, err := conn.DB(db.GetDatabase()).C(db.GetCommentC()).FindId(id).Apply(change, nil); err != nil {
		return storage.ErrCannotVote
	}

	return nil
}

func (db *DB) DownComment(id int) error {
	conn := db.Copy()
	defer conn.Close()

	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"votes": -1}},
		ReturnNew: false,
	}

	if _, err := conn.DB(db.GetDatabase()).C(db.GetCommentC()).FindId(id).Apply(change, nil); err != nil {
		return storage.ErrCannotVote
	}

	return nil
}
