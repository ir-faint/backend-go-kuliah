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
	perms       *helper.PermissionSet
}

func NewStudentService(studentRepo repository.StudentRepository, perms *helper.PermissionSet) *StudentService {
	return &StudentService{
		StudentRepo: studentRepo,
		perms:       perms,
	}
}

func (s *StudentService) ListStudents(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	q := helper.ParseListQuery(c)

	students, total, err := s.StudentRepo.FindAll(ctx, q)
	if err != nil {
		return helper.Internal(err)
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
		return helper.BadRequest("id harus berupa angka positif")
	}

	mahasiswa, err := s.StudentRepo.FindByID(ctx, id)
	if err != nil {
		return errorTranslator(err)
	}

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("pengguna tidak terautentikasi")
	}

	if !CanAccessStudent(current, mahasiswa.OwnerID, s.perms, "student:read:any") {
		return helper.Forbidden("anda tidak memiliki akses ke data mahasiswa ini")
	}

	return helper.Success(c, "mahasiswa ditemukan", mahasiswa)
}

func (s *StudentService) CreateStudent(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("pengguna tidak terautentikasi")
	}

	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}

	if errs := ValidateCreate(req); len(errs) > 0 {
		return helper.Validation(errs)
	}

	baru, err := s.StudentRepo.Create(ctx, model.Student{
		OwnerID:  current.ID,
		NIM:      strings.TrimSpace(req.NIM),
		Name:     strings.TrimSpace(req.Name),
		Grade:    req.Grade,
		IsActive: req.IsActive,
	})
	if err != nil {
		return errorTranslator(err)
	}

	return helper.Created(c, "mahasiswa berhasil dibuat", baru,
		"/api/v1/students/"+strconv.Itoa(baru.ID))
}

func (s *StudentService) ReplaceStudent(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id harus berupa angka positif")
	}

	var req model.UpdateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}

	if errs := ValidateReplace(req); len(errs) > 0 {
		return helper.Validation(errs)
	}

	currentStudent, err := s.StudentRepo.FindByID(ctx, id)
	if err != nil {
		return errorTranslator(err)
	}

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("pengguna tidak terautentikasi")
	}

	if !CanAccessStudent(current, currentStudent.OwnerID, s.perms, "student:update:any") {
		return helper.Forbidden("anda tidak memiliki akses untuk mengubah data mahasiswa ini")
	}

	result, err := s.StudentRepo.Update(ctx, model.Student{
		ID:       id,
		OwnerID:  currentStudent.OwnerID,
		NIM:      strings.TrimSpace(req.NIM),
		Name:     strings.TrimSpace(req.Name),
		Grade:    req.Grade,
		IsActive: req.IsActive,
	})
	if err != nil {
		return errorTranslator(err)
	}

	return helper.Success(c, "mahasiswa berhasil diganti seluruhnya", result)
}

func (s *StudentService) PatchStudent(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id harus berupa angka positif")
	}

	var req model.PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}

	if IsEmptyPatch(req) {
		return helper.BadRequest("tidak ada field yang diubah")
	}

	currentStudent, err := s.StudentRepo.FindByID(ctx, id)
	if err != nil {
		return errorTranslator(err)
	}

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("pengguna tidak terautentikasi")
	}

	if !CanAccessStudent(current, currentStudent.OwnerID, s.perms, "student:update:any") {
		return helper.Forbidden("anda tidak memiliki akses untuk mengubah data mahasiswa ini")
	}

	updatedStudent, errs := ApplyPatch(currentStudent, req)
	if len(errs) > 0 {
		return helper.Validation(errs)
	}

	result, err := s.StudentRepo.Update(ctx, updatedStudent)
	if err != nil {
		return errorTranslator(err)
	}

	return helper.Success(c, "mahasiswa berhasil diperbarui sebagian", result)
}

func (s *StudentService) DeleteStudent(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id harus berupa angka positif")
	}

	if err := s.StudentRepo.Delete(ctx, id); err != nil {
		return errorTranslator(err)
	}

	return helper.NoContent(c)
}

func errorTranslator(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return helper.NotFound("mahasiswa tidak ditemukan")
	case errors.Is(err, repository.ErrDuplicate):
		return helper.Conflict("mahasiswa sudah ada")
	default:
		return helper.Internal(err)
	}
}
