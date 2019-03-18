package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/sashamerkulev/logger"
	"github.com/sashamerkulev/rssservice/domain"
	"github.com/sashamerkulev/rssservice/model"
	"io/ioutil"
	"net/http"
)

func finishUserResponse(w http.ResponseWriter, user model.User, err error, logger logger.UserDbLogger) {
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
	var ur = domain.UpdateUser{Phone: phone, Name: name, Logger: userDbLogger, UserId: userDbLogger.UserId}
	user, err := ur.UpdateUser()
	finishUserResponse(w, user, err, userDbLogger)
}

func usersUploadPhotoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := getAuthorizationToken(r)
	userDbLogger, err := logger.UserDbLogger{DB: dbLogger.DB, UserIP: r.RemoteAddr, UserId: userId}, r.ParseMultipartForm(10<<20)
	if err != nil {
		userDbLogger.Log("ERROR", "PREPARE", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("userPhoto")
	if err != nil {
		userDbLogger.Log("ERROR", "USERSUPLOADPHOTOHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()
	userDbLogger.Log("DEBUG", "USERSUPLOADPHOTOHANDLER", "Uploaded File: "+handler.Filename)
	userDbLogger.Log("DEBUG", "USERSUPLOADPHOTOHANDLER", "File Size: "+fmt.Sprint(handler.Size))
	//userDbLogger.Log("DEBUG", "USERSUPLOADPHOTOHANDLER", "MIME Header: "+handler.Header)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		userDbLogger.Log("ERROR", "USERSUPLOADPHOTOHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var upu = domain.UploadPhotoUser{Logger: userDbLogger, UserId: userDbLogger.UserId, Photo: fileBytes}
	err = upu.UploadPhoto()
	if err != nil {
		userDbLogger.Log("ERROR", "USERSUPLOADPHOTOHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
