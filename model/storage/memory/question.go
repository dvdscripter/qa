package memory

import (
	"html"
	"time"

	"securecodewarrior.com/ddias/heapoverflow/model"
	"securecodewarrior.com/ddias/heapoverflow/model/storage"
)

func (db *DB) FindAllQuestion() ([]model.Question, error) {
	return db.questions, nil
}

func (db *DB) CreateQuestion(q model.Question) (model.Question, error) {
	if _, err := db.FindQuestionByTitle(q.Title); err == nil {
		return model.Question{}, storage.ErrQuestionAlreadyExist
	}

	if _, err := db.FindUser(q.UserID); err != nil {
		return model.Question{}, storage.ErrUserNotFound
	}

	q.ID = len(db.questions) + 1
	q.When = time.Now()
	q.LastEdit = time.Now()
	q.Votes = 0
	q.Title = html.EscapeString(q.Title)
	q.Content = html.EscapeString(q.Content)

	if err := q.Valid(); err != nil {
		return model.Question{}, err
	}

	db.questions = append(db.questions, q)

	return q, nil
}

func (db *DB) UpdateQuestion(q model.Question) (model.Question, error) {
	if err := q.Valid(); err != nil {
		return model.Question{}, err
	}
	for i, question := range db.questions {
		if q.ID == question.ID {
			db.questions[i].Title = html.EscapeString(q.Title)
			db.questions[i].Content = html.EscapeString(q.Content)
			db.questions[i].LastEdit = time.Now()
			return db.questions[i], nil
		}
	}
	return model.Question{}, storage.ErrQuestionNotFound
}

func (db *DB) FindQuestion(id int) (model.Question, error) {
	for _, question := range db.questions {
		if id == question.ID {
			return question, nil
		}
	}
	return model.Question{}, storage.ErrQuestionNotFound
}

func (db *DB) FindQuestionByTitle(title string) (model.Question, error) {
	for _, question := range db.questions {
		if question.Title == title {
			return question, nil
		}
	}
	return model.Question{}, storage.ErrQuestionNotFound
}

func (db *DB) FindQuestionByAuthor(author int) ([]model.Question, error) {
	found := []model.Question{}
	for _, question := range db.questions {
		if question.UserID == author {
			found = append(found, question)
		}
	}
	if len(found) == 0 {
		return nil, storage.ErrQuestionNotFound
	}
	return found, nil
}

func (db *DB) UpQuestion(id int) error {
	question, err := db.FindQuestion(id)
	if err != nil {
		return err
	}

	question.Votes++
	return nil
}

func (db *DB) DownQuestion(id int) error {
	question, err := db.FindQuestion(id)
	if err != nil {
		return err
	}

	question.Votes--
	return nil
}
