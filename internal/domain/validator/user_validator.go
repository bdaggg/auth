package validator

import (
	"errors"
	"regexp"
	"unicode"
)

var (
	ErrInvalidEmail     = errors.New("geçersiz email formatı")
	ErrPasswordTooShort = errors.New("şifre en az 8 karakter olmalıdır")
	ErrPasswordTooWeak  = errors.New("şifre en az bir büyük harf, bir küçük harf ve bir rakam içermelidir")
	ErrPasswordHasSpace = errors.New("şifre boşluk içeremez")
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	var (
		hasUpper  bool
		hasLower  bool
		hasNumber bool
		hasSpace  bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsSpace(char):
			hasSpace = true
		}
	}

	if hasSpace {
		return ErrPasswordHasSpace
	}

	if !hasUpper || !hasLower || !hasNumber {
		return ErrPasswordTooWeak
	}

	return nil
}
