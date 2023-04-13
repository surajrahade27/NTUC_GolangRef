package params

import (
	"campaign-mgmt/app/domain/entities"
	"campaign-mgmt/app/domain/valueobjects"
)

type CampaignProductCreationForm struct {
	//campaign Id
	CampaignID int64 `json:"campaign_id" validate:"required"`
	// List of campaign products
	Products []CampaignProduct `json:"products" validate:"required,dive"`
	// Created by user id
	CreatedBy int64 `json:"created_by" validate:"required"`
}

type CampaignProduct struct {
	ProductID   int64  `json:"product_id" validate:"required"`
	SKUNo       int64  `json:"SKU_no"`
	SerialNo    int    `json:"serial_no"`
	SequenceNo  int    `json:"sequence_no"`
	ProductType string `json:"product_type" validate:"omitempty,oneof=cd ncd"`
}

type UpdateCampaignProduct struct {
	ID          int64  `json:"campaign_product_id"`
	ProductID   int64  `json:"product_id"`
	SKUNo       int64  `json:"SKU_no"`
	SerialNo    int    `json:"serial_no"`
	SequenceNo  int    `json:"sequence_no"`
	ProductType string `json:"product_type" validate:"omitempty,oneof=cd ncd"`
}

func ToCampaignProductEntity(product CampaignProduct, campaignID, userID int64) entities.CampaignProduct {
	return entities.CampaignProduct{
		CampaignID:  valueobjects.CampaignID(campaignID),
		ProductID:   product.ProductID,
		SKUNo:       product.SKUNo,
		SerialNo:    product.SerialNo,
		SequenceNo:  product.SequenceNo,
		ProductType: product.ProductType,
		CreatedBy:   userID,
	}
}

func ToUpdateCampaignProductEntity(product UpdateCampaignProduct, campaignID, userID int64) entities.CampaignProduct {
	return entities.CampaignProduct{
		ID:          valueobjects.CampaignProductID(product.ID),
		CampaignID:  valueobjects.CampaignID(campaignID),
		ProductID:   product.ProductID,
		SKUNo:       product.SKUNo,
		SerialNo:    product.SerialNo,
		SequenceNo:  product.SequenceNo,
		ProductType: product.ProductType,
		UpdatedBy:   userID,
	}
}
