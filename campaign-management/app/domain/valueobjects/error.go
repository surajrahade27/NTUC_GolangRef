package valueobjects

// Error is base error type for domain
type Error string

// ErrNotFound ..
const ErrNotFound Error = "not found"

// ErrValidation ..
const ErrValidation Error = "validation fails"

// Error ..
func (e Error) Error() string {
	return string(e)
}

var (
	ErrCampaignCantGet          Error = "unable to get campaign"
	ErrCampaignCantUpdate       Error = "unable to update campaign"
	ErrCampaignCantCreate       Error = "unable to crate campaign"
	ErrCampaignCantExist        Error = "unable to check existence of campaign"
	ErrCampaignCantGetList      Error = "unable to get campaign list"
	ErrProductCantCreate        Error = "unable to create product(s)"
	ErrProductCantUpdate        Error = "unable to update product(s)"
	ErrStoreCantCreate          Error = "unable to create store(s)"
	ErrStoreCantUpdate          Error = "unable to update store(s)"
	ErrStoreCantDelete          Error = "unable to delete store(s)"
	ErrProductCantDelete        Error = "unable to delete product"
	ErrCampaignStatusCantUpdate Error = "unable to update campaign status"
	ErrStoreCantGet             Error = "unable to get campaign store"
	ErrStoreNotExists           Error = "campaign store not exists"
)
