package seperrors

import (
	"errors"
)

// Errors extracted from the specification document
var ErrTerminalNotFound = errors.New("TerminalNotFound")
var ErrResourceNotFound = errors.New("ResourceNotFound")
var ErrMerchantIpAddressIsInvalid = errors.New("MerchantIpAddressIsInvalid")
var ErrTerminalIsDisabled = errors.New("TerminalIsDisabled")

// Errors infered from what doc explains about the behavior
var ErrXInvalidRequest = errors.New("invalid request")
var ErrXInvalidAction = errors.New("only 'token' is valid for action")
var ErrXInvalidAmount = errors.New("invalid transaction amount")
var ErrXInvalidPhoneNumber = errors.New("invalid phone number")
var ErrXInvalidCardHash = errors.New("hash of card is invalid")
var ErrXInvalidNumberOfCards = errors.New("no more than 10 cards is allowerd")
var ErrXInvalidRedirectURL = errors.New("redirect url does not have correct format")
var ErrXInvalidRedirectURLScheme = errors.New("redirect url does not have correct scheme")
var ErrXEmptyResNum = errors.New("must include resnum")

func GetBankSepErrorCode(err error) int {
	if errors.Is(err, ErrTerminalIsDisabled) {
		return 12
	}
	if errors.Is(err, ErrResourceNotFound) {
		return 404
	}
	if errors.Is(err, ErrMerchantIpAddressIsInvalid) {
		return 8
	}
	if errors.Is(err, ErrTerminalIsDisabled) {
		return 20
	}

	return -1
}
