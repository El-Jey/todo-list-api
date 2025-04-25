package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"

	v1 "todo-list/internal/app/api/v1"
	"todo-list/internal/pkg/store/db"
)

type structValidator struct {
	validate *validator.Validate
}

// Validator needs to implement the Validate method
func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func NewServer(db *db.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
	})

	api := app.Group("/api")

	apiV1 := api.Group("/v1")
	v1.AddRoutes(apiV1, db)

	return app
}
