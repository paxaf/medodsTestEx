package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	APIServer APIServer `mapstructure:"api_server"`
	Database  Database  `mapstructure:"database"`
}

type APIServer struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

func (c *Database) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Name)
}

func MustLoad() (*Config, error) {
	_ = godotenv.Load()
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigFile("./config/config.yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatal(err, "failed to read config")
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation error: %w", err)
	}

	fmt.Println(cfg)
	return &cfg, nil
}

func (c *Config) validate() error {
	if c.APIServer.Host == "" || c.APIServer.Port == "" {
		return fmt.Errorf("API_SERVER_HOST and API_SERVER_PORT are required")
	}

	if c.Database.Host == "" || c.Database.Port == "" {
		return fmt.Errorf("DB_HOST and DB_PORT are required")
	}

	return nil
}
