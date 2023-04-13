package http

import (
	"bytes"
	"campaign-mgmt/app/domain/entities"
	service_mocks "campaign-mgmt/app/domain/services/mocks"
	"campaign-mgmt/app/domain/usecases/mocks"
	"campaign-mgmt/app/domain/valueobjects"
	"campaign-mgmt/app/usecases/dto"
	"campaign-mgmt/app/usecases/params"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestCampaignController_Init(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`{"code": 200,"status": "SUCCESS"}`))
	})
	req, _ := http.NewRequest("GET", "/", nil)
	campaignController := NewCampaignController(nil, nil, nil, nil, nil)
	r := chi.NewRouter()
	campaignController.Init(r)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	body := string(w.Body.String())
	expected := `{"code": 200,"status": "SUCCESS"}`
	if body != expected {
		t.Fatalf("expected:%s got:%s", expected, body)
	}
}

func TestCampaignController_validateCampaignRequest(t *testing.T) {
	appConfig := entities.AppCfg{
		ValidationParam: entities.ValidationParam{
			MaxLeadTime:       20,
			MaxDateDifference: 28,
		},
	}

	mockCampaignUsecase := mocks.NewCampaignUseCases(t)
	mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
	mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
	mockTransactionService := service_mocks.NewTransactionService(t)
	campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
		mockTransactionService, &appConfig)

	t.Run("Request body validation failure : error occured while decoding", func(t *testing.T) {
		var jsonStr = []byte(`{"campaign_status_code": 1,
		"campaign_type": "deli}`)

		req, err := http.NewRequest("POST", "/campaigns", bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}
		_, err = campaignController.validateCampaignRequest(req)
		expectedErr := "unexpected EOF"
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})

	t.Run("Request body validation failure : missing required params", func(t *testing.T) {
		var jsonStr = []byte(`{
			"campaign_status_code": 1,
			"campaign_type": "deli",
			"collection_end_date": "2023-04-05 12:00:00",
			"collection_start_date": "2023-03-05 12:00:00",
			"landing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			"lead_time": 15,
			"listing_description": "test description",
			"listing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			"listing_title": "test screen title",
			"onboarding_description": "test desc",
			"onboarding_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			"onboarding_title": "test campaign",
			"order_end_date": "2023-03-31 12:00:00",
			"order_start_date": "2023-03-01 12:00:00",
			"stores": [
			  83,84,85
			],
			"offer_id": 123,
			"tag_id":456
		  }`)

		req, err := http.NewRequest("POST", "/campaigns", bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}
		_, err = campaignController.validateCampaignRequest(req)
		expectedErr := "Key: 'CampaignCreationForm.Title' Error:Field validation for 'Title' failed on the 'required' tag"
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})


	t.Run("Request body validation failure : invalid dates", func(t *testing.T) {
		var jsonStr = []byte(`{
			"campaign_status_code": 1,
			"campaign_type": "deli",
			"collection_end_date": "2023-04-05 12:00:00",
			"collection_start_date": "2023-03-01 12:00:00",
			"landing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			"lead_time": 20,
			"listing_description": "test description",
			"listing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			"listing_title": "test screen title",
			"onboarding_description": "test desc",
			"onboarding_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			"onboarding_title": "test campaign",
			"order_end_date": "2023-03-31 12:00:00",
			"order_start_date": "2023-03-01 12:00:00",
			"stores": [
			  83,84,85
			],
			"title": "new campaign",
			"offer_id": 123,
			"tag_id":456
		  }`)

		req, err := http.NewRequest("POST", "/campaigns", bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}
		_, err = campaignController.validateCampaignRequest(req)
		expectedErr := "invalid date : Collection start date should be at least 20 days greater than order start date"
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})

	t.Run("Request body validation Success", func(t *testing.T) {
		var jsonStr = []byte(`{
			"campaign_status_code": 1,
			"campaign_type": "deli",
			"collection_end_date": "2023-04-05 12:00:00",
			"collection_start_date": "2023-03-05 12:00:00",
			"landing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			"lead_time": 3,
			"listing_description": "test description",
			"listing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			"listing_title": "test screen title",
			"onboarding_description": "test desc",
			"onboarding_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			"onboarding_title": "test campaign",
			"order_end_date": "2023-03-31 12:00:00",
			"order_start_date": "2023-03-01 12:00:00",
			"stores": [
			  83,84,85
			],
			"title": "new campaign",
			"offer_id": 123,
			"tag_id": 456
		  }`)

		req, err := http.NewRequest("POST", "/campaigns", bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}
		_, err = campaignController.validateCampaignRequest(req)
		ShouldBeNil(err)
		if err != nil {
			t.Errorf("unexpected error : got - %v ; want - nil", err.Error())
		}
	})
}

func TestCampaignController_validateUpdateCampaignRequest(t *testing.T) {
	appConfig := entities.AppCfg{
		ValidationParam: entities.ValidationParam{
			MaxLeadTime:       20,
			MaxDateDifference: 28,
		},
	}

	mockCampaignUsecase := mocks.NewCampaignUseCases(t)
	mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
	mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
	mockTransactionService := service_mocks.NewTransactionService(t)
	campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
		mockTransactionService, &appConfig)

	t.Run("Request body validation failure : error occured while decoding", func(t *testing.T) {
		var jsonStr = []byte(`{"campaign_status_code": 1,
		"campaign_type": "deli}`)

		req, err := http.NewRequest("PUT", "/campaigns/1", bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}
		_, err = campaignController.validateUpdateCampaignRequest(req)
		expectedErr := "unexpected EOF"
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})

	t.Run("Request body validation failure : key validation error", func(t *testing.T) {
		var jsonStr = []byte(`{
			"campaign_status_code": 1,
			"campaign_type": "test",
			"collection_end_date": "2023-04-05 12:00:00",
			"collection_start_date": "2023-03-05 12:00:00",
			"landing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			"lead_time": 15,
			"listing_description": "test description",
			"listing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			"listing_title": "test screen title",
			"onboarding_description": "test desc",
			"onboarding_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			"onboarding_title": "test campaign",
			"order_end_date": "2023-03-31 12:00:00",
			"order_start_date": "2023-03-01 12:00:00"
		  }`)

		req, err := http.NewRequest("PUT", "/campaigns/1", bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}
		_, err = campaignController.validateUpdateCampaignRequest(req)
		expectedErr := "Key: 'CampaignUpdateForm.CampaignType' Error:Field validation for 'CampaignType' failed on the 'oneof' tag"
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})

	t.Run("Request body validation failure : invalid dates", func(t *testing.T) {
		var jsonStr = []byte(`{
			"campaign_status_code": 1,
			"campaign_type": "deli",
			"collection_end_date": "2023-04-05 12:00:00",
			"collection_start_date": "2023-03-01 12:00:00",
			"landing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			"lead_time": 3,
			"listing_description": "test description",
			"listing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			"listing_title": "test screen title",
			"onboarding_description": "test desc",
			"onboarding_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			"onboarding_title": "test campaign",
			"order_end_date": "2023-03-31 12:00:00",
			"order_start_date": "2023-03-01 12:00:00",
			"title": "new campaign"
		  }`)

		req, err := http.NewRequest("PUT", "/campaigns/1", bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}
		_, err = campaignController.validateUpdateCampaignRequest(req)
		expectedErr := "invalid date : Collection start date should be at least 3 days greater than order start date"
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})

	t.Run("Request body validation Success", func(t *testing.T) {
		var jsonStr = []byte(`{
			"campaign_status_code": 1,
			"campaign_type": "deli",
			"collection_end_date": "2023-04-05 12:00:00",
			"collection_start_date": "2023-03-05 12:00:00",
			"landing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			"lead_time": 3,
			"listing_description": "test description",
			"listing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			"listing_title": "test screen title",
			"onboarding_description": "test desc",
			"onboarding_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			"onboarding_title": "test campaign",
			"order_end_date": "2023-03-31 12:00:00",
			"order_start_date": "2023-03-01 12:00:00",
			"offer_id": 123,
			"tag_id":456,
			"title": "new campaign"
		  }`)

		req, err := http.NewRequest("PUT", "/campaigns/1", bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}
		_, err = campaignController.validateUpdateCampaignRequest(req)
		ShouldBeNil(err)
		if err != nil {
			t.Errorf("unexpected error : got - %v ; want - nil", err.Error())
		}
	})
}

