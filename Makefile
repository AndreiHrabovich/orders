tidy:
	go mod tidy
	go mod vendor

run:
	go run cmd/orders/main.go

bump:
	go get -u -v ./...

clean:
	go clean -modcache
	tidy