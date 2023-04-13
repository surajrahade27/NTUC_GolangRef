package http

import (
	"campaign-mgmt/app/domain/services"
	"campaign-mgmt/app/domain/valueobjects"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/usecases"
	"campaign-mgmt/app/usecases/dto"
	"campaign-mgmt/app/usecases/params"
	"campaign-mgmt/app/usecases/util"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type CampaignController struct {
	campaignUseCases        usecases.CampaignUseCases
	campaignStoreUseCases   usecases.CampaignStoreUseCases
	campaignProductUseCases usecases.CampaignProductUseCases
	tx                      services.TransactionService
	appConfig               *entities.AppCfg
}

func NewCampaignController(
	campaignUseCases usecases.CampaignUseCases,
	storeUseCases usecases.CampaignStoreUseCases,
	productUsecases usecases.CampaignProductUseCases,
	transactionService services.TransactionService,
	appConfig *entities.AppCfg) *CampaignController {
	return &CampaignController{
		campaignUseCases:        campaignUseCases,
		campaignStoreUseCases:   storeUseCases,
		campaignProductUseCases: productUsecases,
		tx:                      transactionService,
		appConfig:               appConfig,
	}
}

func (c *CampaignController) Init(r chi.Router) {
	r.Route("/campaigns", func(r chi.Router) {
		r.Post("/", c.CreateCampaign)
		r.Put("/{id}", c.UpdateCampaign)
		r.Get("/{id}", c.GetCampaign)
		r.Get("/", c.GetCampaignList)
		r.Put("/update-status", c.UpdateCampaignStatus)
	})
}

// GetCampaign godoc
//
//	@Summary Get campaign details by id
//	@Description API to get details of particular campaign
//	@Tags campaign
//	@Produce json
//	@Param	id	path int true "Campaign ID"
//	@Param	omit_products query boolean false "Omit Products"
//	@Param	omit_stores query boolean false "Omit Stores"
//	@Success 200 {object} dto.CampaignResponse
//	@Failure 400 {object} dto.Response
//	@Failure 404 {object} dto.Response
//	@Failure 500 {object} dto.Response
//	@Router	/campaigns/{id} [get]
func (c *CampaignController) GetCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	campaignID, campaignIDError := strconv.Atoi(chi.URLParam(r, "id"))
	if campaignIDError != nil {
		dto.BadRequestJSON(w, r, campaignIDError.Error())
	}
	omitProductsOptional := false
	omitStoresOptional := false
	if r.URL.Query().Get("omit_products") != "" {
		omitProducts, omitProductsErr := strconv.ParseBool(r.URL.Query().Get("omit_products"))
		if omitProductsErr != nil {
			dto.BadRequestJSON(w, r, omitProductsErr.Error())
		}
		omitProductsOptional = omitProducts
	}
	if r.URL.Query().Get("omit_stores") != "" {
		omitStores, omitStoresErr := strconv.ParseBool(r.URL.Query().Get("omit_stores"))
		if omitStoresErr != nil {
			dto.BadRequestJSON(w, r, omitStoresErr.Error())
		}
		omitStoresOptional = omitStores
	}

	response, CampaignDataErr := c.campaignUseCases.Get(ctx, int64(campaignID))
	if CampaignDataErr != nil {
		dto.InternalServerErrorJSON(w, r, "Error occured on get campaign details")
	}

	if !omitStoresOptional {
		storeDetails, storeDetailsErr := c.getStores(ctx, int64(campaignID))
		if storeDetailsErr != nil {
			dto.InternalServerErrorJSON(w, r, "Error occured on get store details: "+storeDetailsErr.Error())
		}
		response.CampaignStores = storeDetails
	}
	if !omitProductsOptional {
		productDetails, productDetailsErr := c.getProducts(ctx, int64(campaignID))
		if productDetailsErr != nil {
			dto.InternalServerErrorJSON(w, r, "Error occured on get product details : "+productDetailsErr.Error())
		}
		response.CampaignProducts = productDetails
	}

	render.JSON(w, r, dto.ToCampaignResponse(*response))
}

