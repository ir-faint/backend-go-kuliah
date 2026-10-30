package model

import "time"

type Cursor struct {
	CreatedAt time.Time `json:"created_ad"`
	ID        int       `json:"id"`
}

type CursorMeta struct {
	Limit      int    `json:"limit"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}
