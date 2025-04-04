package main

import (
	"bytes"
	"encoding/json"
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

	logServiceUrl := "http://logger-service/log"
	logData := struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}{}

	logData.Name = "Authentication"
	logData.Data = fmt.Sprintf("%v", payload)

	jsonData, err := json.Marshal(logData)
	if err != nil {
		err = app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			return
		}
	}

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		err = app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			return
		}
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		err = app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			return
		}
	}
	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		return
	}

}
