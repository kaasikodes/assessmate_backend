package httpserver

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type VerifyUserPayload struct {
	Email string `json:"email" validate:"required,email,max=255"`
	Token string `json:"token" validate:"required,min=5,max=200"`
}

func (app *application) verifyHandler(w http.ResponseWriter, r *http.Request) {
	parentTraceCtx, span := app.trace.Start(r.Context(), "verify")

	defer span.End()

	// get the parameters from
	var payload VerifyUserPayload
	if err := readJson(w, r, &payload); err != nil {
		app.logger.WithContext(parentTraceCtx).Error("Error reading user verify payload as json", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.logger.WithContext(parentTraceCtx).Error("Error validating verify payload", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	span.SetAttributes(
		attribute.String("email", payload.Email),
	)
	user, err := app.service.user.VerifyUser(parentTraceCtx, payload.Email, payload.Token)
	if err != nil {
		app.logger.WithContext(parentTraceCtx).Error("unable to verify user", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "User verified successfully!", user)

}
