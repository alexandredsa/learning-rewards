package server

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/alexandredsa/learning-rewards/reward-processor/graph/generated"
	"github.com/alexandredsa/learning-rewards/reward-processor/graph/resolver"
	"github.com/alexandredsa/learning-rewards/reward-processor/internal/repository"
	"github.com/alexandredsa/learning-rewards/reward-processor/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.uber.org/zap"
)

// Config holds the server configuration
type Config struct {
	Port string
}

// Server represents the HTTP server
type Server struct {
	config Config
	log    *zap.Logger
	server *http.Server
}

// New creates a new server instance
func New(cfg Config) *Server {
	return &Server{
		config: cfg,
		log:    logger.Get(),
	}
}

// Start starts the HTTP server
func (s *Server) Start(db *repository.GormRuleRepository) error {
	// Create resolver
	resolver := resolver.NewResolver(db)

	// Create GraphQL server
	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))

	// Configure server with HTTP transport
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	// Enable introspection
	srv.Use(extension.Introspection{})

	// Set custom error presenter
	srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		s.log.Error("GraphQL error", zap.Error(err))
		return gqlerror.Errorf("%s", err.Error())
	})

	// Create router
	router := mux.NewRouter()

	// loggingMiddleware logs the HTTP request details
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			s.log.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.RequestURI),
				zap.Duration("duration", time.Since(start)))
		})
	}

	// recoveryMiddleware recovers from panics and logs the error
	recoveryMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					s.log.Error("panic recovered",
						zap.Any("error", err),
						zap.String("path", r.RequestURI))
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}

	// Apply middleware
	router.Use(loggingMiddleware)
	router.Use(recoveryMiddleware)

	// GraphQL playground for development
	if os.Getenv("ENV") != "production" {
		router.Handle("/", playground.Handler("GraphQL", "/query"))
	}

	// GraphQL endpoint
	router.Handle("/query", srv)

	// Create HTTP server
	s.server = &http.Server{
		Addr:         ":" + s.config.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s.log.Info("Starting server", zap.String("port", s.config.Port))
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.server != nil {
		s.log.Info("Shutting down server...")
		return s.server.Shutdown(ctx)
	}
	return nil
}
