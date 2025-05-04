package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mounis-bhat/rest-api-go/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDb(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "postgres://mounis:3132@localhost:5433/test_db")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// run migrations for our test database
	err = Migrate(db, "../../migrations")
	if err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	_, err = db.Exec("TRUNCATE workouts, workout_entries CASCADE")
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
	return db
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDb(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:           "Morning Run",
				Description:     "A quick morning run",
				DurationMinutes: 30,
				CaloriesBurned:  300,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Running",
						Sets:         1,
						Reps:         utils.IntPtr(2),
						Weight:       utils.Float64Ptr(82.5),
						Notes:        "Felt great",
						OrderIndex:   0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid workout (missing title)",
			workout: &Workout{
				Description:     "A quick morning run",
				DurationMinutes: 30,
				CaloriesBurned:  300,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Running",
						Sets:         1,
						Reps:         utils.IntPtr(2),
						Weight:       utils.Float64Ptr(82.5),
						Notes:        "Felt great",
						OrderIndex:   0,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.workout)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.DurationMinutes, createdWorkout.DurationMinutes)

			retrievedWorkout, err := store.GetWorkoutById(int64(createdWorkout.ID))

			require.NoError(t, err)

			assert.Equal(t, createdWorkout.ID, retrievedWorkout.ID)
			assert.Equal(t, len(createdWorkout.Entries), len(retrievedWorkout.Entries))
			require.Equal(t, len(createdWorkout.Entries), len(retrievedWorkout.Entries), "number of entries should match before comparing entries")

			for i, entry := range createdWorkout.Entries {
				assert.Equal(t, entry.ExerciseName, retrievedWorkout.Entries[i].ExerciseName)
				assert.Equal(t, entry.Sets, retrievedWorkout.Entries[i].Sets)
				assert.Equal(t, entry.Reps, retrievedWorkout.Entries[i].Reps)
			}

		})
	}

}
