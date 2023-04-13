package dto

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSuccessJSON(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/campaigns/1", nil)
	expectedOutput := `{"code":200,"message":"successful request"}`

	SuccessJSON(w, r, "successful request")

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expectedOutput); a != e {
		t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expectedOutput)
	}
}

func TestSuccessJSONResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/campaigns/1", nil)
	expectedOutput := `{"id":101}`

	SuccessJSONResponse(w, r, map[string]interface{}{"id": 101})

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expectedOutput); a != e {
		t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expectedOutput)
	}
}

func TestBadRequestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/campaigns/1", nil)
	expectedOutput := `{"code":400,"message":"Bad Request : invalid id"}`

	BadRequestJSON(w, r, "invalid id")

	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
	if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expectedOutput); a != e {
		t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expectedOutput)
	}
}

func TestConflictErrorJSON(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/campaigns/1", nil)
	expectedOutput := `{"code":409,"message":"Conflict Error : entity already exists"}`

	ConflictErrorJSON(w, r, "entity already exists")

	if status := w.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
	}
	if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expectedOutput); a != e {
		t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expectedOutput)
	}
}

func TestInternalServerErrorJSON(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/campaigns/1", nil)
	expectedOutput := `{"code":500,"message":"Internal Server Error : dummy error"}`

	InternalServerErrorJSON(w, r, "dummy error")

	if status := w.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
	if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expectedOutput); a != e {
		t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expectedOutput)
	}
}

func TestNotFoundJSON(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/campaigns/1", nil)
	expectedOutput := `{"code":404,"message":"Entity Not Found : dummy error"}`

	NotFoundJSON(w, r, "dummy error")

	if status := w.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
	if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expectedOutput); a != e {
		t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expectedOutput)
	}
}
