package service

import (
	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type StudentService struct {
	StudentRepo repository.StudentRepository
}

func NewStudentService(studentRepo repository.StudentRepository) *StudentService {
	return &StudentService{StudentRepo: studentRepo}
}

func (s *StudentService) ListStudents(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	q := helper.ParseListQuery(c)

	students, total, err := s.StudentRepo.FindAll(ctx, q)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "gagal mengambil data mahasiswa")
	}

	totalPages := CountTotalPages(total, q.Limit)

	return helper.SuccessList(c, "daftar mahasiswa berhasil diambil", students, &model.Meta{
		Page:       q.Page,
		Limit:      q.Limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

func (s *StudentService) GetStudent(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	mahasiswa, err := s.StudentRepo.FindByID(ctx, id)
	if err != nil {
		return errorTranslator(c, err, "gagal mengambil data mahasiswa")
	}

	return helper.Success(c, "mahasiswa ditemukan", mahasiswa)
}

func (s *StudentService) CreateStudent(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	if errs := ValidateCreate(req); len(errs) > 0 {
		return helper.FailValidation(c, errs)
	}

	baru, err := s.StudentRepo.Create(ctx, model.Student{
		NIM:      strings.TrimSpace(req.NIM),
		Name:     strings.TrimSpace(req.Name),
		Grade:    req.Grade,
		IsActive: req.IsActive,
	})
	if err != nil {
		// log.Println("ERROR DETAIL CREATE:", err)
		return errorTranslator(c, err, "gagal menyimpan mahasiswa")
	}

	return helper.Created(c, "mahasiswa berhasil dibuat", baru,
		"/api/v1/students/"+strconv.Itoa(baru.ID))
}

func (s *StudentService) ReplaceStudent(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	var req model.UpdateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	if errs := ValidateReplace(req); len(errs) > 0 {
		return helper.FailValidation(c, errs)
	}

	result, err := s.StudentRepo.Update(ctx, model.Student{
		ID:       id,
		NIM:      strings.TrimSpace(req.NIM),
		Name:     strings.TrimSpace(req.Name),
		Grade:    req.Grade,
		IsActive: req.IsActive,
	})
	if err != nil {
		return errorTranslator(c, err, "gagal memperbarui mahasiswa")
	}

	return helper.Success(c, "mahasiswa berhasil diganti seluruhnya", result)
}

func (s *StudentService) PatchStudent(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	var req model.PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	if IsEmptyPatch(req) {
		return helper.Fail(c, fiber.StatusBadRequest, "tidak ada field yang diubah")
	}

	current, err := s.StudentRepo.FindByID(ctx, id)
	if err != nil {
		return errorTranslator(c, err, "gagal mengambil data mahasiswa")
	}

	updatedStudent, errs := ApplyPatch(current, req)
	if len(errs) > 0 {
		return helper.FailValidation(c, errs)
	}

	result, err := s.StudentRepo.Update(ctx, updatedStudent)
	if err != nil {
		return errorTranslator(c, err, "gagal memperbarui mahasiswa")
	}

	return helper.Success(c, "mahasiswa berhasil diperbarui sebagian", result)
}

func (s *StudentService) DeleteStudent(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	if err := s.StudentRepo.Delete(ctx, id); err != nil {
		return errorTranslator(c, err, "gagal menghapus mahasiswa")
	}

	return helper.NoContent(c)
}

func errorTranslator(c *fiber.Ctx, err error, defaultMsg string) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return helper.Fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	case errors.Is(err, repository.ErrDuplicate):
		return helper.Fail(c, fiber.StatusConflict, "mahasiswa sudah ada")
	default:
		return helper.Fail(c, fiber.StatusInternalServerError, defaultMsg)
	}
}
