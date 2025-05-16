package main

import (
	"backend/internal/api"
	"backend/internal/config"
	"backend/internal/service/middleware"
	"backend/internal/service/schedule"
	"backend/internal/service/user"
	"backend/internal/usecase/auth"
	scheduleusecase "backend/internal/usecase/schedule"
	userusecase "backend/internal/usecase/user"
	"backend/pkg/gorm/dbcontext"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"log"
	"net/http"
	"os"
	"time"
)

var version = "1.0.0"

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

func main() {
	flag.Parse()

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	cfg, err := config.Load(*flagConfig, logger)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}

	gormLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlogger.Config{
			LogLevel: gormlogger.Info,
		},
	)

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		logger.Error().Err(err)
		os.Exit(-1)
	}

	logger.Info().Msg("Connected to PostgreSQL")

	connectionDb, err := db.DB()
	if err != nil {
		logger.Error().Err(err)
	}

	defer func() {
		if err := connectionDb.Close(); err != nil {
			logger.Error().Err(err).Msg("")
		}
	}()

	server := http.Server{
		Addr:    "localhost:5000",
		Handler: BuildHandler(logger, dbcontext.New(db), cfg),
	}

	fmt.Println("Listening on localhost:5000")

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func BuildHandler(logger zerolog.Logger, db dbcontext.Sessions, cfg *config.Config) http.Handler {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding"},
			ExposedHeaders:   []string{"Link"},
			MaxAge:           300,
			AllowCredentials: false,
		}),
		ChiLoggerHandler(logger),
	)

	authMiddleware := middleware.Handler(
		cfg.JWTSigningKey,
		user.NewUserRepository(db),
		logger,
	)

	api.RegisterHandlers(router, version)
	authUseCase := auth.NewUsecase(user.NewService(user.NewUserRepository(db)), authMiddleware.Ja, cfg)
	userUseCase := userusecase.NewUsecase(user.NewService(user.NewUserRepository(db)))
	scheduleUseCase := scheduleusecase.NewUsecase(schedule.NewService(schedule.NewScheduleRepository(db)))

	router.Route("/api/v1", func(router chi.Router) {
		api.RegisterAuthHandlers(router, authUseCase, authMiddleware, logger)
		api.RegisterUserHandlers(router, userUseCase, authMiddleware, logger)
		api.RegisterScheduleHandlers(router, scheduleUseCase, authMiddleware, logger)
	})

	return router
}

func ChiLoggerHandler(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return hlog.NewHandler(logger)(
			hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
				hlog.FromRequest(r).Info().
					Str("method", r.Method).
					Stringer("url", r.URL).
					Int("status_code", status).
					Int("response_size_bytes", size).
					Dur("elapsed_ms", duration).
					Msg("incoming request")
			})(
				hlog.UserAgentHandler("http_user_agent")(
					hlog.RequestIDHandler("request_id", "")(next),
				),
			),
		)
	}
}
