package controllers

import (
	"encoding/json"
	"github.com/sashamerkulev/rssservice/domain"
	"github.com/sashamerkulev/rssservice/mysql"
	"net/http"
)

func setupRegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger1 := repository.GetLogger(-1, r.RemoteAddr)

	setupId := r.Form.Get("setupId")
	firebaseId := r.Form.Get("firebaseId")
	var ru = domain.SetupRegister{SetupId: setupId, FirebaseId: firebaseId, Logger: logger1,
		Repository: mysql.SetupRegisterRepositoryImpl{DB: mysql.DB}}
	user, err := ru.RegisterUser()
	if err != nil {
		logger1.Log("ERROR", "FINISH", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func setupFirebaseHandler(w http.ResponseWriter, r *http.Request) {
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

	firebaseId := r.Form.Get("firebaseId")
	setupId := r.Form.Get("setupId")
	var ur = domain.SetupRegister{SetupId: setupId, FirebaseId: firebaseId, Logger: logger1,
		Repository: mysql.SetupRegisterRepositoryImpl{DB: mysql.DB}}
	err = ur.UpdateFirebaseId(userId)
	if err != nil {
		logger1.Log("ERROR", "SETUPFIREBASEHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
