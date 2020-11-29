package auth

type (
	SignUpInfo struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	LoginInfo struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginToken struct {
		Token string `json:"token"`
	}
)
