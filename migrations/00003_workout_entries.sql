-- +goose Up 
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS workouts_entries (
  id BIGSERIAL PRIMARY KEY,
  workout_id BIGINT NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
  exercise_name VARCHAR(255) NOT NULL,
  sets INTEGER NOT NULL,
  reps INTEGER,
  duration_second INTEGER,
  weight DECIMAL(5,2),
  notes TEXT,
  order_index INTEGER NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT valid_workout_entry CHECK (
    (reps IS NOT NULL AND duration_second IS NOT NULL)
    OR
    (reps IS NULL AND duration_second IS NULL)
  )
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workouts_entries;
-- +goose StatementEnd
