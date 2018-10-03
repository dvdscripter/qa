package mongodb

import (
	"math/rand"
	"time"

	"github.com/globalsign/mgo"
	"github.com/pkg/errors"
)

type DB struct {
	userC     string
	commentC  string
	questionC string
	database  string
	*mgo.Session
}

func (db *DB) GetUserC() string {
	return db.userC
}

func (db *DB) GetCommentC() string {
	return db.commentC
}

func (db *DB) GetQuestionC() string {
	return db.questionC
}

func (db *DB) GetDatabase() string {
	return db.database
}

func (db *DB) getID(col string) int {
	conn := db.Copy()
	defer conn.Close()

	for {
		nid := int(rand.Int31())
		c, err := conn.DB(db.database).C(col).FindId(nid).Count()
		if err != nil {
			panic(err)
		}
		if c == 0 {
			return nid
		}
	}
}

func New(URL, database, userC, questionC, commentC string) (*DB, error) {
	rand.Seed(time.Now().UnixNano())
	db, err := mgo.Dial(URL)
	if err != nil {
		return nil, errors.Wrap(err, "cannot init mgo session")
	}
	return &DB{userC, commentC, questionC, database, db}, nil
}

func (db *DB) Close() error {
	return db.Close()
}
