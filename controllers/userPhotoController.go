package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sashamerkulev/rssservice/domain"
	"github.com/sashamerkulev/rssservice/mysql"
	"net/http"
	"strconv"
)

func usersUploadPhotoHandler(w http.ResponseWriter, r *http.Request) {
	userId := GetAuthorizationToken(r)
	if userId == -1 {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	err := r.ParseMultipartForm(10<<20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger1 := repository.GetLogger(userId, r.RemoteAddr)

	form := r.MultipartForm
	files := form.File["File"]
	for _, file := range files {
		if file.Size <= 0 {
			continue
		}
		logger1.Log("DEBUG", "USERSUPLOADPHOTOHANDLER", "Uploading File: "+file.Filename+" File Size: "+fmt.Sprint(file.Size))
		ffile, err := file.Open()
		if err != nil {
			logger1.Log("ERROR", "USERSUPLOADPHOTOHANDLER", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		bytes := make([]byte, file.Size)
		n, err := ffile.Read(bytes)
		if err != nil {
			logger1.Log("ERROR", "USERSUPLOADPHOTOHANDLER", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if int64(n) < file.Size {
			logger1.Log("ERROR", "USERSUPLOADPHOTOHANDLER", "can't read file")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var upu = domain.UserPhoto{Logger: logger1, UserId: userId, Repository: mysql.UserPhotoUploadRepositoryImpl{
			DB:        mysql.DB,
		}}
		err = upu.UploadUserPhoto(bytes)
		if err != nil {
			logger1.Log("ERROR", "USERSUPLOADPHOTOHANDLER", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	var ur = domain.UserInfo{Logger: logger1, UserId: userId, Repository: mysql.UserInfoRepositoryImpl{
		DB:        mysql.DB,
	}}
	user, err := ur.GetUserInfo()
	finishUserResponse(w, user, err, logger1)
}

func authorisedUserDownloadPhotoHandler(w http.ResponseWriter, r *http.Request) {
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

	var ur = domain.UserPhoto{Logger: logger1, UserId: userId,
		Repository: mysql.UserPhotoUploadRepositoryImpl{DB: mysql.DB}}
	bytes, err := ur.GetUserPhoto()
	w.Header().Add("Content-Type", "image/png")
	w.Header().Add("Content-Length", strconv.Itoa(len(bytes)))
	w.Header().Add("filename", fmt.Sprint(ur.UserId)+".png")
	if err != nil {
		logger1.Log("ERROR", "USERSDOWNLOADPHOTOHANDLER", err.Error())
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	n, err := w.Write(bytes)
	if err != nil {
		logger1.Log("ERROR", "USERSDOWNLOADPHOTOHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if n != len(bytes) {
		logger1.Log("ERROR", "USERSDOWNLOADPHOTOHANDLER", "error writing bytes")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func userDownloadPhotoHandler(w http.ResponseWriter, r *http.Request) {
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

	vars := mux.Vars(r)
	var userId2 int64
	userId2Par := vars["userId"]
	if userId2Par == "" {
		w.WriteHeader(http.StatusBadRequest)
		logger1.Log("ERROR", "USERSDOWNLOADPHOTOHANDLER", "userid is required")
		return
	}
	userId2, err = strconv.ParseInt(userId2Par, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger1.Log("ERROR", "USERSDOWNLOADPHOTOHANDLER", "userid is wrong format")
		return
	}

	var ur = domain.UserPhoto{Logger: logger1, UserId: userId2,
		Repository: mysql.UserPhotoUploadRepositoryImpl{DB: mysql.DB}}
	bytes, err := ur.GetUserPhoto()
	w.Header().Add("Content-Type", "image/png")
	w.Header().Add("Content-Length", strconv.Itoa(len(bytes)))
	w.Header().Add("filename", fmt.Sprint(ur.UserId)+".png")
	if err != nil {
		logger1.Log("ERROR", "USERSDOWNLOADPHOTOHANDLER", err.Error())
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	n, err := w.Write(bytes)
	if err != nil {
		logger1.Log("ERROR", "USERSDOWNLOADPHOTOHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if n != len(bytes) {
		logger1.Log("ERROR", "USERSDOWNLOADPHOTOHANDLER", "error writing bytes")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
