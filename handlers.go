package main

import (
	"net/http"

	"github.com/pkg/errors"
	"securecodewarrior.com/ddias/heapoverflow/model"

	"securecodewarrior.com/ddias/heapoverflow/jwt"
)

type appHandler func(w http.ResponseWriter, r *http.Request) (interface{}, error)

func (app *app) CreateUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var user model.User
	err := jsonFromRequest(&user, r)
	if err != nil {
		return nil, err
	}

	payload := jwt.DecodePayload(r)
	if payload.Email != "" {
		return nil, errors.Errorf("Already logged")
	}

	return app.Storage.CreateUser(user)
}

func (app *app) RetrieveUsers(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return app.Storage.FindAllUser()
}

func (app *app) RetrieveUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id, err := idFromRequest("id", r)
	if err != nil {
		return nil, err
	}

	return app.Storage.FindUser(id)
}

func (app *app) DeleteUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id, err := idFromRequest("id", r)
	if err != nil {
		return nil, err
	}

	user, err := app.Storage.FindUser(id)
	if err != nil {
		return nil, err
	}

	payload := jwt.DecodePayload(r)
	if payload.Email != user.Email {
		return nil, errors.Errorf("Cannot delete another user")
	}

	return nil, app.Storage.DeleteUser(id)
}

func (app *app) UpdateUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var user model.User
	err := jsonFromRequest(&user, r)
	if err != nil {
		return nil, err
	}

	payload := jwt.DecodePayload(r)
	if payload.Email != user.Email {
		return nil, errors.Errorf("Cannot update another user")
	}

	id, err := idFromRequest("id", r)
	if err != nil {
		return nil, err
	}

	user.ID = id

	return app.Storage.UpdateUser(user)
}

func (app *app) Login(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var user model.User
	err := jsonFromRequest(&user, r)
	if err != nil {
		return nil, err
	}

	if err := app.Storage.Login(user.Email, user.Password); err != nil {
		return nil, err
	}

	payload := jwt.Payload{
		Email: user.Email,
		Exp:   jwt.DefaultExpiration,
	}

	token := jwt.NewFromFile(payload, app.jwtKeyFile)
	return token.Encode()
}

func (app *app) RetrieveQuestions(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return app.FindAllQuestion()
}

func (app *app) RetrieveQuestion(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id, err := idFromRequest("id", r)
	if err != nil {
		return nil, err
	}

	return app.FindQuestion(id)
}

func (app *app) CreateQuestion(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var question model.Question
	err := jsonFromRequest(&question, r)
	if err != nil {
		return nil, err
	}
	payload := jwt.DecodePayload(r)
	user, err := app.Storage.FindUserByEmail(payload.Email)
	if err != nil {
		return nil, err
	}
	question.UserID = user.ID

	return app.Storage.CreateQuestion(question)
}

func (app *app) UpdateQuestion(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var question model.Question
	err := jsonFromRequest(&question, r)
	if err != nil {
		return nil, err
	}
	id, err := idFromRequest("id", r)
	if err != nil {
		return nil, err
	}

	payload := jwt.DecodePayload(r)
	qstore, err := app.Storage.FindQuestion(id)
	if err != nil {
		return nil, err
	}
	author, err := app.Storage.FindUser(qstore.UserID)
	if err != nil {
		return nil, err
	}
	if payload.Email != author.Email {
		return nil, errors.Errorf("Cannot update another user question")
	}

	question.ID = id

	return app.Storage.UpdateQuestion(question)
}

func (app *app) UpVoteQuestion(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id, err := idFromRequest("id", r)
	if err != nil {
		return nil, err
	}
	payload := jwt.DecodePayload(r)
	question, err := app.Storage.FindQuestion(id)
	if err != nil {
		return nil, err
	}
	author, err := app.Storage.FindUser(question.UserID)
	if err != nil {
		return nil, err
	}
	if payload.Email == author.Email {
		return nil, errors.Errorf("Cannot up vote yourself")
	}

	return nil, app.Storage.UpQuestion(id)
}

