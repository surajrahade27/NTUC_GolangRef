package http

import (
	"bytes"
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/usecases/mocks"
	"campaign-mgmt/app/domain/valueobjects"
	"campaign-mgmt/app/usecases/dto"
	"campaign-mgmt/app/usecases/params"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCampaignStoreController_validateStoresRequest(t *testing.T) {
	mockCampaignUsecase := mocks.NewCampaignUseCases(t)
	mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)

	campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

	t.Run("Request body validation failure : error occured while decoding", func(t *testing.T) {
		var jsonStr = []byte(`{"stores": [1,}`)

		req, err := http.NewRequest("POST", "/campaigns/1/stores", bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}
		_, err = campaignStoreController.validateStoresRequest(req)
		expectedErr := "invalid character '}' looking for beginning of value"
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})

	t.Run("Request body validation failure", func(t *testing.T) {
		var jsonStr = []byte(`{}`)

		req, err := http.NewRequest("POST", "/campaigns/1/stores", bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}
		_, err = campaignStoreController.validateStoresRequest(req)
		expectedErr := "Key: 'CampaignStoresForm.Stores' Error:Field validation for 'Stores' failed on the 'required' tag"
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})

	t.Run("Request body validation Success", func(t *testing.T) {
		var jsonStr = []byte(`
		{
			"stores": [123, 456]
		}
		`)

		req, err := http.NewRequest("POST", "/campaigns/1/stores", bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}
		_, err = campaignStoreController.validateStoresRequest(req)
		ShouldBeNil(err)
		if err != nil {
			t.Errorf("unexpected error : got - %v ; want - nil", err.Error())
		}
	})
}

func TestCampaignStoreController_addStores(t *testing.T) {
	ctx := context.WithValue(context.Background(), "userId", int64(123456))
	campaignID := int64(101)
	t.Run("failure scenario", func(t *testing.T) {
		request := params.CampaignStoresForm{
			Stores: []int64{123},
		}

		storeEntities := []entities.CampaignStore{
			{
				CampaignID: valueobjects.CampaignID(campaignID),
				StoreID:    123,
				CreatedBy:  123456,
			},
		}

		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignStoreUsecase.On("AddStores", ctx, storeEntities).Return(nil, errors.New("db error"))

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

		_, err := campaignStoreController.addStores(ctx, request, int(campaignID), int64(123456))
		ShouldNotBeNil(err)
		ShouldEqual(err.Error(), "db error")
	})

	t.Run("success scenario", func(t *testing.T) {
		request := params.CampaignStoresForm{
			Stores: []int64{123},
		}

		storeEntities := []entities.CampaignStore{
			{
				CampaignID: valueobjects.CampaignID(campaignID),
				StoreID:    123,
				CreatedBy:  123456,
			},
		}

		expectedResult := []*dto.CampaignStores{
			{
				ID:      1,
				StoreID: 123,
			},
		}

		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignStoreUsecase.On("AddStores", ctx, storeEntities).Return(expectedResult, nil)

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

		response, err := campaignStoreController.addStores(ctx, request, int(campaignID), int64(123456))
		ShouldBeNil(err)
		ShouldEqual(response, expectedResult)
	})
}

