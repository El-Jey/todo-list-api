package tasks

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"

	"todo-list/internal/pkg/store/db"
	"todo-list/internal/service/tasks"
)

func Create(db *db.DB) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		task := tasks.NewTask()

		if err := ctx.Bind().JSON(task); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				for _, e := range validationErrors {
					return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"ok":    false,
						"error": e.Error(),
					})
				}
			}

			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "Error parsing request data",
			})
		}

		task, err := tasks.Create(db, task)

		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Error creating task",
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok":      true,
			"message": "Task created successfully",
			"task":    task,
		})
	}
}

func List(db *db.DB) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		tasksList, err := tasks.Read(db, map[string]any{})

		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Error fetching tasks",
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok":    true,
			"tasks": tasksList,
		})
	}
}

func Update(db *db.DB) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		id := ctx.Params("id")

		if id == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "ID is required",
			})
		}

		idInt, err := strconv.Atoi(id)

		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "Invalid ID",
			})
		}

		tasksList, err := tasks.Read(db, map[string]any{"id": idInt})

		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Internal server error",
			})
		}

		if len(tasksList) == 0 {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"ok":    false,
				"error": "Task not found",
			})
		}

		oldTask := tasksList[0]
		modifiedTask := tasks.NewTask()

		if err := ctx.Bind().JSON(modifiedTask); err != nil {
			hasErrors := true

			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				for _, e := range validationErrors {
					if e.Field() == "Title" {
						hasErrors = false
						continue
					}

					return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"ok":    false,
						"error": e.Error(),
					})
				}
			}

			if hasErrors {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"ok":    false,
					"error": "Error parsing request data",
				})
			}
		}

		_, err = oldTask.Update(db, modifiedTask)

		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Error updating task",
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok":      true,
			"message": "Task updated successfully",
		})
	}
}

func Delete(db *db.DB) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		id := ctx.Params("id")

		if id == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "ID is required",
			})
		}

		idInt, err := strconv.Atoi(id)

		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "Invalid ID",
			})
		}

		tasksList, err := tasks.Read(db, map[string]any{"id": idInt})

		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Internal server error",
			})
		}

		if len(tasksList) == 0 {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"ok":    false,
				"error": "Task not found",
			})
		}

		_, err = tasks.Delete(db, map[string]any{"id": idInt})

		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": "Internal server error",
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok":      true,
			"message": "Task deleted successfully",
		})
	}
}
