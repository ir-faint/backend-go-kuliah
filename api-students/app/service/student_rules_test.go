package service

import (
	"testing"

	"api-students/app/model"
)

// Perhatikan: pengujian ini tidak menyalakan server, tidak menyentuh
// database, dan tidak membuat fiber.Ctx.
func TestCountTotalPages(t *testing.T) {
	cases := []struct{ total, limit, want int }{
		{0, 10, 0},
		{1, 10, 1},
		{10, 10, 1},
		{11, 10, 2},
		{137, 20, 7},
	}

	for _, tc := range cases {
		if got := CountTotalPages(tc.total, tc.limit); got != tc.want {
			t.Errorf("total=%d limit=%d: harap %d, dapat %d",
				tc.total, tc.limit, tc.want, got)
		}
	}
}

func TestValidateCreate(t *testing.T) {
	// Case 1: Valid request
	validReq := model.CreateStudentRequest{
		NIM:      "22010101",
		Name:     "Budi Santoso",
		Grade:    3.75,
		IsActive: true,
	}
	if errs := ValidateCreate(validReq); len(errs) != 0 {
		t.Errorf("request POST valid tidak boleh menghasilkan error: %v", errs)
	}

	// Case 2: Invalid request (NIM & Name empty, Grade > 4.00)
	invalidReq := model.CreateStudentRequest{
		NIM:   "   ",
		Name:  "",
		Grade: 5.00,
	}
	errs := ValidateCreate(invalidReq)
	if _, ok := errs["nim"]; !ok {
		t.Error("nim kosong harus menghasilkan error validasi")
	}
	if _, ok := errs["name"]; !ok {
		t.Error("name kosong harus menghasilkan error validasi")
	}
	if _, ok := errs["grade"]; !ok {
		t.Error("grade di luar rentang 0-4.00 harus menghasilkan error validasi")
	}
}

func TestValidateReplace(t *testing.T) {
	// Case 1: Valid PUT request
	validReq := model.UpdateStudentRequest{
		ID:       1,
		NIM:      "22010101",
		Name:     "Budi Santoso Updated",
		Grade:    3.90,
		IsActive: true,
	}
	if errs := ValidateReplace(validReq); len(errs) != 0 {
		t.Errorf("request PUT valid tidak boleh menghasilkan error: %v", errs)
	}

	// Case 2: Invalid PUT request (NIM & Name empty)
	invalidReq := model.UpdateStudentRequest{
		ID:    1,
		NIM:   "",
		Name:  "",
		Grade: 3.50,
	}
	errs := ValidateReplace(invalidReq)
	if msg, ok := errs["nim"]; !ok || msg != "wajib diisi pada PUT" {
		t.Errorf("nim kosong pada PUT harus menghasilkan error 'wajib diisi pada PUT', dapat: %v", errs["nim"])
	}
	if msg, ok := errs["name"]; !ok || msg != "wajib diisi pada PUT" {
		t.Errorf("name kosong pada PUT harus menghasilkan error 'wajib diisi pada PUT', dapat: %v", errs["name"])
	}
}

func TestApplyPatch(t *testing.T) {
	initial := model.Student{
		ID:       1,
		NIM:      "22010101",
		Name:     "Sari",
		Grade:    3.50,
		IsActive: true,
	}
	inactive := false
	newName := "Sari Updated"

	result, errs := ApplyPatch(initial, model.PatchStudentRequest{
		Name:     &newName,
		IsActive: &inactive,
	})

	if len(errs) != 0 {
		t.Fatalf("tidak seharusnya ada error: %v", errs)
	}
	if result.IsActive {
		t.Error("is_active seharusnya berubah menjadi false")
	}
	if result.Name != "Sari Updated" {
		t.Errorf("name seharusnya berubah menjadi 'Sari Updated', dapat %s", result.Name)
	}
	if result.NIM != "22010101" {
		t.Error("field yang tidak dikirim (NIM) seharusnya tidak berubah")
	}
}
