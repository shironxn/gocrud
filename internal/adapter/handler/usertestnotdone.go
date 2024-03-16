package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"gocrud/internal/core/domain"
	"gocrud/internal/mocks"
	"gocrud/internal/util"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUserHandler_Register(t *testing.T) {
	body := domain.UserRegisterRequest{
		Name:     "shiron",
		Email:    "shiron@example.com",
		Password: "password123",
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	expectedUser := &domain.User{
		Model:    gorm.Model{ID: 1},
		Name:     "shiron",
		Email:    "shiron@example.com",
		Password: "password123",
	}

	mockUserService := mocks.NewUserService(t)

	validator, err := util.NewValidator()
	if err != nil {
		t.Fatal(err)
	}

	jwt := util.NewJWT(nil)

	t.Run("success", func(t *testing.T) {
		mockUserService.On("Create", mock.AnythingOfType("domain.UserRegisterRequest")).Return(expectedUser, nil)

		app := fiber.New()
		userHandler := NewUserHandler(mockUserService, validator, jwt)
		app.Post("/api/v1/auth/register", userHandler.Register)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBody domain.SuccessResponse
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "successfully create user", responseBody.Message)
		assert.Equal(t, float64(expectedUser.ID), responseBody.Data.(map[string]interface{})["id"].(float64))

		mockUserService.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockUserService.On("Create", mock.AnythingOfType("domain.UserRegisterRequest")).Return(nil, errors.New("failed to create user"))

		app := fiber.New()
		userHandler := NewUserHandler(mockUserService, validator, jwt)
		app.Post("/api/v1/auth/register", userHandler.Register)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var responseBody fiber.Error
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "failed to create user", responseBody.Message)

		mockUserService.AssertExpectations(t)
	})
}
