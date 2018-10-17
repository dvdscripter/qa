package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"securecodewarrior.com/ddias/heapoverflow/jwt"
	"securecodewarrior.com/ddias/heapoverflow/model"
)

func (app *app) CreateUser(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

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

func (app *app) RetrieveUsers(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

	return app.Storage.FindAllUser()
}

func (app *app) RetrieveUser(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

	id, err := idFromRequest("id", r)
	if err != nil {
		return nil, err
	}

	return app.Storage.FindUser(id)
}

func (app *app) RetrieveUserByEmail(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

	params := mux.Vars(r)
	email, exist := params["email"]
	if !exist {
		return nil, errors.Errorf("Missing email parameter")
	}

	user, err := app.Storage.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	user.Password = ""

	return user, nil
}

func (app *app) DeleteUser(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

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

func (app *app) UpdateUser(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

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

func (app *app) Login(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {

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
