package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/sashamerkulev/rssservice/domain"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"github.com/sashamerkulev/rssservice/mysql"
	"io/ioutil"
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

func usersUpdateHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	name := r.Form.Get("name")
	phone := r.Form.Get("phone")
	var ur = domain.UserUpdate{Phone: phone, Name: name, Logger: logger, UserId: getAuthorizationToken(r), Repository: mysql.UserUpdateRepositoryImpl{}}
	user, err := ur.UpdateUser()
	finishUserResponse(w, user, err, logger)
}

func userInfoHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	var ur = domain.UserInfo{Logger: logger, UserId: getAuthorizationToken(r), Repository: mysql.UserInfoRepositoryImpl{}}
	user, err := ur.GetUserInfo()
	finishUserResponse(w, user, err, logger)
}

func usersUploadPhotoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := getAuthorizationToken(r)
	logger, err := repository.GetLogger(userId, r.RemoteAddr), r.ParseMultipartForm(10<<20)
	if err != nil {
		logger.Log("ERROR", "PREPARE", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("userPhoto")
	if err != nil {
		logger.Log("ERROR", "USERSUPLOADPHOTOHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()
	logger.Log("DEBUG", "USERSUPLOADPHOTOHANDLER", "Uploaded File: "+handler.Filename)
	logger.Log("DEBUG", "USERSUPLOADPHOTOHANDLER", "File Size: "+fmt.Sprint(handler.Size))
	//logger.Log("DEBUG", "USERSUPLOADPHOTOHANDLER", "MIME Header: "+handler.Header)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Log("ERROR", "USERSUPLOADPHOTOHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var upu = domain.UserPhotoUpload{Logger: logger, UserId: getAuthorizationToken(r), Photo: fileBytes, Repository: mysql.UserPhotoUploadRepositoryImpl{}}
	err = upu.UploadUserPhoto()
	if err != nil {
		logger.Log("ERROR", "USERSUPLOADPHOTOHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
