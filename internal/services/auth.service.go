package services

import (
	"workout_tracker/internal/jwt"
	"workout_tracker/internal/models"
	"workout_tracker/internal/repository"
	"workout_tracker/internal/utils"
	"workout_tracker/internal/utils/password"

	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"password"`
}

// SignupDTO defined the /login payload
type SignupDTO struct {
	LoginDTO
	Name string `json:"name" validate:"required,min=3"`
}

// UserResponse todo
type UserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

// AccessResponse todo
type AccessResponse struct {
	Token string `json:"token"`
}

// AuthResponse todo
type AuthResponse struct {
	User *UserResponse   `json:"user"`
	Auth *AccessResponse `json:"auth"`
}

func Login(ctx *fiber.Ctx) error {
	b := new(LoginDTO)

	if err := utils.ParseBodyAndValidate(ctx, b); err != nil {
		return err
	}

	u := &UserResponse{}

	err := repository.FindUserByEmail(u, b.Email).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	if err := password.Verify(u.Password, b.Password); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	t := jwt.Generate(&jwt.TokenPayload{
		ID: u.ID,
	})

	return ctx.JSON(&AuthResponse{
		User: u,
		Auth: &AccessResponse{
			Token: t,
		},
	})
}
func Signup(ctx *fiber.Ctx) error {
	b := new(SignupDTO)

	if err := utils.ParseBodyAndValidate(ctx, b); err != nil {
		return err
	}

	err := repository.FindUserByEmail(&struct{ ID string }{}, b.Email).Error

	// If email already exists, return
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusConflict, "Email already exists")
	}

	user := &models.User{
		Name:     b.Name,
		Password: password.Generate(b.Password),
		Email:    b.Email,
	}

	// Create a user, if error return
	if err := repository.CreateUser(user); err.Error != nil {
		return fiber.NewError(fiber.StatusConflict, err.Error.Error())
	}

	// generate access token
	t := jwt.Generate(&jwt.TokenPayload{
		ID: user.ID,
	})

	return ctx.JSON(&AuthResponse{
		User: &UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		Auth: &AccessResponse{
			Token: t,
		},
	})
}
