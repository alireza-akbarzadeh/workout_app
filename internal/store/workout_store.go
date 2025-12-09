package store

import (
	"database/sql"
)

type Workout struct {
	ID              int            `json:"id"`
	Title           string         `json:"title"`
	Description     *string        `json:"description,omitempty"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  *int           `json:"calories_burned,omitempty"`
	Entries         []WorkoutEntry `json:"entries,omitempty"`
}

type WorkoutEntry struct {
	ID              int      `json:"id"`
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

type WorkoutStore interface {
	CreateWorkout(*Workout) (*Workout, error)
	GetWorkouts() ([]*Workout, error)
	GetWorkoutByID(id int64) (*Workout, error)
	UpdateWorkout(*Workout) error
	DeleteWorkout(id int64) error
}

func (pg *PostgresWorkoutStore) GetWorkouts() ([]*Workout, error) {
	query := `
		SELECT id, title, description, duration_minutes, calories_burned 
		FROM workouts
	`
	rows, err := pg.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workouts []*Workout
	for rows.Next() {
		workout := &Workout{}
		err = rows.Scan(&workout.ID, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned)
		if err != nil {
			return nil, err
		}
		workouts = append(workouts, workout)
	}

	return workouts, nil

}

func (pg *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert workout
	query := `
        INSERT INTO workouts (title, description, duration_minutes, calories_burned)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	err = tx.QueryRow(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)
	if err != nil {
		return nil, err
	}

	// Insert workout entries
	for i := range workout.Entries {
		query := `
            INSERT INTO workouts_entries 
                (workout_id, exercise_name, sets, reps, duration_second, weight, notes, order_index)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
            RETURNING id
        `
		err = tx.QueryRow(
			query,
			workout.ID,
			workout.Entries[i].ExerciseName,
			workout.Entries[i].Sets,
			workout.Entries[i].Reps,
			workout.Entries[i].DurationSeconds,
			workout.Entries[i].Weight,
			workout.Entries[i].Notes,
			workout.Entries[i].OrderIndex,
		).Scan(&workout.Entries[i].ID)
		if err != nil {
			return nil, err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (pg *PostgresWorkoutStore) GetWorkoutByID(id int64) (*Workout, error) {
	workout := &Workout{}
	query := `
		SELECT id, title, description, duration_minutes, calories_burned 
		FROM workouts WHERE id = $1
	`
	err := pg.db.QueryRow(query, id).Scan(&workout.ID, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	entryQuery := `
		SELECT id, exercise_name, sets, reps, duration_second, weight, notes, order_index 
		FROM workouts_entries WHERE workout_id = $1 ORDER BY order_index
	`
	rows, err := pg.db.Query(entryQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry WorkoutEntry
		err = rows.Scan(&entry.ID, &entry.ExerciseName, &entry.Sets, &entry.Reps, &entry.DurationSeconds, &entry.Weight, &entry.Notes, &entry.OrderIndex)
		if err != nil {
			return nil, err
		}
		workout.Entries = append(workout.Entries, entry)
	}

	return workout, nil
}

func (pg *PostgresWorkoutStore) UpdateWorkout(workout *Workout) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update workout table
	query := `UPDATE workouts 
			  SET title=$1, description=$2, duration_minutes=$3, calories_burned=$4 
			  WHERE id=$5`

	result, err := tx.Exec(query,
		workout.Title,
		workout.Description,
		workout.DurationMinutes,
		workout.CaloriesBurned,
		workout.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	// Reset entries
	_, err = tx.Exec(`DELETE FROM workouts_entries WHERE workout_id=$1`, workout.ID)
	if err != nil {
		return err
	}

	// Insert updated entries
	for _, entry := range workout.Entries {

		query := `INSERT INTO workouts_entries
			(workout_id, exercise_name, sets, reps, duration_second, weight, notes, order_index)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

		_, err = tx.Exec(query,
			workout.ID,
			entry.ExerciseName,
			entry.Sets,
			entry.Reps,
			entry.DurationSeconds,
			entry.Weight,
			entry.Notes,
			entry.OrderIndex,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (pg *PostgresWorkoutStore) DeleteWorkout(id int64) error {
	query := `DELETE FROM workouts WHERE id = $1`

	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
