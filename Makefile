MIGRATE = $(HOME)/go/bin/migrate
MIGRATIONS_DIR = internal/db/migrations
DATABASE_URL ?= postgres://doheem_dev_user:simple_pswd@localhost:5432/doheem_dev_db?sslmode=disable

.PHONY: migrate-up migrate-down migrate-create migrate-version

migrate-up:
	$(MIGRATE) -source file://$(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-down:
	$(MIGRATE) -source file://$(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down $(N)

migrate-create:
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) $(NAME)

migrate-version:
	$(MIGRATE) -source file://$(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version
