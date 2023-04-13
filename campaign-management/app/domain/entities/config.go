package entities

type AppCfg struct {
	MYSQLConfig      MYSQLConfig
	PaginationConfig PaginationConfig
	ValidationParam  ValidationParam
}

type MYSQLConfig struct {
	User            string
	Password        string
	Host            string
	Port            string
	Database        string
	DBRetryAttempts int
}

type ValidationParam struct {
	MaxLeadTime       int
	MaxDateDifference int
}

type PaginationConfig struct {
	Limit  int
	Page   int
	Offset int
	Sort   string
	Name   string
	Status int64
}