func (c *CampaignController) getStores(ctx context.Context, campaignID int64) ([]*dto.CampaignStores, error) {
	storeDetails, err := c.campaignStoreUseCases.GetStores(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	return storeDetails, nil
}

func (c *CampaignController) getProducts(ctx context.Context, campaignID int64) ([]*dto.CampaignProducts, error) {
	productDetails, err := c.campaignProductUseCases.GetProducts(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	return productDetails, nil
}

// CreateCampaign godoc
//
//	@Summary Create a campaign
//	@Description API to create new campaign
//	@Tags campaign
//	@Accept json
//	@Produce json
//	@Security ApiKeyAuth
//	@Param	campaign body params.CampaignCreationForm	true "Add campaign details"
//	@Success 200 {object} dto.CampaignDTO
//	@Failure 400 {object} dto.Response
//	@Failure 409 {object} dto.Response
//	@Failure 500 {object} dto.Response
//	@Router	/campaigns [post]
func (c *CampaignController) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var exists bool

	userID, err := util.GetUserID(ctx)
	if err != nil {
		dto.BadRequestJSON(w, r, err.Error())
		return
	}
	campaignRequest, err := c.validateCampaignRequest(r)
	if err != nil {
		dto.BadRequestJSON(w, r, err.Error())
		return
	}

	exists, err = c.campaignUseCases.Exists(ctx, 0, campaignRequest.Title)
	if err != nil {
		dto.InternalServerErrorJSON(w, r, err.Error())
		return
	}
	if exists {
		dto.ConflictErrorJSON(w, r, "campaign with given name already exists")
		return
	}

	response, err := c.create(ctx, campaignRequest, int64(userID))
	if err != nil {
		dto.InternalServerErrorJSON(w, r, err.Error())
		return
	}

	dto.SuccessJSONResponse(w, r, response)
}

func (c *CampaignController) create(ctx context.Context, request params.CampaignCreationForm, userID int64) (*dto.CampaignDTO, error) {
	response := &dto.CampaignDTO{}
	var err error
	c.tx.RunWithTransaction(
		ctx, func(ctx context.Context) error {
			if response, err = c.saveCampaignDetails(ctx, request, userID); err != nil {
				return err
			}
			return nil
		})

	return response, err
}

func (c *CampaignController) saveCampaignDetails(ctx context.Context, request params.CampaignCreationForm, userID int64) (*dto.CampaignDTO, error) {
	campaignEntity, err := params.ToCampaignEntity(request)
	campaignEntity.CreatedBy = userID
	if err != nil {
		return nil, err
	}
	campaignDetails, err := c.campaignUseCases.Create(ctx, campaignEntity)
	if err != nil {
		return nil, err
	}

	if len(request.Stores) != 0 {
		storesDetails, err := c.addStores(ctx, request.Stores, campaignDetails.ID, userID)
		if err != nil {
			return nil, err
		}
		campaignDetails.CampaignStores = storesDetails
	}

	return campaignDetails, nil
}

func (c *CampaignController) addStores(ctx context.Context, storeIDs []int64, campaignID, userID int64) ([]*dto.CampaignStores, error) {
	storeEntities := []entities.CampaignStore{}
	for _, storeID := range storeIDs {
		storeEntities = append(storeEntities, params.ToCampaignStoreEntity(storeID, campaignID, userID))
	}
	storeDetails, err := c.campaignStoreUseCases.AddStores(ctx, storeEntities)
	if err != nil {
		return nil, err
	}
	return storeDetails, nil
}

// func (c *CampaignController) addProducts(ctx context.Context, productRequest []params.CampaignProduct, campaignID, userID int64) ([]*dto.CampaignProducts, error) {
// 	productEntities := []entities.CampaignProduct{}
// 	for _, product := range productRequest {
// 		productEntities = append(productEntities, params.ToCampaignProductEntity(product, campaignID, userID))
// 	}
// 	productDetails, err := c.campaignProductUseCases.AddProducts(ctx, productEntities)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return productDetails, nil
// }

// UpdateCampaign godoc
//
//	@Summary Update campaign details
//	@Description API to update an existing campaign
//	@Tags campaign
//	@Accept json
//	@Produce json
//	@Security ApiKeyAuth
//	@Param	id	path int true "Campaign ID"
//	@Param	campaign body params.CampaignUpdateForm	true "campaign details"
//	@Success 200 {object} dto.Response
//	@Failure 400 {object} dto.Response
//	@Failure 409 {object} dto.Response
//	@Failure 500 {object} dto.Response
//	@Router	/campaigns/{id} [put]
func (c *CampaignController) UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := util.GetUserID(ctx)
	if err != nil {
		dto.BadRequestJSON(w, r, err.Error())
		return
	}

	campaignID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		dto.BadRequestJSON(w, r, fmt.Sprintf(IncorrectCampaignIDErr, err.Error()))
		return
	}

	campaignRequest, err := c.validateUpdateCampaignRequest(r)
	if err != nil {
		dto.BadRequestJSON(w, r, err.Error())
		return
	}

	exists, err := c.campaignUseCases.Exists(ctx, int64(campaignID), "")
	if err != nil {
		dto.InternalServerErrorJSON(w, r, err.Error())
		return
	}
	if !exists {
		dto.BadRequestJSON(w, r, "campaign with given id not exists")
		return
	}

	err = c.update(ctx, int64(campaignID), campaignRequest, int64(userID))
	if err != nil {
		dto.InternalServerErrorJSON(w, r, err.Error())
		return
	}
	dto.SuccessJSON(w, r, fmt.Sprintf("campaign with id %d updated successfully", campaignID))
}

