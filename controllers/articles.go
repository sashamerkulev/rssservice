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
	logger, err := prepareRequest(w, r)
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
	logger.Log("DEBUG", "ARTICLES", r.RequestURI)
	au := domain.UserArticles{
		LastTime:   datetime,
		Logger:     logger,
		UserId:     getAuthorizationToken(r),
		Repository: mysql.UserArticlesRepositoryImpl{},
	}
	articles, err := au.GetUserArticles()
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(articles)
}

func favoriteArticlesHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	logger.Log("DEBUG", "FAVORITEARTICLESHANDLER", r.RequestURI)
	au := domain.UserArticles{
		Logger:     logger,
		UserId:     getAuthorizationToken(r),
		Repository: mysql.UserArticlesRepositoryImpl{},
	}
	articles, err := au.GetUserFavoriteArticles()
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(articles)
}

func articlesLikeHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	articleUser, err := prepareArticleActivity(logger, r)
	if err != nil {
		logger.Log("ERROR", "ARTICLESLIKEHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = articleUser.Like()
	if err != nil {
		logger.Log("ERROR", "ARTICLESLIKEHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func articlesDislikeHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	articleUser, err := prepareArticleActivity(logger, r)
	if err != nil {
		logger.Log("ERROR", "ARTICLESDISLIKEHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = articleUser.Dislike()
	if err != nil {
		logger.Log("ERROR", "ARTICLESDISLIKEHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func prepareArticleActivity(logger logger.Logger, r *http.Request) (domain.UserArticle, error) {
	vars := mux.Vars(r)
	articleId := vars["articleId"]
	id, err := strconv.ParseInt(articleId, 10, 64)
	return domain.UserArticle{ArticleId: id, UserId: getAuthorizationToken(r), Logger: logger, Repository: mysql.UserArticleRepositoryImpl{}}, err
}
