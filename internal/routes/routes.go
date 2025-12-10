package routes

import (
	_ "github.com/alireza-akbarzadeh/fem_project/docs"
	"github.com/alireza-akbarzadeh/fem_project/internal/app"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoute(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/health", app.HealthCheck)

	//swagger
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// workout
	r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutById)
	r.Post("/workouts", app.WorkoutHandler.Insert)
	r.Get("/workouts", app.WorkoutHandler.GetAllWorkouts)
	r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkout)
	r.Delete("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkout)
	// users
	r.Post("/users", app.UserHandler.HandleRegisterUser)
	r.Get("/users", app.UserHandler.HandleGetUserByUsername)
	r.Put("/users", app.UserHandler.HandleUpdateUser)
	r.Get("/users/{id}", app.UserHandler.HandleGetUserByID)
	r.Delete("/users/{id}", app.UserHandler.HandleDeleteUser)

	// tokens
	r.Post("/tokens", app.TokenHandler.CreateToken)

	return r
}
