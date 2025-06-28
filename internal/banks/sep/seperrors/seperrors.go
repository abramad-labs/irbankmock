package seperrors

import "errors"

var ErrTerminalNotFound = errors.New("TerminalNotFound")
var ErrResourceNotFound = errors.New("ResourceNotFound")
var ErrMerchantIpAddressIsInvalid = errors.New("MerchantIpAddressIsInvalid")
var ErrTerminalIsDisabled = errors.New("TerminalIsDisabled")

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