func (c *CampaignController) update(ctx context.Context, campaignID int64, request params.CampaignUpdateForm, userID int64) error {
	var err error
	c.tx.RunWithTransaction(
		ctx, func(ctx context.Context) error {
			if err = c.updateCampaignDetails(ctx, campaignID, request, userID); err != nil {
				return err
			}
			return nil
		})

	return err
}

func (c *CampaignController) validateCampaignRequest(r *http.Request) (params.CampaignCreationForm, error) {
	var campaignRequest params.CampaignCreationForm
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&campaignRequest)
	if err != nil {
		return campaignRequest, err
	}
	defer r.Body.Close()

	validate := validator.New()
	err = validate.Struct(campaignRequest)
	if err != nil {
		return campaignRequest, err
	}

	dates := params.CampaignDates{
		OrderStartDate:      campaignRequest.OrderStartDate,
		OrderEndDate:        campaignRequest.OrderEndDate,
		CollectionStartDate: campaignRequest.CollectionStartDate,
		CollectionEndDate:   campaignRequest.CollectionEndDate,
	}

	err = c.validateCampaignDates(dates, campaignRequest.LeadTime)
	if err != nil {
		return campaignRequest, err
	}
	return campaignRequest, nil
}

func (c *CampaignController) validateUpdateCampaignRequest(r *http.Request) (params.CampaignUpdateForm, error) {
	var campaignRequest params.CampaignUpdateForm
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&campaignRequest)
	if err != nil {
		return campaignRequest, err
	}
	defer r.Body.Close()

	validate := validator.New()
	err = validate.Struct(campaignRequest)
	if err != nil {
		return campaignRequest, err
	}

	dates := params.CampaignDates{
		OrderStartDate:      campaignRequest.OrderStartDate,
		OrderEndDate:        campaignRequest.OrderEndDate,
		CollectionStartDate: campaignRequest.CollectionStartDate,
		CollectionEndDate:   campaignRequest.CollectionEndDate,
	}

	err = c.validateCampaignDates(dates, campaignRequest.LeadTime)
	if err != nil {
		return campaignRequest, err
	}
	return campaignRequest, nil
}

