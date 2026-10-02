package middleware

import (
	"fmt"

	"api-students/helper"

	"github.com/gofiber/fiber/v2"
)

func RequirePermission(perms *helper.PermissionSet, permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := helper.CurrentUser(c)
		if !ok {
			return helper.Fail(c, fiber.StatusUnauthorized, "pengguna tidak terautentikasi")
		}

		if !perms.Can(user.Role, permission) {
			msg := fmt.Sprintf("role %s tidak memiliki hak %s", user.Role, permission)
			return helper.Fail(c, fiber.StatusForbidden, msg)
		}

		return c.Next()
	}
}

// func RequireRole(roles ...string) fiber.Handler {
// 	allowed := make(map[string]struct{}, len(roles))
// 	for _, role := range roles {
// 		allowed[role] = struct{}{}
// 	}

// 	return func(c *fiber.Ctx) error {
// 		user, ok := helper.CurrentUser(c)
// 		if !ok {
// 			return helper.Fail(c, fiber.StatusUnauthorized, "belum terautentikasi")
// 		}

// 		if _, granted := allowed[user.Role]; !granted {
// 			return helper.Fail(c, fiber.StatusForbidden,
// 				"role Anda tidak berhak mengakses endpoint ini")
// 		}

// 		return c.Next()
// 	}
// }
