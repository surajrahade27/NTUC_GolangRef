package middlewares

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
)

type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func Test_validateToken(t *testing.T) {
	t.Run("failure : authorization header is missing", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/campaigns/1", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req.Header.Set("Content-Type", "application/json")
		isValid, context := validateToken(req)
		if isValid != false {
			t.Errorf("unexpected response : got - %v ; want - %v", isValid, true)
		}
		if context != nil {
			t.Errorf("unexpected response : got - %v ; want - %v", context, nil)
		}

	})

	t.Run("legacy token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/campaigns/1", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer abcd")

		isValid, context := validateToken(req)
		if isValid != true {
			t.Errorf("unexpected response : got - %v ; want - %v", isValid, true)
		}
		if context != nil {
			t.Errorf("unexpected response : got - %v ; want - %v", context, nil)
		}
	})
}

func Test_OktaAuthenticator(t *testing.T) {
	t.Run("Authentication Failed", func(t *testing.T) {
		body := bytes.NewBufferString(`{"stores": [123, 456]`)
		req := httptest.NewRequest("POST", "/campaigns/abc/stores", body)
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler := &mockHandler{}
		expectedResponse := `{"code":401,"message":"Not Authorised"}`
		OktaAuthenticator(handler).ServeHTTP(res, req)
		if status := res.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}
		if a, e := strings.TrimSpace(res.Body.String()), strings.TrimSpace(expectedResponse); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expectedResponse)
		}
	})

	t.Run("Skip authentication for OPTIONS request", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/test", nil)
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler := &mockHandler{}
		OktaAuthenticator(handler).ServeHTTP(res, req)
		if status := res.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		if a := strings.TrimSpace(res.Body.String()); a != "" {
			t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), nil)
		}
	})

	t.Run("Successful authentication with legacy auth", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/campaigns", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer abcd")
		res := httptest.NewRecorder()

		handler := &mockHandler{}
		handler.On("ServeHTTP", res, req)
		OktaAuthenticator(handler).ServeHTTP(res, req)
		if status := res.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		if a := strings.TrimSpace(res.Body.String()); a != "" {
			t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), nil)
		}
	})

}
