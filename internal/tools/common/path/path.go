package path

import (
	"os"
	"path/filepath"
	"runtime"
)

type Path struct {
	Root       string
	Migrations string
}

func NewPath() Path {
	p := Path{
		Root:       rootPath(),
		Migrations: migrationsPath(),
	}

	return p
}

func rootPath() string {
	_, b, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(b), "..", "..", "..", "..")
}

func migrationsPath() string {
	return filepath.Join(rootPath(), "db", "migrations")
}

func (p Path) IsFileExists(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}
