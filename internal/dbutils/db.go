package dbutils

import (
	"errors"
	"log"
	"time"

	"github.com/abramad-labs/irbankmock/internal/conf"
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type dbCtxKeyType struct{}

var key dbCtxKeyType

var gormLogger logger.Interface

func InitializeDb() (*gorm.DB, error) {
	dbpath := conf.GetDbPath()
	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	gormLogger = logger.New(log.Default(), logger.Config{
		SlowThreshold:             365 * 24 * time.Hour,
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	})

	return db, err
}

func ContextWithDb(c *fiber.Ctx, db *gorm.DB) *fiber.Ctx {
	c.Locals(key, db)
	return c
}

func GetDb(c *fiber.Ctx) (*gorm.DB, error) {
	db := c.Locals(key).(*gorm.DB)
	if db == nil {
		return nil, errors.New("db is not initialized")
	}
	if conf.IsGormLogDisabled() || gormLogger == nil {
		return db.Session(&gorm.Session{}), nil
	}
	return db.Session(&gorm.Session{Logger: gormLogger}), nil
}
