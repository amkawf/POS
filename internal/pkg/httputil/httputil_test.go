package httputil_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	kitchendomain "pos-backend/internal/kitchen/domain"
	orderdomain "pos-backend/internal/order/domain"
	"pos-backend/internal/pkg/httputil"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestResponseHelpers(t *testing.T) {
	t.Run("BadRequest", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		httputil.BadRequest(c, "TEST_BAD_REQUEST", "invalid payload")

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", w.Code)
		}

		var resp httputil.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal failed: %v", err)
		}
		if resp.Error.Code != "TEST_BAD_REQUEST" || resp.Error.Message != "invalid payload" {
			t.Fatalf("unexpected body: %+v", resp)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		httputil.NotFound(c, "NOT_FOUND", "entity missing")

		if w.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d", w.Code)
		}
	})

	t.Run("HandleError mapping", func(t *testing.T) {
		cases := []struct {
			err          error
			expectedCode int
		}{
			{kitchendomain.ErrTicketNotFound, http.StatusNotFound},
			{errors.New("item not found"), http.StatusNotFound},
			{orderdomain.ErrOrderAlreadyCompleted, http.StatusConflict},
			{orderdomain.ErrInvalidOrderTransition, http.StatusConflict},
			{orderdomain.ErrInvalidOrderType, http.StatusBadRequest},
			{orderdomain.ErrEmptyOrder, http.StatusBadRequest},
			{errors.New("unexpected database error"), http.StatusInternalServerError},
		}

		for _, tc := range cases {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			httputil.HandleError(c, tc.err)
			if w.Code != tc.expectedCode {
				t.Errorf("error %v: expected status %d, got %d", tc.err, tc.expectedCode, w.Code)
			}
		}
	})
}

func TestParamParsing(t *testing.T) {
	testUUID := uuid.New()

	t.Run("ParseUUIDParam", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: testUUID.String()}}

		val, ok := httputil.ParseUUIDParam(c, "id")
		if !ok || val != testUUID {
			t.Fatalf("expected %s, got %s", testUUID, val)
		}

		// Invalid
		w2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(w2)
		c2.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}
		_, ok2 := httputil.ParseUUIDParam(c2, "id")
		if ok2 || w2.Code != http.StatusBadRequest {
			t.Fatalf("expected invalid param to fail with 400")
		}
	})

	t.Run("ParseIntQuery", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		req, _ := http.NewRequest("GET", "/test?limit=25", nil)
		c.Request = req

		val := httputil.ParseIntQuery(c, "limit", 10)
		if val != 25 {
			t.Fatalf("expected 25, got %d", val)
		}

		valFallback := httputil.ParseIntQuery(c, "page", 10)
		if valFallback != 10 {
			t.Fatalf("expected fallback 10, got %d", valFallback)
		}
	})
}