// Get Campaign tests
func TestCampaignController_getCampaign(t *testing.T) {
	appConfig := entities.AppCfg{
		ValidationParam: entities.ValidationParam{
			MaxLeadTime: 20,
		},
	}

	mockCampaignUsecase := mocks.NewCampaignUseCases(t)
	mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
	mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
	mockTransactionService := service_mocks.NewTransactionService(t)
	campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
		mockTransactionService, &appConfig)

	t.Run("Get Campaign request success", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/campaigns/1", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req.Header.Set("Content-Type", "application/json")
		query := req.URL.Query()
		query.Add("omit_stores", "true")
		query.Add("omit_products", "true")
		req.URL.RawQuery = query.Encode()

		res := httptest.NewRecorder()
		response := dto.CampaignDTO{}
		mockCampaignUsecase.On("Get", req.Context(), int64(1)).Return(&response, nil)
		campaignController.GetCampaign(res, req)

		ShouldBeNil(err)
		if err != nil {
			t.Fatal(err)
		}

	})

	t.Run("Get Campaign request success when omit_stores=false and omit_products=false", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/campaigns/1", nil)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req.Header.Set("Content-Type", "application/json")
		query := req.URL.Query()
		query.Add("omit_stores", "false")
		query.Add("omit_products", "false")
		req.URL.RawQuery = query.Encode()

		res := httptest.NewRecorder()
		response := dto.CampaignDTO{}
		mockCampaignUsecase.On("Get", req.Context(), int64(1)).Return(&response, nil)
		var campaignStores []*dto.CampaignStores
		mockCampaignStoreUsecase.On("GetStores", req.Context(), int64(1)).Return(campaignStores, nil)
		var campaignProducts []*dto.CampaignProducts
		mockCampaignProductUsecase.On("GetProducts", req.Context(), int64(1)).Return(campaignProducts, nil)
		campaignController.GetCampaign(res, req)

		ShouldBeNil(err)
		if err != nil {
			t.Fatal(err)
		}

	})

	t.Run("Get Campaign request error when omit_stores=false and omit_products=false", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/campaigns/1", nil)

		query := req.URL.Query()
		query.Add("omit_stores", "false")
		query.Add("omit_products", "false")
		req.URL.RawQuery = query.Encode()

		res := httptest.NewRecorder()
		response := dto.CampaignDTO{}
		mockCampaignUsecase.On("Get", req.Context(), int64(0)).Return(&response, nil)
		var campaignStores []*dto.CampaignStores
		mockCampaignStoreUsecase.On("GetStores", req.Context(), int64(0)).Return(campaignStores, nil)
		var campaignProducts []*dto.CampaignProducts
		mockCampaignProductUsecase.On("GetProducts", req.Context(), int64(0)).Return(campaignProducts, nil)
		campaignController.GetCampaign(res, req)

		ShouldBeNil(err)
		if err != nil {
			t.Fatal(err)
		}

	})

	t.Run("Get Campaign request error when omit_stores=false and omit_products=false are nil", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/campaigns/1", nil)

		res := httptest.NewRecorder()
		response := dto.CampaignDTO{}
		mockCampaignUsecase.On("Get", req.Context(), int64(0)).Return(&response, nil)
		var campaignStores []*dto.CampaignStores
		mockCampaignStoreUsecase.On("GetStores", req.Context(), int64(0)).Return(campaignStores, nil)
		var campaignProducts []*dto.CampaignProducts
		mockCampaignProductUsecase.On("GetProducts", req.Context(), int64(0)).Return(campaignProducts, nil)
		campaignController.GetCampaign(res, req)

		ShouldBeNil(err)
		if err != nil {
			t.Fatal(err)
		}

	})

	t.Run("Get Campaign request error when omit_stores and omit_products are invalid inputs", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/campaigns/1", nil)

		query := req.URL.Query()
		query.Add("omit_stores", "invalid_input")
		query.Add("omit_products", "invalid_input")
		req.URL.RawQuery = query.Encode()

		res := httptest.NewRecorder()
		response := dto.CampaignDTO{}
		mockCampaignUsecase.On("Get", req.Context(), int64(0)).Return(&response, nil)
		var campaignStores []*dto.CampaignStores
		mockCampaignStoreUsecase.On("GetStores", req.Context(), int64(0)).Return(campaignStores, nil)
		var campaignProducts []*dto.CampaignProducts
		mockCampaignProductUsecase.On("GetProducts", req.Context(), int64(0)).Return(campaignProducts, nil)
		campaignController.GetCampaign(res, req)

		ShouldBeNil(err)
		if err != nil {
			t.Fatal(err)
		}

	})
}

