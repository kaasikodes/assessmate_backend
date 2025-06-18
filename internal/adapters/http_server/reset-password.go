package httpserver

import (
	"net/http"

	"go.opentelemetry.io/otel/codes"
)

type ResetPasswordPayload struct {
	Email       string `json:"email" validate:"required,email"`
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"password" validate:"required"`
	UserId      int    `json:"userId" validate:"required"`
}

func (app *application) resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	parentTraceCtx, span := app.trace.Start(r.Context(), "reset passsword")

	defer span.End()

	var payload ResetPasswordPayload
	if err := readJson(w, r, &payload); err != nil {
		app.logger.WithContext(parentTraceCtx).Error("Error reading user reset passsword payload as json", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.logger.WithContext(parentTraceCtx).Error("Error validating reset passsword payload", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	err := app.service.user.ResetPassword(parentTraceCtx, payload.Token, payload.Email, payload.NewPassword, payload.UserId)
	if err != nil {
		app.logger.WithContext(parentTraceCtx).Error("unable to verify user", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, "A password reset link has been sent to you. Please check your mail!", nil); err != nil {
		app.internalServerError(w, r, err)
	}

}
