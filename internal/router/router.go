package router

import (
	"net/http"

	"github.com/avito/internship/pr-service/internal/handler"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type TeamController interface {
	CreateTeam(w http.ResponseWriter, r *http.Request) error
	GetTeamByName(w http.ResponseWriter, r *http.Request) error
}

type UserController interface {
	SetIsActive(w http.ResponseWriter, r *http.Request) error
	GetUsersPullRequests(w http.ResponseWriter, r *http.Request) error
}

type PullRequestController interface {
	Create(w http.ResponseWriter, r *http.Request) error
	Merge(w http.ResponseWriter, r *http.Request) error
	Reassign(w http.ResponseWriter, r *http.Request) error
}

func InitRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	return r
}

func SetupTeamRoutes(r *chi.Mux, c TeamController) {
	r.Route("/team", func(r chi.Router) {
		r.Post("/add", handler.ErrorHandler(c.CreateTeam))
		r.Get("/get", handler.ErrorHandler(c.GetTeamByName))
	})
}

func SetupUserRoutes(r *chi.Mux, c UserController) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/setIsActive", handler.ErrorHandler(c.SetIsActive))
		r.Get("/getReview", handler.ErrorHandler(c.GetUsersPullRequests))
	})
}

func SetupPullRequestRoutes(r *chi.Mux, c PullRequestController) {
	r.Route("/pullRequest", func(r chi.Router) {
		r.Post("/create", handler.ErrorHandler(c.Create))
		r.Post("/merge", handler.ErrorHandler(c.Merge))
		r.Post("/reassign", handler.ErrorHandler(c.Reassign))
	})
}
