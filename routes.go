package main

import (
	"net/http"
)

type route struct {
	pattern string
	method  string
	handler appHandler
	public  bool
}

var routes = []route{
	{"/login", "POST", webapp.Login, true},

	{"/user", "POST", webapp.CreateUser, true},
	{"/user", "GET", webapp.RetrieveUsers, false},
	{"/user/{id:[0-9]+}", "GET", webapp.RetrieveUser, false},
	{"/user/{email}", "GET", webapp.RetrieveUserByEmail, false},
	{"/user/{id:[0-9]+}", "DELETE", webapp.DeleteUser, false},
	{"/user/{id:[0-9]+}", "PUT", webapp.UpdateUser, false},

	{"/question", "POST", webapp.CreateQuestion, false},
	{"/question", "GET", webapp.RetrieveQuestions, false},
	{"/question/{id:[0-9]+}", "GET", webapp.RetrieveQuestion, false},
	{"/question/{id:[0-9]+}", "PUT", webapp.UpdateQuestion, false},
	{"/question/{id:[0-9]+}/vote", "PUT", webapp.UpVoteQuestion, false},
	{"/question/{id:[0-9]+}/vote", "DELETE", webapp.DownVoteQuestion, false},

	{"/question/{id:[0-9]+}/comments", "POST", webapp.CreateQuestionComments, false},
	{"/question/{id:[0-9]+}/comments", "GET", webapp.RetrieveQuestionComments, false},
	{"/question/{id:[0-9]+}/comments/{cid:[0-9]+}", "GET", webapp.RetrieveQuestionComment, false},
	{"/question/{id:[0-9]+}/comments/{cid:[0-9]+}", "PUT", webapp.UpdateQuestionComment, false},
	{"/question/{id:[0-9]+}/comments/{cid:[0-9]+}/vote", "PUT", webapp.UpVoteQuestionComment, false},
	{"/question/{id:[0-9]+}/comments/{cid:[0-9]+}/vote", "DELETE", webapp.DownVoteQuestionComment, false},
}

func (app *app) registerRoutes(logger func(appHandler) http.Handler) {

	for _, route := range app.routes {
		app.router.Handle(route.pattern, logger(route.handler)).Methods(route.method)
	}

}

func (app *app) isPublic(r *http.Request) bool {
	for _, route := range app.routes {
		if route.pattern == r.URL.Path && r.Method == route.method && route.public {
			return true
		}
	}
	return false
}
