package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestHealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	Register("/health", router)

	req, _ := http.NewRequest("GET", "/health/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, w.Code, http.StatusOK)
	require.Equal(t, w.Body.String(), "ok")
}
