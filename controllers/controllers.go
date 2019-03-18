package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sashamerkulev/logger"
	"github.com/sashamerkulev/rssservice/db"
	"net/http"
	"strings"
)

var dbLogger logger.DbLogger

func getAuthorizationToken(r *http.Request) int64 {
	token := r.Header.Get("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1)
	userId, _ := db.GetUserIdByToken(token)
	return userId
}

func prepareRequest(w http.ResponseWriter, r *http.Request) (userDbLogger logger.UserDbLogger, err error) {
	w.Header().Set("Content-Type", "application/json")
	userId := getAuthorizationToken(r)
	log, err := logger.UserDbLogger{DB: dbLogger.DB, UserIP: r.RemoteAddr, UserId: userId}, r.ParseForm()
	if err != nil {
		log.Log("ERROR", "PREPARE", err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}
	return log, err
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println("remoteAddr", r.RemoteAddr)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "[]")
}

func Init(_dbLogger logger.DbLogger) {
	dbLogger = _dbLogger

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/articles", articlesHandler).Methods("GET")
	r.HandleFunc("/articles/{articleId}/like", articlesLikeHandler).Methods("POST")
	r.HandleFunc("/articles/{articleId}/dislike", articlesDislikeHandler).Methods("POST")
	r.HandleFunc("/articles/{articleId}/comments", articlesCommentsHandler).Methods("POST")
	r.HandleFunc("/users", homeHandler).Methods("GET")
	r.HandleFunc("/users/register", usersRegisterHandler).Methods("POST")
	r.HandleFunc("/users/update", usersUpdateHandler).Methods("POST")
	r.HandleFunc("/users/uploadPhoto", usersUploadPhotoHandler).Methods("POST")
	err := http.ListenAndServe(":9000", r)
	if err != nil {
		dbLogger.Log("ERROR", "INIT", err.Error())
	}
}