func (c *CampaignController) validateCampaignDates(dates params.CampaignDates, leadTime int) error {
	var orderStartDate, orderEndDate, collectionStartDate, collectionEndDate time.Time
	var err error
	orderStartDate, err = util.ToDateTime(dates.OrderStartDate)
	if err != nil {
		return err
	}
	orderEndDate, err = util.ToDateTime(dates.OrderEndDate)
	if err != nil {
		return err
	}
	collectionStartDate, err = util.ToDateTime(dates.CollectionStartDate)
	if err != nil {
		return err
	}
	collectionEndDate, err = util.ToDateTime(dates.CollectionEndDate)
	if err != nil {
		return err
	}
	if orderStartDate.After(orderEndDate) {
		return errors.New("invalid date : order start date should be the date before order end date")
	}
	if orderStartDate.After(collectionStartDate) {
		return errors.New("invalid date : order start date should be the date before collection start date")
	}
	if collectionStartDate.After(collectionEndDate) {
		return errors.New("invalid date : collection start date should be the date before collection end date")
	}

	err = c.checkStartDate(orderStartDate, collectionStartDate, leadTime)
	if err != nil {
		return err
	}
	err = c.checkEndDate(orderEndDate, collectionEndDate, leadTime)
	if err != nil {
		return err
	}
	return nil
}

func (c *CampaignController) checkStartDate(orderStartDate, collectionStartDate time.Time, leadTime int) error {
	if !orderStartDate.IsZero() && !collectionStartDate.IsZero() {
		if int64(collectionStartDate.Sub(orderStartDate).Hours()/24) < int64(leadTime) {
			return fmt.Errorf("invalid date : Collection start date should be at least %d days greater than order start date",
				leadTime)
		} else if int64(collectionStartDate.Sub(orderStartDate).Hours()/24) > int64(c.appConfig.ValidationParam.MaxDateDifference) {
			return fmt.Errorf("invalid date : Collection start date should be less than %d days from order start date",
				c.appConfig.ValidationParam.MaxDateDifference)
		}
	}
	return nil
}

func (c *CampaignController) checkEndDate(orderEndDate, collectionEndDate time.Time, leadTime int) error {
	if !orderEndDate.IsZero() && !collectionEndDate.IsZero() {
		if int64(collectionEndDate.Sub(orderEndDate).Hours()/24) < int64(leadTime) {
			return fmt.Errorf("invalid date : Collection end date should be at least %d days greater than order end date",
				leadTime)
		} else if int64(collectionEndDate.Sub(orderEndDate).Hours()/24) > int64(c.appConfig.ValidationParam.MaxDateDifference) {
			return fmt.Errorf("invalid date : Collection end date should be less than %d days from order end date",
				c.appConfig.ValidationParam.MaxDateDifference)
		}
	}
	return nil
}

func (c *CampaignController) updateCampaignDetails(ctx context.Context, campaignID int64, request params.CampaignUpdateForm, userID int64) error {
	campaignEntity, err := params.ToUpdateCampaignEntity(request, campaignID)
	campaignEntity.UpdatedBy = userID
	if err != nil {
		return err
	}
	err = c.campaignUseCases.Update(ctx, campaignEntity)
	if err != nil {
		return err
	}

	err = c.updateStores(ctx, campaignID, request.Stores, userID)
	if err != nil {
		return err
	}

	return nil
}

