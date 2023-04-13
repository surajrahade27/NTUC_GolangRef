package http

import (
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/services"
	"campaign-mgmt/app/domain/usecases"
	"campaign-mgmt/app/usecases/dto"
	"campaign-mgmt/app/usecases/params"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type CampaignProductController struct {
	campaignProductUseCases usecases.CampaignProductUseCases
	tx                      services.TransactionService
	appConfig               *entities.AppCfg
}

func NewCampaignProductController(
	campaignProductUseCases usecases.CampaignProductUseCases,
	transactionService services.TransactionService,
	appConfig *entities.AppCfg) *CampaignProductController {
	return &CampaignProductController{
		campaignProductUseCases: campaignProductUseCases,
		tx:                      transactionService,
		appConfig:               appConfig,
	}
}
func (c *CampaignProductController) Init(r chi.Router) {
	r.Route("/campaigns/{campaign_id}/products", func(r chi.Router) {
		r.Delete("/{id}", c.DeleteProduct)
		r.Delete("/", c.DeleteAllProduct)
	})
	r.Post("/campaigns/products", c.AddProducts)

}

// CreateCampaignProducts godoc
//
//	@Summary Create a campaign products
//	@Description API to create new campaign products
//	@Tags campaign products
//	@Accept json
//	@Produce json
//	@Param	campaign body params.CampaignProductCreationForm	true "Add campaign products details"
//	@Success 200 {object} []dto.CampaignProducts
//	@Failure 400 {object} dto.Response
//	@Failure 409 {object} dto.Response
//	@Failure 500 {object} dto.Response
//	@Router	/campaigns/products [post]
func (c *CampaignProductController) AddProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	campaignRequest, err := c.validateCampaignRequest(r)
	if err != nil {
		dto.BadRequestJSON(w, r, err.Error())
		return
	}

	response, err := c.create(ctx, campaignRequest)
	if err != nil {
		dto.InternalServerErrorJSON(w, r, err.Error())
		return
	}

	dto.SuccessJSONResponse(w, r, response)
}

func (c *CampaignProductController) create(ctx context.Context, request params.CampaignProductCreationForm) ([]*dto.CampaignProducts, error) {
	response := []*dto.CampaignProducts{}
	var err error
	c.tx.RunWithTransaction(
		ctx, func(ctx context.Context) error {
			if response, err = c.addProducts(ctx, request.Products, request.CampaignID, request.CreatedBy); err != nil {
				return err
			}
			return nil
		})
	return response, err
}

func (c *CampaignProductController) validateCampaignRequest(r *http.Request) (params.CampaignProductCreationForm, error) {
	var campaignProductRequest params.CampaignProductCreationForm
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&campaignProductRequest)
	if err != nil {
		return campaignProductRequest, err
	}
	defer r.Body.Close()

	validate := validator.New()
	err = validate.Struct(campaignProductRequest)
	if err != nil {
		return campaignProductRequest, err
	}

	return campaignProductRequest, nil
}

func (c *CampaignProductController) addProducts(ctx context.Context, productRequest []params.CampaignProduct, campaignID, userID int64) ([]*dto.CampaignProducts, error) {
	productEntities := []entities.CampaignProduct{}
	for _, product := range productRequest {
		productEntities = append(productEntities, params.ToCampaignProductEntity(product, campaignID, userID))
	}
	return c.campaignProductUseCases.AddProducts(ctx, productEntities)
}

// DeleteProduct godoc
//
//	@Summary Delete particular campaign product by product id
//	@Description API to delete particular product under a specified campaign
//	@Tags campaign products
//	@Produce json
//	@Param	campaign_id	path int true "Campaign ID"
//	@Param	id	path int true "Product ID"
//	@Success 200 {object} dto.Response
//	@Failure 400 {object} dto.Response
//	@Failure 404 {object} dto.Response
//	@Failure 500 {object} dto.Response
//	@Router	/campaigns/{campaign_id}/products/{id} [delete]
func (c *CampaignProductController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	campaignID, campaignIDError := strconv.Atoi(chi.URLParam(r, "campaign_id"))
	if campaignIDError != nil {
		dto.BadRequestJSON(w, r, "please provide correct campaign_id ")
	}
	productID, productIDError := strconv.Atoi(chi.URLParam(r, "id"))
	if productIDError != nil {
		dto.BadRequestJSON(w, r, "please provide correct id")
	}
	err := c.campaignProductUseCases.DeleteByCampaignId(ctx, int64(campaignID), int64(productID))
	if err != nil {
		dto.BadRequestJSON(w, r, "product not deleted")
	}
	dto.SuccessJSON(w, r, fmt.Sprintf("product with id %d deleted successfully", productID))
}

// DeleteProduct godoc
//
//	@Summary Delete particular campaign products
//	@Description API to delete all products under a specified campaign
//	@Tags campaign products
//	@Produce json
//	@Param	campaign_id	path int true "Campaign ID"
//	@Success 200 {object} dto.Response
//	@Failure 400 {object} dto.Response
//	@Failure 404 {object} dto.Response
//	@Failure 500 {object} dto.Response
//	@Router	/campaigns/{campaign_id}/products [delete]
func (c *CampaignProductController) DeleteAllProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	campaignID, campaignIDError := strconv.Atoi(chi.URLParam(r, "campaign_id"))
	if campaignIDError != nil {
		dto.BadRequestJSON(w, r, "please provide correct campaign_id ")
	}
	err := c.campaignProductUseCases.DeleteAllByCampaignId(ctx, int64(campaignID))
	if err != nil {
		dto.BadRequestJSON(w, r, "products not deleted")
	}
	dto.SuccessJSON(w, r, fmt.Sprintf("all products for campaign id %d deleted successfully", campaignID))
}
