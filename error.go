package viesparser

import "errors"

var (
	ErrorInvalidOption          = errors.New("invalid option")
	ErrorMissingAddress         = errors.New("missing address")
	ErrorMissingCountryCode     = errors.New("missing country code")
	ErrorUnsupportedCountryCode = errors.New("unsupported country code")
)
