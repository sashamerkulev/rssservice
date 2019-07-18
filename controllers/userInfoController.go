package controllers

import (
	"encoding/json"
	"github.com/sashamerkulev/rssservice/domain"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"github.com/sashamerkulev/rssservice/mysql"
	"net/http"
)

func finishUserResponse(w http.ResponseWriter, user model.User, err error, logger logger.Logger) {
	if err != nil {
		logger.Log("ERROR", "FINISH", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func userInfoHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	var ur = domain.UserInfo{Logger: logger, UserId: getAuthorizationToken(r),
		Repository: mysql.UserInfoRepositoryImpl{DB: mysql.DB}}
	user, err := ur.GetUserInfo()
	finishUserResponse(w, user, err, logger)
}
