package service

import (
	"api-students/helper"
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
			err := helper.CheckPasswordStrength(tt.password)
			if (err != "") != tt.wantErr {
				t.Errorf("checkPasswordStrength(%q) error = %v, wantErr %v", tt.password, err, tt.wantErr)
			}
		})
	}
}
