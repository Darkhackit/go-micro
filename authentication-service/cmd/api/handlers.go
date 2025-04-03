package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		err := app.errorJSON(w, err, http.StatusBadRequest)
		if err != nil {
			return
		}
	}
	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		err = app.errorJSON(w, err, http.StatusNotFound)
		if err != nil {
			return
		}
	}
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		err = app.errorJSON(w, errors.New("invalid credentials"), http.StatusNotFound)
		if err != nil {
			return
		}
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user  %s", user.Email),
		Data:    user,
	}

	fmt.Println(payload)

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		return
	}

}
