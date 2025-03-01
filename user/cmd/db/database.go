package db

import (
	"fmt"
	"os"
	"time"
	"user_service/types"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(c *gin.Context) (*gorm.DB, error) {
	dbModel := []interface{}{&types.User{}}

	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dsn := fmt.Sprintf("host=postgres user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Jakarta", dbUser, dbPass, dbName)
	for i := 0; i < 5; i++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			for _, v := range dbModel {
				db.AutoMigrate(v)
			}
			return db, nil
		}

		logrus.Warnf("Database connection failed. Retrying in 5 seconds... (%d/5)", i+1)
		time.Sleep(5 * time.Second)
	}
	return nil, fmt.Errorf("failed to connect to database")
}
