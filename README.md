<div align="center">
  <img
    alt="sliced piece of pie"
    src="./assets/slice-it.jpg"
    height="250px"
  />
</div>
<h1 align="center">Welcome to the Slice It API 👋</h1>
<p align="center">
  <a href="https://golang.org/dl" target="_blank">
    <img alt="Using go version 1.14" src="https://img.shields.io/badge/go-1.14-9cf.svg" />
  </a>
  <a href="https://travis-ci.com/bradford-hamilton/slice-it-api" target="_blank">
    <img alt="Using go version 1.14" src="https://travis-ci.com/bradford-hamilton/slice-it-api.svg?branch=master" />
  </a>
  <a href="https://goreportcard.com/report/github.com/bradford-hamilton/slice-it-api" target="_blank">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/bradford-hamilton/slice-it-api/pkg" />
  </a>
  <a href="https://godoc.org/github.com/bradford-hamilton/slice-it-api/pkg" target="_blank">
    <img alt="godoc" src="https://godoc.org/github.com/bradford-hamilton/slice-it-api/pkg?status.svg" />
  </a>
  <a href="#" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
</p>

## Dependencies
- Be sure to have Go (1.14+)
- Be sure to have postgres running locally
___
## Database
- Set up development db with
  ```
  createdb slice_it_api_dev
  ```
- Run migration:
  ```
  psql -U slice_it_user -d slice_it_api_dev -a -f internal/storage/migrations/schema.sql
  ```
___
## Usage
### Development
```
go run cmd/server/main.go
```
___
## Testing
Standard:
```
go test ./...
```

For some nice color output:
```
make test
```
___
## Deployment

Build docker image after code changes:
```
docker build \
  --build-arg SLICE_IT_API_SERVER_PORT={server_port} \
  --build-arg SLICE_IT_API_DB_HOST={host} \
  --build-arg SLICE_IT_API_DB_PORT={db_port} \
  --build-arg SLICE_IT_API_DB_USER={db_username} \
  --build-arg SLICE_IT_API_DB_PASSWORD={db_password} \
  --build-arg SLICE_IT_API_DB_NAME={db_name} \
  --build-arg SLICE_IT_API_SSL_MODE={ssl_mode} \
  -t bradfordhamilton/slice-it-api:latest .
```

Push image:
```
docker push bradfordhamilton/slice-it-api:latest
```

Then cd into `build/terraform` and run the following commands.
```
terraform plan -out=tfplan -input=false .
```
```
terraform apply -input=false tfplan
```
