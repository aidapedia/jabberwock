# initializes the project
init:
	@chmod +x ./scripts/init.sh
	@./scripts/init.sh

# runs wire
wire:
	@./scripts/wire.sh

# runs database migrations
database-up:
	go run ./pkg/migration/main.go -method=up

database-down:
	go run ./pkg/migration/main.go -method=down

database-reset:
	go run ./pkg/migration/main.go -method=reset

generate_key:
	@chmod +x ./scripts/generate_private_key.sh
	@./scripts/generate_private_key.sh