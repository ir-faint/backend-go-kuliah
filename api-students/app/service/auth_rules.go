package service

import (
	"api-students/app/model"
	"api-students/helper"
	"net/mail"
	"regexp"
	"strings"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,50}$`)
	letterRegex   = regexp.MustCompile(`[a-zA-Z]`)
	digitRegex    = regexp.MustCompile(`[0-9]`)

	weakPasswords = map[string]bool{
		"password":    true,
		"12345678":    true,
		"qwerty123":   true,
		"password123": true,
		"admin123":    true,
	}
)

func isValidUsername(username string) bool {
	return usernameRegex.MatchString(username)
}

func ValidateRegister(req model.RegisterRequest) map[string]string {
	errs := make(map[string]string)

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" {
		errs["username"] = "username wajib diisi"
	} else if !isValidUsername(req.Username) {
		errs["username"] = "username hanya boleh berisi huruf, angka, underscore, strip (3-50 karakter)"
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		errs["email"] = "email wajib diisi"
	} else if _, err := mail.ParseAddress(req.Email); err != nil {
		errs["email"] = "format email tidak valid"
	}

	if req.Password == "" {
		errs["password"] = "password wajib diisi"
	} else if err := helper.CheckPasswordStrength(req.Password); err != "" {
		errs["password"] = err
	}

	return errs
}

func ValidateLogin(req model.LoginRequest) map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(req.Username) == "" {
		errs["username"] = "username wajib diisi"
	}

	if req.Password == "" {
		errs["password"] = "password wajib diisi"
	}

	return errs
}
