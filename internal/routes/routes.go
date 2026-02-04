package routes

import (
	"github.com/DiegoBM/goWorkout/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	// Group endpoints that require user information (either anonymous or logged-in)
	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)

		// Workout endpoints
		r.Get("/workouts/{id}", app.Middleware.ProtectedEndpoint(app.WorkoutHandler.HandleGetWorkoutByID))
		r.Post("/workouts", app.Middleware.ProtectedEndpoint(app.WorkoutHandler.HandleCreateWorkout))
		r.Put("/workouts/{id}", app.Middleware.ProtectedEndpoint(app.WorkoutHandler.HandleUpdateWorkoutByID))
		r.Delete("/workouts/{id}", app.Middleware.ProtectedEndpoint(app.WorkoutHandler.HandleDeleteWorkout))
	})

	// Healthcheck
	r.Get("/health", app.HealthCheck)

	// User endpoints
	r.Post("/users", app.UserHandler.HandleRegisterUser)

	// Token endpoints
	r.Post("/tokens/authentication", app.TokenHandler.HandleCreateToken)

	return r
}
