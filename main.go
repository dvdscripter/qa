package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"securecodewarrior.com/ddias/heapoverflow/model/storage/mongodb"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"securecodewarrior.com/ddias/heapoverflow/model/storage"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type app struct {
	storage.Storage
	jwtKeyFile string
	routes     []route
	router     *mux.Router
}

var webapp app

func main() {

	// cert := flag.String("cert", "server.crt", "public certificate")
	// key := flag.String("key", "server.key", "private certificate")
	jwtKey := flag.String("jwt", "jwt.key", "file with jwt key")
	// openssl rand -out jwt.key -hex 256

	flag.Parse()

	limits := &limit{rate.NewLimiter(10, 10)}

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

	// db, err := sql.New("database.db")
	db, err := mongodb.New("localhost", "go-qa-forum", "users", "questions",
		"comments")
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	defer db.Close()
	// db.DropTableIfExists("users", "questions", "comments")
	// if err := db.AutoMigrate(&model.User{}, &model.Question{},
	// &model.Comment{}).Error; err != nil {
	// 	log.Fatalf("Error migrating db %s\n", err)
	// }

	webapp = app{db, *jwtKey, routes, mux.NewRouter()}
	webapp.registerRoutes(middleJSONLogger)
	webapp.router.Use(
		limits.toLimit,
		webapp.Validate,
	)

	srv := http.Server{
		Addr: "localhost:8000",
		Handler: handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedHeaders([]string{"X-Requested-With",
				"Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT",
				"OPTIONS", "DELETE"}),
		)(webapp.router),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	// openssl genrsa -out server.key 2048
	// openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
	log.Fatal(srv.ListenAndServe())
	// log.Fatal(srv.ListenAndServe())
}
