package controllers

import (
	"encoding/json"
	"github.com/sashamerkulev/rssservice/domain"
	"net/http"
)

func setupSourcesHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	var ur = domain.SetupSources{Logger: logger}
	sources, err := ur.GetSources()
	if err != nil {
		logger.Log("ERROR", "SETUPSOURCESHANDLER", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sources)
}
