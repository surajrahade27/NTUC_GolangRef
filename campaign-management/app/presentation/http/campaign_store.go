package http

import (
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/usecases"
	"campaign-mgmt/app/usecases/dto"
	"campaign-mgmt/app/usecases/params"
	"campaign-mgmt/app/usecases/util"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

const (
	IncorrectCampaignIDErr = "incorrect campaign id value, err : %v"
	CampaignNotExistsErr   = "campaign with id %d not exists"
)

type CampaignStoreController struct {
	campaignUseCases      usecases.CampaignUseCases
	campaignStoreUseCases usecases.CampaignStoreUseCases
}

func NewCampaignStoreController(campaignUsecases usecases.CampaignUseCases,
	campaignStoreUseCases usecases.CampaignStoreUseCases) *CampaignStoreController {
	return &CampaignStoreController{
		campaignUseCases:      campaignUsecases,
		campaignStoreUseCases: campaignStoreUseCases,
	}
}

func (c *CampaignStoreController) Init(r chi.Router) {
	r.Route("/campaigns/{campaign_id}/stores", func(r chi.Router) {
		r.Delete("/", c.DeleteStores)
		r.Delete("/{id}", c.DeleteStore)
		r.Post("/", c.AddStores)
	})
}

// DeleteStores godoc
//
//	@Summary Delete all stores under partilcular campaign
//	@Description API to delete all stores under specified campaign
//	@Tags campaign stores
//	@Produce json
//	@Security ApiKeyAuth
//	@Param	campaign_id	path int true "Campaign ID"
//	@Success 200 {object} dto.Response
//	@Failure 400 {object} dto.Response
//	@Failure 404 {object} dto.Response
//	@Failure 500 {object} dto.Response
//	@Router	/campaigns/{campaign_id}/stores [delete]
func (c *CampaignStoreController) DeleteStores(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := util.GetUserID(ctx)
	if err != nil {
		dto.BadRequestJSON(w, r, err.Error())
		return
	}
	campaignID, err := strconv.Atoi(chi.URLParam(r, "campaign_id"))
	if err != nil {
		dto.BadRequestJSON(w, r, fmt.Sprintf(IncorrectCampaignIDErr, err.Error()))
		return
	}
	exists, err := c.campaignUseCases.Exists(ctx, int64(campaignID), "")
	if err != nil {
		dto.InternalServerErrorJSON(w, r, err.Error())
		return
	}

	if !exists {
		dto.BadRequestJSON(w, r, fmt.Sprintf(CampaignNotExistsErr, campaignID))
		return
	}

	err = c.campaignStoreUseCases.DeleteStores(ctx, int64(campaignID), int64(userID))
	if err != nil {
		dto.InternalServerErrorJSON(w, r, err.Error())
		return
	}

	dto.SuccessJSON(w, r, fmt.Sprintf("all campaign stores with campaign id %d deleted successfully", campaignID))
}

// DeleteStore godoc
//
//	@Summary Delete specified store with given store id under partilcular campaign
//	@Description API to delete particular store under specified campaign
//	@Tags campaign stores
//	@Produce json
//	@Security ApiKeyAuth
//	@Param	campaign_id	path int true "Campaign ID"
//	@Param	id	path int true "Campaign Store ID"
//	@Success 200 {object} dto.Response
//	@Failure 400 {object} dto.Response
//	@Failure 404 {object} dto.Response
//	@Failure 500 {object} dto.Response
//	@Router	/campaigns/{campaign_id}/stores/{id} [delete]
func (c *CampaignStoreController) DeleteStore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := util.GetUserID(ctx)
	if err != nil {
		dto.BadRequestJSON(w, r, err.Error())
		return
	}
	campaignID, err := strconv.Atoi(chi.URLParam(r, "campaign_id"))
	if err != nil {
		dto.BadRequestJSON(w, r, fmt.Sprintf(IncorrectCampaignIDErr, err.Error()))
		return
	}
	storeID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		dto.BadRequestJSON(w, r, fmt.Sprintf("incorrect store id value, err : %v", err.Error()))
		return
	}

	exists, err := c.campaignUseCases.Exists(ctx, int64(campaignID), "")
	if err != nil {
		dto.InternalServerErrorJSON(w, r, err.Error())
		return
	}

	if !exists {
		dto.BadRequestJSON(w, r, fmt.Sprintf(CampaignNotExistsErr, campaignID))
		return
	}

	err = c.campaignStoreUseCases.DeleteStore(ctx, int64(campaignID), int64(storeID), int64(userID))
	if err != nil {
		dto.InternalServerErrorJSON(w, r, err.Error())
		return
	}

	dto.SuccessJSON(w, r, fmt.Sprintf("store with id %d deleted successfully", storeID))
}

// AddStores godoc
//
//	@Summary add stores for specific campaign
//	@Description API to insert new stores under given campaign id
//	@Tags campaign stores
//	@Accept json
//	@Produce json
//	@Security ApiKeyAuth
//	@Param	campaign_id	path int true "Campaign ID"
//	@Param	stores body params.CampaignStoresForm true "Store Details"
//	@Success 200 {object} dto.CampaignStoresDTO
//	@Failure 400 {object} dto.Response
//	@Failure 409 {object} dto.Response
//	@Failure 500 {object} dto.Response
//	@Router	/campaigns/{campaign_id}/stores [post]
func (c *CampaignStoreController) AddStores(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := util.GetUserID(ctx)
	if err != nil {
		dto.BadRequestJSON(w, r, err.Error())
		return
	}
	campaignID, err := strconv.Atoi(chi.URLParam(r, "campaign_id"))
	if err != nil {
		dto.BadRequestJSON(w, r, fmt.Sprintf(IncorrectCampaignIDErr, err.Error()))
		return
	}

	request, err := c.validateStoresRequest(r)
	if err != nil {
		dto.BadRequestJSON(w, r, err.Error())
		return
	}

	var exists bool
	exists, err = c.campaignUseCases.Exists(ctx, int64(campaignID), "")
	if err != nil {
		dto.InternalServerErrorJSON(w, r, err.Error())
		return
	}

	if !exists {
		dto.BadRequestJSON(w, r, fmt.Sprintf(CampaignNotExistsErr, campaignID))
		return
	}

	stores, err := c.addStores(ctx, request, campaignID, int64(userID))
	if err != nil {
		dto.InternalServerErrorJSON(w, r, err.Error())
		return
	}

	response := dto.CampaignStoresDTO{
		CampaignID: int64(campaignID),
		Stores:     stores,
	}

	dto.SuccessJSONResponse(w, r, response)
}

func (c *CampaignStoreController) addStores(ctx context.Context, request params.CampaignStoresForm, campaignID int, userID int64) ([]*dto.CampaignStores, error) {
	storeEntities := []entities.CampaignStore{}
	for _, storeID := range request.Stores {
		storeEntities = append(storeEntities, params.ToCampaignStoreEntity(storeID, int64(campaignID), userID))
	}
	storeDetails, err := c.campaignStoreUseCases.AddStores(ctx, storeEntities)
	if err != nil {
		return nil, err
	}
	return storeDetails, nil
}

func (c *CampaignStoreController) validateStoresRequest(r *http.Request) (params.CampaignStoresForm, error) {
	var request params.CampaignStoresForm
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return request, err
	}
	defer r.Body.Close()

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		return request, err
	}

	return request, nil
}
