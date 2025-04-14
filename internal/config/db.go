package config

type DB struct {
	Host            string
	Port            string
	User            string
	Password        string `yaml:"password,omitempty"`
	Database        string
	MigrationsTable string `yaml:"migrations_table"`
}
