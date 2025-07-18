package config

import (
	provider "PocGo/internal/configuration/providers"
	envConfig "github.com/joho/godotenv"
	"log"
	setter "os"
	"strconv"
)

type Config struct {
	Database *provider.DatabaseConfig
	Routine  *provider.RoutineConfig
}

func LoadConfig(env string) *Config {
	err := envConfig.Load(".env." + env)
	if err != nil {
		log.Println("Erro ao carregar arquivo : ", err)
	}

	connStr := setter.Getenv("DB_CONNECTION_STRING")
	maxOpers, _ := strconv.Atoi(setter.Getenv("DB_MAX_OPEN_CONNECTIONS"))
	maxIdle, _ := strconv.Atoi(setter.Getenv("DB_MAX_IDLE_CONNECTIONS"))

	day, _ := strconv.Atoi(setter.Getenv("RT_INCREMENT_DAY"))
	hour, _ := strconv.Atoi(setter.Getenv("RT_HOUR"))
	minute, _ := strconv.Atoi(setter.Getenv("RT_MINUTE"))
	second, _ := strconv.Atoi(setter.Getenv("RT_SECOND"))
	millisecond, _ := strconv.Atoi(setter.Getenv("RT_MILLISECOND"))

	return &Config{
		Database: &provider.DatabaseConfig{
			ConnectionString: connStr,
			MaxOpenConns:     maxOpers,
			MaxIdleConns:     maxIdle,
		},
		Routine: &provider.RoutineConfig{
			IncrementDay: day,
			Hour:         hour,
			Minute:       minute,
			Second:       second,
			Millisecond:  millisecond,
		},
	}
}
