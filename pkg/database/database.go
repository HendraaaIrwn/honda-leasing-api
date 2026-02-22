package database

import (
	"fmt"
	"log"
	"time"

	configs "github.com/HendraaaIrwn/honda-leasing-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

// Databese is kept as a backward-compatible alias.
type Databese = Database

func GetDB(db *Database) *gorm.DB {
	if db.DB == nil {
		panic("Database connection is not initialized")
	}
	return db.DB
}

func CloseDB(db *Database) error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("error getting sql.DB from gorm.DB: %w", err)
	}
	return sqlDB.Close()
}

func GenerateDSN(cfg configs.DatabaseConfig) string {
	sslMode := cfg.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	timezone := cfg.Timezone
	if timezone == "" {
		timezone = "Asia/Jakarta"
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		sslMode,
		timezone,
	)
}

func AutoMigrate(db *Database, model ...any) error {
	if db.DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	log.Printf("Starting auto migration for models: %v", model)

	if err := db.DB.AutoMigrate(model...); err != nil {
		return fmt.Errorf("error during auto migration: %w", err)
	}
	return nil
}

func InitAutoMigrate(db *Database, model ...any) {
	schemas := []string{"mst", "account", "dealer", "leasing", "payment"}
	for _, schema := range schemas {
		if err := db.DB.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)).Error; err != nil {
			log.Fatalf("Failed creating schema %s: %v", schema, err)
		}
	}

	if err := AutoMigrate(db, model...); err != nil {
		log.Fatalf("Auto migration failed: %v", err)
	}
	log.Println("Auto migration completed successfully")
}
func InitDB(cfg *configs.Config) (*Database, error) {
	dsn := GenerateDSN(cfg.Database)

	log.Printf("Connecting to database with DSN: %s@%s:%s/%s",
		cfg.Database.User,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name)

	gomrmConfig := &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info), // Enable detailed logging for debugging
	}

	if cfg.Environment == "development" {
		gomrmConfig.Logger = logger.Default.LogMode(logger.Info) // Enable detailed logging in development
	} else {
		gomrmConfig.Logger = logger.Default.LogMode(logger.Silent) // Disable logging in production
	}

	db, err := gorm.Open(postgres.Open(dsn), gomrmConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error getting sql.DB from gorm.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	log.Println("Database connection established successfully")
	return &Database{DB: db}, nil
}

func SetupDB() (*Database, error) {
	cfg, err := configs.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return InitDB(cfg)
}
