# Run e2e tests
.PHONY: e2e_tests unit_tests pr_reviewer_up
e2e_tests:
	docker compose down && \
	docker compose -f docker-compose.yml -f docker-compose.e2e.yml --env-file testing.env --profile runner up \
		--build --force-recreate -d --renew-anon-volumes testing >> build_log.txt && \
	docker compose -f docker-compose.yml -f docker-compose.e2e.yml --env-file testing.env logs -f testing && \
    docker compose -f docker-compose.yml -f docker-compose.e2e.yml --env-file testing.env --profile runner down -v --remove-orphans

unit_tests:
	go test ./...

pr_reviewer_up:
	docker compose up
