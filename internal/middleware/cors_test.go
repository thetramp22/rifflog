package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORS_AllowsLocalOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Origin", "http://127.0.0.1:5173")

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = request

	t.Setenv(
		"CORS_ALLOWED_ORIGINS",
		"http://127.0.0.1:5173,https://rifflog.scottstarks.dev",
	)

	middleware := CORS()
	middleware(context)

	allowOrigin := recorder.Header().Get("Access-Control-Allow-Origin")
	expectedAllowOrigin := "http://127.0.0.1:5173"

	if allowOrigin != expectedAllowOrigin {
		t.Errorf("Origins should match: got %v; want %v", allowOrigin, expectedAllowOrigin)
	}
}

func TestCORS_AllowsProductionOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Origin", "https://rifflog.scottstarks.dev")

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = request

	t.Setenv(
		"CORS_ALLOWED_ORIGINS",
		"http://127.0.0.1:5173,https://rifflog.scottstarks.dev",
	)

	middleware := CORS()
	middleware(context)

	allowOrigin := recorder.Header().Get("Access-Control-Allow-Origin")
	expectedAllowOrigin := "https://rifflog.scottstarks.dev"

	if allowOrigin != expectedAllowOrigin {
		t.Errorf("Origins should match: got %v; want %v", allowOrigin, expectedAllowOrigin)
	}
}

func TestCORS_RejectUnknownOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Origin", "http://example.com")

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = request

	t.Setenv(
		"CORS_ALLOWED_ORIGINS",
		"http://127.0.0.1:5173,https://rifflog.scottstarks.dev",
	)

	middleware := CORS()
	middleware(context)

	allowOrigin := recorder.Header().Get("Access-Control-Allow-Origin")
	expectedAllowOrigin := ""

	if allowOrigin != expectedAllowOrigin {
		t.Errorf("Origins should match: got %v; want %v", allowOrigin, expectedAllowOrigin)
	}
}

func TestCORS_AllowsLocalPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)

	request := httptest.NewRequest(http.MethodOptions, "/", nil)
	request.Header.Set("Origin", "http://127.0.0.1:5173")

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = request

	t.Setenv(
		"CORS_ALLOWED_ORIGINS",
		"http://127.0.0.1:5173,https://rifflog.scottstarks.dev",
	)

	middleware := CORS()
	middleware(context)

	allowOrigin := recorder.Header().Get("Access-Control-Allow-Origin")
	expectedAllowOrigin := "http://127.0.0.1:5173"

	status := recorder.Code
	expectedStatus := 204

	allowMethods := recorder.Header().Get("Access-Control-Allow-Methods")
	expectedAllowMethods := "GET, POST, PUT, DELETE, OPTIONS"

	allowHeaders := recorder.Header().Get("Access-Control-Allow-Headers")
	expectedAllowHeaders := "Content-Type, Authorization"

	if allowOrigin != expectedAllowOrigin {
		t.Errorf("Origins should match: got %v; want %v", allowOrigin, expectedAllowOrigin)
	}

	if status != expectedStatus {
		t.Errorf("Status should match: got %v; want %v", status, expectedStatus)
	}

	if allowMethods != expectedAllowMethods {
		t.Errorf("Methods should match: got %v; want %v", allowMethods, expectedAllowMethods)
	}

	if !strings.Contains(allowHeaders, expectedAllowHeaders) {
		t.Errorf("Allowed headers must contain: %v, got %v", expectedAllowHeaders, allowHeaders)
	}
}

func TestCORS_RejectUnknownPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Origin", "http://example.com")
	request.Method = "OPTIONS"

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = request

	t.Setenv(
		"CORS_ALLOWED_ORIGINS",
		"http://127.0.0.1:5173,https://rifflog.scottstarks.dev",
	)

	middleware := CORS()
	middleware(context)

	allowOrigin := recorder.Header().Get("Access-Control-Allow-Origin")
	expectedAllowOrigin := ""

	status := recorder.Code
	expectedStatus := 204

	allowMethods := recorder.Header().Get("Access-Control-Allow-Methods")
	expectedAllowMethods := ""

	if allowOrigin != expectedAllowOrigin {
		t.Errorf("Origins should match: got %v; want %v", allowOrigin, expectedAllowOrigin)
	}

	if status != expectedStatus {
		t.Errorf("Status should match: got %v; want %v", status, expectedStatus)
	}

	if allowMethods != expectedAllowMethods {
		t.Errorf("Methods should match: got %v; want %v", allowMethods, expectedAllowMethods)
	}
}
