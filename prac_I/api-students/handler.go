package main

import (
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var students []Student
var nextStudentID = 1

func findStudentIndex(id int) int {
	for i := range students {
		if students[i].ID == id {
			return i
		}
	}
	return -1
}

func isNIMTaken(nim string, excludeID int) bool {
	for _, s := range students {
		if s.NIM == nim && s.ID != excludeID {
			return true
		}
	}
	return false
}

func listStudents(c *fiber.Ctx) error {
	q := parseListQuery(c)
	var filtered []Student

	for _, s := range students {
		if q.IsActive != nil && s.IsActive != *q.IsActive {
			continue
		}
		if q.Search != "" && !strings.Contains(strings.ToLower(s.Name), strings.ToLower(q.Search)) {
			continue
		}
		filtered = append(filtered, s)
	}

	sort.SliceStable(filtered, func(i, j int) bool {
		var less bool
		switch q.Sort {
		case "nim":
			less = filtered[i].NIM < filtered[j].NIM
		case "name":
			less = filtered[i].Name < filtered[j].Name
		case "grade":
			less = filtered[i].Grade < filtered[j].Grade
		default:
			less = filtered[i].ID < filtered[j].ID
		}
		if q.Order == "desc" {
			return !less
		}
		return less
	})

	total := len(filtered)
	totalPages := (total + q.Limit - 1) / q.Limit
	if totalPages == 0 {
		totalPages = 1
	}

	start := (q.Page - 1) * q.Limit
	if start > total {
		start = total
	}
	end := start + q.Limit
	if end > total {
		end = total
	}

	return okList(c, "daftar mahasiswa berhasil diambil", filtered[start:end], &Meta{
		Page: q.Page, Limit: q.Limit, Total: total, TotalPages: totalPages,
	})
}

func getStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	i := findStudentIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	}
	return ok(c, "mahasiswa ditemukan", students[i])
}

func createStudent(c *fiber.Ctx) error {
	var req CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	errs := map[string]string{}
	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	if req.NIM == "" {
		errs["nim"] = "wajib diisi"
	} else if isNIMTaken(req.NIM, 0) {
		return fail(c, fiber.StatusConflict, "NIM sudah terdaftar")
	}
	if req.Name == "" {
		errs["name"] = "wajib diisi"
	}
	if req.Grade < 0 || req.Grade > 100 {
		errs["grade"] = "nilai harus berada di rentang 0-100"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	baru := Student{
		ID:       nextStudentID,
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		IsActive: req.IsActive,
	}
	nextStudentID++
	students = append(students, baru)

	return created(c, "mahasiswa berhasil dibuat", baru, "/api/v1/students/"+strconv.Itoa(baru.ID))
}

func replaceStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	i := findStudentIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	}

	var req UpdateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	errs := map[string]string{}
	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	if req.NIM == "" {
		errs["nim"] = "wajib diisi pada PUT"
	} else if isNIMTaken(req.NIM, id) {
		return fail(c, fiber.StatusConflict, "NIM sudah terdaftar")
	}
	if req.Name == "" {
		errs["name"] = "wajib diisi pada PUT"
	}
	if req.Grade < 0 || req.Grade > 100 {
		errs["grade"] = "nilai harus berada di rentang 0-100"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	students[i].NIM = req.NIM
	students[i].Name = req.Name
	students[i].Grade = req.Grade
	students[i].IsActive = req.IsActive

	return ok(c, "data mahasiswa berhasil diganti seluruhnya", students[i])
}

func patchStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	i := findStudentIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	}

	var req PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	if req.NIM == nil && req.Name == nil && req.Grade == nil && req.IsActive == nil {
		return fail(c, fiber.StatusBadRequest, "tidak ada field yang dikirim untuk diubah")
	}

	errs := map[string]string{}

	if req.NIM != nil {
		trimmed := strings.TrimSpace(*req.NIM)
		if trimmed == "" {
			errs["nim"] = "tidak boleh kosong"
		} else if isNIMTaken(trimmed, id) {
			return fail(c, fiber.StatusConflict, "NIM sudah terdaftar")
		}
	}
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" {
			errs["name"] = "tidak boleh kosong"
		}
	}
	if req.Grade != nil {
		if *req.Grade < 0 || *req.Grade > 100 {
			errs["grade"] = "nilai harus berada di rentang 0-100"
		}
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	if req.NIM != nil {
		students[i].NIM = strings.TrimSpace(*req.NIM)
	}
	if req.Name != nil {
		students[i].Name = strings.TrimSpace(*req.Name)
	}
	if req.Grade != nil {
		students[i].Grade = *req.Grade
	}
	if req.IsActive != nil {
		students[i].IsActive = *req.IsActive
	}

	return ok(c, "data mahasiswa berhasil diperbarui sebagian", students[i])
}

func deleteStudent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	i := findStudentIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	}

	students = append(students[:i], students[i+1:]...)
	return noContent(c)
}
