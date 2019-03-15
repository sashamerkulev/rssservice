package controllers

import (
	"encoding/json"
	"github.com/sashamerkulev/logger"
	"github.com/sashamerkulev/rssservice/domain"
	"github.com/sashamerkulev/rssservice/model"
	"net/http"
)

func finishUserResponse(w http.ResponseWriter, user model.User, err error, logger logger.UserDbLogger) {
	logger.UserId = user.UserId
	if err != nil {
		logger.Log("ERROR", "FINISH", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func usersRegisterHandler(w http.ResponseWriter, r *http.Request) {
	userDbLogger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	deviceId := r.Form.Get("deviceId")
	firebaseId := r.Form.Get("firebaseId")
	var ru = domain.RegisterUser{DeviceId: deviceId, FirebaseId: firebaseId, Logger: userDbLogger}
	user, err := ru.RegisterUser()
	finishUserResponse(w, user, err, userDbLogger)
}

func usersUpdateHandler(w http.ResponseWriter, r *http.Request) {
	userDbLogger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	name := r.Form.Get("name")
	phone := r.Form.Get("phone")
	token := r.Header.Get("Authorization")
	var ur = domain.UpdateUser{Phone: phone, Name: name, Logger: userDbLogger, UserToken: token}
	user, err := ur.UpdateUser()
	finishUserResponse(w, user, err, userDbLogger)
}
