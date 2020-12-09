package handlers

type ErrorMessage struct {
	Message string `json:"message"`
}

var (
	// Generic
	ErrorEmptyBody = &ErrorMessage{Message: "empty body"}
	// Auth/Account
	ErrorName                   = &ErrorMessage{Message: "first and/or last name missing"}
	ErrorEmail                  = &ErrorMessage{Message: "invalid email address"}
	ErrorMissingPasswordOrEmail = &ErrorMessage{Message: "both email and password are required for login"}
	ErrorNonExistentUser        = &ErrorMessage{Message: "user with this email does not exist"}
	ErrorEmailConflict          = &ErrorMessage{Message: "user with this email already exists"}
	ErrorWrongPassword          = &ErrorMessage{Message: "wrong password"}
	// Party
	ErrorPartyNameTaken = &ErrorMessage{Message: "party with the same name, belonging to the same user already exists"}
	// Wallet
	ErrorWalletNameTaken = &ErrorMessage{Message: "wallet with the same name, belonging to the same user already exists"}
	// Transaction
	ErrorRequiredAmount   = &ErrorMessage{Message: "cannot create new transaction with an amount of 0"}
	ErrorRequiredWalletID = &ErrorMessage{Message: "a valid wallet id must be specified to register a new transaction"}
	ErrorRequiredPartyID  = &ErrorMessage{Message: "a valid wallet id must be specified to register a new transaction"}
	ErrorWalletNotFound   = &ErrorMessage{Message: "wallet with specified id not found"}
	ErrorBadWalletID      = &ErrorMessage{Message: "wallet with specified id belongs to another user"}
	ErrorPartyNotFound    = &ErrorMessage{Message: "party with specified id not found"}
	ErrorBadPartyID       = &ErrorMessage{Message: "party with specified id belongs to another user"}
)
