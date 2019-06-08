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
	params := r.URL.Query()["lastArticleReadDate"]
	// loc, _ := time.LoadLocation("UTC")
	var datetime time.Time
	if len(params) > 0 {
		datetime = domain.StringToDate(params[0])
		logger.Log("DEBUG", "ARTICLESHANDLER", datetime.String())
	} else {
		datetime = domain.StringToDate("")
	}

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
	logger.Log("DEBUG", "ARTICLES", r.RequestURI+" ("+datetime.String()+")"+" count: "+strconv.FormatInt(int64(len(articles)), 10))
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
	article, err := articleUser.Like()
	if err != nil {
		logger.Log("ERROR", "ARTICLESLIKEHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(article)
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
	article, err := articleUser.Dislike()
	if err != nil {
		logger.Log("ERROR", "ARTICLESDISLIKEHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(article)
}

func articleHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	articleUser, err := prepareArticleActivity(logger, r)
	if err != nil {
		logger.Log("ERROR", "articleLikeHandler", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	article, err := articleUser.GetUserArticle()
	if err != nil {
		logger.Log("ERROR", "articleLikeHandler", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(article)
}

func prepareArticleActivity(logger logger.Logger, r *http.Request) (domain.UserArticle, error) {
	vars := mux.Vars(r)
	articleId := vars["articleId"]
	id, err := strconv.ParseInt(articleId, 10, 64)
	return domain.UserArticle{ArticleId: id, UserId: getAuthorizationToken(r), Logger: logger, Repository: mysql.UserArticleRepositoryImpl{}}, err
}
