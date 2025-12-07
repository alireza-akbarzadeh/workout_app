package api

import (
	"fmt"
	"net/http"

	"strconv"

	"github.com/go-chi/chi/v5"
)

type WorkoutHnadler struct{}

func NewWorkoutHandler() *WorkoutHnadler {
	return &WorkoutHnadler{}
}

func (wh *WorkoutHnadler) Get(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutId := chi.URLParam(r, "id")
	if paramsWorkoutId == "" {
		http.NotFound(w, r)
		return
	}

	workoutID, err := strconv.ParseInt(paramsWorkoutId, 10, 64)

	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "this is th eworkkout id %d\n", workoutID)
}

func (wh *WorkoutHnadler) Insert(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "create a workout")
}
