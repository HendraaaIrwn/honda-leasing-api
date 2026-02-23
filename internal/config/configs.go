package configs

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// Config maps the root structure of configs.development.toml.
type Config struct {
	Environment string         `mapstructure:"ENVIRONMENT" toml:"ENVIRONMENT"`
	Server      ServerConfig   `mapstructure:"SERVER" toml:"SERVER"`
	Database    DatabaseConfig `mapstructure:"DATABASE" toml:"DATABASE"`
	JWT         JWTConfig      `mapstructure:"JWT" toml:"JWT"`
	Storage     StorageConfig  `mapstructure:"STORAGE" toml:"STORAGE"`
	CORS        CORSConfig     `mapstructure:"CORS" toml:"CORS"`
}

type ServerConfig struct {
	Address      string   `mapstructure:"ADDRESS" toml:"ADDRESS"`
	TrustedProxy []string `mapstructure:"TRUSTED_PROXY" toml:"TRUSTED_PROXY"`
	BasePath     string   `mapstructure:"BASE_PATH" toml:"BASE_PATH"`
	ReadTimeout  int      `mapstructure:"READ_TIMEOUT" toml:"READ_TIMEOUT"`
	WriteTimeout int      `mapstructure:"WRITE_TIMEOUT" toml:"WRITE_TIMEOUT"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"HOST" toml:"HOST"`
	Port            string `mapstructure:"PORT" toml:"PORT"`
	User            string `mapstructure:"USER" toml:"USER"`
	Password        string `mapstructure:"PASSWORD" toml:"PASSWORD"`
	Name            string `mapstructure:"NAME" toml:"NAME"`
	SSLMode         string `mapstructure:"SSL_MODE" toml:"SSL_MODE"`
	Timezone        string `mapstructure:"TIMEZONE" toml:"TIMEZONE"`
	MaxOpenConns    int    `mapstructure:"MAX_OPEN_CONNS" toml:"MAX_OPEN_CONNS"`
	MaxIdleConns    int    `mapstructure:"MAX_IDLE_CONNS" toml:"MAX_IDLE_CONNS"`
	ConnMaxLifetime int    `mapstructure:"CONN_MAX_LIFETIME" toml:"CONN_MAX_LIFETIME"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"SECRET" toml:"SECRET"`
	ExpiryHours int    `mapstructure:"EXPIRY_HOURS" toml:"EXPIRY_HOURS"`
}

type StorageConfig struct {
	UploadPath  string              `mapstructure:"UPLOAD_PATH" toml:"UPLOAD_PATH"`
	MaxFileSize int64               `mapstructure:"MAX_FILE_SIZE" toml:"MAX_FILE_SIZE"`
	PublicURL   string              `mapstructure:"PUBLIC_URL" toml:"PUBLIC_URL"`
	Employees   StorageSubdirConfig `mapstructure:"EMPLOYEES" toml:"EMPLOYEES"`
}

type AllowedTypesConfig struct {
	Images    []string `mapstructure:"IMAGES" toml:"IMAGES"`
	Documents []string `mapstructure:"DOCUMENTS" toml:"DOCUMENTS"`
}

type StorageSubdirConfig struct {
	MaxSize      int64    `mapstructure:"MAX_SIZE" toml:"MAX_SIZE"`
	AllowedTypes []string `mapstructure:"ALLOWED_TYPES" toml:"ALLOWED_TYPES"`
	Subdirectory string   `mapstructure:"SUBDIRECTORY" toml:"SUBDIRECTORY"`
}

type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"ALLOWED_ORIGINS" toml:"ALLOWED_ORIGINS"`
	AllowedMethods   []string `mapstructure:"ALLOWED_METHODS" toml:"ALLOWED_METHODS"`
	AllowedHeaders   []string `mapstructure:"ALLOWED_HEADERS" toml:"ALLOWED_HEADERS"`
	AllowCredentials bool     `mapstructure:"ALLOW_CREDENTIALS" toml:"ALLOW_CREDENTIALS"`
}

func loadConfig() (*Config, error) {
	//buat default environment, bisa di override dengan env var ENVIRONMENT
	env := "development"
	if env == "" {
		env = "development"
	}

	configDir, err := getConfigDir()
	if err != nil {
		log.Printf("Error getting config directory: %v", err)
		configDir = "."
		return nil, err
	}

	connfigPath := configDir + "/configs." + env + ".toml"
	log.Printf("Loading configuration from: %s", connfigPath)

	viper.SetConfigFile(connfigPath)
	viper.AutomaticEnv()

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("Error unmarshalling config: %v", err)
		return nil, err
	} else {
		log.Printf("Configuration loaded successfully: %+v", config)
	}

	return &config, nil
}

func LoadConfig() (*Config, error) {
	return loadConfig()
}

func getConfigDir() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine caller information")
	}

	return filepath.Dir(filename), nil
}

func setDefaults() {
	viper.SetDefault("SERVER.ADDRESS", ":8080")
	viper.SetDefault("SERVER.TRUSTED_PROXY", []string{})
	viper.SetDefault("SERVER.BASE_PATH", "")
	viper.SetDefault("SERVER.READ_TIMEOUT", 15)
	viper.SetDefault("SERVER.WRITE_TIMEOUT", 15)

	viper.SetDefault("DATABASE.HOST", "localhost")
	viper.SetDefault("DATABASE.PORT", "5432")
	viper.SetDefault("DATABASE.USER", "postgres")
	viper.SetDefault("DATABASE.PASSWORD", "password")
	viper.SetDefault("DATABASE.NAME", "myapp")
	viper.SetDefault("DATABASE.SSL_MODE", "disable")
	viper.SetDefault("DATABASE.TIMEZONE", "Asia/Jakarta")
	viper.SetDefault("DATABASE.MAX_OPEN_CONNS", 25)
	viper.SetDefault("DATABASE.MAX_IDLE_CONNS", 25)
	viper.SetDefault("DATABASE.CONN_MAX_LIFETIME", 5)

	viper.SetDefault("JWT.SECRET", "your-secret-key")
	viper.SetDefault("JWT.EXPIRY_HOURS", 24)

	viper.SetDefault("STORAGE.UPLOAD_PATH", "./uploads")
	viper.SetDefault("STORAGE.MAX_FILE_SIZE", 10*1024*1024) // 10 MB
	viper.SetDefault("STORAGE.PUBLIC_URL", "http://localhost:8080/uploads")
	viper.SetDefault("STORAGE.EMPLOYEES.MAX_SIZE", 5*1024*1024) // 5 MB
	viper.SetDefault("STORAGE.EMPLOYEES.ALLOWED_TYPES", []string{"image/jpeg", "image/png"})
	viper.SetDefault("STORAGE.EMPLOYEES.SUBDIRECTORY", "employees")

	viper.SetDefault("CORS.ALLOWED_ORIGINS", []string{"*"})
	viper.SetDefault("CORS.ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE"})
	viper.SetDefault("CORS.ALLOWED_HEADERS", []string{"Content-Type", "Authorization"})
	viper.SetDefault("CORS.ALLOW_CREDENTIALS", true)
}
