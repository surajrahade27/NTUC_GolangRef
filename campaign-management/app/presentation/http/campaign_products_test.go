package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCampaignProductController_AddProducts(t *testing.T) {
	// Setup
	// t.Run("Add Products Executed Sucessfully", func(t *testing.T) {
	// 	controller := NewCampaignProductController(nil, nil, nil)

	// 	payload := map[string]interface{}{
	// 		"campaign_id": 1,
	// 		"created_by":  11,
	// 		"products": []entities.CampaignProduct{
	// 			{
	// 				ProductID:   123,
	// 				SKUNo:       1111,
	// 				SerialNo:    1,
	// 				SequenceNo:  1,
	// 				ProductType: "cd",
	// 				CampaignID:  valueobjects.CampaignID(1),
	// 				CreatedBy:   int64(11),
	// 			},
	// 			{
	// 				ProductID:   456,
	// 				SKUNo:       2222,
	// 				SerialNo:    2,
	// 				SequenceNo:  2,
	// 				ProductType: "cd",
	// 				CampaignID:  valueobjects.CampaignID(1),
	// 				CreatedBy:   int64(11),
	// 			},
	// 		},
	// 	}

	// 	payloadBytes, _ := json.Marshal(payload)
	// 	req, _ := http.NewRequest("POST", "/campaigns/products", bytes.NewBuffer(payloadBytes))
	// 	rr := httptest.NewRecorder()

	// 	// Execute
	// 	controller.AddProducts(rr, req)

	// 	// Assert
	// 	assert.Equal(t, http.StatusOK, rr.Code)
	// })

	t.Run("Add Products Executed Failed due to error", func(t *testing.T) {
		controller := NewCampaignProductController(nil, nil, nil)

		payload := map[string]interface{}{}

		payloadBytes, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/campaigns/products", bytes.NewBuffer(payloadBytes))
		rr := httptest.NewRecorder()

		// Execute
		controller.AddProducts(rr, req)

		// Assert
		assert.Equal(t, rr.Body.String(), "{\"code\":400,\"message\":\"Bad Request : Key: 'CampaignProductCreationForm.CampaignID' Error:Field validation for 'CampaignID' failed on the 'required' tag\\nKey: 'CampaignProductCreationForm.Products' Error:Field validation for 'Products' failed on the 'required' tag\\nKey: 'CampaignProductCreationForm.CreatedBy' Error:Field validation for 'CreatedBy' failed on the 'required' tag\"}\n")
	})

}

// func TestAddProducts(t *testing.T) {
// 	// Create a new instance of the controller
// 	controller := &CampaignProductController{}
// 	ctx := context.Background()
// 	// Create a new request with sample data
// 	payload := map[string]interface{}{
// 		"campaign_id": 1,
// 		"created_by":  11,
// 		"products": []entities.CampaignProduct{
// 			{
// 				ProductID:  123,
// 				SKUNo:      1111,
// 				SerialNo:   1,
// 				SequenceNo: 1,
// 				ProductType: "cd",
// 				CampaignID: valueobjects.CampaignID(1),
// 				CreatedBy:  int64(11),
// 			},
// 			{
// 				ProductID:  456,
// 				SKUNo:      2222,
// 				SerialNo:   2,
// 				SequenceNo: 2,
// 				ProductType: "cd",
// 				CampaignID: valueobjects.CampaignID(1),
// 				CreatedBy:  int64(11),
// 			},
// 		},
// 	}

// 	payloadBytes, _ := json.Marshal(payload)
// 	req, err := http.NewRequest("POST", "/campaigns/products", bytes.NewBuffer(payloadBytes))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Create a new response recorder
// 	rr := httptest.NewRecorder()
// 	txs := &mock1.TransactionService{}
// 	productMock := &mocks.CampaignProductUseCases{}
// 	controller.tx = txs
// 	controller.campaignProductUseCases = productMock
// 	txs.On("RunWithTransaction", ctx, mock.Anything).
// 		Return(func(ctx context.Context, fn func(context.Context) error) error {
// 			return fn(ctx)
// 		})
// 	productMock.On("AddProducts", ctx, mock.Anything).Return(nil, nil)

// 	r := chi.NewRouter()
// 	controller.Init(r)
// 	r.ServeHTTP(rr, req)
// 	controller.create(ctx, params.CampaignProductCreationForm{})

// 	// Check that the status code is as expected
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}
// }
