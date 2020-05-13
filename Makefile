# for local dev
test_bot:
	docker-compose -f docker-compose.test.yml up -d bot_test && docker attach khti_timetable_bot_test && docker-compose -f docker-compose.test.yml down
coverage_bot:
	go tool cover -html bot/coverage.out
test_vk:
	docker-compose -f docker-compose.test.yml up -d vk_library_test && docker attach khti_timetable_bot_test_vk_library && docker-compose -f docker-compose.test.yml down
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
