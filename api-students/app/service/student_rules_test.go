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
