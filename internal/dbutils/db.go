package dbutils

import (
	"github.com/abramad-labs/irbankmock/internal/conf"
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type dbCtxKeyType struct{}

var key dbCtxKeyType

func InitializeDb() (*gorm.DB, error) {
	dbpath := conf.GetDbPath()
	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, err
}

func ContextWithDb(c *fiber.Ctx, db *gorm.DB) *fiber.Ctx {
	c.Locals(key, db)
	return c
}

func GetDb(c *fiber.Ctx) *gorm.DB {
	return c.Locals(key).(*gorm.DB)
}
