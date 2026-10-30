package service

import (
	"testing"

	"api-students/app/model"
	"api-students/helper"
)

func TestCanAccessStudent(t *testing.T) {
	rawPerms := map[string][]string{
		"admin": {"student:read:any", "student:update:any"},
		"staff": {"student:read:any"},
		"user":  {},
	}
	perms := helper.NewPermissionSet(rawPerms)

	userOwner := model.AuthUser{ID: 10, Username: "student_user", Role: "user", IsActive: true}
	userOther := model.AuthUser{ID: 20, Username: "other_user", Role: "user", IsActive: true}
	staffUser := model.AuthUser{ID: 30, Username: "staff_user", Role: "staff", IsActive: true}
	adminUser := model.AuthUser{ID: 40, Username: "admin_user", Role: "admin", IsActive: true}

	studentOwnerID := 10

	// 1. Pemilik data (selalu diizinkan meski role user tidak punya permission eksplisit)
	if !CanAccessStudent(userOwner, studentOwnerID, perms, "student:read:any") {
		t.Error("pemilik data harus dapat mengakses data milik sendiri")
	}
	if !CanAccessStudent(userOwner, studentOwnerID, perms, "student:update:any") {
		t.Error("pemilik data harus dapat mengubah data milik sendiri")
	}

	// 2. User lain bukan pemilik data (harus ditolak)
	if CanAccessStudent(userOther, studentOwnerID, perms, "student:read:any") {
		t.Error("user bukan pemilik data TIDAK boleh mengakses data milik pengguna lain")
	}
	if CanAccessStudent(userOther, studentOwnerID, perms, "student:update:any") {
		t.Error("user bukan pemilik data TIDAK boleh mengubah data milik pengguna lain")
	}

	// 3. Staf (student:read:any diizinkan, student:update:any ditolak)
	if !CanAccessStudent(staffUser, studentOwnerID, perms, "student:read:any") {
		t.Error("staf harus dapat membaca data mahasiswa milik siapa saja")
	}
	if CanAccessStudent(staffUser, studentOwnerID, perms, "student:update:any") {
		t.Error("staf TIDAK boleh mengubah data mahasiswa jika tidak memiliki izin student:update:any")
	}

	// 4. Admin (semua diizinkan)
	if !CanAccessStudent(adminUser, studentOwnerID, perms, "student:read:any") {
		t.Error("admin harus dapat membaca data mahasiswa milik siapa saja")
	}
	if !CanAccessStudent(adminUser, studentOwnerID, perms, "student:update:any") {
		t.Error("admin harus dapat mengubah data mahasiswa milik siapa saja")
	}
}
