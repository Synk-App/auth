package controller

import (
	"database/sql"
	"net/http"
	"synk/gateway/app/model"
)

type Users struct {
	model *model.Users
}

type HandleShowResponse struct {
	Resource ResponseHeader    `json:"resource"`
	Data     []model.UsersList `json:"user"`
}

func NewUsers(db *sql.DB) *Users {
	users := Users{
		model: model.NewUsers(db),
	}

	return &users
}

func (u *Users) HandleShow(w http.ResponseWriter, r *http.Request) {
	SetJsonContentType(w)

	var user []model.UsersList

	userId := r.URL.Query().Get("user_id")

	response := HandleShowResponse{
		Resource: ResponseHeader{
			Ok: true,
		},
		Data: user,
	}

	if userId == "" {
		response.Resource.Ok = false
		response.Resource.Error = "fields user_id into query string is required"

		WriteErrorResponse(w, response, "/users", response.Resource.Error, http.StatusBadRequest)

		return
	}

	userList, userErr := u.model.List(userId)

	if userErr != nil {
		response.Resource.Ok = false
		response.Resource.Error = userErr.Error()

		WriteErrorResponse(w, response, "/users", "error on user show", http.StatusInternalServerError)

		return
	}

	if len(userList) > 0 {
		user = userList
	}

	response.Data = user

	WriteSuccessResponse(w, response)
}
