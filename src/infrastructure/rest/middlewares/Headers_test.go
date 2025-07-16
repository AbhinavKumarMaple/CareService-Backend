package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCommonHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(CommonHeaders)

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	router.ServeHTTP(w, req)

	headers := w.Header()

	expectedHeaders := map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Methods":     "POST, OPTIONS, DELETE, GET, PUT",
		"X-Frame-Options":                  "SAMEORIGIN",
		"Cache-Control":                    "no-cache, no-store",
		"Pragma":                           "no-cache",
		"Expires":                          "0",
	}

	for key, expectedValue := range expectedHeaders {
		actualValue := headers.Get(key)
		if actualValue != expectedValue {
			t.Errorf("Header %s: expected %s, got %s", key, expectedValue, actualValue)
		}
	}

	allowHeaders := headers.Get("Access-Control-Allow-Headers")
	expectedAllowHeaders := "Content-Type, Depth, User-Agent, X-File-Size, X-Requested-With, If-Modified-Since, X-File-CompanyName, Cache-Control"
	if allowHeaders != expectedAllowHeaders {
		t.Errorf("Access-Control-Allow-Headers: expected %s, got %s", expectedAllowHeaders, allowHeaders)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
