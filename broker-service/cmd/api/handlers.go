package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LogPayload struct {
	Name string `json:"name"`
	Date string `json:"date"`
}
type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
	}
	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.LogItem(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		err := app.errorJSON(w, fmt.Errorf("invalid action"))
		if err != nil {
			return
		}
	}
}

func (app *Config) LogItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, err := json.Marshal(entry)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
	}
	logServiceUrl := "http://logger-service/log"
	req, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)

	if response.StatusCode != http.StatusCreated {
		err := app.errorJSON(w, fmt.Errorf("invalid status code: %d", response.StatusCode))
		if err != nil {
			return
		}
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Log item successfully created"

	err = app.writeJSON(w, response.StatusCode, payload)
	if err != nil {
		return
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
	}
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		err = app.errorJSON(w, err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		err = app.errorJSON(w, err)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)

	fmt.Println(response.Body, response.StatusCode)

	if response.StatusCode == http.StatusUnauthorized {
		err := app.errorJSON(w, fmt.Errorf("authentication failed"))
		if err != nil {
			return
		}
	} else if response.StatusCode != http.StatusOK {
		err := app.errorJSON(w, fmt.Errorf("authentication failed"))
		if err != nil {
			return
		}
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		err = app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}

	if jsonFromService.Error {
		err := app.errorJSON(w, err, http.StatusUnauthorized)
		if err != nil {
			return
		}
	}
	var payload jsonResponse

	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		return
	}

}

func (app *Config) sendMail(w http.ResponseWriter, mail MailPayload) {
	jsonData, err := json.MarshalIndent(mail, "", "\t")
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
	}
	mailServiceUrl := "http://mail-service/send"
	req, err := http.NewRequest("POST", mailServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		err := app.errorJSON(w, err)
		if err != nil {
			return
		}
		return
	}
	fmt.Println(response.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		err := app.errorJSON(w, fmt.Errorf("invalid status code: %d", response.StatusCode))
		if err != nil {
			return
		}
		return
	}
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Mail successfully sent"

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		return
	}

}
