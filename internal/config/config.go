package config

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/jinzhu/configor"

	tools "todo-list/internal/tools/common"
)

var (
	config Config
	once   sync.Once
)

type Config struct {
	DB     DB
	Server Server
}

func Get() (Config, error) {
	tools := tools.NewTools()

	var err error
	var configBytes []byte

	once.Do(func() {
		configPath := filepath.Join(tools.Path.Root, "configs", "config.yml")

		err = configor.New(&configor.Config{Environment: "dev", ErrorOnUnmatchedKeys: true}).Load(&config, configPath)

		if err == nil {
			configBytes, err = json.MarshalIndent(config, "", "  ")
		}

		fmt.Println("Configuration:", string(configBytes))
	})

	return config, err
}
