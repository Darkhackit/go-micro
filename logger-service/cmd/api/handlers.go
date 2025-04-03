package main

import (
	"github.com/Darkhackit/go-micro-logger/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var request JSONPayload
	err := app.readJSON(w, r, &request)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
	}
	event := data.LogEntry{
		Name: request.Name,
		Data: request.Data,
	}

	err = app.Models.LogEntry.Insert(event)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
	}
	response := jsonResponse{
		Error:   false,
		Message: "Log entry created",
	}
	err = app.writeJSON(w, http.StatusCreated, response)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
	}
}
