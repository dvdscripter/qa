package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"securecodewarrior.com/ddias/heapoverflow/model"
	"securecodewarrior.com/ddias/heapoverflow/model/storage"
	"securecodewarrior.com/ddias/heapoverflow/model/storage/sql"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type app struct {
	storage.Storage
	jwtKeyFile string
}

func main() {

	cert := flag.String("cert", "server.crt", "public certificate")
	key := flag.String("key", "server.key", "private certificate")
	jwtKey := flag.String("jwt", "jwt.key", "file with jwt key")
	// openssl rand -out jwt.key -hex 256

	flag.Parse()

	limits := &limit{
		max:      100,
		interval: 10 * time.Second,
	}
	go limits.start()
	// app := app{
	// 	storage.Storage{
	// 		UserStorage:     db,
	// 		QuestionStorage: db,
	// 		CommentStorage:  db,
	// 	},
	// 	*jwtKey,
	// }

	// db := memory.New()
	// app := app{db, *jwtKey}

	db, err := sql.New("database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// db.DropTableIfExists("users", "questions", "comments")
	if err := db.AutoMigrate(&model.User{}, &model.Question{}, &model.Comment{}).Error; err != nil {
		log.Fatalf("Error migrating db %s\n", err)
	}
	app := app{db, *jwtKey}

	if _, err := app.Storage.CreateUser(model.User{
		Email:    "zim@email.com",
		Nick:     "zim",
		Password: "secret1",
	}); err != nil {
		log.Fatal("Cannot create user zero")
	}

	router := mux.NewRouter()

	rUser := router.PathPrefix("/user").Subrouter()
	rUser.Handle("", middleJSONLogger(app.CreateUser)).Methods("POST")
	rUser.Handle("", app.Validate(middleJSONLogger(app.RetrieveUsers))).Methods("GET")
	rUser.Handle("/{id:[0-9]+}", app.Validate(middleJSONLogger(app.RetrieveUser))).Methods("GET")
	rUser.Handle("/{id:[0-9]+}", app.Validate(middleJSONLogger(app.DeleteUser))).Methods("DELETE")
	rUser.Handle("/{id:[0-9]+}", app.Validate(middleJSONLogger(app.UpdateUser))).Methods("PUT")

	rQuestion := router.PathPrefix("/question").Subrouter()
	rQuestion.Handle("", middleJSONLogger(app.CreateQuestion)).Methods("POST")
	rQuestion.Handle("", middleJSONLogger(app.RetrieveQuestions)).Methods("GET")
	rQuestion.Handle("/{id:[0-9]+}", middleJSONLogger(app.RetrieveQuestion)).Methods("GET")
	rQuestion.Handle("/{id:[0-9]+}", middleJSONLogger(app.UpdateQuestion)).Methods("PUT")
	rQuestion.Handle("/{id:[0-9]+}/vote", middleJSONLogger(app.UpVoteQuestion)).Methods("PUT")
	rQuestion.Handle("/{id:[0-9]+}/vote", middleJSONLogger(app.DownVoteQuestion)).Methods("DELETE")

	rQuestion.Handle("/{id:[0-9]+}/comments", middleJSONLogger(app.RetrieveQuestionComments)).Methods("GET")
	rQuestion.Handle("/{id:[0-9]+}/comments", middleJSONLogger(app.CreateQuestionComments)).Methods("POST")
	rQuestion.Handle("/{id:[0-9]+}/comments/{cid:[0-9]+}", middleJSONLogger(app.RetrieveQuestionComment)).Methods("GET")
	rQuestion.Handle("/{id:[0-9]+}/comments/{cid:[0-9]+}", middleJSONLogger(app.UpdateQuestionComment)).Methods("PUT")
	rQuestion.Handle("/{id:[0-9]+}/comments/{cid:[0-9]+}/vote", middleJSONLogger(app.UpVoteQuestionComment)).Methods("PUT")
	rQuestion.Handle("/{id:[0-9]+}/comments/{cid:[0-9]+}/vote", middleJSONLogger(app.DownVoteQuestionComment)).Methods("DELETE")

	rQuestion.Use(app.Validate)

	router.Handle("/login", middleJSONLogger(app.Login)).Methods("POST")

	router.Use(
		limits.toLimit,
	)

	srv := http.Server{
		Addr:         "localhost:8000",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	// openssl genrsa -out server.key 2048
	// openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
	log.Fatal(srv.ListenAndServeTLS(*cert, *key))

}
