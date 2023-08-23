package handlers

import (
	"fmt"
	"strconv"
	"strings"
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

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf(
		`<tr>
                <td>%s</td>
                <td>%d</td>
                <td>%d</td>
                <td>%s</td>
                <td>%s</td>
            </tr>`,
		date.Format("02/01/2006"), uint(reps_int), uint(weight_int), payload.ExerciseName, payload.ExerciseType,
	))
	//return ctx.JSON(Response{
	//Success: true,
	//Data:    newSet,
	//})

	return ctx.SendString(builder.String())

}

func LoadSets(ctx *fiber.Ctx) error {

	t := ctx.Cookies("jwt")
	claims, err := jwt.Verify(t)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}
	userID := claims.ID

	var sets []models.Set
	result := database.DB.Instance.Where("user_id = ?", userID).Order("date DESC").Find(&sets)

	if result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error fetching sets from database")
	}

	var builder strings.Builder
	for _, set := range sets {
		builder.WriteString(fmt.Sprintf(
			`<tr>
                <td>%s</td>
                <td>%d</td>
                <td>%d</td>
                <td>%s</td>
                <td>%s</td>
            </tr>`,
			set.Date.Format("02/01/2006"), set.Reps, set.Weight, set.ExerciseName, set.ExerciseType,
		))
	}
	return ctx.SendString(builder.String())
}
