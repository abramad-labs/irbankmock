package sep

import (
	"net/url"
	"regexp"
	"slices"

	"github.com/abramad-labs/irbankmock/internal/banks/sep/seperrors"
)

const SepMinimimTokenExpiry = 20
const SepMaximumTokenExpiry = 3600

func IsValidPhoneNumber(number string) bool {
	re := regexp.MustCompile(`^(09\d{9}|9\d{9})$`)
	return re.MatchString(number)
}

func ValidateURL(rawURL string) error {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return seperrors.ErrXInvalidRedirectURL
	}
	if u.Host == "" {
		return seperrors.ErrXInvalidRedirectURL
	}

	allowedSchemes := []string{"http", "https"}
	if !slices.Contains(allowedSchemes, u.Scheme) {
		return seperrors.ErrXInvalidRedirectURLScheme
	}
	return nil
}

func ClampTokenExpiryMinute(minute int) int {
	if minute <= SepMinimimTokenExpiry {
		return SepMinimimTokenExpiry
	}
	if minute >= SepMaximumTokenExpiry {
		return SepMaximumTokenExpiry
	}
	return minute
}
