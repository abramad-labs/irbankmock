package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

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

	var staticNext func(*fiber.Ctx) bool
	if !app.Config().CaseSensitive {
		staticNext = func(c *fiber.Ctx) bool {
			path := c.Path()
			if strings.HasPrefix(path, registry.RegistryBanksPrefix) {
				c.Path(strings.ToLower(path))
			}
			return false
		}
	}

	app.Static("/", conf.GetWebAppPath(), fiber.Static{
		Browse: false,
		Next:   staticNext,
	})
	app.Use(func(c *fiber.Ctx) error {
		c = dbutils.ContextWithDb(c, db)
		return c.Next()
	})
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "${locals:requestid} ${status} - ${method} ${path}\u200b\n",
	}))

	rootGroup := app.Group("/")
	registry.ConfigAppRouters(rootGroup.(*fiber.Group))

	PrintAllRoutes(app)
	listenaddr := conf.GetListenAddress()
	err = app.Listen(listenaddr)
	if err != nil {
		log.Fatal(err)
	}
}

func PrintAllRoutes(app *fiber.App) {
	routes := app.GetRoutes()

	pathMap := make(map[string][]string)
	for _, r := range routes {
		pathMap[r.Path] = append(pathMap[r.Path], r.Method)
	}

	paths := make([]string, 0, len(pathMap))
	for p := range pathMap {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Println()
	fmt.Fprintln(w, "PATH\tMETHODS")

	for _, p := range paths {
		methods := pathMap[p]
		sort.Strings(methods)
		fmt.Fprintf(w, "%s\t%s\n", p, methods)
	}

	w.Flush()
}
