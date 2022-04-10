package main

import (
	"net/http"
)

func (app *Application) statusHandler(w http.ResponseWriter, r *http.Request) {
	currentStatus := AppState{
		Status:      "OK",
		Environment: app.config.env,
		Version:     version,
	}

	err := app.writeJSON(w, http.StatusOK, currentStatus, "")
	if err != nil {
		app.logger.Println("statusHandler: " + err.Error())
		//app.writeError(w, http.StatusInternalServerError, err)
	}

}
