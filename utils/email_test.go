package utils_test

import (
	"testing"

	"expense-api/utils"
)

func TestIsEmailValid(t *testing.T) {
	t.Run("Valid email cases", func(t *testing.T) {
		testCases := []struct {
			desc  string
			email string
		}{
			{
				desc:  "Username is a single alpha char",
				email: "a@domain.com",
			},
			{
				desc:  "Username is a numerical char",
				email: "1@domain.com",
			},
			{
				desc:  "Username is a single special char",
				email: "!@domain.com",
			},
			{
				desc:  "Username is consisting of alpha-numerical and speical chars",
				email: "!test123_@domain.com",
			},
			{
				desc:  "Domain is a single alpha char",
				email: "!test123_@a",
			},
			{
				desc:  "Domain is a single alpha char followed by dot",
				email: "!test123_@a.com",
			},
			{
				desc:  "Domain is a single alpha char followed by dot",
				email: "!test123_@a.net",
			},
			{
				desc:  "Valid email",
				email: "test123@test.com",
			},
		}

		for _, tC := range testCases {
			t.Run(tC.desc, func(t *testing.T) {
				if !utils.IsEmailValid(tC.email) {
					t.Errorf("Email: %s should be valid", tC.email)
				}
			})
		}
	})

	t.Run("Invalid email cases", func(t *testing.T) {
		testCases := []struct {
			desc  string
			email string
		}{
			{
				desc:  "Empty string",
				email: "",
			},
			{
				desc:  "Missing @",
				email: "asd",
			},
			{
				desc:  "Missing username",
				email: "@asd",
			},
			{
				desc:  "Missing domain",
				email: "asd@",
			},
			{
				desc:  "Only @",
				email: "@",
			},
			{
				desc:  "Invalid domain",
				email: "test@!",
			},
			{
				desc:  "Invalid domain",
				email: "test@!?.com",
			},
			{
				desc:  "Invalid domain",
				email: "test@!?.{}+=",
			},
		}

		for _, tC := range testCases {
			t.Run(tC.desc, func(t *testing.T) {
				if utils.IsEmailValid(tC.email) {
					t.Errorf("Email: %s shouldn't be valid", tC.email)
				}
			})
		}
	})
}
