package routes

import (
	"github.com/alireza-akbarzadeh/fem_project/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoute(app *app.Applicaiton) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", app.HealthCheck)

	// workout
	r.Get("/workout/{id}", app.WorkoutHnadler.Get)
	r.Post("/workouts", app.WorkoutHnadler.Insert)

	return r
}
