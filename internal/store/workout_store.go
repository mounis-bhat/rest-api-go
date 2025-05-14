package store

import (
	"database/sql"
	"fmt"
	"time"
)

type Workout struct {
	ID              int            `json:"id"`
	UserID          int64          `json:"user_id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"`
	CaloriesBurned  int            `json:"calories_burned"` // in kcal
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	Entries         []WorkoutEntry `json:"entries"`
}

type WorkoutEntry struct {
	ID              int      `json:"id"`
	ExerciseName    string   `json:"exercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"`
	Weight          *float64 `json:"weight"` // in kg
	Notes           string   `json:"notes"`
	OrderIndex      int      `json:"order_index"`
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

type WorkoutStore interface {
	CreateWorkout(workout *Workout) (*Workout, error)
	GetWorkoutById(id int64) (*Workout, error)
	UpdateWorkout(workout *Workout) error
	DeleteWorkout(id int64) error
	GetAllWorkouts() ([]*Workout, error)
	GetWorkoutOwner(id int64) (int, error)
}

func (s *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	if workout.Title == "" {
		return nil, fmt.Errorf("workout title is required")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO workouts (user_id, title, description, duration_minutes, calories_burned)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`

	err = tx.QueryRow(query, workout.UserID, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID, &workout.CreatedAt, &workout.UpdatedAt)
	if err != nil {
		return nil, err
	}

	insertedEntries := make([]WorkoutEntry, 0, len(workout.Entries))
	for _, entry := range workout.Entries {
		query = `INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
		err = tx.QueryRow(query, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)
		if err != nil {
			return nil, err
		}
		insertedEntries = append(insertedEntries, entry)
	}
	workout.Entries = insertedEntries

	return workout, tx.Commit()
}

func (s *PostgresWorkoutStore) GetWorkoutById(id int64) (*Workout, error) {
	query := `SELECT id, user_id, title, description, duration_minutes, calories_burned, created_at, updated_at
		FROM workouts WHERE id = $1`
	workout := &Workout{}
	err := s.db.QueryRow(query, id).Scan(&workout.ID, &workout.UserID, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned, &workout.CreatedAt, &workout.UpdatedAt)
	if err != nil {
		return nil, err
	}

	query = `SELECT id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index
		FROM workout_entries WHERE workout_id = $1`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		entry := WorkoutEntry{}
		err := rows.Scan(&entry.ID, &entry.ExerciseName, &entry.Sets, &entry.Reps, &entry.DurationSeconds, &entry.Weight, &entry.Notes, &entry.OrderIndex)
		if err != nil {
			return nil, err
		}
		workout.Entries = append(workout.Entries, entry)
	}

	return workout, nil
}

func (s *PostgresWorkoutStore) UpdateWorkout(workout *Workout) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE workouts SET title = $1, description = $2, duration_minutes = $3, calories_burned = $4, updated_at = NOW()
		WHERE id = $5`

	result, err := tx.Exec(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned, workout.ID)
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

	for _, entry := range workout.Entries {
		query = `UPDATE workout_entries SET exercise_name = $1, sets = $2, reps = $3, duration_seconds = $4, weight = $5, notes = $6, order_index = $7
			WHERE id = $8`
		_, err := tx.Exec(query, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex, entry.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *PostgresWorkoutStore) DeleteWorkout(id int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `DELETE FROM workout_entries WHERE workout_id = $1`
	_, err = tx.Exec(query, id)
	if err != nil {
		return err
	}

	query = `DELETE FROM workouts WHERE id = $1`
	result, err := tx.Exec(query, id)
	if err != nil {
		return err
	}

	// Check if any row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit()
}

func (s *PostgresWorkoutStore) GetAllWorkouts() ([]*Workout, error) {
	query := `SELECT id, user_id, title, description, duration_minutes, calories_burned, created_at, updated_at
		FROM workouts`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workouts := []*Workout{}
	for rows.Next() {
		workout := &Workout{}
		err := rows.Scan(&workout.ID, &workout.UserID, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned, &workout.CreatedAt, &workout.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Load entries for this workout
		entriesQuery := `SELECT id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index
			FROM workout_entries WHERE workout_id = $1`
		entryRows, err := s.db.Query(entriesQuery, workout.ID)
		if err != nil {
			return nil, err
		}
		defer entryRows.Close()

		workout.Entries = []WorkoutEntry{}
		for entryRows.Next() {
			entry := WorkoutEntry{}
			err := entryRows.Scan(&entry.ID, &entry.ExerciseName, &entry.Sets, &entry.Reps, &entry.DurationSeconds, &entry.Weight, &entry.Notes, &entry.OrderIndex)
			if err != nil {
				return nil, err
			}
			workout.Entries = append(workout.Entries, entry)
		}

		if err = entryRows.Err(); err != nil {
			return nil, err
		}

		workouts = append(workouts, workout)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return workouts, nil
}

func (s *PostgresWorkoutStore) GetWorkoutOwner(id int64) (int, error) {
	query := `SELECT user_id FROM workouts WHERE id = $1`
	var userID int
	err := s.db.QueryRow(query, id).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
