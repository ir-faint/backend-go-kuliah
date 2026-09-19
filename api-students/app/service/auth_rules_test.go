package service

import (
	"api-students/app/model"
	"testing"
)

func TestCheckPasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Valid strong password", "Secret123!", false},
		{"Too short", "Sec1!", true},
		{"No digits", "SecretPassword", true},
		{"No letters", "123456789", true},
		{"Weak common password", "password123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkPasswordStrength(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkPasswordStrength(%q) error = %v, wantErr %v", tt.password, err, tt.wantErr)
			}
		})
	}
}

func TestIsValidUsername(t *testing.T) {
	tests := []struct {
		username string
		want     bool
	}{
		{"john_doe", true},
		{"user123", true},
		{"ab", false},          // too short
		{"user@domain", false}, // invalid char @
		{"valid-name", true},
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			got := isValidUsername(tt.username)
			if got != tt.want {
				t.Errorf("isValidUsername(%q) = %v, want %v", tt.username, got, tt.want)
			}
		})
	}
}

func TestValidateRegister(t *testing.T) {
	t.Run("Valid registration payload", func(t *testing.T) {
		req := model.RegisterRequest{
			Username: "validuser",
			Email:    "user@example.com",
			Password: "StrongPassword123!",
		}
		errs := ValidateRegister(req)
		if len(errs) != 0 {
			t.Errorf("ValidateRegister() expected no errors, got %v", errs)
		}
	})

	t.Run("Invalid registration payload with missing fields and weak password", func(t *testing.T) {
		req := model.RegisterRequest{
			Username: "u",
			Email:    "invalid-email",
			Password: "123",
		}
		errs := ValidateRegister(req)
		if len(errs) != 3 {
			t.Errorf("ValidateRegister() expected 3 errors, got %v", len(errs))
		}
		if _, ok := errs["username"]; !ok {
			t.Errorf("expected username error")
		}
		if _, ok := errs["email"]; !ok {
			t.Errorf("expected email error")
		}
		if _, ok := errs["password"]; !ok {
			t.Errorf("expected password error")
		}
	})
}