func (c *CampaignController) updateStores(ctx context.Context, campaignID int64, stores []int64, userID int64) error {
	var newStores []entities.CampaignStore
	for _, storeID := range stores {
		_, err := c.campaignStoreUseCases.GetByStoreID(ctx, campaignID, storeID)
		if err != nil && errors.As(err, &valueobjects.ErrStoreNotExists) {
			storeEntity := params.ToCampaignStoreEntity(storeID, campaignID, userID)
			newStores = append(newStores, storeEntity)
		} else if err != nil {
			return err
		}

		// this is additional code for future reference.
		// currently there are no additional fields to
		// update, that's why there is no need to update the store entity
		// var existingStores []entities.CampaignStore
		// _ = storeDetails
		// if storeDetails.ID > 0 {
		// 	store := params.UpdateCampaignStore{
		// 		ID:      storeDetails.ID,
		// 		StoreID: storeDetails.StoreID,
		// 	}
		// 	existingStoreEntity := params.ToUpdateCampaignStoreEntity(store, campaignID, userID)
		// 	existingStores = append(existingStores, existingStoreEntity)
		// }

	}

	// this is additional code for future reference. currently there are no additional fields to
	// update, that's why there is no need to update the store entity
	// if len(existingStores) > 0 {
	// 	err := c.campaignStoreUseCases.UpdateStores(ctx, existingStores)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	dbStoreIDS := []int64{}
	dbStores, err := c.campaignStoreUseCases.GetStores(ctx, campaignID)
	if err != nil {
		return err
	}
	for _, store := range dbStores {
		dbStoreIDS = append(dbStoreIDS, store.StoreID)
	}
	storesToDelete := util.Difference(dbStoreIDS, stores)

	if len(newStores) > 0 {
		_, err := c.campaignStoreUseCases.AddStores(ctx, newStores)
		if err != nil {
			return err
		}
	}

	for _, id := range storesToDelete {
		err := c.campaignStoreUseCases.DeleteByStoreID(ctx, campaignID, id, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetCampaignList godoc
//
//	@Summary Get list of all campaigns
//	@Description API to get details of all campaigns
//	@Tags campaign
//	@Produce json
//	@Param	page query int false "Page Number"
//	@Param	limit query int false "Limit"
//	@Param	sort query string false "Sort Type [created_at asc/created_at desc]"
//	@Param	name query string false "Campaign Name"
//	@Param	status query string false "Campaign Status [InActive/Active/Scheduled]"
//	@Success 200 {object} dto.CampaignListResponse
//	@Failure 400 {object} dto.Response
//	@Failure 404 {object} dto.Response
//	@Failure 500 {object} dto.Response
//	@Router	/campaigns [get]
func (c *CampaignController) GetCampaignList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pagination := c.generatePaginationFromRequest(r)
	paginationData := params.ToPaginationEntity(pagination)
	response, err := c.campaignUseCases.GetList(ctx, paginationData)
	if err != nil {
		dto.InternalServerErrorJSON(w, r, err.Error())
		return
	}
	render.JSON(w, r, response)
}

func (c *CampaignController) generatePaginationFromRequest(r *http.Request) params.Pagination {
	// Initializing default
	paginationConfig := c.appConfig.PaginationConfig
	limit := paginationConfig.Limit
	page := paginationConfig.Page
	sort := paginationConfig.Sort
	query := r.URL.Query()
	name := paginationConfig.Name
	var statusCode int64

	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
			break
		case "page":
			page, _ = strconv.Atoi(queryValue)
			break
		case "sort":
			sort = queryValue
			break
		case "name":
			name = queryValue
		case "status":
			switch queryValue {
			case "InActive":
				statusCode = valueobjects.CampaignStatusInActive.Code()
				break
			case "Active":
				statusCode = valueobjects.CampaignStatusActive.Code()
				break
			case "Scheduled":
				statusCode = valueobjects.CampaignStatusScheduled.Code()
				break
			}
			break
		}
	}

	return params.Pagination{
		Limit:  limit,
		Page:   page,
		Sort:   sort,
		Name:   name,
		Status: statusCode,
	}
}

// UpdateCampaignStatus godoc
//
//	@Summary Update status of campaign
//	@Description API to update the status of campaign
//	@Tags campaign
//	@Produce json
//	@Security ApiKeyAuth
//	@Success 200 {object} dto.Response
//	@Failure 400 {object} dto.Response
//	@Failure 409 {object} dto.Response
//	@Failure 500 {object} dto.Response
//	@Router	/campaigns/update-status [put]
func (c *CampaignController) UpdateCampaignStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := c.updateStatus(ctx)
	if err != nil {
		dto.InternalServerErrorJSON(w, r, err.Error())
		return
	}
	dto.SuccessJSON(w, r, "campaigns status updated successfully")
}

func (c *CampaignController) updateStatus(ctx context.Context) error {
	var err error
	c.tx.RunWithTransaction(
		ctx, func(ctx context.Context) error {
			if err = c.campaignUseCases.UpdateStatus(ctx); err != nil {
				return err
			}
			return nil
		})
	return err
}
