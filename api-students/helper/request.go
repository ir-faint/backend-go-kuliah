package helper

import (
	"api-students/app/model"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var allowedSortFields = map[string]bool{
	"id":    true,
	"nim":   true,
	"name":  true,
	"grade": true,
}

func ParseListQuery(c *fiber.Ctx) model.ListQuery {
	q := model.ListQuery{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: strings.TrimSpace(c.Query("search")),
		Sort:   c.Query("sort", "id"),
		Order:  strings.ToLower(c.Query("order", "asc")),
	}

	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 10
	}
	if q.Limit > 100 {
		q.Limit = 100
	}
	if !allowedSortFields[q.Sort] {
		q.Sort = "id"
	}
	if q.Order != "desc" {
		q.Order = "asc"
	}

	if raw := c.Query("is_active"); raw != "" {
		if v, err := strconv.ParseBool(raw); err == nil {
			q.IsActive = &v
		}
	}

	return q
}

// ParseCursorQuery membaca parameter pagination cursor dari query string URL
func ParseCursorQuery(c *fiber.Ctx) (model.CursorQuery, error) {
	q := model.CursorQuery{
		Limit:  c.QueryInt("limit", 10),
		Search: strings.TrimSpace(c.Query("search")),
	}

	if q.Limit < 1 {
		q.Limit = 10
	}
	if q.Limit > 100 {
		q.Limit = 100
	}

	// Baca cursor jika dikirim
	if rawCursor := strings.TrimSpace(c.Query("cursor")); rawCursor != "" {
		cur, err := DecodeCursor(rawCursor)
		if err != nil {
			return q, BadRequest("cursor tidak valid")
		}
		q.After = &cur
	}

	// Baca filter is_active jika dikirim
	if raw := strings.TrimSpace(c.Query("is_active")); raw != "" {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return q, BadRequest("is_active harus bernilai boolean (true/false)")
		}
		q.IsActive = &v
	}

	return q, nil
}

func RequestContext(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.UserContext(), 5*time.Second)
}

func ParamID(c *fiber.Ctx) (int, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}
