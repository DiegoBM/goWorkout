package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DiegoBM/goWorkout/internal/api"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Our stores will go here

	app := &Application{
		Logger:         logger,
		WorkoutHandler: api.NewWorkoutHandler(),
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
