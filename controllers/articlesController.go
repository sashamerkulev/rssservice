package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sashamerkulev/rssservice/domain"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/mysql"
	"net/http"
	"strconv"
	"time"
)

func articlesHandler(w http.ResponseWriter, r *http.Request) {
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

	logger1 := repository.GetLogger(userId, r.RemoteAddr)

	params := r.URL.Query()["lastArticleReadDate"]
	var datetime time.Time
	if len(params) > 0 {
		datetime = domain.StringToDate(params[0])
		logger1.Log("DEBUG", "ARTICLESHANDLER", datetime.String())
	} else {
		datetime = domain.StringToDate("")
	}

	userArticles := domain.UserArticles{
		LastTime: datetime,
		Logger:   logger1,
		UserId:   userId,
		Repository: mysql.UserArticlesRepositoryImpl{
			DB:        mysql.DB,
			TableName: "article",
		},
	}
	articles, err := userArticles.GetUserArticles()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger1.Log("ERROR", "ARTICLESHANDLER", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(articles)
}

func articlesLikeHandler(w http.ResponseWriter, r *http.Request) {
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

	logger1 := repository.GetLogger(userId, r.RemoteAddr)

	userArticle, err := prepareUserArticle(logger1, r, userId)
	if err != nil {
		logger1.Log("ERROR", "ARTICLESLIKEHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	article, err := userArticle.Like()
	if err != nil {
		logger1.Log("ERROR", "ARTICLESLIKEHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(article)
}

func articlesDislikeHandler(w http.ResponseWriter, r *http.Request) {
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

	logger1 := repository.GetLogger(userId, r.RemoteAddr)

	userArticle, err := prepareUserArticle(logger1, r, userId)
	if err != nil {
		logger1.Log("ERROR", "ARTICLESDISLIKEHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	article, err := userArticle.Dislike()
	if err != nil {
		logger1.Log("ERROR", "ARTICLESDISLIKEHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(article)
}

func articleHandler(w http.ResponseWriter, r *http.Request) {
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

	logger1 := repository.GetLogger(userId, r.RemoteAddr)

	userArticle, err := prepareUserArticle(logger1, r, userId)
	if err != nil {
		logger1.Log("ERROR", "articleLikeHandler", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	article, err := userArticle.GetUserArticle()
	if err != nil {
		logger1.Log("ERROR", "articleLikeHandler", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(article)
}

func prepareUserArticle(logger logger.Logger, r *http.Request, userId int64) (domain.UserArticle, error) {
	vars := mux.Vars(r)
	articleId := vars["articleId"]
	id, err := strconv.ParseInt(articleId, 10, 64)
	return domain.UserArticle{ArticleId: id, UserId: userId, Logger: logger,
		Repository: mysql.UserArticleRepositoryImpl{
			DB:        mysql.DB,
			TableName: "article",
		}}, err
}
