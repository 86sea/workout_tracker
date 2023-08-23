package handlers

import (
	"strconv"
	"time"
	"workout_tracker/internal/database"
	"workout_tracker/internal/jwt"
	"workout_tracker/internal/models"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type SetPayload struct {
	Reps         string `json:"Reps"`
	Weight       string `json:"Weight"`
	ExerciseName string `json:"ExerciseName"`
	ExerciseType string `json:"ExerciseType"`
}

func NewSet(ctx *fiber.Ctx) error {

	t := ctx.Cookies("jwt")
	claims, err := jwt.Verify(t)
	if err != nil {
		panic(err)
	}
	userID := claims.ID

	date := time.Now()
	payload := new(SetPayload)
	if err := ctx.BodyParser(payload); err != nil {
		return err
	}
	reps_int, err := strconv.Atoi(payload.Reps)
	if err != nil {
		panic(err)
	}
	weight_int, err := strconv.Atoi(payload.Weight)
	if err != nil {
		panic(err)
	}

	newSet := models.Set{
		Date:         date,
		Reps:         uint(reps_int),
		Weight:       uint(weight_int),
		ExerciseName: payload.ExerciseName,
		ExerciseType: payload.ExerciseType,
		UserID:       userID,
	}

	result := database.DB.Instance.Create(&newSet)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error adding set to database")
	}

	return ctx.JSON(Response{
		Success: true,
		Data:    newSet,
	})

}
