package sep

type BankSepGetTerminalsResponse struct {
	Terminals []BankSepTerminal `json:"terminals"`
}

type BankSepCreateTerminalRequest struct {
	Name string
}

type BankSepCreateTerminalResponse struct {
	ID       uint64
	Name     string
	Username string
	Password string
}

type BankSepManagementError struct {
	Error   bool
	Message string
}

type BankSepTransactionRequest struct {
	// action to proceed e.g. "token" for getting a payment token
	Action string `json:"action"`

	// the merchant/termianl ID
	TerminalId int64 `json:"terminalId"`

	// amount of payment in IRR
	Amount int64 `json:"amount"`

	// a unique number generated in merchant side to prevent double-spending and can be used for inquery
	ResNum string `json:"resNum"`

	// where to redirect the buyer after the transaction finished
	RedirectURL string `json:"redirectURL"`

	// optional fee of transaction, usually used for business partnership programs
	Wage *int64 `json:"wage,omitempty"`

	// amount that is reduced from the customer card. this parameter is ignore by the
	// irbankmock service.
	AffectiveAmount *int `json:"affectiveAmount,omitempty"`

	// optional buyer phone number64 used to store and retrieve card info and auto-fill
	// the payment form
	CellNumber *string `json:"cellNumber,omitempty"`

	// validity duration of this token in range 20 to 3600 minutes
	TokenExpiryInMin int `json:"tokenExpiryInMin"`

	// optional hash of the card number. irbankmock ignores this value.
	HashedCardNumber *string `json:"hashedCardNumber,omitempty"`

	// if provided, you should pass this key to be able to receive the receipt
	TxnRandomSessionKey *int64 `json:"txnRandomSessionKey,omitempty"`
}

type BankSepTransactionResponse struct {
	// 1 for ok case and -1 in case of failure
	Status    int
	Token     string
	ErrorCode string
	ErrorDesc string
}

type BankSepCustomerPayRequest struct {
	// card number
	Pan string

	// internet payment password
	Pin2       string
	ExpireDate string
	Cvv2       string
	Email      *string
}

// this payload is sent back to the merchant's redirect url to take furthur action
type BankSepRedirectedPayload struct {
	Token string
	// receipt number used for validation and verification of the payment
	RefNum string
}

type BankSepGetReceiptBaseRequest struct {
	TerminalNumber int64

	// only required if this value is provided while requesting token.
	// must match the token's random session key.
	TxnRandomSessionKey *int64
	Rrn                 *int64
}

type BankSepGetReceiptViaRefNumRequest struct {
	BankSepGetReceiptBaseRequest
	RefNum string
}

type BankSepGetReceiptViaTokenRequest struct {
	BankSepGetReceiptBaseRequest
	Token string
}

type BankSepValidationError struct {
	FieldName     string
	ErrorMessages []string
}

type PaymentReceiptStatus int

const (
	PaymentReceiptStatusInProgress PaymentReceiptStatus = iota
	PaymentReceiptStatusCanceledByUser
	PaymentReceiptStatusOK
	PaymentReceiptStatusFailed
)

type PaymentReceiptState string

const PaymentReceiptStateInProgress = PaymentReceiptState("InProgress")
const PaymentReceiptStateCanceledByUser = PaymentReceiptState("CanceledByUser")
const PaymentReceiptStateOK = PaymentReceiptState("OK")
const PaymentReceiptStateFailed = PaymentReceiptState("Failed")
const PaymentReceiptStateUnknown = PaymentReceiptState("Unknown")

func (prs PaymentReceiptStatus) GetState() PaymentReceiptState {
	switch prs {
	case PaymentReceiptStatusInProgress:
		return PaymentReceiptStateInProgress
	case PaymentReceiptStatusCanceledByUser:
		return PaymentReceiptStateCanceledByUser
	case PaymentReceiptStatusOK:
		return PaymentReceiptStateOK
	case PaymentReceiptStatusFailed:
		return PaymentReceiptStateFailed
	}
	return PaymentReceiptStateUnknown
}

type BankSepPaymentReceipt struct {
	State            PaymentReceiptState
	Status           PaymentReceiptStatus
	TerminalId       int64
	Token            string
	RefNum           string
	ResNum           string
	TraceNo          int64
	Amount           int64
	AffectiveAmount  int64
	Rrn              int64
	HashedCardNumber string
}

type BankSepGetReceiptResponse struct {
	HasError         bool
	Data             BankSepPaymentReceipt
	ValidationErrors []BankSepValidationError
	ErrorCode        int32
	ErrorMessage     string
}
