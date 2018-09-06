package storage

import (
	"errors"

	"securecodewarrior.com/ddias/heapoverflow/model"
)

type Storage interface {
	UserStorage
	QuestionStorage
	CommentStorage
}

var (
	ErrUserNotFound         = errors.New("User not found")
	ErrQuestionNotFound     = errors.New("Question not found")
	ErrCommentNotFound      = errors.New("Comment not found")
	ErrQuestionAlreadyExist = errors.New("Question already exist")
	ErrCannotVote           = errors.New("Cannot change votes")
)

type UserStorage interface {
	FindAllUser() ([]model.User, error)
	CreateUser(model.User) (model.User, error)
	UpdateUser(model.User) (model.User, error)
	DeleteUser(int) error

	FindUser(int) (model.User, error)

	FindUserByNick(string) (model.User, error)
	FindUserByEmail(string) (model.User, error)
	Login(string, string) error
}

type QuestionStorage interface {
	FindAllQuestion() ([]model.Question, error)
	CreateQuestion(model.Question) (model.Question, error)
	UpdateQuestion(model.Question) (model.Question, error)

	FindQuestion(int) (model.Question, error)

	FindQuestionByTitle(string) (model.Question, error)
	FindQuestionByAuthor(int) ([]model.Question, error)
	UpQuestion(int) error
	DownQuestion(int) error
}

type CommentStorage interface {
	CreateComment(model.Comment) (model.Comment, error)
	UpdateComment(model.Comment) (model.Comment, error)

	FindComment(int) (model.Comment, error)

	FindCommentByAuthor(int) ([]model.Comment, error)
	FindCommentByQuestion(int) ([]model.Comment, error)
	UpComment(int) error
	DownComment(int) error
}
