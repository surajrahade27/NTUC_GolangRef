# Campaign Management

### Environment Variables

- export DB_USER=username
- export DB_PASSWORD=password
- export DB_HOST=hostname
- export DB_NAME=campaign_management
- export DB_PORT=5432

### Set Environment Variables
```
source <sh file path>
```
### Code Coverage
run following to generate html doc for code coverage
```
go test -coverprofile=coverage.out ./... ; go tool cover -html=coverage.out
```

### Run following to Generating Mocks 
- go install github.com/vektra/mockery/v2@latest
- mockery --version
- //go:generate mockery --name Campaigns --filename < filename.go >

### Swagger
run following to generate Swagger documentation
```
swag init -g cmd/api/main.go
```

### Web URL for swagger documentation (local)
```
http://{host_name}/swagger/index.html
```