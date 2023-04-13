package dto

import "campaign-mgmt/app/domain/entities"

type CampaignProductsDTO struct {
	CampaignID int64              `json:"campaign_id"`
	Products   []CampaignProducts `json:"products"`
}

type CampaignProducts struct {
	ID          int64  `json:"campaign_product_id"`
	ProductID   int64  `json:"product_id"`
	SKUNo       int64  `json:"sku_no"`
	SerialNo    int    `json:"serial_no"`
	SequenceNo  int    `json:"sequence_no"`
	ProductType string `json:"product_type"`
}

func ToCampaignProductDTO(productEntity entities.CampaignProduct) *CampaignProducts {
	return &CampaignProducts{
		ID:          productEntity.ID.ToInt64(),
		ProductID:   productEntity.ProductID,
		SKUNo:       productEntity.SKUNo,
		SerialNo:    productEntity.SerialNo,
		SequenceNo:  productEntity.SequenceNo,
		ProductType: productEntity.ProductType,
	}
}
