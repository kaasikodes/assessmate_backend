package httpserver

import (
	"context"
	"time"

	usermanagment "github.com/kaasikodes/assessmate_backend/internal/core/application/services/user-managment"
	jwtport "github.com/kaasikodes/assessmate_backend/internal/ports/outbound/jwt"
)

const (
	ExpiresAtVerificationToken = time.Hour * 24 * 5
	AccessTokenDuration        = time.Duration(time.Hour * 24 * 3)
)

type ContextKeyUser struct{}
type ContextKeyClaims struct{}

type paginatedResponse struct {
	Total  int   `json:"total"`
	Result []any `json:"result"`
}

func createPaginatedResponse(result []any, total int) paginatedResponse {
	return paginatedResponse{
		Total:  total,
		Result: result,
	}

}

func (app *application) isProduction() bool {
	return app.config.env == "production"
}
func getUserFromContext(ctx context.Context) (*usermanagment.User, bool) {
	user, ok := ctx.Value(ContextKeyUser{}).(*usermanagment.User)
	return user, ok
}

func getClaimsFromContext(ctx context.Context) (*jwtport.CustomClaims, bool) {
	claims, ok := ctx.Value(ContextKeyClaims{}).(*jwtport.CustomClaims)
	return claims, ok
}
