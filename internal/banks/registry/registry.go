package registry

import (
	"github.com/gofiber/fiber/v2"
)

type bankEntry struct {
	Name   string
	Action func(g fiber.Router)
}

var banks []bankEntry

func RegisterBank(name string, action func(g fiber.Router)) {
	banks = append(banks, bankEntry{
		Name:   name,
		Action: action,
	})
}

func Cleanup() {
	banks = nil
}

func ConfigAppRouters(app *fiber.App) {
	for _, entry := range banks {
		grp := app.Group("/api/" + entry.Name)
		entry.Action(grp)
	}
}
