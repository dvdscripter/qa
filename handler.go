package main

import (
	"net/http"
)

type appHandler func(w http.ResponseWriter, r *http.Request) (interface{},
	error)
