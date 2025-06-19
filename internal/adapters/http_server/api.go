package httpserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/kaasikodes/assessmate_backend/env"
	email_adapter "github.com/kaasikodes/assessmate_backend/internal/adapters/email"
	jwttoken "github.com/kaasikodes/assessmate_backend/internal/adapters/jwt"
	log_adapter "github.com/kaasikodes/assessmate_backend/internal/adapters/logger"
	"github.com/kaasikodes/assessmate_backend/internal/adapters/store"
	usermanagment "github.com/kaasikodes/assessmate_backend/internal/core/application/services/user-managment"
	"github.com/kaasikodes/assessmate_backend/internal/db"
	"github.com/kaasikodes/assessmate_backend/internal/ports/outbound/email"
	jwtport "github.com/kaasikodes/assessmate_backend/internal/ports/outbound/jwt"
	"github.com/kaasikodes/assessmate_backend/internal/ports/outbound/logger"
	user_repo "github.com/kaasikodes/assessmate_backend/internal/ports/outbound/user"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type config struct {
	addr        string
	grpcAddr    string
	db          dbConfig
	env         string
	apiURL      string
	frontendUrl string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}
type application struct {
	config  config
	logger  logger.Logger
	metrics *metrics
	jwt     jwtport.JwtMaker
	trace   trace.Tracer

	// service
	service Service
}

type Service struct {
	user usermanagment.UserManagementService
}

func (app *application) mount(reg *prometheus.Registry) http.Handler {
	app.logger.Info("api mounted ...")
	r := chi.NewRouter()
	// Add the metrics middleware
	r.Use(app.metricsMiddleware)

	r.Get("/healthz", app.healthzHandler)
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}).ServeHTTP(w, r)
		// promhttp.Handler().ServeHTTP(w, r)
	})
	r.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			// TODO: add rate limiting for auth required endpoints to prevent abuse
			r.Post("/register", app.registerHandler) // customer(happy path), vendor
			r.Post("/verify", app.verifyHandler)
			r.Post("/login", app.loginHandler)
			r.Post("/forgot-password", app.forgotPasswordHandler)
			r.Post("/reset-password", app.resetPasswordHandler)
			r.Post("/resend-verification", app.resendVerificationHandler)
			// oauth providers
			// r.Route("/oauth", func(r chi.Router) {
			// 	r.Get("/github/login", app.githubOauthLoginHandler)
			// 	r.Get("/github/callback", app.githubOauthCallbackHandler)
			// })

			r.Group(func(r chi.Router) {
				r.Use(app.authMiddleware)
				r.Get("/me", app.retriveAuthAccountHandler)
			})
		})

	})

	return r

}
func (app *application) run(mux http.Handler) error {

	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Info("App running starting to run on .....", app.config.addr)

	err := server.ListenAndServe()
	app.logger.Info("App running stopping to run on .....", app.config.addr)

	if err != nil {
		return err
	}

	return nil

}
func createUserMgtService(repo user_repo.UserRepository, jwt jwtport.JwtMaker, emailClient email.EmailClient, logger logger.Logger) (*usermanagment.UserManagementService, error) {

	service := usermanagment.NewUserManagementService(repo, jwt, emailClient, logger)
	return service, nil

}
func Start() error {

	cfg := config{
		grpcAddr:    env.GetString("GRPC_ADDR", ":4010"),
		addr:        env.GetString("ADDR", ":3010"),
		apiURL:      env.GetString("API_URL", "localhost:9010"),
		frontendUrl: env.GetString("FRONTEND_URL", "localhost:3000"),
		env:         env.GetString("ENV", "development"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", ""),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
	// TODO: Refactor app to use the options wrapper pattern
	//logging
	logCfg := logger.LogConfig{
		LogFilePath:       "logs/auth-service.log",
		Format:            logger.DefaultLogFormat,
		PrimaryIdentifier: serviceIdentifier,
	}
	logger := log_adapter.New(logCfg)
	// logger := logger.NewZapLogger(logCfg)
	// metrics
	metricsReg := prometheus.NewRegistry()
	metrics := NewMetrics(metricsReg)
	// trace
	tracing := otel.Tracer("assessmate.com/trace")

	//jwt

	jwt := jwttoken.NewJwtMaker(env.GetString("JWT_SECRET", ""))
	//email
	email := email_adapter.NewEmailNotificationService(email_adapter.MailConfig{
		Host:      env.GetString("SMTP_HOST", "sandbox.smtp.mailtrap.io"),
		Port:      env.GetInt("SMTP_PORT", 2525),
		FromEmail: env.GetString("SMTP_FROM_EMAIL", "hello@assessmate.com"),
		Username:  env.GetString("SMTP_USERNAME", ""),
		Password:  env.GetString("SMTP_PASSWORD", ""),
	}, logger)
	//storage
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxOpenConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	persistentStorage := store.NewUserRepository(db, logger)
	// service
	userMgtService, err := createUserMgtService(persistentStorage, jwt, email, logger)
	if err != nil {
		return fmt.Errorf("error creating user management service: %w", err)
	}

	app := &application{
		config:  cfg,
		logger:  logger,
		metrics: metrics,
		trace:   tracing,
		jwt:     jwt,
		service: Service{
			user: *userMgtService,
		},
	}
	mux := app.mount(metricsReg)

	return app.run(mux)

}
