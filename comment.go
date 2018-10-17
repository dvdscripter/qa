package main

import (
	"net/http"

	"github.com/pkg/errors"
	"securecodewarrior.com/ddias/heapoverflow/jwt"
	"securecodewarrior.com/ddias/heapoverflow/model"
)

func (app *app) RetrieveQuestionComments(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

	id, err := idFromRequest("id", r)
	if err != nil {
		return nil, err
	}

	return app.Storage.FindCommentByQuestion(id)
}

func (app *app) RetrieveQuestionComment(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

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

func (app *app) CreateQuestionComments(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

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

func (app *app) UpdateQuestionComment(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

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

func (app *app) UpVoteQuestionComment(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

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

func (app *app) DownVoteQuestionComment(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

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
