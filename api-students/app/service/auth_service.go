package service

import (
	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const refreshTokenBytes = 32

type AuthService struct {
	UserRepo   repository.UserRepository
	TokenRepo  repository.TokenRepository
	JWT        *helper.JWTManager
	RefreshTTL time.Duration
	perms      *helper.PermissionSet
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	jwtManager *helper.JWTManager,
	refreshTTL time.Duration,
	perms *helper.PermissionSet,
) *AuthService {
	return &AuthService{
		UserRepo:   userRepo,
		TokenRepo:  tokenRepo,
		JWT:        jwtManager,
		RefreshTTL: refreshTTL,
		perms:      perms,
	}
}

func (s *AuthService) issueTokenPair(ctx context.Context, user model.User) (model.TokenPair, error) {
	accessToken, err := s.JWT.GenerateAccess(user)
	if err != nil {
		return model.TokenPair{}, err
	}

	rawRefreshToken, err := helper.RandomToken(refreshTokenBytes)
	if err != nil {
		return model.TokenPair{}, err
	}

	tokenHash := helper.SHA256Hex(rawRefreshToken)
	err = s.TokenRepo.Save(ctx, model.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.RefreshTTL),
	})
	if err != nil {
		return model.TokenPair{}, err
	}

	return model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.JWT.AccessTTL().Seconds()),
	}, nil
}

func (s *AuthService) Register(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}

	if errs := ValidateRegister(req); len(errs) > 0 {
		return helper.Validation(errs)
	}

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return helper.Internal(err)
	}

	user, err := s.UserRepo.Create(ctx, model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user", // Strictly hardcoded to prevent mass assignment
		IsActive: true,
	})

	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return helper.Conflict("username atau email sudah terdaftar")
		}
		return helper.Internal(err)
	}

	return helper.Created(c, "registrasi berhasil", user, "/api/v1/users/"+strconv.Itoa(user.ID))
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}

	if errs := ValidateLogin(req); len(errs) > 0 {
		return helper.Validation(errs)
	}

	user, err := s.UserRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			helper.VerifyDummyPassword() // Prevent timing attack / user enumeration
			return helper.Unauthorized("username atau password salah")
		}
		return helper.Internal(err)
	}

	if !user.IsActive {
		return helper.Unauthorized("akun pengguna tidak aktif")
	}

	if !helper.VerifyPassword(user.Password, req.Password) {
		return helper.Unauthorized("username atau password salah")
	}

	pair, err := s.issueTokenPair(ctx, user)
	if err != nil {
		return helper.Internal(err)
	}

	return helper.Success(c, "login berhasil", pair)
}

func (s *AuthService) Refresh(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}

	if req.RefreshToken == "" {
		return helper.BadRequest("refresh token wajib diisi")
	}

	tokenHash := helper.SHA256Hex(req.RefreshToken)
	tokenRecord, err := s.TokenRepo.FindActive(ctx, tokenHash)
	if err != nil {
		return helper.Unauthorized("refresh token tidak valid atau telah kedaluwarsa")
	}

	// Revoke old refresh token (token rotation)
	_ = s.TokenRepo.Revoke(ctx, tokenHash)

	user, err := s.UserRepo.FindByID(ctx, tokenRecord.UserID)
	if err != nil || !user.IsActive {
		return helper.Unauthorized("pengguna tidak valid atau tidak aktif")
	}

	pair, err := s.issueTokenPair(ctx, user)
	if err != nil {
		return helper.Internal(err)
	}

	return helper.Success(c, "token berhasil diperbarui", pair)
}

func (s *AuthService) Logout(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.RefreshRequest
	if err := c.BodyParser(&req); err == nil && req.RefreshToken != "" {
		tokenHash := helper.SHA256Hex(req.RefreshToken)
		_ = s.TokenRepo.Revoke(ctx, tokenHash)
	}

	return helper.Success(c, "logout berhasil", nil)
}

func (s *AuthService) Me(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	authUser, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("unauthorized")
	}

	user, err := s.UserRepo.FindByID(ctx, authUser.ID)
	if err != nil {
		return helper.NotFound("user tidak ditemukan")
	}

	permissions := s.perms.PermissionsOf(user.Role)

	return helper.Success(c, "informasi pengguna berhasil diambil", fiber.Map{
		"user":        user,
		"permissions": permissions,
	})
}

func (s *AuthService) AssignRole(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	targetID, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id user harus berupa angka positif")
	}

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("pengguna tidak terautentikasi")
	}

	var req model.AssignRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}

	if errs := ValidateAssignRole(current, targetID, req, s.perms); len(errs) > 0 {
		return helper.Validation(errs)
	}

	updatedUser, err := s.UserRepo.UpdateRole(ctx, targetID, strings.TrimSpace(req.Role))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.NotFound("user tidak ditemukan")
		}
		return helper.Internal(err)
	}

	return helper.Success(c, "role user berhasil diperbarui", updatedUser)
}

func (s *AuthService) DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id user harus berupa angka positif")
	}

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("pengguna tidak terautentikasi")
	}

	// Punya permission menghapus tidak berarti boleh menghapus dirinya sendiri.
	if current.ID == id {
		return helper.Forbidden("tidak boleh menghapus akun sendiri")
	}

	if err := s.UserRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.NotFound("user tidak ditemukan")
		}
		return helper.Internal(err)
	}

	return helper.NoContent(c)
}
