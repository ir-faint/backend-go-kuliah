package helper

import (
	"reflect"
	"testing"
)

func TestPermissionSet(t *testing.T) {
	raw := map[string][]string{
		"admin": {"student:list", "student:read:any", "student:create", "student:update:any", "student:delete"},
		"staff": {"student:list", "student:read:any", "student:create"},
		"user":  {},
	}

	ps := NewPermissionSet(raw)

	// Test Can
	if !ps.Can("admin", "student:delete") {
		t.Error("admin harus memiliki izin student:delete")
	}
	if !ps.Can("staff", "student:create") {
		t.Error("staff harus memiliki izin student:create")
	}
	if ps.Can("staff", "student:delete") {
		t.Error("staff TIDAK boleh memiliki izin student:delete")
	}
	if ps.Can("user", "student:list") {
		t.Error("user TIDAK boleh memiliki izin student:list")
	}
	if ps.Can("unknown_role", "student:list") {
		t.Error("role tidak dikenal TIDAK boleh memiliki izin apa pun")
	}

	// Test Nil receiver (fail-closed)
	var nilPS *PermissionSet
	if nilPS.Can("admin", "student:delete") {
		t.Error("nil PermissionSet harus mengembalikan false untuk Can")
	}
	if len(nilPS.PermissionsOf("admin")) != 0 {
		t.Error("nil PermissionSet harus mengembalikan slice kosong untuk PermissionsOf")
	}
	if len(nilPS.KnownRoles()) != 0 {
		t.Error("nil PermissionSet harus mengembalikan slice kosong untuk KnownRoles")
	}
	if nilPS.IsKnownRole("admin") {
		t.Error("nil PermissionSet harus mengembalikan false untuk IsKnownRole")
	}

	// Test PermissionsOf (terurut)
	staffPerms := ps.PermissionsOf("staff")
	expectedStaffPerms := []string{"student:create", "student:list", "student:read:any"}
	if !reflect.DeepEqual(staffPerms, expectedStaffPerms) {
		t.Errorf("diharapkan izin staff %v, didapat %v", expectedStaffPerms, staffPerms)
	}

	userPerms := ps.PermissionsOf("user")
	if len(userPerms) != 0 {
		t.Errorf("diharapkan izin user kosong, didapat %v", userPerms)
	}

	// Test KnownRoles (terurut)
	roles := ps.KnownRoles()
	expectedRoles := []string{"admin", "staff", "user"}
	if !reflect.DeepEqual(roles, expectedRoles) {
		t.Errorf("diharapkan known roles %v, didapat %v", expectedRoles, roles)
	}

	// Test IsKnownRole
	if !ps.IsKnownRole("admin") || !ps.IsKnownRole("user") {
		t.Error("admin dan user harus dikenali sebagai role terdaftar")
	}
	if ps.IsKnownRole("guest") {
		t.Error("guest tidak boleh dikenali sebagai role terdaftar")
	}
}