func TestCampaignStoreController_AddStores(t *testing.T) {
	t.Run("failure due to incorrect user id", func(t *testing.T) {
		body := bytes.NewBufferString(`{"stores": [123, 456]`)
		req := httptest.NewRequest("POST", "/campaigns/abc/stores", body)
		req.Header.Set("Content-Type", "application/json")

		res := httptest.NewRecorder()
		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

		campaignStoreController.AddStores(res, req)

		if status := res.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
		expected := `{"code":400,"message":"Bad Request : invalid user id"}`
		if a, e := strings.TrimSpace(res.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
		}
	})

	t.Run("failure due to incorrect campaign id", func(t *testing.T) {
		body := bytes.NewBufferString(`{"stores": [123, 456]`)
		req := httptest.NewRequest("POST", "/campaigns/abc/stores", body)
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))

		res := httptest.NewRecorder()
		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

		campaignStoreController.AddStores(res, req)

		if status := res.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		expected := `{"code":400,"message":"Bad Request : incorrect campaign id value, err : strconv.Atoi: parsing \"\": invalid syntax"}`

		if a, e := strings.TrimSpace(res.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
		}
	})

	t.Run("failure due to request validation error", func(t *testing.T) {
		var jsonStr = []byte(`{"stores": ["123", "456"]}`)

		req, _ := http.NewRequest("POST", "/campaigns/1/stores", bytes.NewBuffer(jsonStr))
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("campaign_id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

		campaignStoreController.AddStores(w, req)

		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		expected := `{"code":400,"message":"Bad Request : json: cannot unmarshal string into Go struct field CampaignStoresForm.stores of type int64"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to error occured while checking campaign existence", func(t *testing.T) {
		var jsonStr = []byte(`{"stores": [123, 456]}`)

		req, _ := http.NewRequest("POST", "/campaigns/1/stores", bytes.NewBuffer(jsonStr))
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("campaign_id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

		mockCampaignUsecase.On("Exists", req.Context(), int64(1), "").Return(false, errors.New("db error"))

		w := httptest.NewRecorder()
		campaignStoreController.AddStores(w, req)

		if status := w.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		expected := `{"code":500,"message":"Internal Server Error : db error"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to campaign with given id not exists", func(t *testing.T) {
		var jsonStr = []byte(`{"stores": [123, 456]}`)

		req, _ := http.NewRequest("POST", "/campaigns/1/stores", bytes.NewBuffer(jsonStr))
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("campaign_id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

		mockCampaignUsecase.On("Exists", req.Context(), int64(1), mock.Anything).Return(false, nil)

		w := httptest.NewRecorder()
		campaignStoreController.AddStores(w, req)

		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		expected := `{"code":400,"message":"Bad Request : campaign with id 1 not exists"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to error occured while adding stores", func(t *testing.T) {
		var jsonStr = []byte(`{"stores": [123, 456]}`)

		req, _ := http.NewRequest("POST", "/campaigns/1/stores", bytes.NewBuffer(jsonStr))
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("campaign_id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

		mockCampaignUsecase.On("Exists", req.Context(), int64(1), "").Return(true, nil)

		storeEntities := []entities.CampaignStore{
			{
				StoreID:    123,
				CampaignID: 1,
				CreatedBy:  12345,
			},
			{
				StoreID:    456,
				CampaignID: 1,
				CreatedBy:  12345,
			},
		}
		mockCampaignStoreUsecase.On("AddStores", req.Context(), storeEntities).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		campaignStoreController.AddStores(w, req)

		if status := w.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		expected := `{"code":500,"message":"Internal Server Error : db error"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("Success : stores added successfully", func(t *testing.T) {
		var jsonStr = []byte(`{"stores": [123, 456]}`)

		req, _ := http.NewRequest("POST", "/campaigns/1/stores", bytes.NewBuffer(jsonStr))
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("campaign_id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

		mockCampaignUsecase.On("Exists", req.Context(), int64(1), "").Return(true, nil)

		storeEntities := []entities.CampaignStore{
			{
				StoreID:    123,
				CampaignID: valueobjects.CampaignID(1),
				CreatedBy:  12345,
			},
			{
				StoreID:    456,
				CampaignID: valueobjects.CampaignID(1),
				CreatedBy:  12345,
			},
		}

		response := []*dto.CampaignStores{
			{
				ID:      1,
				StoreID: 123,
			},
			{
				ID:      2,
				StoreID: 456,
			},
		}
		mockCampaignStoreUsecase.On("AddStores", req.Context(), storeEntities).Return(response, nil)

		w := httptest.NewRecorder()
		campaignStoreController.AddStores(w, req)

		if status := w.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := `{"campaign_id":1,"stores":[{"campaign_store_id":1,"store_id":123},{"campaign_store_id":2,"store_id":456}]}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})
}

func TestCampaignStoreController_DeleteStores(t *testing.T) {
	mockCampaignUsecase := mocks.NewCampaignUseCases(t)
	mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)

	campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

	t.Run("failure due to incorrect user id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/campaigns/aaa/stores", nil)
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		campaignStoreController.DeleteStores(res, req)

		if status := res.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		expected := `{"code":400,"message":"Bad Request : invalid user id"}`

		if a, e := strings.TrimSpace(res.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
		}
	})

	t.Run("failure due to incorrect campaign id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/campaigns/aaa/stores", nil)
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		campaignStoreController.DeleteStores(res, req)

		if status := res.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		expected := `{"code":400,"message":"Bad Request : incorrect campaign id value, err : strconv.Atoi: parsing \"\": invalid syntax"}`

		if a, e := strings.TrimSpace(res.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
		}
	})

	t.Run("failure due to error occured while checking campaign existence", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/campaigns/1/stores", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("campaign_id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

		mockCampaignUsecase.On("Exists", req.Context(), int64(1), "").Return(false, errors.New("db error"))

		w := httptest.NewRecorder()
		campaignStoreController.DeleteStores(w, req)

		if status := w.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		expected := `{"code":500,"message":"Internal Server Error : db error"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to campaign with given id not exists", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/campaigns/1/stores", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("campaign_id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

		mockCampaignUsecase.On("Exists", req.Context(), int64(1), mock.Anything).Return(false, nil)

		w := httptest.NewRecorder()
		campaignStoreController.DeleteStores(w, req)

		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		expected := `{"code":400,"message":"Bad Request : campaign with id 1 not exists"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to error occured while deleting stores", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/campaigns/1/stores", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("campaign_id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 123))
		req.Header.Set("Content-Type", "application/json")

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

		mockCampaignUsecase.On("Exists", req.Context(), int64(1), mock.Anything).Return(true, nil)
		mockCampaignStoreUsecase.On("DeleteStores", req.Context(), int64(1), int64(123)).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		campaignStoreController.DeleteStores(w, req)

		if status := w.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		expected := `{"code":500,"message":"Internal Server Error : db error"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("success : deleted all stores for given campaign", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/campaigns/1/stores", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("campaign_id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 123))
		req.Header.Set("Content-Type", "application/json")

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

		mockCampaignUsecase.On("Exists", req.Context(), int64(1), mock.Anything).Return(true, nil)
		mockCampaignStoreUsecase.On("DeleteStores", req.Context(), int64(1), int64(123)).Return(nil)

		w := httptest.NewRecorder()
		campaignStoreController.DeleteStores(w, req)

		if status := w.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := `{"code":200,"message":"all campaign stores with campaign id 1 deleted successfully"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})
}

func TestCampaignStoreController_DeleteStore(t *testing.T) {
	mockCampaignUsecase := mocks.NewCampaignUseCases(t)
	mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)

	campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

	t.Run("failure due to incorrect user id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/campaigns/1/stores/123", nil)
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		campaignStoreController.DeleteStore(res, req)

		if status := res.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		expected := `{"code":400,"message":"Bad Request : invalid user id"}`

		if a, e := strings.TrimSpace(res.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
		}
	})

	t.Run("failure due to incorrect campaign id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/campaigns/aaa/stores/123", nil)
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		campaignStoreController.DeleteStore(res, req)

		if status := res.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		expected := `{"code":400,"message":"Bad Request : incorrect campaign id value, err : strconv.Atoi: parsing \"\": invalid syntax"}`

		if a, e := strings.TrimSpace(res.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
		}
	})

	t.Run("failure due to incorrect store id", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/campaigns/1/stores/aaa", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("campaign_id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		campaignStoreController.DeleteStore(res, req)

		if status := res.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		expected := `{"code":400,"message":"Bad Request : incorrect store id value, err : strconv.Atoi: parsing \"\": invalid syntax"}`

		if a, e := strings.TrimSpace(res.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", res.Body.String(), expected)
		}
	})

	t.Run("failure due to error occured while checking campaign existence", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/campaigns/1/stores/123", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("campaign_id", "1")
		ctx.URLParams.Add("id", "123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

		mockCampaignUsecase.On("Exists", req.Context(), int64(1), "").Return(false, errors.New("db error"))

		w := httptest.NewRecorder()
		campaignStoreController.DeleteStore(w, req)

		if status := w.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		expected := `{"code":500,"message":"Internal Server Error : db error"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to campaign with given id not exists", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/campaigns/1/stores/123", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("campaign_id", "1")
		ctx.URLParams.Add("id", "123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

		mockCampaignUsecase.On("Exists", req.Context(), int64(1), mock.Anything).Return(false, nil)

		w := httptest.NewRecorder()
		campaignStoreController.DeleteStore(w, req)

		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		expected := `{"code":400,"message":"Bad Request : campaign with id 1 not exists"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to error occured while deleting store", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/campaigns/1/stores/987", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("campaign_id", "1")
		ctx.URLParams.Add("id", "987")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 123))
		req.Header.Set("Content-Type", "application/json")

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

		mockCampaignUsecase.On("Exists", req.Context(), int64(1), mock.Anything).Return(true, nil)
		mockCampaignStoreUsecase.On("DeleteStore", req.Context(), int64(1), int64(987), int64(123)).Return(errors.New("dummy error"))

		w := httptest.NewRecorder()
		campaignStoreController.DeleteStore(w, req)

		if status := w.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		expected := `{"code":500,"message":"Internal Server Error : dummy error"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("success : deleted given store for given campaign", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/campaigns/1/stores/987", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("campaign_id", "1")
		ctx.URLParams.Add("id", "987")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 123))
		req.Header.Set("Content-Type", "application/json")

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		campaignStoreController := NewCampaignStoreController(mockCampaignUsecase, mockCampaignStoreUsecase)

		mockCampaignUsecase.On("Exists", req.Context(), int64(1), mock.Anything).Return(true, nil)
		mockCampaignStoreUsecase.On("DeleteStore", req.Context(), int64(1), int64(987), int64(123)).Return(nil)

		w := httptest.NewRecorder()
		campaignStoreController.DeleteStore(w, req)

		if status := w.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expected := `{"code":200,"message":"store with id 987 deleted successfully"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

}

func TestCampaignStoreController_Init(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`{"code": 200,"message": "all campaign stores with campaign id 1 deleted successfully"}`))
	})
	req, _ := http.NewRequest("DELETE", "campaigns/1/stores", nil)
	campaignStoreController := NewCampaignStoreController(nil, nil)
	r := chi.NewRouter()
	campaignStoreController.Init(r)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	body := string(w.Body.String())
	expected := `{"code": 200,"message": "all campaign stores with campaign id 1 deleted successfully"}`
	if body != expected {
		t.Fatalf("expected:%s got:%s", expected, body)
	}
}
