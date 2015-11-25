package controller

import (
	"cubar.com/lib/authentication"
	"cubar.com/lib/util"
	"cubar.com/model"
	"cubar.com/service"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {

	account := &model.Account{}

	if err := util.DecodeJson(r.Body, account); err != nil {
		packError(w, err, ERROR_CODE, "register, user structure decode error")
		return
	}

	if err := account.Register(); err != nil {
		packError(w, err, ERROR_CODE, "register, register error")
		return
	}

	packResult(w, "ok")
}

func Login(w http.ResponseWriter, r *http.Request) {

	login := &model.Login{}
	if err := util.DecodeJson(r.Body, login); err != nil {
		packError(w, err, ERROR_CODE, "")
		return
	}

	responseStatus, token := service.Login(login)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {

	if ok, _ := authentication.RequireTokenAuthentication(w, r); !ok {
		return
	}

	login := &model.Login{}
	if err := util.DecodeJson(r.Body, login); err != nil {
		packError(w, err, ERROR_CODE, "")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(service.RefreshToken(login))
}

func Logout(w http.ResponseWriter, r *http.Request) {

	if ok, _ := authentication.RequireTokenAuthentication(w, r); !ok {
		return
	}

	err := service.Logout(r)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
