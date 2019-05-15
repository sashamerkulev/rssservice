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
	val := r.Form.Get("lastArticleReadDate")
	loc, _ := time.LoadLocation("UTC")
	var datetime time.Time
	if val == "" {
		datetime = time.Date(2019, 1, 1, 0, 0, 0, 0, loc)
	} else {
		//2019-05-03T16:54:07
		datetime, err = time.Parse("2006-01-02T15:04:05", val)
		if err != nil {
			logger.Log("ERROR", "ARTICLES", err.Error())
		}
		datetime = datetime.Add(time.Duration(-3) * time.Hour) // TODO
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
	logger.Log("DEBUG", "ARTICLES", r.RequestURI+" ("+val+")"+" count: "+strconv.FormatInt(int64(len(articles)), 10))
}

func usersActivityArticlesHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	logger.Log("DEBUG", "USERSACTIVITYARTICLESHANDLER", r.RequestURI)
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
