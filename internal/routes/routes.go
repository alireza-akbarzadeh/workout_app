package routes

import (
	"github.com/alireza-akbarzadeh/fem_project/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoute(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", app.HealthCheck)

	// workout
	r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutById)
	r.Post("/workouts", app.WorkoutHandler.Insert)
	r.Get("/workouts", app.WorkoutHandler.GetAllWorkouts)
	r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkout)
	r.Delete("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkout)

	return r
}
