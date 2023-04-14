## setup
- `export PATH=$PATH:/usr/local/go/bin`
- `go install github.com/gofiber/fiber/v2@latest`
- `go mod init <app_name>`
- `go get github.com/gofiber/fiber/v2`

## run
- `go run main.go`

## connect to db:
- `docker exec -it mysql mysql -uappuser -psecretpassword`
- `USE app;` - choose database
- `SHOW TABLES;` - list tables
