package main

import "net/http"

func (app *application) handleCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}
	if err := app.writeResponse(w, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
