package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sashamerkulev/logger"
	"github.com/sashamerkulev/rssservice/domain"
	"net/http"
	"strconv"
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
		LastTime: datetime,
		Logger:   userDbLogger,
		UserId:   userDbLogger.UserId,
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
	userDbLogger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	articleUser, err := prepareArticleActivity(userDbLogger, r)
	if err != nil {
		userDbLogger.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	comments := r.Form.Get("comments")
	err = articleUser.Comment(comments)
	if err != nil {
		userDbLogger.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func articlesLikeHandler(w http.ResponseWriter, r *http.Request) {
	userDbLogger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	articleUser, err := prepareArticleActivity(userDbLogger, r)
	if err != nil {
		userDbLogger.Log("ERROR", "ARTICLESLIKEHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = articleUser.Like()
	if err != nil {
		userDbLogger.Log("ERROR", "ARTICLESLIKEHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func articlesDislikeHandler(w http.ResponseWriter, r *http.Request) {
	userDbLogger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	articleUser, err := prepareArticleActivity(userDbLogger, r)
	if err != nil {
		userDbLogger.Log("ERROR", "ARTICLESDISLIKEHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = articleUser.Dislike()
	if err != nil {
		userDbLogger.Log("ERROR", "ARTICLESDISLIKEHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func prepareArticleActivity(logger logger.UserDbLogger, r *http.Request) (domain.ArticleUserLike, error) {
	vars := mux.Vars(r)
	articleId := vars["articleId"]
	id, err := strconv.ParseInt(articleId, 10, 64)
	return domain.ArticleUserLike{ArticleId: id, UserId: logger.UserId, Logger: logger}, err
}
