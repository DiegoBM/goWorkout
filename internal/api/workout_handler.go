package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/DiegoBM/goWorkout/internal/middleware"
	"github.com/DiegoBM/goWorkout/internal/store"
	"github.com/DiegoBM/goWorkout/internal/utils"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

func (h *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	workout, err := h.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		h.logger.Printf("ERROR: getWorkoutByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (h *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout

	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		h.logger.Printf("ERROR: decodingCreateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	// Assign the currently logged in user
	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		h.logger.Printf("ERROR: getUser: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "you need to log in to create a workout"})
		return
	}

	workout.UserID = currentUser.ID

	newWorkout, err := h.workoutStore.CreateWorkout(&workout)
	if err != nil {
		h.logger.Printf("ERROR: createWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"workout": newWorkout})
}

func (h *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	workout, err := h.workoutStore.GetWorkoutByID(workoutID)
	if err == sql.ErrNoRows {
		h.logger.Printf("ERROR: getWorkoutByID: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout does not exist"})
		return
	}
	if err != nil {
		h.logger.Printf("ERROR: getWorkoutByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser || currentUser.ID != workout.UserID {
		h.logger.Printf("ERROR: getUser: %v", err)
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "you are not authorized to modify this workout"})
		return
	}

	var updateWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)
	if err != nil {
		h.logger.Printf("ERROR: decodingUpdateRequest: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	if updateWorkoutRequest.Title != nil {
		workout.Title = *updateWorkoutRequest.Title
	}

	if updateWorkoutRequest.Description != nil {
		workout.Description = *updateWorkoutRequest.Description
	}

	if updateWorkoutRequest.DurationMinutes != nil {
		workout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}

	if updateWorkoutRequest.CaloriesBurned != nil {
		workout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}

	if updateWorkoutRequest.Entries != nil {
		workout.Entries = updateWorkoutRequest.Entries
	}

	err = h.workoutStore.UpdateWorkout(workout)
	if err != nil {
		h.logger.Printf("ERROR: updateWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (h *WorkoutHandler) HandleDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		h.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	userID, err := h.workoutStore.GetWorkoutOwner(workoutID)
	if err == sql.ErrNoRows {
		h.logger.Printf("ERROR: GetWorkoutOwner: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout does not exist"})
		return
	}
	if err != nil {
		h.logger.Printf("ERROR: GetWorkoutOwner: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser || currentUser.ID != userID {
		h.logger.Printf("ERROR: getUser: %v", err)
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "you are not authorized to modify this workout"})
		return
	}

	err = h.workoutStore.DeleteWorkout(workoutID)
	if err == sql.ErrNoRows {
		h.logger.Printf("ERROR: deleteWorkoutNoRows: %v", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}

	if err != nil {
		h.logger.Printf("ERROR: deleteWorkout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, utils.Envelope{"success": "workout deleted"})
}
