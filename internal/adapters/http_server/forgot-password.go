package httpserver

import (
	"net/http"

	"go.opentelemetry.io/otel/codes"
)

type ForgotPasswordPayload struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

func (app *application) forgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	parentTraceCtx, span := app.trace.Start(r.Context(), "forgot passsword")

	defer span.End()

	var payload ForgotPasswordPayload
	if err := readJson(w, r, &payload); err != nil {
		app.logger.WithContext(parentTraceCtx).Error("Error reading user forgot passsword payload as json", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.logger.WithContext(parentTraceCtx).Error("Error validating forgot passsword payload", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	data, err := app.service.user.ForgotPassword(parentTraceCtx, payload.Email)
	if err != nil {
		app.logger.WithContext(parentTraceCtx).Error("unable to perform fgt password service action", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, "A password reset link has been sent to you. Please check your mail!", data); err != nil {
		app.internalServerError(w, r, err)
	}

}
