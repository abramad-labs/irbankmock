package registry

import (
	"github.com/gofiber/fiber/v2"
)

type routerPrefixType struct{}

var routerPrefixKey routerPrefixType

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

func ConfigAppRouters(app *fiber.Group) {
	for _, entry := range banks {
		grpPath := "/banks/" + entry.Name
		grp := app.Group(grpPath)
		grp.Use(func(c *fiber.Ctx) error {
			c.Locals(routerPrefixKey, grpPath)
			return c.Next()
		})
		entry.Action(grp)
	}
}

func GetRouterPrefix(c *fiber.Ctx) string {
	return c.Locals(routerPrefixKey).(string)
}
