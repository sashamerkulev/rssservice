package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sashamerkulev/rssservice/domain"
	"net/http"
	"strings"
)

var repository domain.MainRepository

func GetAuthorizationToken(r *http.Request) int64 {
	token := r.Header.Get("Authorization")
	token = strings.Replace(strings.ToLower(token), "bearer ", "", -1)
	userId, _ := repository.GetUserIdByToken(token)
	return userId
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	userId := GetAuthorizationToken(r)
	if userId == -1 {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "[]")
}

func Init(_repository domain.MainRepository) *mux.Router {
	repository = _repository
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET", "POST", "PUT", "DELETE")
	r.HandleFunc("/articles", articlesHandler).Methods("GET")
	r.HandleFunc("/articles/{articleId}", articleHandler).Methods("GET")
	r.HandleFunc("/articles/{articleId}/like", articlesLikeHandler).Methods("PUT")
	r.HandleFunc("/articles/{articleId}/dislike", articlesDislikeHandler).Methods("PUT")
	r.HandleFunc("/articles/{articleId}/comments", articlesCommentsHandler).Methods("GET")
	r.HandleFunc("/articles/{articleId}/comments", articlesAddCommentsHandler).Methods("POST")

	r.HandleFunc("/articles/comments", homeHandler).Methods("GET", "POST", "PUT", "DELETE")
	r.HandleFunc("/articles/comments/{commentId}", articlesDeleteCommentsHandler).Methods("DELETE")
	r.HandleFunc("/articles/comments/{commentId}/like", articlesCommentsLikeHandler).Methods("PUT")
	r.HandleFunc("/articles/comments/{commentId}/dislike", articlesCommentsDislikeHandler).Methods("PUT")

	r.HandleFunc("/setup", homeHandler).Methods("GET", "POST", "PUT", "DELETE")
	r.HandleFunc("/setup/register", setupRegisterHandler).Methods("POST")
	r.HandleFunc("/setup/firebase", setupFirebaseHandler).Methods("POST")
	r.HandleFunc("/setup/sources", setupSourcesHandler).Methods("GET")

	r.HandleFunc("/users", homeHandler).Methods("GET", "POST", "PUT", "DELETE")
	r.HandleFunc("/users/info", userInfoHandler).Methods("GET")
	r.HandleFunc("/users/update", usersUpdateHandler).Methods("PUT")
	r.HandleFunc("/users/uploadPhoto", usersUploadPhotoHandler).Methods("PUT")
	r.HandleFunc("/users/downloadPhoto", authorisedUserDownloadPhotoHandler).Methods("GET")
	r.HandleFunc("/users/{userId}/downloadPhoto", userDownloadPhotoHandler).Methods("GET")
	return r
}
