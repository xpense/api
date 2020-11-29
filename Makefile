.PHONY: run, start-db, stop-db, generate-mocks, test

run:
	air

start-db:
	docker-compose -f docker-compose.dev.yml up -d

stop-db:
	docker-compose -f docker-compose.dev.yml down

start-test-db:
	docker-compose -f docker-compose.test.yml up -d

stop-test-db:
	docker-compose -f docker-compose.test.yml down

generate-mocks:
	# Repository
	mockery --name Repository --filename repository_spy.go --dir internal/repository --output test/spies --outpkg spies --structname RepositorySpy
	# Password Hasher
	mockery --name PasswordHasher --filename password_hasher_spy.go --dir internal/utils --output test/spies --outpkg spies --structname PasswordHasherSpy
	# JWT Service
	mockery --name JWTService --filename jwt_service_spy.go --dir internal/middleware/auth --output test/spies --outpkg spies --structname JWTServiceSpy

test:
	go test ./... -short

test-all:
	godotenv go test ./...

test-integration:
	godotenv go test ./test/integration