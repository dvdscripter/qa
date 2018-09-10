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

func (db *DB) FindAllQuestion() ([]model.Question, error) {
	conn := db.Copy()
	defer conn.Close()

	var questions []model.Question
	if err := conn.DB(db.GetDatabase()).C(db.GetQuestionC()).Find(nil).All(&questions); err != nil {
		return nil, err
	}
	return questions, nil
}

func (db *DB) CreateQuestion(q model.Question) (model.Question, error) {
	conn := db.Copy()
	defer conn.Close()

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
		return model.Question{}, model.ErrInvalidQuestion
	}

	nid, err := conn.DB(db.GetDatabase()).C(db.GetQuestionC()).Find(nil).Count()
	if err != nil {
		return model.Question{}, errors.Wrap(err, "cannot get new ID")
	}

	q.ID = nid

	if err := conn.DB(db.GetDatabase()).C(db.GetQuestionC()).Insert(&q); err != nil {
		return model.Question{}, errors.Wrap(err, "cannot create new question")
	}

	return q, nil
}

func (db *DB) UpdateQuestion(q model.Question) (model.Question, error) {
	conn := db.Copy()
	defer conn.Close()

	if err := q.Valid(); err != nil {
		return model.Question{}, model.ErrInvalidQuestion
	}

	question, err := db.FindQuestion(q.ID)
	if err != nil {
		return model.Question{}, storage.ErrQuestionNotFound
	}

	question.Title = html.EscapeString(q.Title)
	question.Content = html.EscapeString(q.Content)
	question.LastEdit = time.Now()
	if err := conn.DB(db.GetDatabase()).C(db.GetQuestionC()).UpdateId(question.ID, &question); err != nil {
		return model.Question{}, errors.Wrap(err, "cannot update question")
	}

	return question, nil
}

func (db *DB) FindQuestion(id int) (model.Question, error) {
	conn := db.Copy()
	defer conn.Close()

	var question model.Question

	if err := conn.DB(db.GetDatabase()).C(db.GetQuestionC()).FindId(id).One(&question); err != nil {
		return model.Question{}, storage.ErrQuestionNotFound
	}

	return question, nil
}

func (db *DB) FindQuestionByTitle(title string) (model.Question, error) {
	conn := db.Copy()
	defer conn.Close()

	var question model.Question

	if err := conn.DB(db.GetDatabase()).C(db.GetQuestionC()).Find(bson.M{"title": title}).One(&question); err != nil {
		return model.Question{}, storage.ErrQuestionNotFound
	}

	return question, nil
}

func (db *DB) FindQuestionByAuthor(author int) ([]model.Question, error) {
	conn := db.Copy()
	defer conn.Close()

	var questions []model.Question

	if err := conn.DB(db.GetDatabase()).C(db.GetQuestionC()).Find(bson.M{"user_id": author}).All(&questions); err != nil {
		return nil, storage.ErrQuestionNotFound
	}

	return questions, nil
}

func (db *DB) UpQuestion(id int) error {
	conn := db.Copy()
	defer conn.Close()

	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"votes": 1}},
		ReturnNew: false,
	}

	if _, err := conn.DB(db.GetDatabase()).C(db.GetQuestionC()).FindId(id).Apply(change, nil); err != nil {
		return storage.ErrCannotVote
	}

	return nil
}

func (db *DB) DownQuestion(id int) error {
	conn := db.Copy()
	defer conn.Close()

	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"votes": -1}},
		ReturnNew: false,
	}

	if _, err := conn.DB(db.GetDatabase()).C(db.GetQuestionC()).FindId(id).Apply(change, nil); err != nil {
		return storage.ErrCannotVote
	}

	return nil
}
