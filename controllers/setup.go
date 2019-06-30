package controllers

import (
	"encoding/json"
	"github.com/sashamerkulev/rssservice/domain"
	"github.com/sashamerkulev/rssservice/mysql"
	"github.com/sashamerkulev/rssservice/mysql/data"
	"net/http"
)

func setupRegisterHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	setupId := r.Form.Get("setupId")
	firebaseId := r.Form.Get("firebaseId")
	var ru = domain.SetupRegister{SetupId: setupId, FirebaseId: firebaseId, Logger: logger, Repository: mysql.UserRegisterRepositoryImpl{}}
	user, err := ru.RegisterUser()
	finishUserResponse(w, user, err, logger)
}

func setupSourcesHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	var ur = domain.SetupSources{Logger: logger, Repository: data.MainRepositoryImpl{}}
	sources, err := ur.Repository.GetSources()
	if err != nil {
		logger.Log("ERROR", "SETUPSOURCESHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sources)
}

func setupFirebaseHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	firebaseId := r.Form.Get("firebaseId")
	var ur = domain.SetupRegister{FirebaseId: firebaseId, Logger: logger, Repository: mysql.UserRegisterRepositoryImpl{}}
	err = ur.UpdateFirebaseId(getAuthorizationToken(r))
	if err != nil {
		logger.Log("ERROR", "SETUPFIREBASEHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
