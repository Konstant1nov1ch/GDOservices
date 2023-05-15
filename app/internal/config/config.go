package config

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug       bool `env:"IS_DEBUG" env-default:"false"`
	IsDevelopment bool `env:"IS_DEV" env-default:"true"`
	Listen        struct {
		Type       string `env:"LISTEN_TYPE" env-default:"port" env-description:"'port' or 'sock'. if 'sock' then env 'SOCKET_FILE' is required"`
		BindIP     string `env:"BIND_IP" env-default:"0.0.0.0"`
		Port       string `env:"PORT" env-default:"10000"`
		SocketFile string `env:"SOCKET_FILE" env-default:"app.sock"`
	}
	AppConfig struct {
		LogLevel  string `env:"LOG_LEVEL" env-default:"trace"`
		AdminUser struct {
			Email    string `env:"ADMIN_EMAIL" env-default:"admin"`
			Password string `env:"ADMIN_PWD" env-default:"admin"`
		}
	}
	PostgreSQL struct {
		Username string `env:"PSQL_USERNAME" env-default:"konstantin"`
		Password string `env:"PSQL_PASSWORD" env-default:"konstantin"`
		Host     string `env:"PSQL_HOST" env-default:"localhost"`
		Port     string `env:"PSQL_PORT" env-default:"5432"`
		Database string `env:"PSQL_DATABASE" env-default:"todo"`
	}
	DB *sql.DB // Добавляем поле для хранения соединения с базой данных
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		log.Print("gather config")

		instance = &Config{}

		if err := cleanenv.ReadEnv(instance); err != nil {
			helpText := "Monolith Notes System"
			help, _ := cleanenv.GetDescription(instance, &helpText)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}

func (c *Config) GetDB() (*sql.DB, error) {
	cfg := GetConfig()
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgreSQL.Host, cfg.PostgreSQL.Port, cfg.PostgreSQL.Username, cfg.PostgreSQL.Password, cfg.PostgreSQL.Database))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
