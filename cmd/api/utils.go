package main

import (
	"encoding/json"
	"net/http"
)

func (app *Application) writeJSON(w http.ResponseWriter, status int, data interface{}, wrap string) error {
	wrapper := struct {
		Data interface{} `json:"data"`
	}{
		Data: data,
	}
	wrapper.Data = data
	if wrap != "" {
		wrapper.Data = map[string]interface{}{
			wrap: wrapper.Data,
		}
	}
	js, err := json.Marshal(wrapper.Data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) writeError(w http.ResponseWriter, status int, err error) error {
	return app.writeJSON(w, status, map[string]interface{}{
		"error": err.Error(),
	}, "")
}
