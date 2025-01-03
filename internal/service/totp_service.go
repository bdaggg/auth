package service

import (
	"crypto/rand"
	"encoding/base32"

	"github.com/pquerna/otp/totp"
)

type TOTPService struct {
	issuer string
}

func NewTOTPService(issuer string) *TOTPService {
	return &TOTPService{
		issuer: issuer,
	}
}

func (s *TOTPService) GenerateSecret() (string, error) {
	secret := make([]byte, 20)
	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(secret), nil
}

func (s *TOTPService) GenerateQRCode(email, secret string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.issuer,
		AccountName: email,
		Secret:      []byte(secret),
	})
	if err != nil {
		return "", err
	}
	return key.URL(), nil
}

func (s *TOTPService) ValidateCode(secret, code string) bool {
	return totp.Validate(code, secret)
}
