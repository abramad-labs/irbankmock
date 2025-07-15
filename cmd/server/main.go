package main

import (
	"log"

	_ "github.com/abramad-labs/irbankmock/internal/banks"
	"github.com/abramad-labs/irbankmock/internal/banks/registry"
	fibererror "github.com/abramad-labs/irbankmock/internal/usererror/fiber"
	_ "go.uber.org/automaxprocs"

	"github.com/abramad-labs/irbankmock/internal/conf"
	"github.com/abramad-labs/irbankmock/internal/dbutils"
	"github.com/abramad-labs/irbankmock/internal/dbutils/migration"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	db, err := dbutils.InitializeDb()
	if err != nil {
		log.Fatalf("failed to init sqlite db: %s", err.Error())
	}

	if conf.ShouldAutoMigrate() {
		migrator := db.Migrator()
		err = migration.ApplyMigrations(migrator)
		if err != nil {
			log.Fatalf("migration failed: %s", err.Error())
		}
	}

	app := fiber.New(fiber.Config{
		CaseSensitive: false,
		ErrorHandler:  fibererror.FiberUserErrorHandling,
	})

	app.Static("/", conf.GetWebAppPath(), fiber.Static{
		Browse: false,
	})
	app.Use(func(c *fiber.Ctx) error {
		c = dbutils.ContextWithDb(c, db)
		return c.Next()
	})
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "${locals:requestid} ${status} - ${method} ${path}\u200b\n",
	}))

	registry.ConfigAppRouters(app)

	listenaddr := conf.GetListenAddress()
	err = app.Listen(listenaddr)
	if err != nil {
		log.Fatal(err)
	}
}
