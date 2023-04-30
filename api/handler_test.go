package api

import "testing"

func TestIsValidLogin(t *testing.T) {
	tests := []struct {
		name  string
		pass  string
		valid bool
	}{
		{
			"NoCharacter",
			"",
			false,
		},
		{
			"CyrillicCharacter",
			"Гофер",
			false,
		},
		{
			"EmptyStringAndWhitespace",
			" \n\t\r\v\f ",
			false,
		},
		{
			"InvalidCharacter",
			"Uua?aaaa",
			false,
		},
		{
			"LongName",
			"Uu10123456789012345",
			false,
		},
		{
			"LessThanRequiredMinimumLength",
			"Uu",
			false,
		},
		{
			"ValidLogin",
			"Uu.12-3_4",
			true,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				if test.valid != isValidLogin(test.pass) {
					t.Fatal("invalid login")
				}
			},
		)
	}
}

func TestIsValidName(t *testing.T) {
	tests := []struct {
		name  string
		pass  string
		valid bool
	}{
		{
			"NoCharacter",
			"",
			false,
		},
		{
			"CyrillicCharacter",
			"Гофер",
			false,
		},
		{
			"EmptyStringAndWhitespace",
			" \n\t\r\v\f ",
			false,
		},
		{
			"InvalidCharacter",
			"Uua?aaaa",
			false,
		},
		{
			"LongName",
			"Goooooooooooooooooooooooooooooooooooopher",
			false,
		},
		{
			"LessThanRequiredMinimumLength",
			"Uu",
			false,
		},
		{
			"ValidName",
			"Gopher",
			true,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				if test.valid != isValidName(test.pass) {
					t.Fatal("invalid name")
				}
			},
		)
	}
}

func TestIsValPassword(t *testing.T) {
	tests := []struct {
		name  string
		pass  string
		valid bool
	}{
		{
			"NoCharacter",
			"",
			false,
		},
		{
			"EmptyStringAndWhitespace",
			" \n\t\r\v\f ",
			false,
		},
		{
			"MixtureOfEmptyStringAndWhitespace",
			"U u\n1\t?\r1\v2\f34",
			false,
		},
		{
			"MissingUpperCaseString",
			"uu112345",
			false,
		},
		{
			"MissingLowerCaseString",
			"UU112345",
			false,
		},
		{
			"MissingNumber",
			"Uuaaaaa",
			false,
		},
		{
			"LessThanRequiredMinimumLength",
			"Uu123",
			false,
		},
		{
			"HaveCyrillicChars",
			"Uu123456ы",
			false,
		},
		{
			"ValidPassword",
			"Uu123456",
			true,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				if test.valid != isValidPassword(test.pass) {
					t.Fatal("invalid password")
				}
			},
		)
	}
}
