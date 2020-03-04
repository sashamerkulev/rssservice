package controllers

import (
	"encoding/json"
	"github.com/sashamerkulev/rssservice/domain"
	"net/http"
)

func setupSourcesHandler(w http.ResponseWriter, r *http.Request) {
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

	ss, err := repository.GetSources()
	var ur = domain.SetupSources{Logger: logger1, Sources: ss}
	sources, err := ur.GetSources()
	if err != nil {
		logger1.Log("ERROR", "SETUPSOURCESHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sources)
}
