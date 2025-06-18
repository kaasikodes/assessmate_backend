package httpserver

import "net/http"

const version = "0.0.0"
const serviceIdentifier = "auth-service"

func (app *application) healthzHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "ok",
		"environment": app.config.env,
		"version":     version,
		"service":     "Authentication",
	}
	if err := app.jsonResponse(w, http.StatusOK, "Health status retrieved successfully!", data); err != nil {
		app.internalServerError(w, r, err)
	}

}
