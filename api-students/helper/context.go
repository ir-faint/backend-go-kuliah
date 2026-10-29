package helper

import (
	"api-students/app/model"

	"github.com/gofiber/fiber/v2"
)

func CurrentUser(c *fiber.Ctx) (model.AuthUser, bool) {
	u, ok := c.Locals("authUser").(model.AuthUser)
	return u, ok
}

// RequestID mengambil request ID unik dari context Fiber yang di-generate oleh middleware requestid.
func RequestID(c *fiber.Ctx) string {
	if id, ok := c.Locals("requestid").(string); ok {
		return id
	}
	return ""
}
