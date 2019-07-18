package controllers

import (
	"github.com/sashamerkulev/rssservice/domain"
	"github.com/sashamerkulev/rssservice/mysql"
	"net/http"
)

func usersUpdateHandler(w http.ResponseWriter, r *http.Request) {
	logger, err := prepareRequest(w, r)
	if err != nil {
		return
	}
	name := r.Form.Get("name")
	phone := r.Form.Get("phone")
	var ur = domain.UserUpdate{Phone: phone, Name: name, Logger: logger, UserId: getAuthorizationToken(r),
		Repository: mysql.UserUpdateRepositoryImpl{DB: mysql.DB}}
	user, err := ur.UpdateUser()
	finishUserResponse(w, user, err, logger)
}
