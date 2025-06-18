package httpserver

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"

	"go.opentelemetry.io/otel/codes"
)

type RegisterUserPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Name     string `json:"name" validate:"required,max=100"`
	Password string `json:"password" validate:"required,min=5,max=17"`
}

func (app *application) registerHandler(w http.ResponseWriter, r *http.Request) {
	registerTraceCtx, span := app.trace.Start(r.Context(), "register")

	defer span.End()

	// get the parameters from
	var payload RegisterUserPayload
	if err := readJson(w, r, &payload); err != nil {
		app.logger.WithContext(registerTraceCtx).Error("Error reading user registration payload as json", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.logger.WithContext(registerTraceCtx).Error("Error validating registration payload", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	ctx, cancel := context.WithTimeout(registerTraceCtx, time.Second*5)
	defer cancel()
	span.SetAttributes(
		attribute.String("email", payload.Email),
		attribute.String("name", payload.Name),
	)
	user, err := app.service.user.Register(ctx, payload.Email, payload.Name, payload.Password)
	if err != nil {
		app.logger.WithContext(registerTraceCtx).Error("Error registering user", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		app.badRequestResponse(w, r, err)
		return
	}
	app.jsonResponse(w, http.StatusCreated, "New account created successfully, please check email for a verification link!", user)
	app.logger.WithContext(registerTraceCtx).Info("New account Registeration was a success ...", user.Email)

}
