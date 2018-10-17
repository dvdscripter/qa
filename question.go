package main

import (
	"net/http"

	"github.com/pkg/errors"
	"securecodewarrior.com/ddias/heapoverflow/jwt"
	"securecodewarrior.com/ddias/heapoverflow/model"
)

func (app *app) RetrieveQuestions(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

	return app.FindAllQuestion()
}

func (app *app) RetrieveQuestion(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

	id, err := idFromRequest("id", r)
	if err != nil {
		return nil, err
	}

	return app.FindQuestion(id)
}

func (app *app) CreateQuestion(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

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

func (app *app) UpdateQuestion(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

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

func (app *app) UpVoteQuestion(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

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

func (app *app) DownVoteQuestion(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

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
