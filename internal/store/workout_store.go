package store

import "database/sql"

type Workout struct {
	Id              int            `json:"id"`
	Title           string         `json:"title"`
	Description     *string        `json:"description,omitempty"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  *int           `json:"calories_burned,omitempty"`
	Entries         []WorkoutEntry `json:"entries,omitempty"`
}

type WorkoutEntry struct {
	Id              int      `json:"id"`
	WorkoutID       int      `json:"workout_id"`
	ExerciseName    string   `json:"exercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps,omitempty"`
	DurationSeconds *int     `json:"duration_seconds,omitempty"`
	Weight          *float64 `json:"weight,omitempty"`
	Notes           *string  `json:"notes,omitempty"`
	OrderIndex      int      `json:"order_index"`
	CreatedAt       string   `json:"created_at,omitempty"`
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}
