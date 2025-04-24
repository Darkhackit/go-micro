package main

import (
	"log"
	"net/http"
)

func (app *Config) sendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		err = app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			return
		}
		return
	}
	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Println(err)
		err = app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			return
		}
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: "Sent to: " + requestPayload.To,
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		return
	}
}
