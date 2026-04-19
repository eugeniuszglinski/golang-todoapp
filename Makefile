include .env
export


export PROJECT_ROOT=$(shell pwd)

# target "env-up"
env-up:
	@docker compose up -d todoapp-postgres

# target "env-up"
env-down:
	@docker compose down todoapp-postgres

# @ makes you commands invisible in terminal
env-cleanup:
	@read -p "Do you want to cleanup all volume files in an environment? There is a risk of losing data. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
	  docker compose down todoapp-postgres port-forwarder && \
	  rm -rf ${PROJECT_ROOT}/out/pgdata && \
	  echo "Environment cleanup successfully completed"; \
	else \
	  echo "Environment cleanup was rejected"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

# for long running services use "compose up", otherwise use "compose run"
# $(seq) takes argument given to MAKEFILE target
migrate-create:
	@if [ -z "$(seq)" ]; then \
  		echo "missing required parameter: seq. Example: make migrate-create seq=init"; \
  		exit 1; \
  	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "missing required parameter: action. Example: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

log-cleanup:
	@read -p "Do you want to cleanup all log files? There is a risk of losing logs. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
	  rm -rf ${PROJECT_ROOT}/out/${LOGGER_FOLDER} && \
	  echo "Logs cleanup successfully completed"; \
	else \
	  echo "Logs cleanup was rejected"; \
	fi

todoapp-run:
	@export LOGGER_FOLDER=$(PROJECT_ROOT)/out/${LOGGER_FOLDER} && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/todoapp/main.go

todoapp-deploy:
	docker compose up -d --build todoapp

ps:
	@docker compose ps