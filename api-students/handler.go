package main

import (
	"api-students/app/model"
	"api-students/app/repository"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type StudentHandler struct {
	StudentRepo repository.StudentRepository
}

func NewStudentHandler(studentRepo repository.StudentRepository) *StudentHandler {
	return &StudentHandler{StudentRepo: studentRepo}
}

func errorHandler(c *fiber.Ctx, err error, defaultMsg string) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return fail(c, fiber.StatusNotFound, "User tidak ditemukan")
	case errors.Is(err, repository.ErrDuplicate):
		return fail(c, fiber.StatusConflict, "User sudah ada")
	default:
		return fail(c, fiber.StatusInternalServerError, defaultMsg)
	}
}

func (h *StudentHandler) listStudents(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	q := parseListQuery(c)

	users, total, err := h.StudentRepo.FindAll(ctx, q)
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "gagal mengambil data mahasiswa")
	}

	totalPages := 0
	if q.Limit > 0 {
		totalPages = (total + q.Limit - 1) / q.Limit
	}

	return okList(c, "daftar mahasiswa berhasil diambil", users, &model.Meta{
		Page: q.Page, Limit: q.Limit, Total: total, TotalPages: totalPages,
	})
}

func (h *StudentHandler) getStudent(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	mahasiswa, err := h.StudentRepo.FindByID(ctx, id)
	if err != nil {
		return errorHandler(c, err, "gagal mengambil data mahasiswa")
	}

	return ok(c, "mahasiswa ditemukan", mahasiswa)

}

func (h *StudentHandler) createStudent(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	errs := map[string]string{}
	if req.NIM == "" {
		errs["nim"] = "wajib diisi"
	}
	if req.Name == "" {
		errs["name"] = "wajib diisi"
	}
	if req.Grade < 0 || req.Grade > 4.00 {
		errs["grade"] = "nilai harus berada di rentang 0.00-4.00"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	baru, err := h.StudentRepo.Create(ctx, model.Student{
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		IsActive: req.IsActive,
	})
	if err != nil {
		log.Println("ERROR DETAIL CREATE:", err)
		return errorHandler(c, err, "gagal menyimpan mahasiswa")
	}

	return created(c, "mahasiswa berhasil dibuat", baru,
		"/api/v1/students/"+strconv.Itoa(baru.ID))
}

func (h *StudentHandler) replaceStudent(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	var req model.UpdateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	errs := map[string]string{}
	if strings.TrimSpace(req.NIM) == "" {
		errs["nim"] = "wajib diisi pada PUT"
	}
	if strings.TrimSpace(req.Name) == "" {
		errs["name"] = "wajib diisi pada PUT"
	}
	if req.Grade < 0 || req.Grade > 4.00 {
		errs["grade"] = "nilai harus berada di rentang 0.00-4.00"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	hasil, err := h.StudentRepo.Update(ctx, model.Student{
		ID: id, NIM: req.NIM, Name: req.Name, Grade: req.Grade, IsActive: req.IsActive,
	})
	if err != nil {
		return errorHandler(c, err, "gagal memperbarui mahasiswa")
	}

	return ok(c, "mahasiswa berhasil diganti seluruhnya", hasil)
}

func (h *StudentHandler) patchStudent(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	var req model.PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	if req.NIM == nil && req.Name == nil && req.Grade == nil && req.IsActive == nil {
		return fail(c, fiber.StatusBadRequest, "tidak ada field yang diubah")
	}

	saatIni, err := h.StudentRepo.FindByID(ctx, id)
	if err != nil {
		return errorHandler(c, err, "gagal mengambil data mahasiswa")
	}

	if req.NIM != nil {
		if strings.TrimSpace(*req.NIM) == "" {
			return failValidation(c, map[string]string{"nim": "tidak boleh kosong"})
		}
		saatIni.NIM = *req.NIM
	}
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return failValidation(c, map[string]string{"name": "tidak boleh kosong"})
		}
		saatIni.Name = *req.Name
	}
	if req.Grade != nil {
		if *req.Grade < 0 || *req.Grade > 4.00 {
			return failValidation(c, map[string]string{"grade": "nilai harus berada di rentang 0.00-4.00"})
		}
		saatIni.Grade = *req.Grade
	}
	if req.IsActive != nil {
		saatIni.IsActive = *req.IsActive
	}

	hasil, err := h.StudentRepo.Update(ctx, saatIni)
	if err != nil {
		return errorHandler(c, err, "gagal memperbarui mahasiswa")
	}

	return ok(c, "mahasiswa berhasil diperbarui sebagian", hasil)
}

func (h *StudentHandler) deleteStudent(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	if err := h.StudentRepo.Delete(ctx, id); err != nil {
		return errorHandler(c, err, "gagal menghapus mahasiswa")
	}

	return noContent(c)
}
