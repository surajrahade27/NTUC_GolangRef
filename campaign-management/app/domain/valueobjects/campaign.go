package valueobjects

type (
	CampaignID         int64
	CampaignProductID  int64
	CampaignStoreID    int64
	DailyTimeSlotID    int64
	SpecificTimeSlotID int64
	CampaignType       string
	CampaignStatusCode int64
)

const (
	CampaignTypePreOrder CampaignType = "preorder" // should it be deli ?
)

func (c CampaignID) ToInt64() int64 {
	return int64(c)
}

func (c CampaignProductID) ToInt64() int64 {
	return int64(c)
}

func (c CampaignStoreID) ToInt64() int64 {
	return int64(c)
}

func (c DailyTimeSlotID) ToInt64() int64 {
	return int64(c)
}

func (c SpecificTimeSlotID) ToInt64() int64 {
	return int64(c)
}
func (d CampaignType) String() string {
	return string(d)
}
