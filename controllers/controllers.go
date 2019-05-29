package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sashamerkulev/rssservice/domain"
	"github.com/sashamerkulev/rssservice/logger"
	"net/http"
	"strings"
)

var repository domain.MainRepository

func getAuthorizationToken(r *http.Request) int64 {
	token := r.Header.Get("Authorization")
	token = strings.Replace(strings.ToLower(token), "bearer ", "", -1)
	userId, _ := repository.GetUserIdByToken(token)
	return userId
}

func prepareRequest(w http.ResponseWriter, r *http.Request) (logger logger.Logger, err error) {
	w.Header().Set("Content-Type", "application/json")
	userId := getAuthorizationToken(r)
	log, err := repository.GetLogger(userId, r.RemoteAddr), r.ParseForm()
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

func Init(_repository domain.MainRepository) {
	repository = _repository
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET", "POST", "PUT", "DELETE")
	r.HandleFunc("/articles", articlesHandler).Methods("POST")
	r.HandleFunc("/articles/{articleId}", articleHandler).Methods("GET")
	r.HandleFunc("/articles/{articleId}/like", articlesLikeHandler).Methods("PUT")
	r.HandleFunc("/articles/{articleId}/dislike", articlesDislikeHandler).Methods("PUT")
	r.HandleFunc("/articles/{articleId}/comments", articlesCommentsHandler).Methods("GET")
	r.HandleFunc("/articles/{articleId}/comments", articlesAddCommentsHandler).Methods("POST")
	r.HandleFunc("/articles/comments/{commentId}", articlesDeleteCommentsHandler).Methods("DELETE")
	r.HandleFunc("/articles/comments/{commentId}/like", articlesCommentsLikeHandler).Methods("PUT")
	r.HandleFunc("/articles/comments/{commentId}/dislike", articlesCommentsDislikeHandler).Methods("PUT")
	r.HandleFunc("/setup/register", setupRegisterHandler).Methods("POST")
	r.HandleFunc("/setup/sources", setupSourcesHandler).Methods("GET")
	r.HandleFunc("/users", homeHandler).Methods("GET", "POST", "PUT", "DELETE")
	r.HandleFunc("/users/info", userInfoHandler).Methods("GET")
	r.HandleFunc("/users/update", usersUpdateHandler).Methods("PUT")
	r.HandleFunc("/users/uploadPhoto", usersUploadPhotoHandler).Methods("PUT")
	r.HandleFunc("/users/downloadPhoto", usersDownloadPhotoHandler).Methods("GET")
	err := http.ListenAndServe(":9000", r)
	if err != nil {
		repository.GetLogger(-1, "").Log("ERROR", "INIT", err.Error())
	}
}
