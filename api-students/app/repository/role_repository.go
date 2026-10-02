package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepository interface {
	LoadPermissions(ctx context.Context) (map[string][]string, error)
}

type roleRepository struct {
	pool *pgxpool.Pool
}

func NewRoleRepository(pool *pgxpool.Pool) RoleRepository {
	return &roleRepository{pool: pool}
}

func (r *roleRepository) LoadPermissions(ctx context.Context) (map[string][]string, error) {
	sqlText := `
		SELECT r.name, COALESCE(rp.permission_name, '')
		FROM roles r
		LEFT JOIN role_permissions rp ON r.name = rp.role_name
		ORDER BY r.name, rp.permission_name
	`

	rows, err := r.pool.Query(ctx, sqlText)
	if err != nil {
		return nil, fmt.Errorf("memuat data permission role: %w", err)
	}
	defer rows.Close()

	result := make(map[string][]string)

	for rows.Next() {
		var roleName, permName string
		if err := rows.Scan(&roleName, &permName); err != nil {
			return nil, fmt.Errorf("membaca baris permission role: %w", err)
		}

		if _, exists := result[roleName]; !exists {
			result[roleName] = []string{}
		}

		if permName != "" {
			result[roleName] = append(result[roleName], permName)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("membaca hasil query permission role: %w", err)
	}

	return result, nil
}
