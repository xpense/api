#!/bin/sh

# Repository
mockery --name Repository --filename another_repository_spy.go --dir repository --output router/test/spies --outpkg spies --structname RepositorySpy

# Password Hasher
mockery --name PasswordHasher --filename password_hasher_spy.go --dir utils --output router/test/spies --outpkg spies --structname PasswordHasherSpy