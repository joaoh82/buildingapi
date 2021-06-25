package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joaoh82/buildingapi/api/handlers"
	"github.com/joaoh82/buildingapi/api/middlewares"
	"github.com/joaoh82/buildingapi/api/presenters"
	"github.com/joaoh82/buildingapi/entities"
	"github.com/joaoh82/buildingapi/repositories"
	"github.com/joaoh82/buildingapi/services"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const PORT = 3000

func main() {
	loggerOutput := zerolog.ConsoleWriter{Out: os.Stderr}
	logger := zerolog.New(loggerOutput)

	exporter, err := jaeger.NewRawExporter(
		jaeger.WithAgentEndpoint(jaeger.WithAgentHost("localhost"), jaeger.WithAgentPort("6831")),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed connecting to apm exporter")
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.ServiceNameKey.String("TodoAPI"),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// DB Connection
	const dsn = "postgresql://postgres:postgres@localhost:5432/postgres"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("An error occurred connection to db: %v\n", err)
		os.Exit(1)
	}

	db.AutoMigrate(entities.Todo{})

	// Presenters
	presenters := presenters.NewPresenters(logger)

	// Repositories
	todoRepoditory := repositories.NewTodoRepository(db)

	// Services
	todoService := services.NewTodoService(todoRepoditory)

	r := chi.NewRouter()

	r.Route("/v1", func(v1 chi.Router) {
		v1.Use(middleware.RealIP)
		v1.Use(middlewares.RequestId)
		v1.Use(middlewares.Tracer)
		v1.Use(middlewares.Logger(logger))

		v1.Use(cors.Default().Handler)

		// This middleware will catch and treat panics
		v1.Use(middlewares.Recover(logger))

		v1.Mount("/todos", handlers.TodoRouter(logger, presenters, todoService))
	})

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", PORT),
		Handler: r,

		// These values are here to make sure that the server does not hang
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		// This value is extremely important, it prevents us from suffering a Slowloris attack
		IdleTimeout: time.Second * 60,
	}

	// Creating a channel that listens on incoming interrupt signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+C
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)

	// Graceful shutdown
	go func() {
		// Wait for a new signal on channel
		<-signalChan
		// Signal received, shutdown the server
		fmt.Println("shutting down...")

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		srv.Shutdown(ctx)

		// Check if context timeouts, in worst case call cancel via defer
		select {
		case <-time.After(21 * time.Second):
			fmt.Println("Not all connection are done")
		case <-ctx.Done():
		}
	}()

	log.Printf("Server started at %d", PORT)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal().Err(err).Msg("server crashed")
	}
}
