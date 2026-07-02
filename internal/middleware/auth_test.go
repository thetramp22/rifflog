package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/thetramp22/rifflog/internal/auth"
	"github.com/thetramp22/rifflog/internal/config"
)

type testApp struct {
	router         *gin.Engine
	jwtService     *auth.JWTService
	authMiddleware *AuthMiddleware
}

type errorResponse struct {
	Error string `json:"error"`
}

func testSetup() *testApp {
	router := gin.Default()
	jwtService := auth.NewJWTService(config.JWTSecret())
	authMiddleware := NewAuthMiddleware(jwtService)

	protected := router.Group("/api", authMiddleware.Authenticate)
	{
		protected.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

	}

	return &testApp{
		router:         router,
		jwtService:     jwtService,
		authMiddleware: authMiddleware,
	}
}

func TestAuthNoHeader(t *testing.T) {
	app := testSetup()

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()

	app.router.ServeHTTP(w, req)

	asserErrorResponse(t, w, http.StatusUnauthorized, "authorization header is missing")
}

func TestAuthInvalidHeaderFormat(t *testing.T) {
	app := testSetup()

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Authorization", "potato")
	w := httptest.NewRecorder()

	app.router.ServeHTTP(w, req)

	asserErrorResponse(t, w, http.StatusUnauthorized, "authorization header is malformed")
}

func TestAuthValidJWT(t *testing.T) {
	app := testSetup()

	userID := 22
	token, err := app.jwtService.GenerateToken(userID)
	if err != nil {
		t.Fatalf("Error generating token: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	w := httptest.NewRecorder()

	app.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("expected 200, got %v, body=%s", status, w.Body.String())
	}
}

func TestAuthInvalidJWT(t *testing.T) {
	app := testSetup()
	service2 := auth.NewJWTService("secret-two")

	userID := 22
	invalidToken, err := service2.GenerateToken(userID)
	if err != nil {
		t.Fatalf("Error generating token: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", invalidToken))
	w := httptest.NewRecorder()

	app.router.ServeHTTP(w, req)

	asserErrorResponse(t, w, http.StatusUnauthorized, "invalid token")
}

func asserErrorResponse(t *testing.T, w *httptest.ResponseRecorder, wantStatus int, wantErr string) {
	if status := w.Code; status != wantStatus {
		t.Fatalf("expected 401, got %v, body=%s", status, w.Body.String())
	}

	var got errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if diff := cmp.Diff(got.Error, wantErr); diff != "" {
		t.Errorf("values mismatch (-want +got):\n%s", diff)
	}
}
