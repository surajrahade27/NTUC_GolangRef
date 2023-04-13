package params

import (
	"campaign-mgmt/app/domain/entities"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_ToPaginationEntity(t *testing.T) {
	t.Run("test conversion : success scenario", func(t *testing.T) {
		request := Pagination{
			Limit:  20,
			Page:   1,
			Sort:   "asc",
			Name:   "campaign",
			Status: 1,
		}
		expectedResponse := entities.PaginationConfig{
			Limit:  20,
			Page:   1,
			Sort:   "asc",
			Name:   "campaign",
			Status: 1,
		}

		response := ToPaginationEntity(request)
		ShouldEqual(response, expectedResponse)
	})
}
