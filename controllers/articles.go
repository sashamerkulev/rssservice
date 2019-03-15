package controllers

import (
	"encoding/json"
	"github.com/sashamerkulev/rssservice/domain"
	"net/http"
	"time"
)

func articlesHandler(w http.ResponseWriter, r *http.Request) {
	userDbLogger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	val := r.Form.Get("datetime")
	var datetime time.Time
	if val == "" {
		loc, _ := time.LoadLocation("UTC")
		datetime = time.Date(2019, 1, 1, 0, 0, 0, 0, loc)
	} else {
		datetime, _ = time.Parse("2006-01-02 15:04:05", val)
	}
	userDbLogger.Log("DEBUG", "ARTICLES", r.RequestURI)

	au := domain.ArticleUser{
		LastTime:  datetime,
		Logger:    userDbLogger,
		UserToken: "token",
	}
	articles, err := au.GetArticleUser()
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(articles)
}

func articlesCommentsHandler(w http.ResponseWriter, r *http.Request) {
	_, err := prepareRequest(w, r)
	if err != nil {
		return
	}
}

func articlesLikeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := prepareRequest(w, r)
	if err != nil {
		return
	}
}

func articlesDislikeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := prepareRequest(w, r)
	if err != nil {
		return
	}
}
