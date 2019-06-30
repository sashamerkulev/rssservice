package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sashamerkulev/rssservice/domain"
	"github.com/sashamerkulev/rssservice/logger"
	"github.com/sashamerkulev/rssservice/model"
	"github.com/sashamerkulev/rssservice/mysql"
	"net/http"
	"strconv"
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
		logger.Log("ERROR", "USERSUPLOADPHOTOHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	form := r.MultipartForm
	files := form.File["File"]
	for _, file := range files {
		if file.Size <= 0 {
			continue
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
		var upu = domain.UserPhoto{Logger: logger, UserId: getAuthorizationToken(r), Repository: mysql.UserPhotoUploadRepositoryImpl{}}
		err = upu.UploadUserPhoto(bytes)
		if err != nil {
			logger.Log("ERROR", "USERSUPLOADPHOTOHANDLER", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	//w.WriteHeader(http.StatusOK)
	var ur = domain.UserInfo{Logger: logger, UserId: userId, Repository: mysql.UserInfoRepositoryImpl{}}
	user, err := ur.GetUserInfo()
	finishUserResponse(w, user, err, logger)
}

func authorisedUserDownloadPhotoHandler(w http.ResponseWriter, r *http.Request) {
	userId := getAuthorizationToken(r)
	logger, err := repository.GetLogger(userId, r.RemoteAddr), r.ParseForm()
	if err != nil {
		logger.Log("ERROR", "PREPARE", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var ur = domain.UserPhoto{Logger: logger, UserId: getAuthorizationToken(r), Repository: mysql.UserPhotoUploadRepositoryImpl{}}
	bytes, err := ur.GetUserPhoto()
	w.Header().Add("Content-Type", "image/png")
	w.Header().Add("Content-Length", strconv.Itoa(len(bytes)))
	w.Header().Add("filename", fmt.Sprint(ur.UserId)+".png")
	w.WriteHeader(http.StatusOK)
	n, err := w.Write(bytes)
	if err != nil {
		logger.Log("ERROR", "USERSDOWNLOADPHOTOHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if n != len(bytes) {
		logger.Log("ERROR", "USERSDOWNLOADPHOTOHANDLER", "error writing bytes")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func userDownloadPhotoHandler(w http.ResponseWriter, r *http.Request) {
	userId := getAuthorizationToken(r)
	logger, err := repository.GetLogger(userId, r.RemoteAddr), r.ParseForm()
	if err != nil {
		logger.Log("ERROR", "PREPARE", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	var userId2 int64
	userId2Par := vars["userId"]
	if userId2Par != "" {
		userId2, _ = strconv.ParseInt(userId2Par, 10, 64)
	}

	var ur = domain.UserPhoto{Logger: logger, UserId: userId2, Repository: mysql.UserPhotoUploadRepositoryImpl{}}
	bytes, err := ur.GetUserPhoto()
	w.Header().Add("Content-Type", "image/png")
	w.Header().Add("Content-Length", strconv.Itoa(len(bytes)))
	w.Header().Add("filename", fmt.Sprint(ur.UserId)+".png")
	w.WriteHeader(http.StatusOK)
	n, err := w.Write(bytes)
	if err != nil {
		logger.Log("ERROR", "USERSDOWNLOADPHOTOHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if n != len(bytes) {
		logger.Log("ERROR", "USERSDOWNLOADPHOTOHANDLER", "error writing bytes")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
