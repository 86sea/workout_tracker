package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"workout_tracker/internal/database"
	"workout_tracker/internal/jwt"
	"workout_tracker/internal/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

type UpdateSetPayload struct {
	Reps         string `json:"Reps"`
	Weight       string `json:"Weight"`
	ExerciseName string `json:"ExerciseName"`
	ExerciseType string `json:"ExerciseType"`
	SetID        uint   `json:"SetID"`
}

func BuildRow(date time.Time, reps, weight int, exerciseName, exerciseType string, setID uint) string {
	return fmt.Sprintf(
		`<tr id="row-%d">
            <td>%s</td>
            <td>%d</td>
            <td>%d</td>
            <td>%s</td>
            <td>%s</td>
            <td>
                <div class="dropdown">
                    <button class="btn btn-secondary dropdown-toggle" type="button" data-bs-toggle="dropdown" aria-expanded="false">
                        Options 
                    </button>
                    <ul class="dropdown-menu">
                        <li><a class="dropdown-item" href="#" 
                            hx-get="/user/updateform?setId=%d" 
                            hx-target="#row-%d" 
                            hx-swap="innerHTML"
                            >
                            Update
                            </a></li>
                        <li>
                            <a class="dropdown-item" 
                               href="#" 
                               hx-post="/user/deleteset?setID=%d" 
                               hx-target="#row-%d" 
                               hx-swap="delete" 
                               hx-confirm="Are you sure you want to delete this?">
                            Delete</a>
                        </li>
                    </ul>
                </div>
            </td>
        </tr>`,
		setID,
		date.Format("02/01/2006"), reps, weight, exerciseName, exerciseType,
		setID, setID, setID, setID,
	)
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
	builder.WriteString(BuildRow(date, reps_int, weight_int, payload.ExerciseName, payload.ExerciseType, newSet.ID))
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
		builder.WriteString(BuildRow(set.Date, int(set.Reps), int(set.Weight), set.ExerciseName, set.ExerciseType, set.ID))
	}
	return ctx.SendString(builder.String())
}

func UpdateSet(ctx *fiber.Ctx) error {

	t := ctx.Cookies("jwt")
	claims, err := jwt.Verify(t)
	if err != nil {
		panic(err)
	}
	userID := claims.ID

	date := time.Now()
	payload := new(UpdateSetPayload)
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
	println(payload.SetID)
	if err := database.DB.Instance.First(&newSet, int(payload.SetID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).SendString("Set not found")
		}
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if err := database.DB.Instance.Save(&newSet).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error updating set in the database")
	}

	var builder strings.Builder
	builder.WriteString(BuildRow(date, reps_int, weight_int, payload.ExerciseName, payload.ExerciseType, newSet.ID))
	return ctx.SendString(builder.String())
}

func UpdateForm(ctx *fiber.Ctx) error {

	setID := ctx.Query("setId")

	if setID == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Set ID not provided")
	}
	targetRowId := "row-" + setID
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf(
		`<div><tr id="%s">
    <form id="set-form" hx-post="/user/updateset" hx-trigger="submit" hx-target="#%s" hx-swap="outerHTML">
        <input type="hidden" name="SetID" value="%s">
        <label class="form-label">Reps:</label> <input class="form-control" type="number" name="Reps">
        <label class="form-label">Weight:</label> <input class="form-control" type="number" name="Weight">
        <label class="form-label">Name:</label>   <input class="form-control" type="text" name="ExerciseName">
        <label class="form-label">Type:</label><br>

    <div class="form-check">
            <input class="form-check-input" type="radio" name="ExerciseType" id="dumbbells" value="dumbbells">
            <label class="form-check-label" for="dumbbells">Dumb Bells</label>
        </div>
        
        <div class="form-check">
            <input class="form-check-input" type="radio" name="ExerciseType" id="barbells" value="barbells">
            <label class="form-check-label" for="barbells">Barbells</label>
        </div>
        
        <div class="form-check">
            <input class="form-check-input" type="radio" name="ExerciseType" id="machine" value="machine">
            <label class="form-check-label" for="machine">Machine</label>
        </div>
        
        <div class="form-check">
            <input class="form-check-input" type="radio" name="ExerciseType" id="bodyweight" value="bodyweight">
            <label class="form-check-label" for="bodyweight">Bodyweight</label>
        </div>
     <label class="form-label">Submit:</label>   <input type="submit" class="btn btn-primary" value="Submit">

</form>
  </tr><div>`, targetRowId, targetRowId, setID,
	))
	return ctx.SendString(builder.String())
}

func DeleteSet(ctx *fiber.Ctx) error {
	setID := ctx.Query("setID")

    println(setID)
	if setID == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Set ID not provided")
	}

    if err := database.DB.Instance.Delete(&models.Set{}, setID).Error; err != nil {
        return ctx.Status(fiber.StatusInternalServerError).SendString("Error deleting set from databse")
    }

    return ctx.SendStatus(fiber.StatusOK)
}
