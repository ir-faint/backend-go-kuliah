package helper

import (
	"api-students/app/model"

	"github.com/gofiber/fiber/v2"
)

func CurrentUser(c *fiber.Ctx) (model.AuthUser, bool) {
	u, ok := c.Locals("authUser").(model.AuthUser)
	return u, ok
}
