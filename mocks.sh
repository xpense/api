#!/bin/sh

# Repository
mockery --name Repository --filename repository_spy.go --dir internal/repository --output test/spies --outpkg spies --structname RepositorySpy

# Password Hasher
mockery --name PasswordHasher --filename password_hasher_spy.go --dir internal/utils --output test/spies --outpkg spies --structname PasswordHasherSpy

# JWT Service
mockery --name JWTService --filename jwt_service_spy.go --dir internal/middleware/auth --output test/spies --outpkg spies --structname JWTServiceSpy