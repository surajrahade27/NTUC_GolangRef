package params

import (
	"campaign-mgmt/app/domain/entities"
)

type Pagination struct {
	Limit  int    `json:"limit"`
	Page   int    `json:"page"`
	Sort   string `json:"sort"`
	Name   string `json:"name"`
	Status int64  `json:"status"`
}

func ToPaginationEntity(paginationData Pagination) entities.PaginationConfig {
	return entities.PaginationConfig{
		Limit:  paginationData.Limit,
		Sort:   paginationData.Sort,
		Page:   paginationData.Page,
		Name:   paginationData.Name,
		Status: paginationData.Status,
	}
}
