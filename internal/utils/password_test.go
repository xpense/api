package utils_test

import (
	"expense-api/internal/utils"
	"testing"
)

func TestIsPasswordStrong(t *testing.T) {
	t.Run("Weak password cases", func(t *testing.T) {
		testCases := []struct {
			desc     string
			password string
			err      error
		}{
			{
				desc:     "Doesn't contain lower char",
				password: "123WHATEVER!@{}",
				err:      utils.ErrorPasswordLowerChar,
			},
			{
				desc:     "Doesn't contain upper char",
				password: "123whatever!@{}",
				err:      utils.ErrorPasswordUpperChar,
			},
			{
				desc:     "Doesn't contain numerical char",
				password: "Whatever!{}",
				err:      utils.ErrorPasswordDigitChar,
			},
			{
				desc:     "Doesn't contain special char",
				password: "123Whatever",
				err:      utils.ErrorPasswordSpecialChar,
			},
			{
				desc:     "Valid password is too short",
				password: "1W!",
				err:      utils.ErrorPasswordLength,
			},
		}

		for _, tC := range testCases {
			t.Run(tC.desc, func(t *testing.T) {
				_, got := utils.IsPasswordStrong(tC.password)
				want := tC.err

				if got != want {
					t.Errorf("Got error: '%s'; Want error: '%s'", got.Error(), want.Error())
				}
			})
		}
	})

	t.Run("Strong password cases", func(t *testing.T) {
		testCases := []struct {
			desc     string
			password string
		}{
			{
				desc:     "Valid password containing lower, upper, numerical and special char of len 8",
				password: "1Whateve!",
			},
			{
				desc:     "Valid password containing lower, upper, numerical and special char with a len > 8",
				password: "123Whatever{}!",
			},
		}

		for _, tC := range testCases {
			t.Run(tC.desc, func(t *testing.T) {
				strong, err := utils.IsPasswordStrong(tC.password)

				if !strong {
					t.Errorf("Password: '%s' should be strong enough", tC.password)
				}

				if err != nil {
					t.Errorf("Password: '%s' should be strong enough; Error Message: '%s'", tC.password, err.Error())
				}
			})
		}
	})
}
