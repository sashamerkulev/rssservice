package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sashamerkulev/rssservice/domain"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/mysql"
	"net/http"
	"strconv"
)

func articlesCommentsHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	articleUser, err := prepareArticleCommentActivity(logger, r)
	if err != nil {
		logger.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	comments, err := articleUser.GetComments()
	if err != nil {
		logger.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}

func articlesDeleteCommentsHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	articleUser, err := prepareArticleCommentActivity(logger, r)
	if err != nil {
		logger.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = articleUser.DeleteComment()
	if err != nil {
		logger.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func articlesAddCommentsHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	articleUser, err := prepareArticleCommentActivity(logger, r)
	if err != nil {
		logger.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	comments := r.Form.Get("comments")
	comment, err := articleUser.AddComment(comments)
	if err != nil {
		logger.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)
}

func articlesCommentsLikeHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	articleUser, err := prepareArticleCommentActivity(logger, r)
	if err != nil {
		logger.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	comment, err := articleUser.LikeComment()
	if err != nil {
		logger.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)
}

func articlesCommentsDislikeHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	articleUser, err := prepareArticleCommentActivity(logger, r)
	if err != nil {
		logger.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	comment, err := articleUser.DislikeComment()
	if err != nil {
		logger.Log("ERROR", "ARTICLESCOMMENTSHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)
}

func prepareArticleCommentActivity(logger logger.Logger, r *http.Request) (domain.ArticleComment, error) {
	vars := mux.Vars(r)
	articleId := vars["articleId"]
	aId, err := strconv.ParseInt(articleId, 10, 64)
	commentId := vars["commentId"]
	cId, err := strconv.ParseInt(commentId, 10, 64)
	return domain.ArticleComment{ArticleId: aId, UserId: getAuthorizationToken(r), CommentId: cId, Logger: logger, Repository: mysql.ArticleCommentRepositoryImpl{}}, err
}
