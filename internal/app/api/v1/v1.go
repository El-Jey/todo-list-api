package v1

import (
	"context"

	"github.com/gofiber/fiber/v3"

	tasksHandler "todo-list/internal/app/api/v1/handler/tasks"
	"todo-list/internal/pkg/store/db"
)

func AddRoutes(group fiber.Router, db *db.DB, ctx context.Context) {
	group.Get("tasks", tasksHandler.List(db))
	group.Post("tasks", tasksHandler.Create(db))
	group.Put("tasks/:id", tasksHandler.Update(db))
	group.Delete("tasks/:id", tasksHandler.Delete(db))
}
