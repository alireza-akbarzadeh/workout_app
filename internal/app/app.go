package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alireza-akbarzadeh/fem_project/internal/api"
	"github.com/alireza-akbarzadeh/fem_project/internal/store"
	"github.com/alireza-akbarzadeh/fem_project/migrations"
)

type Applicaiton struct {
	Logger         *log.Logger
	WorkoutHnadler *api.WorkoutHnadler
	DB             *sql.DB
}

func NewApplication() (*Applicaiton, error) {
	pgDb, err := store.Open()
	if err != nil {
		return nil, err
	}
	err = store.MigrateFS(pgDb, migrations.Fs, ".")
	if err != nil {
		return nil, err
	}
	// out stores will go here
	// our handlers will go here
	workoutHandler := api.NewWorkoutHandler()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := &Applicaiton{
		Logger:         logger,
		WorkoutHnadler: workoutHandler,
		DB:             pgDb,
	}

	return app, nil
}

func (app *Applicaiton) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "status is available\n")
}
