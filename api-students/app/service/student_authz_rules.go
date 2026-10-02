package service

import (
	"api-students/app/model"
	"api-students/helper"
	"strings"
)

// CanAccessStudent memeriksa apakah pengguna dapat mengakses/mengubah data mahasiswa.
// Akses diberikan jika pengguna adalah pemilik data (current.ID == ownerID)
// ATAU jika role pengguna memiliki izin (permission) yang ditentukan.
func CanAccessStudent(current model.AuthUser, ownerID int, perms *helper.PermissionSet, anyPermission string) bool {
	if current.ID == ownerID {
		return true
	}
	return perms.Can(current.Role, anyPermission)
}

func ValidateAssignRole(current model.AuthUser, targetID int, req model.AssignRoleRequest, perms *helper.PermissionSet) map[string]string {
	errs := map[string]string{}

	role := strings.TrimSpace(req.Role)
	if role == "" {
		errs["role"] = "wajib diisi"
		return errs
	}

	if !perms.IsKnownRole(role) {
		errs["role"] = "role tidak dikenal, pilih salah satu dari: " +
			strings.Join(perms.KnownRoles(), ", ")
	}

	if current.ID == targetID {
		errs["role"] = "tidak boleh mengubah role diri sendiri"
	}

	return errs
}
