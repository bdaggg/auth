package handlers

import (
	"auth-service/internal/domain/entity"
	"auth-service/internal/domain/oauth"
	"auth-service/internal/service"

	"github.com/gofiber/fiber/v2"
)

func Register(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input service.RegisterInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Geçersiz istek formatı",
			})
		}

		user, err := authService.Register(c.Context(), input)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(user)
	}
}

func Login(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input service.LoginInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Geçersiz istek formatı",
			})
		}

		tokens, err := authService.Login(c.Context(), input)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(tokens)
	}
}

func RefreshToken(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Query("token")
		if token == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Token gerekli",
			})
		}

		tokens, err := authService.RefreshToken(c.Context(), token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(tokens)
	}
}

func ForgotPassword(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			Email string `json:"email"`
		}
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Geçersiz istek formatı",
			})
		}
		return authService.InitiatePasswordReset(c.Context(), input.Email)
	}
}

func ResetPassword(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			Token       string `json:"token"`
			NewPassword string `json:"new_password"`
		}
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Geçersiz istek formatı",
			})
		}
		return authService.ResetPassword(c.Context(), input.Token, input.NewPassword)
	}
}

func VerifyEmail(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Query("token")
		if token == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Token gerekli",
			})
		}
		return authService.VerifyEmail(c.Context(), token)
	}
}

func GoogleLogin(provider oauth.Provider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		url := provider.GetAuthURL()
		return c.Redirect(url)
	}
}

func GoogleCallback(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		code := c.Query("code")
		if code == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Authorization code required",
			})
		}
		tokens, err := authService.HandleOAuthCallback(c.Context(), code)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(tokens)
	}
}

func ListUsers(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		offset := c.QueryInt("offset", 0)
		limit := c.QueryInt("limit", 10)

		users, err := authService.ListUsers(c.Context(), offset, limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(users)
	}
}

func ChangeUserRole(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Params("id")
		var input struct {
			Role entity.Role `json:"role"`
		}
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Geçersiz istek formatı",
			})
		}
		return authService.ChangeUserRole(c.Context(), userID, input.Role)
	}
}
