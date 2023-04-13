package main

import (
	"campaign-mgmt/app/domain/entities"
	repo "campaign-mgmt/app/infrastructure/mysql"
	"campaign-mgmt/app/middlewares"
	presentation "campaign-mgmt/app/presentation/http"
	"campaign-mgmt/app/usecases"
	"campaign-mgmt/docs"
	_ "campaign-mgmt/docs"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	logger "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlRepoServices struct {
	CampaignRepoService          *repo.CampaignService
	CampaignProductRepoService   *repo.CampaignProductService
	CampaignStoreRepoService     *repo.CampaignStoreService
	CampaignCodeService          *repo.CampaignCodeService
	StoreDailyTimeSlotService    *repo.StoreDailyTimeSlotService
	StoreSpecificTimeSlotService *repo.StoreSpecificTimeSlotService
	TransactionService           *repo.TransactionService
}

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @title           Campaign Management Swagger API
// @version         2.0
// @description     This is a Campaign Management Server, Which provides set of APIs to create and manage campaigns, for adding products and stores against particular campaign and also APIs for time slot management.
func main() {
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(false)
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASSWORD", "Password@1")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_NAME", "campaigns")
	os.Setenv("OKTA_AUDIENCE", "DBP_Back_Office")
	os.Setenv("OKTA_ISSUER", "https://ntucenterprise.oktapreview.com/oauth2/aus11bm2mb1uNA9FN1d7")
	logger.Info("Starting campaign management service..")

	conf, err := InitConfig()
	if err != nil {
		logger.Fatalf("Unable to initialise config, err : %v", err)
	}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.OktaAuthenticator)

	r.Use(cors.Handler(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	db, err := newDBConnection(conf.MYSQLConfig)
	if err != nil {
		logger.Fatalf("Error occurred while initiating database connection : %v", err)
	}

	var defaultCampaignStatusDBEntry = []map[string]interface{}{
		{"StatusCode": 1, "StatusValue": "InActive"},
		{"StatusCode": 2, "StatusValue": "Active"},
		{"StatusCode": 3, "StatusValue": "Scheduled"},
	}
	repos := registerRepoServices(db, defaultCampaignStatusDBEntry)

	campaignUseCase := usecases.NewCampaignUseCase(repos.CampaignRepoService)
	storeUseCase := usecases.NewCampaignStoreUseCase(repos.CampaignStoreRepoService)
	productUseCase := usecases.NewCampaignProductUseCase(repos.CampaignProductRepoService)

	campaignHandler := presentation.NewCampaignController(campaignUseCase, storeUseCase, productUseCase, repos.TransactionService, conf)
	campaignHandler.Init(r)
	productHandler := presentation.NewCampaignProductController(productUseCase, repos.TransactionService, conf)
	productHandler.Init(r)
	storeHandler := presentation.NewCampaignStoreController(campaignUseCase, storeUseCase)
	storeHandler.Init(r)

	logger.Info("Campaign management server started")
	logger.Info("visit http://localhost:8080/swagger/index.html  for swagger documentation")
	http.ListenAndServe(":8080", r)
}

func InitConfig() (*entities.AppCfg, error) {
	var conf entities.AppCfg

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Campaign Management Swagger API"
	docs.SwaggerInfo.Description = "This is a Campaign Management Server," +
		" Which provides set of APIs to create and manage campaigns, for adding products and stores " +
		"against particular campaign and also APIs for time slot management. "
	docs.SwaggerInfo.Version = "2.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	conf.MYSQLConfig = entities.MYSQLConfig{
		User:            os.Getenv("DB_USER"),
		Password:        os.Getenv("DB_PASSWORD"),
		Host:            os.Getenv("DB_HOST"),
		Port:            os.Getenv("DB_PORT"),
		Database:        os.Getenv("DB_NAME"),
		DBRetryAttempts: 3,
	}
	conf.PaginationConfig = entities.PaginationConfig{
		Limit:  20,
		Offset: 0,
		Page:   1,
		Sort:   "created_at asc",
		Name:   "",
	}

	conf.ValidationParam = entities.ValidationParam{
		MaxLeadTime:       20,
		MaxDateDifference: 28,
	}
	return &conf, nil
}

func newDBConnection(mysqlConf entities.MYSQLConfig) (*gorm.DB, error) {
	var err error
	var connection *gorm.DB
	for i := 0; i < mysqlConf.DBRetryAttempts; i++ {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			mysqlConf.User,
			mysqlConf.Password,
			mysqlConf.Host,
			mysqlConf.Port,
			mysqlConf.Database,
		)

		// tlsConf := createTLSConf()
		// err = gomysql.RegisterTLSConfig("custom", &tlsConf)
		// if err != nil {
		// 	logger.Fatalf("Error %s when RegisterTLSConfig\n", err)
		// 	return nil, err
		// }

		connection, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			// Get generic database object sql.DB to use its functions
			sqlDB, err := connection.DB()
			if err == nil {
				// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
				sqlDB.SetMaxIdleConns(10)

				// SetMaxOpenConns sets the maximum number of open connections to the database.
				sqlDB.SetMaxOpenConns(100)

				// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
				sqlDB.SetConnMaxLifetime(time.Hour)
			}
			return connection, nil
		}
		time.Sleep(time.Millisecond * 500)
	}
	return nil, err
}

func registerRepoServices(db *gorm.DB, defaultCampaignStatusDBEntry []map[string]interface{}) *MysqlRepoServices {
	var repos MysqlRepoServices
	repos.CampaignRepoService = repo.NewCampaignService(db)
	if err := repos.CampaignRepoService.Migrate(); err != nil {
		logger.Fatal(err)
	}
	repos.CampaignProductRepoService = repo.NewCampaignProductService(db)
	if err := repos.CampaignProductRepoService.Migrate(); err != nil {
		logger.Fatal(err)
	}
	repos.CampaignStoreRepoService = repo.NewCampaignStoreService(db)
	if err := repos.CampaignStoreRepoService.Migrate(); err != nil {
		logger.Fatal(err)
	}

	repos.CampaignCodeService = repo.NewCampaignCodeService(db)
	if err := repos.CampaignCodeService.Migrate(defaultCampaignStatusDBEntry); err != nil {
		logger.Fatal(err)
	}
	repos.StoreDailyTimeSlotService = repo.NewStoreDailyTimeSlotService(db)
	if err := repos.StoreDailyTimeSlotService.Migrate(); err != nil {
		logger.Fatal(err)
	}
	repos.StoreSpecificTimeSlotService = repo.NewStoreSpecificTimeSlotService(db)
	if err := repos.StoreSpecificTimeSlotService.Migrate(); err != nil {
		logger.Fatal(err)
	}
	repos.TransactionService = repo.NewTransactionService(db)
	return &repos
}

func createTLSConf() tls.Config {
	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile("./server-ca.pem")
	if err != nil {
		logger.Fatal(err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		logger.Fatal("Failed to append PEM.")
	}
	clientCert := make([]tls.Certificate, 0, 1)

	certs, err := tls.LoadX509KeyPair("./client-cert.pem", "./client-key.pem")
	if err != nil {
		logger.Fatal(err)
	}

	clientCert = append(clientCert, certs)

	return tls.Config{
		RootCAs:            rootCertPool,
		Certificates:       clientCert,
		InsecureSkipVerify: true,
	}
}
