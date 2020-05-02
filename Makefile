# for local dev
test_bot:
	cd bot && go test -cover -coverprofile=coverage.out ./... && go tool cover -func coverage.out
coverage_bot:
	go tool cover -html bot/coverage.out
test_vk:
	cd vk && go test -cover -coverprofile=coverage.out -coverpkg=./  ./... && go tool cover -func coverage.out
coverage_vk:
	go tool cover -html vk/coverage.out
fmt_check:
	gofmt -l -s .

# for prod env
upd:
	docker-compose up -d
clear:
	docker-compose down -v
stop:
	docker-compose stop
