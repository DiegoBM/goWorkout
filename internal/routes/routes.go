package routes

import (
	"github.com/DiegoBM/goWorkout/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	// Healthcheck
	r.Get("/health", app.HealthCheck)

	// Workout endpoints
	r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
	r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkoutByID)
	r.Delete("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkout)

	// User endpoints
	r.Post("/users", app.UserHandler.HandleRegisterUser)

	// Token endpoints
	r.Post("/tokens/authentication", app.TokenHandler.HandleCreateToken)

	return r
}