func (app *app) DownVoteQuestion(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id, err := idFromRequest("id", r)
	if err != nil {
		return nil, err
	}
	payload := jwt.DecodePayload(r)
	question, err := app.Storage.FindQuestion(id)
	if err != nil {
		return nil, err
	}
	author, err := app.Storage.FindUser(question.UserID)
	if err != nil {
		return nil, err
	}
	if payload.Email == author.Email {
		return nil, errors.Errorf("Cannot up vote yourself")
	}

	return nil, app.Storage.DownQuestion(id)
}

func (app *app) RetrieveQuestionComments(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id, err := idFromRequest("id", r)
	if err != nil {
		return nil, err
	}

	return app.Storage.FindCommentByQuestion(id)
}

func (app *app) RetrieveQuestionComment(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id, err := idFromRequest("id", r)
	if err != nil {
		return nil, err
	}
	cid, err := idFromRequest("cid", r)
	if err != nil {
		return nil, err
	}

	if _, err := app.Storage.FindQuestion(id); err != nil {
		return nil, err
	}

	return app.Storage.FindComment(cid)
}

func (app *app) CreateQuestionComments(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var comment model.Comment
	if err := jsonFromRequest(&comment, r); err != nil {
		return nil, err
	}
	id, err := idFromRequest("id", r)
	if err != nil {
		return nil, err
	}
	payload := jwt.DecodePayload(r)
	user, err := app.Storage.FindUserByEmail(payload.Email)
	if err != nil {
		return nil, err
	}
	comment.UserID = user.ID
	comment.QuestionID = id

	return app.Storage.CreateComment(comment)
}

func (app *app) UpdateQuestionComment(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var comment model.Comment
	if err := jsonFromRequest(&comment, r); err != nil {
		return nil, err
	}
	id, err := idFromRequest("id", r)
	if err != nil {
		return nil, err
	}
	cid, err := idFromRequest("cid", r)
	if err != nil {
		return nil, err
	}

	payload := jwt.DecodePayload(r)
	if _, err := app.Storage.FindQuestion(id); err != nil {
		return nil, err
	}
	cstore, err := app.Storage.FindComment(cid)
	if err != nil {
		return nil, err
	}
	author, err := app.Storage.FindUser(cstore.UserID)
	if err != nil {
		return nil, err
	}
	if payload.Email != author.Email {
		return nil, errors.Errorf("Cannot update another user comment")
	}

	comment.ID = cid
	comment.QuestionID = id

	return app.Storage.UpdateComment(comment)
}

func (app *app) UpVoteQuestionComment(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id, err := idFromRequest("id", r)
	if err != nil {
		return nil, err
	}
	cid, err := idFromRequest("cid", r)
	if err != nil {
		return nil, err
	}

	if _, err := app.Storage.FindQuestion(id); err != nil {
		return nil, err
	}
	comment, err := app.Storage.FindComment(cid)
	if err != nil {
		return nil, err
	}
	author, err := app.Storage.FindUser(comment.UserID)
	if err != nil {
		return nil, err
	}
	payload := jwt.DecodePayload(r)
	if payload.Email == author.Email {
		return nil, errors.Errorf("Cannot up vote yourself")
	}

	return nil, app.Storage.UpComment(cid)
}

func (app *app) DownVoteQuestionComment(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	id, err := idFromRequest("id", r)
	if err != nil {
		return nil, err
	}
	cid, err := idFromRequest("cid", r)
	if err != nil {
		return nil, err
	}

	if _, err := app.Storage.FindQuestion(id); err != nil {
		return nil, err
	}
	comment, err := app.Storage.FindComment(cid)
	if err != nil {
		return nil, err
	}
	author, err := app.Storage.FindUser(comment.UserID)
	if err != nil {
		return nil, err
	}
	payload := jwt.DecodePayload(r)
	if payload.Email == author.Email {
		return nil, errors.Errorf("Cannot up vote yourself")
	}

	return nil, app.Storage.DownComment(cid)
}
