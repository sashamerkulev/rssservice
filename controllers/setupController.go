package controllers

import (
	"encoding/json"
	"github.com/sashamerkulev/rssservice/domain"
	"github.com/sashamerkulev/rssservice/mysql"
	"net/http"
)

func setupRegisterHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	setupId := r.Form.Get("setupId")
	firebaseId := r.Form.Get("firebaseId")
	var ru = domain.SetupRegister{SetupId: setupId, FirebaseId: firebaseId, Logger: logger,
		Repository: mysql.SetupRegisterRepositoryImpl{DB: mysql.DB}}
	user, err := ru.RegisterUser()
	if err != nil {
		logger.Log("ERROR", "FINISH", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func setupFirebaseHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	firebaseId := r.Form.Get("firebaseId")
	setupId := r.Form.Get("setupId")
	var ur = domain.SetupRegister{SetupId: setupId, FirebaseId: firebaseId, Logger: logger,
		Repository: mysql.SetupRegisterRepositoryImpl{DB: mysql.DB}}
	err = ur.UpdateFirebaseId(getAuthorizationToken(r))
	if err != nil {
		logger.Log("ERROR", "SETUPFIREBASEHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
