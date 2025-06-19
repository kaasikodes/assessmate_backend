package httpserver

import (
	"net/http"

	"go.opentelemetry.io/otel/codes"
)

type ResendVerificationPayload struct {
	Email string `json:"email" validate:"required,email"`
}

func (app *application) resendVerificationHandler(w http.ResponseWriter, r *http.Request) {
	parentTraceCtx, span := app.trace.Start(r.Context(), "resending verification")

	defer span.End()

	var payload ResendVerificationPayload
	if err := readJson(w, r, &payload); err != nil {
		app.logger.WithContext(parentTraceCtx).Error("Error reading user resending verification payload as json", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.logger.WithContext(parentTraceCtx).Error("Error validating resending verification payload", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	user, err := app.service.user.CreateAndSendVerificationTokenForExistingUser(parentTraceCtx, payload.Email)
	if err != nil {
		app.logger.WithContext(parentTraceCtx).Error("unable to verify user", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, "Verification link has been sent to your mail!", user); err != nil {
		app.internalServerError(w, r, err)
	}

}
