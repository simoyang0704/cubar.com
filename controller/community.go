package controller

import (
	"cubar.com/lib/authentication"
	"cubar.com/lib/util"
	"cubar.com/model"
	"github.com/gorilla/mux"
	"net/http"
)

func GetCommunity(w http.ResponseWriter, r *http.Request) {

	if ok, _ := authentication.RequireTokenAuthentication(w, r); !ok {
		return
	}

	vars := mux.Vars(r)
	communityId := util.StringToInt(vars["community_id"])

	community := &model.Community{
		CommunityId: communityId,
	}

	if exists, err := community.Get(); err != nil {
		packError(w, err, ERROR_CODE, "")
		return
	} else if !exists {
		packResult(w, "")
		return
	}

	packResult(w, community)
}

func AddCommunity(w http.ResponseWriter, r *http.Request) {

	var loggedUser *model.User
	var ok bool
	if ok, loggedUser = authentication.RequireTokenAuthentication(w, r); !ok {
		return
	}

	community := &model.Community{}
	if err := util.DecodeJson(r.Body, community); err != nil {
		packError(w, err, ERROR_CODE, "add community, struct error")
		return
	}

	community.CreatedUserId = loggedUser.UserId

	if err := community.Add(); err != nil {
		packError(w, err, ERROR_CODE, "")
		return
	}

	packResult(w, "ok")
}

func UpdateCommunity(w http.ResponseWriter, r *http.Request) {

	var loggedUser *model.User
	var ok bool
	if ok, loggedUser = authentication.RequireTokenAuthentication(w, r); !ok {
		return
	}

	community := &model.Community{}
	if err := util.DecodeJson(r.Body, community); err != nil {
		packError(w, err, ERROR_CODE, "update community, struct error")
		return
	}

	community.UpdatedUserId = loggedUser.UserId

	if err := community.Update(); err != nil {
		packError(w, err, ERROR_CODE, "")
		return
	}

	packResult(w, "ok")
}

func DeleteCommunity(w http.ResponseWriter, r *http.Request) {

	var loggedUser *model.User
	var ok bool
	if ok, loggedUser = authentication.RequireTokenAuthentication(w, r); !ok {
		return
	}

	vars := mux.Vars(r)
	communityId := util.StringToInt(vars["community_id"])

	community := &model.Community{
		CommunityId: communityId,
	}

	community.DeletedUserId = loggedUser.UserId

	if err := community.Delete(); err != nil {
		packError(w, err, ERROR_CODE, "")
		return
	}

	packResult(w, "ok")
}

func QueryCommunities(w http.ResponseWriter, r *http.Request) {

	if ok, _ := authentication.RequireTokenAuthentication(w, r); !ok {
		return
	}

	communityQuery := &model.CommunityQuery{}
	if err := util.DecodeJson(r.Body, communityQuery); err != nil {
		packError(w, err, ERROR_CODE, "query communitied, struct error")
		return
	}

	communities := communityQuery.Query()

	packResult(w, communities)
}

func QueryUserCommunities(w http.ResponseWriter, r *http.Request) {

	var loggedUser *model.User
	var ok bool
	if ok, loggedUser = authentication.RequireTokenAuthentication(w, r); !ok {
		return
	}

	vars := mux.Vars(r)
	status := vars["status"]

	communityQuery := &model.CommunityQuery{}
	if err := util.DecodeJson(r.Body, communityQuery); err != nil {
		packError(w, err, ERROR_CODE, "query communitied, struct error")
		return
	}

	statusCode := 0
	if status == "joined" {
		statusCode = model.COMMUNITY_USER_STATUS_JOINED
	} else if status == "apply" {
		statusCode = model.COMMUNITY_USER_STATUS_APPLY
	}

	if statusCode == 0 {
		packResult(w, nil)
		return
	}

	communities := communityQuery.QueryByUserAndStatus(loggedUser, statusCode)
	packResult(w, communities)
}
