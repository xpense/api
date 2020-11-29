package handlers

import "errors"

var (
	// Generic
	ErrorEmptyBody = errors.New("empty body")
	// Auth/Account
	ErrorName                   = errors.New("first and/or last name missing")
	ErrorEmail                  = errors.New("invalid email address")
	ErrorMissingPasswordOrEmail = errors.New("both email and password are required for login")
	ErrorNonExistentUser        = errors.New("user with this email does not exist")
	ErrorEmailConflict          = errors.New("user with this email already exists")
	ErrorWrongPassword          = errors.New("wrong password")
	// Party
	ErrorPartyNameTaken = errors.New("party with the same name, belonging to the same user already exists")
	// Wallet
	ErrorWalletNameTaken = errors.New("wallet with the same name, belonging to the same user already exists")
	// Transaction
	ErrorRequiredAmount   = errors.New("cannot create new transaction with an amount of 0")
	ErrorRequiredWalletID = errors.New("a valid wallet id must be specified to register a new transaction")
	ErrorRequiredPartyID  = errors.New("a valid wallet id must be specified to register a new transaction")
	ErrorWalletNotFound   = errors.New("wallet with specified id not found")
	ErrorBadWalletID      = errors.New("wallet with specified id belongs to another user")
	ErrorPartyNotFound    = errors.New("party with specified id not found")
	ErrorBadPartyID       = errors.New("party with specified id belongs to another user")
)
