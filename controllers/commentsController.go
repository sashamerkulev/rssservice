package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sashamerkulev/rssservice/domain"
	"github.com/sashamerkulev/rssservice/errors"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/mysql"
	"net/http"
	"strconv"
	"time"
)

func articlesCommentsHandler(w http.ResponseWriter, r *http.Request) {
	userId := getAuthorizationToken(r)
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

	articleUser, err := prepareArticleCommentGet(logger1, r, userId)
	if err != nil {
		logger1.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	comments, err := articleUser.GetComments()
	if err != nil {
		logger1.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}

func articlesDeleteCommentsHandler(w http.ResponseWriter, r *http.Request) {
	userId := getAuthorizationToken(r)
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

	articleUser, err := prepareArticleCommentActivity(logger1, r, userId)
	if err != nil {
		logger1.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = articleUser.DeleteComment()
	if err != nil {
		logger1.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func articlesAddCommentsHandler(w http.ResponseWriter, r *http.Request) {
	userId := getAuthorizationToken(r)
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

	articleUser, err := prepareArticleCommentAdd(logger1, r, userId)
	if err != nil {
		logger1.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	comments := r.Form.Get("comments")
	comment, err := articleUser.AddComment(comments)
	if err != nil {
		logger1.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)
}

func articlesCommentsLikeHandler(w http.ResponseWriter, r *http.Request) {
	userId := getAuthorizationToken(r)
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

	articleUser, err := prepareArticleCommentActivity(logger1, r, userId)
	if err != nil {
		logger1.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	comment, err := articleUser.LikeComment()
	if err != nil {
		logger1.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)
}

func articlesCommentsDislikeHandler(w http.ResponseWriter, r *http.Request) {
	userId := getAuthorizationToken(r)
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

	articleUser, err := prepareArticleCommentActivity(logger1, r, userId)
	if err != nil {
		logger1.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	comment, err := articleUser.DislikeComment()
	if err != nil {
		logger1.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)
}

func prepareArticleCommentActivity(logger logger.Logger, r *http.Request, userId int64) (domain.ArticleComment, error) {
	vars := mux.Vars(r)
	commentId := vars["commentId"]
	if commentId == "" {
		logger.Log("ERROR", "ARTICLESCOMMENTSHANDLER", errors.EmptyParametersError.Error())
		return domain.ArticleComment{}, errors.EmptyParametersError
	}
	cId, err := strconv.ParseInt(commentId, 10, 64)
	if err != nil {
		logger.Log("ERROR", "ARTICLESCOMMENTSHANDLER", "can't parse "+commentId)
		return domain.ArticleComment{}, errors.WrongCommentIdParameterError
	}
	return domain.ArticleComment{UserId: userId, CommentId: cId, Logger: logger,
		Repository: mysql.ArticleCommentRepositoryImpl{
			DB:        mysql.DB,
			TableName: "article",
		}}, nil
}

func prepareArticleCommentGet(logger logger.Logger, r *http.Request, userId int64) (domain.ArticleCommentGet, error) {
	vars := mux.Vars(r)
	articleId := vars["articleId"]
	params := r.URL.Query()["lastCommentsReadDate"]
	if articleId == "" {
		logger.Log("ERROR", "ARTICLESCOMMENTSHANDLER", errors.EmptyParametersError.Error())
		return domain.ArticleCommentGet{}, errors.EmptyParametersError
	}
	aId, err := strconv.ParseInt(articleId, 10, 64)
	if err != nil {
		logger.Log("ERROR", "ARTICLESCOMMENTSHANDLER", "can't parse "+articleId)
		return domain.ArticleCommentGet{}, errors.WrongArticleIdParameterError
	}
	var datetime time.Time
	if len(params) > 0 {
		datetime = domain.StringToDate(params[0])
		logger.Log("DEBUG", "ARTICLESCOMMENTSHANDLER", datetime.String())
	} else {
		datetime = domain.StringToDate("")
	}
	return domain.ArticleCommentGet{UserId: userId, ArticleId: aId, Logger: logger, LastArticleReadDate: datetime,
		Repository: mysql.ArticleCommentRepositoryImpl{
			DB:        mysql.DB,
			TableName: "article",
		}}, nil
}

func prepareArticleCommentAdd(logger logger.Logger, r *http.Request, userId int64) (domain.ArticleCommentAdd, error) {
	vars := mux.Vars(r)
	articleId := vars["articleId"]
	if articleId == "" {
		logger.Log("ERROR", "ARTICLESCOMMENTSHANDLER", errors.EmptyParametersError.Error())
		return domain.ArticleCommentAdd{}, errors.EmptyParametersError
	}
	aId, err := strconv.ParseInt(articleId, 10, 64)
	if err != nil {
		logger.Log("ERROR", "ARTICLESCOMMENTSHANDLER", "can't parse "+articleId)
		return domain.ArticleCommentAdd{}, errors.WrongArticleIdParameterError
	}
	return domain.ArticleCommentAdd{UserId: userId, ArticleId: aId, Logger: logger,
		Repository: mysql.ArticleCommentRepositoryImpl{
			DB:        mysql.DB,
			TableName: "article",
		}}, nil
}
