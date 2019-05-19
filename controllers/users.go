package controllers

import (
	"encoding/json"
	"fmt"
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
	form := r.MultipartForm
	files := form.File["File"]
	for _, file := range files {
		if file.Size < 0 {
			return
		}
		logger.Log("DEBUG", "USERSUPLOADPHOTOHANDLER", "Uploading File: "+file.Filename+" File Size: "+fmt.Sprint(file.Size))
		ffile, err := file.Open()
		if err != nil {
			logger.Log("ERROR", "USERSUPLOADPHOTOHANDLER", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		bytes := make([]byte, file.Size)
		n, err := ffile.Read(bytes)
		if err != nil {
			logger.Log("ERROR", "USERSUPLOADPHOTOHANDLER", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if int64(n) < file.Size {
			logger.Log("ERROR", "USERSUPLOADPHOTOHANDLER", "can't read file")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var upu = domain.UserPhotoUpload{Logger: logger, UserId: getAuthorizationToken(r), Photo: bytes, Repository: mysql.UserPhotoUploadRepositoryImpl{}}
		err = upu.UploadUserPhoto()
		if err != nil {
			logger.Log("ERROR", "USERSUPLOADPHOTOHANDLER", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	//logger.Log("DEBUG", "USERSUPLOADPHOTOHANDLER", "MIME Header: "+handler.Header)
	//if err != nil {
	//	logger.Log("ERROR", "USERSUPLOADPHOTOHANDLER", err.Error())
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	w.WriteHeader(http.StatusOK)
}
