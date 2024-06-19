deploy-dev:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/tgstat cmd/tgstat/main.go
	ssh tgstat "systemctl stop tgstat-api.service"
	scp bin/tgstat "tgstat:/var/www/tgstat/api/bin"
	ssh tgstat "systemctl start tgstat-api.service"
