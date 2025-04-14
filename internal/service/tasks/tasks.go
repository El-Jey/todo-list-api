package tasks

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"todo-list/internal/pkg/store/db"
)

type Task struct {
	Id          int       `json:"id,omitempty"`
	Title       string    `json:"title" validate:"required" message:"Title is required"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status,omitempty" validate:"oneof=new in_progress done" message:"Status not supported"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func NewTask() *Task {
	return &Task{
		Id:          0,
		Title:       "",
		Description: "",
		Status:      "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func Create(db *db.DB, t *Task) (*Task, error) {
	rows, err := db.Insert(
		"tasks",
		[]string{"title", "description", "status"},
		[]any{t.Title, t.Description, t.Status},
	)

	if err != nil {
		return nil, err
	}

	task, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[Task])
	if err != nil {
		fmt.Printf("Create task error: %v", err)
		return nil, err
	}

	return task, nil
}

func Update(db *db.DB, task *Task) (int64, error) {
	affected, err := db.Update(
		"tasks",
		[]string{"title", "description", "status", "updated_at"},
		[]any{task.Title, task.Description, task.Status, task.UpdatedAt},
		map[string]any{"id": task.Id},
	)

	if err != nil {
		return 0, err
	}

	return affected, nil
}

func Read(db *db.DB, where map[string]any) ([]*Task, error) {
	rows, err := db.Select("tasks", []string{}, where)

	if err != nil {
		return nil, err
	}

	tasksList, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Task])
	if err != nil {
		fmt.Printf("Fetch all tasks error: %v", err)
		return nil, err
	}

	return tasksList, nil
}

func (t *Task) Update(db *db.DB, modifiedTask *Task) (int64, error) {
	if modifiedTask.Title == "" {
		modifiedTask.Title = t.Title
	}

	modifiedTask.CreatedAt = t.CreatedAt

	modifiedTask.Id = t.Id

	return Update(db, modifiedTask)
}

func Delete(db *db.DB, where map[string]any) (int64, error) {
	affected, err := db.Delete("tasks", where)

	if err != nil {
		return 0, err
	}

	return affected, nil
}