func TestCampaignController_CreateCampaign(t *testing.T) {
	appConfig := entities.AppCfg{
		ValidationParam: entities.ValidationParam{
			MaxLeadTime:       20,
			MaxDateDifference: 28,
		},
	}

	t.Run("failure due to invalid user id error", func(t *testing.T) {
		var jsonStr = []byte(`{"campaign_status_code": 1,
		"campaign_type": "deli}`)

		req, _ := http.NewRequest("POST", "/campaigns", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		campaignController.CreateCampaign(w, req)

		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		expected := `{"code":400,"message":"Bad Request : invalid user id"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to request validation error", func(t *testing.T) {
		var jsonStr = []byte(`{"campaign_status_code": 1,
		"campaign_type": "deli}`)

		req, _ := http.NewRequest("POST", "/campaigns", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))

		w := httptest.NewRecorder()

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		campaignController.CreateCampaign(w, req)

		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		expected := `{"code":400,"message":"Bad Request : unexpected EOF"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to error occured while checking campaign existence", func(t *testing.T) {
		campaignRequest := bytes.NewBuffer([]byte(`{
			"campaign_status_code": 1,
			"campaign_type": "deli",
			"collection_end_date": "2023-04-05 12:00:00",
			"collection_start_date": "2023-03-05 12:00:00",
			"landing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			"lead_time": 3,
			"listing_description": "test description",
			"listing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			"listing_title": "test screen title",
			"onboarding_description": "test desc",
			"onboarding_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			"onboarding_title": "test campaign",
			"order_end_date": "2023-03-31 12:00:00",
			"order_start_date": "2023-03-01 12:00:00",
			"stores": [
			  83,84,85
			],
			"title": "new campaign",
			"offer_id": 123,
			"tag_id": 456
		  }`))
		req, _ := http.NewRequest("POST", "/campaigns", campaignRequest)
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		w := httptest.NewRecorder()

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)

		mockCampaignUsecase.On("Exists", req.Context(), int64(0), "new campaign").Return(false, errors.New("db error"))

		campaignController.CreateCampaign(w, req)

		if status := w.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		expected := `{"code":500,"message":"Internal Server Error : db error"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to campaign with given name already exists", func(t *testing.T) {
		campaignRequest := bytes.NewBuffer([]byte(`{
			"campaign_status_code": 1,
			"campaign_type": "deli",
			"collection_end_date": "2023-04-05 12:00:00",
			"collection_start_date": "2023-03-05 12:00:00",
			"landing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			"lead_time": 3,
			"listing_description": "test description",
			"listing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			"listing_title": "test screen title",
			"onboarding_description": "test desc",
			"onboarding_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			"onboarding_title": "test campaign",
			"order_end_date": "2023-03-31 12:00:00",
			"order_start_date": "2023-03-01 12:00:00",
			"stores": [
			  83,84,85
			],
			"title": "new campaign",
			"offer_id": 123,
			"tag_id": 456
		  }`))
		req, _ := http.NewRequest("POST", "/campaigns", campaignRequest)
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		w := httptest.NewRecorder()

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)

		mockCampaignUsecase.On("Exists", req.Context(), int64(0), "new campaign").Return(true, nil)

		campaignController.CreateCampaign(w, req)

		if status := w.Code; status != http.StatusConflict {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
		}

		expected := `{"code":409,"message":"Conflict Error : campaign with given name already exists"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to error occured while creating campaign", func(t *testing.T) {
		campaignRequest := bytes.NewBuffer([]byte(`{
			"campaign_status_code": 1,
			"campaign_type": "deli",
			"collection_end_date": "2023-04-05 12:00:00",
			"collection_start_date": "2023-03-05 12:00:00",
			"landing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			"lead_time": 3,
			"listing_description": "test description",
			"listing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			"listing_title": "test screen title",
			"onboarding_description": "test desc",
			"onboarding_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			"onboarding_title": "test campaign",
			"order_end_date": "2023-03-31 12:00:00",
			"order_start_date": "2023-03-01 12:00:00",
			"stores": [
			  83,84,85
			],
			"title": "new campaign",
			"offer_id": 123,
			"tag_id": 456
		  }`))
		req, _ := http.NewRequest("POST", "/campaigns", campaignRequest)
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)

		mockCampaignUsecase.On("Exists", req.Context(), int64(0), "new campaign").Return(false, nil)

		campaignEntity := entities.Campaign{
			Title:               "new campaign",
			StatusCode:          int64(1),
			CampaignType:        "deli",
			ListingTitle:        "test screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			OnboardTitle:        "test campaign",
			OnboardDesc:         "test desc",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OrderStartDate:      time.Date(2023, time.March, 1, 12, 0, 0, 0, time.UTC),
			OrderEndDate:        time.Date(2023, time.March, 31, 12, 0, 0, 0, time.UTC),
			CollectionStartDate: time.Date(2023, time.March, 5, 12, 0, 0, 0, time.UTC),
			CollectionEndDate:   time.Date(2023, time.April, 5, 12, 0, 0, 0, time.UTC),
			OfferID:             123,
			TagID:               456,
			CreatedBy:           12345,
			LeadTime:            3,
		}

		mockCampaignUsecase.On("Create", req.Context(), campaignEntity).Return(nil, errors.New("db error"))

		mockTransactionService.On("RunWithTransaction", req.Context(), mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		campaignController.CreateCampaign(w, req)

		if status := w.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		expected := `{"code":500,"message":"Internal Server Error : db error"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to error occured while adding stores", func(t *testing.T) {
		campaignRequest := bytes.NewBuffer([]byte(`{
			"campaign_status_code": 1,
			"campaign_type": "deli",
			"collection_end_date": "2023-04-05 12:00:00",
			"collection_start_date": "2023-03-05 12:00:00",
			"landing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			"lead_time": 3,
			"listing_description": "test description",
			"listing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			"listing_title": "test screen title",
			"onboarding_description": "test desc",
			"onboarding_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			"onboarding_title": "test campaign",
			"order_end_date": "2023-03-31 12:00:00",
			"order_start_date": "2023-03-01 12:00:00",
			"stores": [
			  83
			],
			"title": "new campaign",
			"offer_id": 123,
			"tag_id": 456
		  }`))
		req, _ := http.NewRequest("POST", "/campaigns", campaignRequest)
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		mockCampaignUsecase.On("Exists", req.Context(), int64(0), "new campaign").Return(false, nil)
		campaignEntity := entities.Campaign{
			Title:               "new campaign",
			StatusCode:          int64(1),
			CampaignType:        "deli",
			ListingTitle:        "test screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			OnboardTitle:        "test campaign",
			OnboardDesc:         "test desc",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OrderStartDate:      time.Date(2023, time.March, 1, 12, 0, 0, 0, time.UTC),
			OrderEndDate:        time.Date(2023, time.March, 31, 12, 0, 0, 0, time.UTC),
			CollectionStartDate: time.Date(2023, time.March, 5, 12, 0, 0, 0, time.UTC),
			CollectionEndDate:   time.Date(2023, time.April, 5, 12, 0, 0, 0, time.UTC),
			OfferID:             123,
			TagID:               456,
			CreatedBy:           12345,
			LeadTime:            3,
		}
		campaignDTO := dto.CampaignDTO{
			ID:                  int64(1),
			Title:               "new campaign",
			StatusCode:          1,
			CampaignType:        "deli",
			ListingTitle:        "test screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			OnboardTitle:        "test campaign",
			OnboardDesc:         "test desc",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OrderStartDate:      "2023-03-01 12:00:00",
			OrderEndDate:        "2023-03-31 12:00:00",
			CollectionStartDate: "2023-03-05 12:00:00",
			CollectionEndDate:   "2023-04-05 12:00:00",
			OfferID:             123,
			TagID:               456,
			LeadTime:            3,
		}
		storeEntities := []entities.CampaignStore{
			{
				CampaignID: 1,
				StoreID:    83,
				CreatedBy:  12345,
			},
		}
		mockCampaignUsecase.On("Create", req.Context(), campaignEntity).Return(&campaignDTO, nil)
		mockCampaignStoreUsecase.On("AddStores", req.Context(), storeEntities).Return(nil, errors.New("db error"))
		mockTransactionService.On("RunWithTransaction", req.Context(), mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})
		campaignController.CreateCampaign(w, req)
		if status := w.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
		expected := `{"code":500,"message":"Internal Server Error : db error"}`
		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("Success : campaign created successfully", func(t *testing.T) {
		campaignRequest := bytes.NewBuffer([]byte(`{
			"campaign_status_code": 1,
			"campaign_type": "deli",
			"collection_end_date": "2023-04-05 12:00:00",
			"collection_start_date": "2023-03-05 12:00:00",
			"landing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			"lead_time": 3,
			"listing_description": "test description",
			"listing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			"listing_title": "test screen title",
			"onboarding_description": "test desc",
			"onboarding_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			"onboarding_title": "test campaign",
			"order_end_date": "2023-03-31 12:00:00",
			"order_start_date": "2023-03-01 12:00:00",
			"stores": [
			  83
			],
			"title": "new campaign",
			"offer_id": 123,
			"tag_id": 456
		  }`))
		req, _ := http.NewRequest("POST", "/campaigns", campaignRequest)
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		mockCampaignUsecase.On("Exists", req.Context(), int64(0), "new campaign").Return(false, nil)
		campaignEntity := entities.Campaign{
			Title:               "new campaign",
			StatusCode:          int64(1),
			CampaignType:        "deli",
			ListingTitle:        "test screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			OnboardTitle:        "test campaign",
			OnboardDesc:         "test desc",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OrderStartDate:      time.Date(2023, time.March, 1, 12, 0, 0, 0, time.UTC),
			OrderEndDate:        time.Date(2023, time.March, 31, 12, 0, 0, 0, time.UTC),
			CollectionStartDate: time.Date(2023, time.March, 5, 12, 0, 0, 0, time.UTC),
			CollectionEndDate:   time.Date(2023, time.April, 5, 12, 0, 0, 0, time.UTC),
			OfferID:             123,
			TagID:               456,
			CreatedBy:           12345,
			LeadTime:            3,
		}
		campaignDTO := dto.CampaignDTO{
			ID:                  int64(1),
			Title:               "new campaign",
			Name:                "new campaign",
			StatusCode:          1,
			CampaignType:        "deli",
			ListingTitle:        "test screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			OnboardTitle:        "test campaign",
			OnboardDesc:         "test desc",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OrderStartDate:      "2023-03-01 12:00:00",
			OrderEndDate:        "2023-03-31 12:00:00",
			CollectionStartDate: "2023-03-05 12:00:00",
			CollectionEndDate:   "2023-04-05 12:00:00",
			OfferID:             123,
			TagID:               456,
			LeadTime:            3,
			IsCampaignPublished: false,
		}
		storeEntities := []entities.CampaignStore{
			{
				CampaignID: 1,
				StoreID:    83,
				CreatedBy:  12345,
			},
		}
		storesDTO := []*dto.CampaignStores{
			{
				ID:      1,
				StoreID: 83,
			},
		}
		mockCampaignUsecase.On("Create", req.Context(), campaignEntity).Return(&campaignDTO, nil)
		mockCampaignStoreUsecase.On("AddStores", req.Context(), storeEntities).Return(storesDTO, nil)
		mockTransactionService.On("RunWithTransaction", req.Context(), mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})
		campaignController.CreateCampaign(w, req)
		if status := w.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		expected := `{"id":1,"campaign_title":"new campaign","name":"new campaign","campaign_status_code":1,"campaign_type":"deli","listing_title":"test screen title","listing_description":"test description","listing_image_path":"https://preprod-media.nedigital.sg/fairprice/images/img2.jpg","onboarding_title":"test campaign","onboarding_description":"test desc","onboard_image_path":"https://preprod-media.nedigital.sg/fairprice/images/img3.jpg","landing_image_path":"https://preprod-media.nedigital.sg/fairprice/images/img1.jpg","order_start_date":"2023-03-01 12:00:00","order_end_date":"2023-03-31 12:00:00","collection_start_date":"2023-03-05 12:00:00","collection_end_date":"2023-04-05 12:00:00","lead_time":3,"offer_id":123,"tag_id":456,"is_campaign_published":false,"campaign_stores":[{"campaign_store_id":1,"store_id":83}]}`
		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})
}

func TestCampaignController_saveCampaignDetails(t *testing.T) {
	appConfig := entities.AppCfg{
		ValidationParam: entities.ValidationParam{
			MaxLeadTime:       20,
			MaxDateDifference: 28,
		},
	}
	t.Run("failure due to error occured while campaign entity conversion", func(t *testing.T) {
		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		request := params.CampaignCreationForm{
			Title:               "new campaign",
			StatusCode:          1,
			CampaignType:        "deli",
			ListingTitle:        "test screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			OnboardTitle:        "test campaign",
			OnboardDesc:         "test desc",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OrderEndDate:        "2023-03-31 12:00:00",
			OrderStartDate:      "2023-03-01 12:00:00",
			CollectionEndDate:   "2023-04-05 12:00:00",
			CollectionStartDate: "2023-13-02 12:00:00",
			OfferID:             123,
			TagID:               456,
			LeadTime:            3,
		}
		ctx := context.WithValue(context.Background(), "userId", int64(12345))

		_, err := campaignController.saveCampaignDetails(ctx, request, int64(12345))
		expectedErr := `parsing time "2023-13-02 12:00:00": month out of range`
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})
}

func TestCampaignController_validateCampaignDates(t *testing.T) {
	appConfig := entities.AppCfg{
		ValidationParam: entities.ValidationParam{
			MaxLeadTime:       20,
			MaxDateDifference: 28,
		},
	}
	t.Run("failure due invalid order start date", func(t *testing.T) {
		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		dates := params.CampaignDates{
			OrderStartDate:      "2023-03-01",
			OrderEndDate:        "2023-03-31 12:00:00",
			CollectionStartDate: "2023-03-02 12:00:00",
			CollectionEndDate:   "2023-04-05 12:00:00",
		}
		err := campaignController.validateCampaignDates(dates, 10)
		expectedErr := `parsing time "2023-03-01" as "2006-01-02 15:04:05": cannot parse "" as "15"`
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})
	t.Run("failure due invalid order end date", func(t *testing.T) {
		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		dates := params.CampaignDates{
			OrderStartDate:      "2023-03-31 12:00:00",
			OrderEndDate:        "2023-03-31",
			CollectionStartDate: "2023-03-02 12:00:00",
			CollectionEndDate:   "2023-04-05 12:00:00",
		}
		err := campaignController.validateCampaignDates(dates, 10)
		expectedErr := `parsing time "2023-03-31" as "2006-01-02 15:04:05": cannot parse "" as "15"`
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})
	t.Run("failure due invalid collection start date", func(t *testing.T) {
		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		dates := params.CampaignDates{
			OrderStartDate:      "2023-03-03 12:00:00",
			OrderEndDate:        "2023-03-31 12:00:00",
			CollectionStartDate: "2023-03-02 12:00",
			CollectionEndDate:   "2023-04-05 12:00:00",
		}
		err := campaignController.validateCampaignDates(dates, 10)
		expectedErr := `parsing time "2023-03-02 12:00" as "2006-01-02 15:04:05": cannot parse "" as ":"`
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})
	t.Run("failure due invalid collection end date", func(t *testing.T) {
		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		dates := params.CampaignDates{
			OrderStartDate:      "2023-03-03 12:00:00",
			OrderEndDate:        "2023-03-31 12:00:00",
			CollectionStartDate: "2023-03-02 12:00:00",
			CollectionEndDate:   "2023-04-77 12:00:00",
		}
		err := campaignController.validateCampaignDates(dates, 10)
		expectedErr := `parsing time "2023-04-77 12:00:00": day out of range`
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})
	t.Run("failure due order start date is after order end date", func(t *testing.T) {
		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		dates := params.CampaignDates{
			OrderStartDate:      "2023-03-03 12:00:00",
			OrderEndDate:        "2023-03-02 12:00:00",
			CollectionStartDate: "2023-03-02 12:00:00",
			CollectionEndDate:   "2023-04-02 12:00:00",
		}
		err := campaignController.validateCampaignDates(dates, 10)
		expectedErr := `invalid date : order start date should be the date before order end date`
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})
	t.Run("failure due order start date is after collection start date", func(t *testing.T) {
		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		dates := params.CampaignDates{
			OrderStartDate:      "2023-03-03 12:00:00",
			OrderEndDate:        "2023-03-31 12:00:00",
			CollectionStartDate: "2023-03-02 12:00:00",
			CollectionEndDate:   "2023-04-02 12:00:00",
		}
		err := campaignController.validateCampaignDates(dates, 10)
		expectedErr := `invalid date : order start date should be the date before collection start date`
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})
	t.Run("failure due collection start date is after collection end date", func(t *testing.T) {
		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		dates := params.CampaignDates{
			OrderStartDate:      "2023-03-03 12:00:00",
			OrderEndDate:        "2023-03-31 12:00:00",
			CollectionStartDate: "2023-03-15 12:00:00",
			CollectionEndDate:   "2023-03-05 12:00:00",
		}
		err := campaignController.validateCampaignDates(dates, 10)
		expectedErr := `invalid date : collection start date should be the date before collection end date`
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})
	t.Run("failure due collection start date and order start date validation not matched", func(t *testing.T) {
		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		dates := params.CampaignDates{
			OrderStartDate:      "2023-03-01 12:00:00",
			OrderEndDate:        "2023-03-31 12:00:00",
			CollectionStartDate: "2023-03-30 12:00:00",
			CollectionEndDate:   "2023-04-11 12:00:00",
		}
		err := campaignController.validateCampaignDates(dates, 10)
		expectedErr := `invalid date : Collection start date should be less than 28 days from order start date`
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})
	t.Run("failure due collection end date and order end date validation not matched", func(t *testing.T) {
		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		dates := params.CampaignDates{
			OrderStartDate:      "2023-03-01 12:00:00",
			OrderEndDate:        "2023-03-31 12:00:00",
			CollectionStartDate: "2023-03-28 12:00:00",
			CollectionEndDate:   "2023-04-30 12:00:00",
		}
		err := campaignController.validateCampaignDates(dates, 10)
		expectedErr := `invalid date : Collection end date should be less than 28 days from order end date`
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})
	t.Run("failure due collection end date less than lead time", func(t *testing.T) {
		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		dates := params.CampaignDates{
			OrderStartDate:      "2023-03-01 12:00:00",
			OrderEndDate:        "2023-03-31 12:00:00",
			CollectionStartDate: "2023-03-28 12:00:00",
			CollectionEndDate:   "2023-04-09 12:00:00",
		}
		err := campaignController.validateCampaignDates(dates, 10)
		expectedErr := `invalid date : Collection end date should be at least 10 days greater than order end date`
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})
	t.Run("success scenario", func(t *testing.T) {
		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		dates := params.CampaignDates{
			OrderStartDate:      "2023-03-03 12:00:00",
			OrderEndDate:        "2023-03-31 12:00:00",
			CollectionStartDate: "2023-03-15 12:00:00",
			CollectionEndDate:   "2023-04-28 12:00:00",
		}
		err := campaignController.validateCampaignDates(dates, 3)
		ShouldBeNil(err)
		if err != nil {
			t.Errorf("unexpected error : got - %v ; want - nil", err.Error())
		}
	})
}

func TestCampaignController_UpdateCampaign(t *testing.T) {
	appConfig := entities.AppCfg{
		ValidationParam: entities.ValidationParam{
			MaxLeadTime: 20,
			MaxDateDifference: 28,
		},
	}

	t.Run("failure due to incorrect user id", func(t *testing.T) {
		var jsonStr = []byte(`{
			"campaign_status_code": 1,
			"campaign_type" : "deli",
			"title": "my new campaign 113",
			"stores":[11111, 22222, 33333, 44444]
		  }`)

		req, _ := http.NewRequest("PUT", "/campaigns/1", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		campaignController.UpdateCampaign(w, req)

		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		expected := `{"code":400,"message":"Bad Request : invalid user id"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to incorrect campaign id", func(t *testing.T) {
		var jsonStr = []byte(`{
			"campaign_status_code": 1,
			"campaign_type" : "deli",
			"title": "my new campaign 113",
			"stores":[11111, 22222, 33333, 44444]
		  }`)

		req, _ := http.NewRequest("PUT", "/campaigns/1", bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		w := httptest.NewRecorder()

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		campaignController.UpdateCampaign(w, req)

		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		expected := `{"code":400,"message":"Bad Request : incorrect campaign id value, err : strconv.Atoi: parsing \"\": invalid syntax"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to request validation error", func(t *testing.T) {
		var jsonStr = []byte(`{"campaign_status_code": 1,
		"campaign_type": "deli}`)

		req, _ := http.NewRequest("PUT", "/campaigns/1", bytes.NewBuffer(jsonStr))
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))

		w := httptest.NewRecorder()

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		campaignController.UpdateCampaign(w, req)

		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		expected := `{"code":400,"message":"Bad Request : unexpected EOF"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to error occured while checking campaign existence", func(t *testing.T) {
		campaignRequest := bytes.NewBuffer([]byte(`{
			"campaign_status_code": 1,
			"campaign_type": "deli",
			"collection_end_date": "2023-04-05 12:00:00",
			"collection_start_date": "2023-03-05 12:00:00",
			"landing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			"lead_time": 3,
			"listing_description": "test description",
			"listing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			"listing_title": "test screen title",
			"onboarding_description": "test desc",
			"onboarding_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			"onboarding_title": "test campaign",
			"order_end_date": "2023-03-31 12:00:00",
			"order_start_date": "2023-03-01 12:00:00",
			"stores": [
			  83,84,85
			],
			"title": "new campaign",
			"offer_id": 123,
			"tag_id": 456
		  }`))
		req, _ := http.NewRequest("PUT", "/campaigns/1", campaignRequest)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		w := httptest.NewRecorder()

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)

		mockCampaignUsecase.On("Exists", req.Context(), int64(1), "").Return(false, errors.New("db error"))

		campaignController.UpdateCampaign(w, req)

		if status := w.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		expected := `{"code":500,"message":"Internal Server Error : db error"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to campaign with given id not exists", func(t *testing.T) {
		campaignRequest := bytes.NewBuffer([]byte(`{
			"campaign_status_code": 1,
			"campaign_type": "deli",
			"collection_end_date": "2023-04-05 12:00:00",
			"collection_start_date": "2023-03-05 12:00:00",
			"landing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			"lead_time": 3,
			"listing_description": "test description",
			"listing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			"listing_title": "test screen title",
			"onboarding_description": "test desc",
			"onboarding_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			"onboarding_title": "test campaign",
			"order_end_date": "2023-03-31 12:00:00",
			"order_start_date": "2023-03-01 12:00:00",
			"stores": [
			  83,84,85
			],
			"title": "new campaign",
			"offer_id": 123,
			"tag_id": 456
		  }`))
		req, _ := http.NewRequest("PUT", "/campaigns/1", campaignRequest)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)

		mockCampaignUsecase.On("Exists", req.Context(), int64(1), "").Return(false, nil)

		campaignController.UpdateCampaign(w, req)

		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		expected := `{"code":400,"message":"Bad Request : campaign with given id not exists"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to error occured while updating campaign", func(t *testing.T) {
		campaignRequest := bytes.NewBuffer([]byte(`{
			"campaign_status_code": 1,
			"campaign_type": "deli",
			"collection_end_date": "2023-04-05 12:00:00",
			"collection_start_date": "2023-03-05 12:00:00",
			"landing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			"lead_time": 3,
			"listing_description": "test description",
			"listing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			"listing_title": "test screen title",
			"onboarding_description": "test desc",
			"onboarding_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			"onboarding_title": "test campaign",
			"order_end_date": "2023-03-31 12:00:00",
			"order_start_date": "2023-03-01 12:00:00",
			"stores": [
			  83,84,85
			],
			"title": "new campaign",
			"offer_id": 123,
			"tag_id": 456
		  }`))
		req, _ := http.NewRequest("PUT", "/campaigns/1", campaignRequest)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)

		mockCampaignUsecase.On("Exists", req.Context(), int64(1), "").Return(true, nil)

		campaignEntity := entities.Campaign{
			ID:                  valueobjects.CampaignID(1),
			Title:               "new campaign",
			StatusCode:          int64(1),
			CampaignType:        "deli",
			ListingTitle:        "test screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			OnboardTitle:        "test campaign",
			OnboardDesc:         "test desc",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OrderStartDate:      time.Date(2023, time.March, 1, 12, 0, 0, 0, time.UTC),
			OrderEndDate:        time.Date(2023, time.March, 31, 12, 0, 0, 0, time.UTC),
			CollectionStartDate: time.Date(2023, time.March, 5, 12, 0, 0, 0, time.UTC),
			CollectionEndDate:   time.Date(2023, time.April, 5, 12, 0, 0, 0, time.UTC),
			OfferID:             123,
			TagID:               456,
			UpdatedBy:           12345,
			LeadTime:            3,
		}

		mockCampaignUsecase.On("Update", req.Context(), campaignEntity).Return(errors.New("db error"))

		mockTransactionService.On("RunWithTransaction", req.Context(), mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		campaignController.UpdateCampaign(w, req)

		if status := w.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		expected := `{"code":500,"message":"Internal Server Error : db error"}`

		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to error occured while getting store details", func(t *testing.T) {
		campaignRequest := bytes.NewBuffer([]byte(`{
			"campaign_status_code": 1,
			"campaign_type": "deli",
			"collection_end_date": "2023-04-05 12:00:00",
			"collection_start_date": "2023-03-05 12:00:00",
			"landing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			"lead_time": 3,
			"listing_description": "test description",
			"listing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			"listing_title": "test screen title",
			"onboarding_description": "test desc",
			"onboarding_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			"onboarding_title": "test campaign",
			"order_end_date": "2023-03-31 12:00:00",
			"order_start_date": "2023-03-01 12:00:00",
			"stores": [
			  83
			],
			"title": "new campaign",
			"offer_id": 123,
			"tag_id": 456
		  }`))

		req, _ := http.NewRequest("PUT", "/campaigns/1", campaignRequest)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		mockCampaignUsecase.On("Exists", req.Context(), int64(1), "").Return(true, nil)
		campaignEntity := entities.Campaign{
			ID:                  valueobjects.CampaignID(1),
			Title:               "new campaign",
			StatusCode:          int64(1),
			CampaignType:        "deli",
			ListingTitle:        "test screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			OnboardTitle:        "test campaign",
			OnboardDesc:         "test desc",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OrderStartDate:      time.Date(2023, time.March, 1, 12, 0, 0, 0, time.UTC),
			OrderEndDate:        time.Date(2023, time.March, 31, 12, 0, 0, 0, time.UTC),
			CollectionStartDate: time.Date(2023, time.March, 5, 12, 0, 0, 0, time.UTC),
			CollectionEndDate:   time.Date(2023, time.April, 5, 12, 0, 0, 0, time.UTC),
			OfferID:             123,
			TagID:               456,
			UpdatedBy:           12345,
			LeadTime:            3,
		}
		mockCampaignUsecase.On("Update", req.Context(), campaignEntity).Return(nil)
		mockCampaignStoreUsecase.On("GetByStoreID", req.Context(), int64(1), int64(83)).Return(dto.CampaignStores{},
			errors.New("db error"))
		mockTransactionService.On("RunWithTransaction", req.Context(), mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})
		campaignController.UpdateCampaign(w, req)
		if status := w.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
		expected := `{"code":500,"message":"Internal Server Error : db error"}`
		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to error occured while adding new stores", func(t *testing.T) {
		campaignRequest := bytes.NewBuffer([]byte(`{
			"campaign_status_code": 1,
			"campaign_type": "deli",
			"collection_end_date": "2023-04-05 12:00:00",
			"collection_start_date": "2023-03-05 12:00:00",
			"landing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			"lead_time": 3,
			"listing_description": "test description",
			"listing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			"listing_title": "test screen title",
			"onboarding_description": "test desc",
			"onboarding_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			"onboarding_title": "test campaign",
			"order_end_date": "2023-03-31 12:00:00",
			"order_start_date": "2023-03-01 12:00:00",
			"stores": [
			  83
			],
			"title": "new campaign",
			"offer_id": 123,
			"tag_id": 456
		  }`))

		req, _ := http.NewRequest("PUT", "/campaigns/1", campaignRequest)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		mockCampaignUsecase.On("Exists", req.Context(), int64(1), "").Return(true, nil)
		campaignEntity := entities.Campaign{
			ID:                  valueobjects.CampaignID(1),
			Title:               "new campaign",
			StatusCode:          int64(1),
			CampaignType:        "deli",
			ListingTitle:        "test screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			OnboardTitle:        "test campaign",
			OnboardDesc:         "test desc",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OrderStartDate:      time.Date(2023, time.March, 1, 12, 0, 0, 0, time.UTC),
			OrderEndDate:        time.Date(2023, time.March, 31, 12, 0, 0, 0, time.UTC),
			CollectionStartDate: time.Date(2023, time.March, 5, 12, 0, 0, 0, time.UTC),
			CollectionEndDate:   time.Date(2023, time.April, 5, 12, 0, 0, 0, time.UTC),
			OfferID:             123,
			TagID:               456,
			UpdatedBy:           12345,
			LeadTime:            3,
		}
		storeEntities := []entities.CampaignStore{
			{
				CampaignID: 1,
				StoreID:    83,
				CreatedBy:  12345,
			},
		}
		mockCampaignUsecase.On("Update", req.Context(), campaignEntity).Return(nil)
		mockCampaignStoreUsecase.On("GetByStoreID", req.Context(), int64(1), int64(83)).Return(dto.CampaignStores{},
			fmt.Errorf("%w", valueobjects.ErrStoreNotExists))
		mockCampaignStoreUsecase.On("GetStores", req.Context(), int64(1)).Return([]*dto.CampaignStores{}, nil)
		mockCampaignStoreUsecase.On("AddStores", req.Context(), storeEntities).Return(nil, errors.New("db error"))
		mockTransactionService.On("RunWithTransaction", req.Context(), mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})
		campaignController.UpdateCampaign(w, req)
		if status := w.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
		expected := `{"code":500,"message":"Internal Server Error : db error"}`
		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to error occured while adding getting all stores", func(t *testing.T) {
		campaignRequest := bytes.NewBuffer([]byte(`{
			"campaign_status_code": 1,
			"campaign_type": "deli",
			"collection_end_date": "2023-04-05 12:00:00",
			"collection_start_date": "2023-03-05 12:00:00",
			"landing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			"lead_time": 3,
			"listing_description": "test description",
			"listing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			"listing_title": "test screen title",
			"onboarding_description": "test desc",
			"onboarding_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			"onboarding_title": "test campaign",
			"order_end_date": "2023-03-31 12:00:00",
			"order_start_date": "2023-03-01 12:00:00",
			"stores": [
			  83
			],
			"title": "new campaign",
			"offer_id": 123,
			"tag_id": 456
		  }`))

		req, _ := http.NewRequest("PUT", "/campaigns/1", campaignRequest)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		mockCampaignUsecase.On("Exists", req.Context(), int64(1), "").Return(true, nil)
		campaignEntity := entities.Campaign{
			ID:                  valueobjects.CampaignID(1),
			Title:               "new campaign",
			StatusCode:          int64(1),
			CampaignType:        "deli",
			ListingTitle:        "test screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			OnboardTitle:        "test campaign",
			OnboardDesc:         "test desc",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OrderStartDate:      time.Date(2023, time.March, 1, 12, 0, 0, 0, time.UTC),
			OrderEndDate:        time.Date(2023, time.March, 31, 12, 0, 0, 0, time.UTC),
			CollectionStartDate: time.Date(2023, time.March, 5, 12, 0, 0, 0, time.UTC),
			CollectionEndDate:   time.Date(2023, time.April, 5, 12, 0, 0, 0, time.UTC),
			OfferID:             123,
			TagID:               456,
			UpdatedBy:           12345,
			LeadTime:            3,
		}
		mockCampaignUsecase.On("Update", req.Context(), campaignEntity).Return(nil)
		mockCampaignStoreUsecase.On("GetByStoreID", req.Context(), int64(1), int64(83)).Return(dto.CampaignStores{},
			fmt.Errorf("%w", valueobjects.ErrStoreNotExists))
		mockCampaignStoreUsecase.On("GetStores", req.Context(), int64(1)).Return([]*dto.CampaignStores{}, errors.New("db error"))
		mockTransactionService.On("RunWithTransaction", req.Context(), mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})
		campaignController.UpdateCampaign(w, req)
		if status := w.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
		expected := `{"code":500,"message":"Internal Server Error : db error"}`
		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("failure due to error occured while deleting old stores", func(t *testing.T) {
		campaignRequest := bytes.NewBuffer([]byte(`{
			"campaign_status_code": 1,
			"campaign_type": "deli",
			"collection_end_date": "2023-04-05 12:00:00",
			"collection_start_date": "2023-03-05 12:00:00",
			"landing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			"lead_time": 3,
			"listing_description": "test description",
			"listing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			"listing_title": "test screen title",
			"onboarding_description": "test desc",
			"onboarding_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			"onboarding_title": "test campaign",
			"order_end_date": "2023-03-31 12:00:00",
			"order_start_date": "2023-03-01 12:00:00",
			"stores": [
			  83
			],
			"title": "new campaign",
			"offer_id": 123,
			"tag_id": 456
		  }`))

		req, _ := http.NewRequest("PUT", "/campaigns/1", campaignRequest)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		mockCampaignUsecase.On("Exists", req.Context(), int64(1), "").Return(true, nil)
		campaignEntity := entities.Campaign{
			ID:                  valueobjects.CampaignID(1),
			Title:               "new campaign",
			StatusCode:          int64(1),
			CampaignType:        "deli",
			ListingTitle:        "test screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			OnboardTitle:        "test campaign",
			OnboardDesc:         "test desc",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OrderStartDate:      time.Date(2023, time.March, 1, 12, 0, 0, 0, time.UTC),
			OrderEndDate:        time.Date(2023, time.March, 31, 12, 0, 0, 0, time.UTC),
			CollectionStartDate: time.Date(2023, time.March, 5, 12, 0, 0, 0, time.UTC),
			CollectionEndDate:   time.Date(2023, time.April, 5, 12, 0, 0, 0, time.UTC),
			OfferID:             123,
			TagID:               456,
			UpdatedBy:           12345,
			LeadTime:            3,
		}
		storeEntities := []entities.CampaignStore{
			{
				CampaignID: 1,
				StoreID:    83,
				CreatedBy:  12345,
			},
		}
		createStoresDTO := []*dto.CampaignStores{
			{
				ID:      2,
				StoreID: 83,
			},
		}
		getStoresDTO := []*dto.CampaignStores{
			{
				ID:      1,
				StoreID: 82,
			},
		}
		mockCampaignUsecase.On("Update", req.Context(), campaignEntity).Return(nil)
		mockCampaignStoreUsecase.On("GetByStoreID", req.Context(), int64(1), int64(83)).Return(dto.CampaignStores{},
			fmt.Errorf("%w", valueobjects.ErrStoreNotExists))
		mockCampaignStoreUsecase.On("GetStores", req.Context(), int64(1)).Return(getStoresDTO, nil)
		mockCampaignStoreUsecase.On("AddStores", req.Context(), storeEntities).Return(createStoresDTO, nil)
		mockCampaignStoreUsecase.On("DeleteByStoreID", req.Context(), int64(1), int64(82), int64(12345)).Return(errors.New("db error"))
		mockTransactionService.On("RunWithTransaction", req.Context(), mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})
		campaignController.UpdateCampaign(w, req)
		if status := w.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
		expected := `{"code":500,"message":"Internal Server Error : db error"}`
		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("Success : campaign update successfully", func(t *testing.T) {
		campaignRequest := bytes.NewBuffer([]byte(`{
			"campaign_status_code": 1,
			"campaign_type": "deli",
			"collection_end_date": "2023-04-05 12:00:00",
			"collection_start_date": "2023-03-05 12:00:00",
			"landing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			"lead_time": 3,
			"listing_description": "test description",
			"listing_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			"listing_title": "test screen title",
			"onboarding_description": "test desc",
			"onboarding_image_path": "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			"onboarding_title": "test campaign",
			"order_end_date": "2023-03-31 12:00:00",
			"order_start_date": "2023-03-01 12:00:00",
			"stores": [
			  84
			],
			"title": "new campaign",
			"offer_id": 123,
			"tag_id": 456
		  }`))
		req, _ := http.NewRequest("PUT", "/campaigns/1", campaignRequest)
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req = req.WithContext(context.WithValue(req.Context(), "userId", 12345))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		mockCampaignUsecase.On("Exists", req.Context(), int64(1), "").Return(true, nil)
		campaignEntity := entities.Campaign{
			ID:                  valueobjects.CampaignID(1),
			Title:               "new campaign",
			StatusCode:          int64(1),
			CampaignType:        "deli",
			ListingTitle:        "test screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			OnboardTitle:        "test campaign",
			OnboardDesc:         "test desc",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OrderStartDate:      time.Date(2023, time.March, 1, 12, 0, 0, 0, time.UTC),
			OrderEndDate:        time.Date(2023, time.March, 31, 12, 0, 0, 0, time.UTC),
			CollectionStartDate: time.Date(2023, time.March, 5, 12, 0, 0, 0, time.UTC),
			CollectionEndDate:   time.Date(2023, time.April, 5, 12, 0, 0, 0, time.UTC),
			OfferID:             123,
			TagID:               456,
			UpdatedBy:           12345,
			LeadTime:            3,
		}
		storeEntities := []entities.CampaignStore{
			{
				CampaignID: 1,
				StoreID:    84,
				CreatedBy:  12345,
			},
		}
		createStoresDTO := []*dto.CampaignStores{
			{
				ID:      2,
				StoreID: 84,
			},
		}
		getStoresDTO := []*dto.CampaignStores{
			{
				ID:      1,
				StoreID: 83,
			},
		}
		mockCampaignUsecase.On("Update", req.Context(), campaignEntity).Return(nil)
		mockCampaignStoreUsecase.On("GetByStoreID", req.Context(), int64(1), int64(84)).Return(dto.CampaignStores{},
			fmt.Errorf("%w", valueobjects.ErrStoreNotExists))
		mockCampaignStoreUsecase.On("GetStores", req.Context(), int64(1)).Return(getStoresDTO, nil)
		mockCampaignStoreUsecase.On("AddStores", req.Context(), storeEntities).Return(createStoresDTO, nil)
		mockCampaignStoreUsecase.On("DeleteByStoreID", req.Context(), int64(1), int64(83), int64(12345)).Return(nil)
		mockTransactionService.On("RunWithTransaction", req.Context(), mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		campaignController.UpdateCampaign(w, req)

		if status := w.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		expected := `{"code":200,"message":"campaign with id 1 updated successfully"}`
		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})
}

func TestCampaignController_updateCampaignDetails(t *testing.T) {
	appConfig := entities.AppCfg{
		ValidationParam: entities.ValidationParam{
			MaxLeadTime: 20,
			MaxDateDifference: 28,
		},
	}
	t.Run("failure due to error occured while campaign entity conversion", func(t *testing.T) {
		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		request := params.CampaignUpdateForm{
			Title:               "new campaign",
			StatusCode:          1,
			CampaignType:        "deli",
			ListingTitle:        "test screen title",
			ListingDesc:         "test description",
			ListingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img2.jpg",
			OnboardTitle:        "test campaign",
			OnboardDesc:         "test desc",
			OnboardImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img3.jpg",
			LandingImagePath:    "https://preprod-media.nedigital.sg/fairprice/images/img1.jpg",
			OrderEndDate:        "2023-03-31 12:00:00",
			OrderStartDate:      "2023-03-01 12:00:00",
			CollectionEndDate:   "2023-04-05 12:00:00",
			CollectionStartDate: "2023-13-02 12:00:00",
			OfferID:             123,
			TagID:               456,
			LeadTime:            3,
		}
		ctx := context.WithValue(context.Background(), "userId", int64(12345))
		err := campaignController.updateCampaignDetails(ctx, int64(1), request, int64(12345))
		expectedErr := `parsing time "2023-13-02 12:00:00": month out of range`
		ShouldNotBeNil(err)
		if err.Error() != expectedErr {
			t.Errorf("unexpected error : got - %v ; want - %v", err.Error(), expectedErr)
		}
	})
}

func TestCampaignController_GetCampaignList(t *testing.T) {
	appConfig := entities.AppCfg{
		ValidationParam: entities.ValidationParam{
			MaxLeadTime: 20,
			MaxDateDifference: 28,
		},
	}
	mockCampaignUsecase := mocks.NewCampaignUseCases(t)
	mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
	mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
	mockTransactionService := service_mocks.NewTransactionService(t)
	campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
		mockTransactionService, &appConfig)

	t.Run("Get Campaign List request success for InActive", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/campaigns", nil)
		query := req.URL.Query()
		query.Add("limit", "10")
		query.Add("page", "1")
		query.Add("sort", "created_at desc")
		query.Add("status", "InActive")
		query.Add("name", "campaign1")
		req.URL.RawQuery = query.Encode()
		res := httptest.NewRecorder()
		response := dto.CampaignListResponse{}
		mockCampaignUsecase.On("GetList", req.Context(), entities.PaginationConfig{
			Limit: 10, Page: 1, Sort: "created_at desc", Name: "campaign1", Status: 1}).Return(&response, nil)
		campaignController.GetCampaignList(res, req)
		ShouldBeNil(err)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Get Campaign List request success for Active status", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/campaigns", nil)
		query := req.URL.Query()
		query.Add("limit", "10")
		query.Add("page", "1")
		query.Add("sort", "created_at desc")
		query.Add("status", "Active")
		query.Add("name", "campaign1")
		req.URL.RawQuery = query.Encode()
		res := httptest.NewRecorder()
		response := dto.CampaignListResponse{}
		mockCampaignUsecase.On("GetList", req.Context(), entities.PaginationConfig{
			Limit: 10, Page: 1, Sort: "created_at desc", Name: "campaign1", Status: 2}).Return(&response, nil)
		campaignController.GetCampaignList(res, req)
		ShouldBeNil(err)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Get Campaign List request success for Scheduled status", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/campaigns", nil)
		query := req.URL.Query()
		query.Add("limit", "10")
		query.Add("page", "1")
		query.Add("sort", "created_at desc")
		query.Add("status", "Scheduled")
		query.Add("name", "campaign1")
		req.URL.RawQuery = query.Encode()
		res := httptest.NewRecorder()
		response := dto.CampaignListResponse{}
		mockCampaignUsecase.On("GetList", req.Context(), entities.PaginationConfig{
			Limit: 10, Page: 1, Sort: "created_at desc", Name: "campaign1", Status: 3}).Return(&response, nil)
		campaignController.GetCampaignList(res, req)
		ShouldBeNil(err)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestCampaignController_UpdateCampaignStatus(t *testing.T) {
	appConfig := entities.AppCfg{
		ValidationParam: entities.ValidationParam{
			MaxLeadTime: 20,
			MaxDateDifference: 28,
		},
	}

	t.Run("failure due to error occured while updating campaign status", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/campaigns/update-status", nil)
		ctx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		mockCampaignUsecase.On("UpdateStatus", req.Context()).Return(errors.New("db error"))
		mockTransactionService.On("RunWithTransaction", req.Context(), mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})
		campaignController.UpdateCampaignStatus(w, req)
		if status := w.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
		expected := `{"code":500,"message":"Internal Server Error : db error"}`
		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})

	t.Run("Success : campaign status updated successfully", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/campaigns/update-status", nil)
		ctx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockCampaignUsecase := mocks.NewCampaignUseCases(t)
		mockCampaignStoreUsecase := mocks.NewCampaignStoreUseCases(t)
		mockCampaignProductUsecase := mocks.NewCampaignProductUseCases(t)
		mockTransactionService := service_mocks.NewTransactionService(t)
		campaignController := NewCampaignController(mockCampaignUsecase, mockCampaignStoreUsecase, mockCampaignProductUsecase,
			mockTransactionService, &appConfig)
		mockCampaignUsecase.On("UpdateStatus", req.Context()).Return(nil)
		mockTransactionService.On("RunWithTransaction", req.Context(), mock.Anything).
			Return(func(ctx context.Context, fn func(context.Context) error) error {
				return fn(ctx)
			})

		campaignController.UpdateCampaignStatus(w, req)

		if status := w.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		expected := `{"code":200,"message":"campaigns status updated successfully"}`
		if a, e := strings.TrimSpace(w.Body.String()), strings.TrimSpace(expected); a != e {
			t.Errorf("handler returned unexpected body: got %v want %v", w.Body.String(), expected)
		}
	})
}
