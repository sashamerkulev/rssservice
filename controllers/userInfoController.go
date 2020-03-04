package controllers

import (
	"encoding/json"
	"github.com/sashamerkulev/rssservice/domain"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/models"
	"github.com/sashamerkulev/rssservice/mysql"
	"net/http"
)

func userInfoHandler(w http.ResponseWriter, r *http.Request) {
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

	var ur = domain.UserInfo{Logger: logger1, UserId: userId,
		Repository: mysql.UserInfoRepositoryImpl{DB: mysql.DB}}
	user, err := ur.GetUserInfo()

	finishUserResponse(w, user, err, logger1)
}

func usersUpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	name := r.Form.Get("name")
	phone := r.Form.Get("phone")
	var ur = domain.UserUpdate{Phone: phone, Name: name, Logger: logger1, UserId: userId,
		Repository: mysql.UserUpdateRepositoryImpl{DB: mysql.DB}}
	user, err := ur.UpdateUser()

	finishUserResponse(w, user, err, logger1)
}

func finishUserResponse(w http.ResponseWriter, user models.User, err error, logger logger.Logger) {
	if err != nil {
		logger.Log("ERROR", "FINISH", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
