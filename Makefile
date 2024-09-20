-include .env
export

PERSONS_SERVICE_MIGRATIONS_DIR:="migrations/persons-service"

.PHONY: persons-service-migrate-up
persons-service-migrate-up:
	goose -dir $(PERSONS_SERVICE_MIGRATIONS_DIR) postgres "${PERSONS_SERVICE_POSTGRESQL_DSN}" up

.PHONY: persons-service-migrate-down
persons-service-migrate-down:
	goose -dir $(PERSONS_SERVICE_MIGRATIONS_DIR) postgres "${PERSONS_SERVICE_POSTGRESQL_DSN}" down

.PHONY: create-persons-service-migration
create-persons-service-migration:
ifeq ($(name),)
	@echo "You forgot to add migration name, example:\nmake create-migration name=create_users_table"
else
	goose -dir $(PERSONS_SERVICE_MIGRATIONS_DIR) create $(name) sql
endif