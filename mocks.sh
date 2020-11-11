#!/bin/sh

# Repository
mockery --name Repository --filename another_repository_spy.go --dir repository --output router/test --outpkg test --structname RepositorySpy

# Password Hasher
mockery --name PasswordHasher --filename password_hasher_spy.go --dir utils --output router/test --outpkg test --structname PasswordHasherSpy