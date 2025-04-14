package tools

import (
	"todo-list/internal/tools/common/path"
)

type Tools struct {
	Path path.Path
}

func NewTools() Tools {
	t := Tools{
		Path: path.NewPath(),
	}

	return t
}
