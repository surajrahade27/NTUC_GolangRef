package valueobjects

// CampaignStatusCode ..
func (c CampaignStatusCode) Code() int64 {
	return int64(c)
}

var (
	CampaignStatusInActive  CampaignStatusCode = 1
	CampaignStatusActive    CampaignStatusCode = 2
	CampaignStatusScheduled CampaignStatusCode = 3
)
