package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/mr-emerald-wolf/yantra-backend/models"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(config *Config) {
	var err error
	// dsn := viper.GetString("DATABASE_URL")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", viper.GetString("POSTGRES_HOST"), viper.GetString("POSTGRES_USER"), viper.GetString("POSTGRES_PASSWORD"), viper.GetString("POSTGRES_DB"), viper.GetString("POSTGRES_PORT"))

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database! \n", err.Error())
		os.Exit(1)
	}

	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	DB.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running Migrations")
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Ngo{})
	DB.AutoMigrate(&models.Request{})
	DB.AutoMigrate(&models.Volunteer{})
	DB.AutoMigrate(&models.Event{})
	DB.AutoMigrate(&models.Blog{})

	log.Println("🚀 Connected Successfully to the Database")
}
