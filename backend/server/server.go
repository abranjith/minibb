package server

import (
	"context"
	"embed"
	"net/http"
	"os"
	"time"

	"minibb/backend/handlers"
	"minibb/backend/models"
	"minibb/backend/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"golang.org/x/time/rate"
)

type Server struct {
	router   *chi.Mux
	handlers *handlers.Handler
	webFS    embed.FS
}

func New(db models.DB, webFS embed.FS) *Server {
	r := chi.NewRouter()
	h := handlers.New(db)

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS for development
	if os.Getenv("ENV") == "development" {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
	}

	// Rate limiting
	rateLimiter := utils.NewRateLimiter(rate.Every(time.Second), 10)
	r.Use(rateLimiter.Middleware)

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Boards
		r.Get("/boards", h.GetBoards)
		r.Get("/board/{slug}", h.GetBoardBySlug)
		r.Post("/board", h.CreateBoard) // Admin only in production

		// Topics
		r.Get("/topics/board/{id}", h.GetTopicsByBoardID)
		r.Get("/topic/{id}", h.GetTopicByID)
		r.Post("/topic", h.CreateTopic)

		// Posts
		r.Get("/posts/topic/{id}", h.GetPostsByTopicID)
		r.Post("/post", h.CreatePost)

		// Post Replies
		r.Get("/posts/{id}/replies", h.GetRepliesByPostID)
		r.Post("/posts/{id}/replies", h.CreatePostReply)
		r.Delete("/replies/{id}", h.DeletePostReply)
	})

	return &Server{
		router:   r,
		handlers: h,
		webFS:    webFS,
	}
}

func (s *Server) Start(ctx context.Context) error {
	// Add frontend serving routes
	if os.Getenv("ENV") != "development" {
		// Production: serve embedded files
		s.router.Handle("/*", http.FileServer(http.FS(s.webFS)))
	} else {
		// Development: serve a simple message, frontend runs on separate port
		s.router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("<h1>MiniBB Backend</h1><p>Frontend should be running on <a href='http://localhost:5173'>http://localhost:5173</a></p>"))
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: s.router,
	}

	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	return srv.ListenAndServe()
}
