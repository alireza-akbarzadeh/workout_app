package store

import (
	"database/sql"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}
	_, err = db.Exec(`TRUNCATE workouts, workouts_entries CASCADE`)
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
	return db
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name        string
		input       *Workout
		expectedErr bool
	}{
		{
			name: "valid workout",
			input: &Workout{
				Title:           "Morning Routine",
				Description:     "A quick morning workout",
				DurationMinutes: 30,
				CaloriesBurned:  200,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Push-ups",
						Sets:         3,
						Reps:         IntPtr(15),
						OrderIndex:   1,
						Weight:       Float64Ptr(133.5),
						Notes:        strPtr("warp up properly"),
					},
					{
						ExerciseName:    "Jogging",
						Sets:            1,
						DurationSeconds: IntPtr(600),
						OrderIndex:      2,
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "workout with invalid entries",
			input: &Workout{
				Title:           "full body",
				Description:     "compete workout",
				DurationMinutes: 90,
				CaloriesBurned:  500,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "plank",
						Sets:         3,
						Reps:         IntPtr(60),
						Notes:        strPtr("keep form"),
						OrderIndex:   1,
					},
					{
						ExerciseName:    "squats",
						Sets:            4,
						Reps:            IntPtr(12),
						DurationSeconds: IntPtr(60),
						Notes:           strPtr("keep form"),
						Weight:          Float64Ptr(90),
						OrderIndex:      2,
					},
				},
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createWorkout, err := store.CreateWorkout(tt.input)
			if (err != nil) != tt.expectedErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.input.Title, createWorkout.Title)
			assert.Equal(t, tt.input.Description, createWorkout.Description)
			assert.Equal(t, tt.input.DurationMinutes, createWorkout.DurationMinutes)
			assert.Equal(t, tt.input.CaloriesBurned, createWorkout.CaloriesBurned)

			retrieved, err := store.GetWorkoutByID(int64(createWorkout.ID))

			require.NoError(t, err)
			assert.Equal(t, createWorkout.ID, retrieved.ID)
			assert.Equal(t, len(tt.input.Entries), len(retrieved.Entries))
			for i := range retrieved.Entries {
				assert.Equal(t, tt.input.Entries[i].ExerciseName, retrieved.Entries[i].ExerciseName)
				assert.Equal(t, tt.input.Entries[i].Sets, retrieved.Entries[i].Sets)
				assert.Equal(t, tt.input.Entries[i].OrderIndex, retrieved.Entries[i].OrderIndex)
			}
		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func strPtr(s string) *string { return &s }
