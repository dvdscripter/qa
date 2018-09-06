package memory

import (
	"securecodewarrior.com/ddias/heapoverflow/model"
)

type DB struct {
	users     []model.User
	questions []model.Question
	comments  []model.Comment
}

func New() *DB {
	return &DB{}
}
