package middlewares

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockResponseWriter struct {
	*httptest.ResponseRecorder
}

func (m *MockResponseWriter) CloseNotify() <-chan bool {
	return make(chan bool)
}

func (m *MockResponseWriter) Flush() {
}

func (m *MockResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}

func (m *MockResponseWriter) Size() int {
	return len(m.Body.Bytes())
}

func (m *MockResponseWriter) Status() int {
	return m.Code
}

func (m *MockResponseWriter) WriteHeaderNow() {
}

func (m *MockResponseWriter) Written() bool {
	return m.Code != 0
}

func (m *MockResponseWriter) WriteString(string) (int, error) {
	return 0, nil
}

func (m *MockResponseWriter) WriteHeader(code int) {
	m.Code = code
}

func (m *MockResponseWriter) Pusher() http.Pusher {
	return nil
}

func TestGinBodyLogMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(GinBodyLogMiddleware)

	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test response"})
	})

	requestBody := `{"test": "data"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expectedResponse := `{"message":"test response"}`
	if !strings.Contains(w.Body.String(), expectedResponse) {
		t.Errorf("Expected response to contain %s, got %s", expectedResponse, w.Body.String())
	}
}

func TestGinBodyLogMiddleware_EmptyBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(GinBodyLogMiddleware)

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	if req.Body == nil {
		req.Body = io.NopCloser(bytes.NewBuffer([]byte("")))
	}

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestBodyLogWriter_Write(t *testing.T) {
	mockWriter := &MockResponseWriter{
		ResponseRecorder: httptest.NewRecorder(),
	}

	blw := &bodyLogWriter{
		ResponseWriter: mockWriter,
		body:           bytes.NewBufferString(""),
	}

	testData := []byte("test response data")

	bytesWritten, err := blw.Write(testData)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if bytesWritten != len(testData) {
		t.Errorf("Expected %d bytes written, got %d", len(testData), bytesWritten)
	}

	if blw.body.String() != string(testData) {
		t.Errorf("Expected body to contain %s, got %s", string(testData), blw.body.String())
	}

	if mockWriter.Body.String() != string(testData) {
		t.Errorf("Expected response writer to contain %s, got %s", string(testData), mockWriter.Body.String())
	}
}

func TestGinBodyLogMiddleware_LargeBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(GinBodyLogMiddleware)

	router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "large body test"})
	})

	largeBody := strings.Repeat("a", 5000)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(largeBody))

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
