package sep

import "time"

type BankSepGetTerminalsResponseEndpoints struct {
	PaymentGateway     string `json:"paymentGateway"`
	PaymentToken       string `json:"paymentToken"`
	Receipt            string `json:"receipt"`
	VerifyTransaction  string `json:"verifyTransaction"`
	ReverseTransaction string `json:"reverseTransaction"`
}

type BankSepGetTerminalsResponse struct {
	Terminals []*BankSepTerminalResponse            `json:"terminals"`
	Endpoints *BankSepGetTerminalsResponseEndpoints `json:"endpoints"`
}

type BankSepCreateTerminalRequest struct {
	Name string `json:"name"`
}

type BankSepTerminalResponse struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type BankSepManagementError struct {
	Error   bool
	Message string
}

type BankSepTransactionRequest struct {
	// action to proceed e.g. "token" for getting a payment token
	Action string `json:"action"`

	// the merchant/termianl ID
	TerminalId string `json:"terminalId"`

	// amount of payment in IRR
	Amount int64 `json:"amount"`

	// reservation number is a unique number generated in merchant side to
	// prevent double-spending and can be used for inquery
	ResNum string `json:"resNum"`

	// optional resnum, used for reporting
	ResNum1 *string `json:"resNum1"`
	// optional resnum, used for reporting
	ResNum2 *string `json:"resNum2"`
	// optional resnum, used for reporting
	ResNum3 *string `json:"resNum3"`
	// optional resnum, used for reporting
	ResNum4 *string `json:"resNum4"`

	// where to redirect the buyer after the transaction finished
	RedirectURL string `json:"redirectURL"`

	// optional fee of transaction, usually used for business partnership programs
	Wage *int64 `json:"wage,omitempty"`

	// amount that is reduced from the customer card. this parameter is ignore by the
	// irbankmock service.
	AffectiveAmount *int64 `json:"affectiveAmount,omitempty"`

	// optional buyer phone number64 used to store and retrieve card info and auto-fill
	// the payment form
	CellNumber *string `json:"cellNumber,omitempty"`

	// validity duration of this token in range 20 to 3600 minutes
	TokenExpiryInMin int `json:"tokenExpiryInMin"`

	// optional md5 hash of the card number for input and sha256 for output.
	// forces user to pick these cards only.
	// separate with one of the |;, characters to send multiple hashes.
	// maximum 10 cards allowed.
	HashedCardNumber *string `json:"hashedCardNumber,omitempty"`

	// if provided, you should pass this key to be able to receive the receipt
	TxnRandomSessionKey *int64 `json:"txnRandomSessionKey,omitempty"`
}

type BankSepTransactionResponse struct {
	// 1 for ok case and -1 in case of failure
	Status    int    `json:"status"`
	Token     string `json:"token,omitempty"`
	ErrorCode string `json:"errorCode,omitempty"`
	ErrorDesc string `json:"errorDesc,omitempty"`
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

type BankSepGetReceiptRequest struct {
	BankSepGetReceiptBaseRequest
	Token  *string
	RefNum *string
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
	ValidationErrors []*BankSepValidationError
	ErrorCode        int32
	ErrorMessage     string
}

type BankSepPublicTokenInfoResponse struct {
	TerminalName string    `json:"terminalName"`
	TerminalId   uint64    `json:"terminalId"`
	Website      string    `json:"website"`
	Amount       int64     `json:"amount"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

type BankSepCancelOrFailTokenRequest struct {
	Token string `json:"token"`
}

type BankSepSubmitTokenRequest struct {
	Token        string `json:"token"`
	CardNumber   string `json:"cardNumber"`
	Cvv          int32  `json:"cvv"`
	ExpiryMonth  int32  `json:"expiryMonth"`
	ExpiryYear   int32  `json:"expiryYear"`
	CardPassword string `json:"cardPassword"`
	Captcha      string `json:"captcha"`
}

type BankSepTokenFinalizeResponseCallbackData struct {
	MID              string `json:"MID"`
	TerminalId       string `json:"terminalId"`
	State            string `json:"state"`
	Status           string `json:"status"`
	Rrn              string `json:"rrn"`
	RefNum           string `json:"refNum"`
	ResNum           string `json:"resNum"`
	TraceNo          string `json:"traceNo"`
	Amount           string `json:"amount"`
	AffectiveAmount  string `json:"affectiveAmount"`
	Wage             string `json:"wage"`
	SecurePan        string `json:"securePan"`
	HashedCardNumber string `json:"hashedCardNumber"`
	Token            string `json:"token"`
}

type BankSepTokenFinalizeResponse struct {
	RedirectURL  string                                    `json:"redirectURL"`
	CallbackData *BankSepTokenFinalizeResponseCallbackData `json:"callbackData"`
}

type BankSepVerificationRequest struct {
	RefNum         string
	TerminalNumber int64
}

type BankSepTransactionDetailResponse struct {
	RRN             string
	RefNum          string
	MaskedPan       string
	HashedPan       string
	TerminalNumber  int32
	OriginalAmount  int64
	AffectiveAmount int64
	StraceDate      time.Time
	StraceNo        int64
}

type BankSepVerificationResponse struct {
	TransactionDetail *BankSepTransactionDetailResponse
	ResultCode        int32
	ResultDescription string
	Success           bool
}

type BankSepReverseRequest BankSepVerificationRequest
type BankSepReverseResponse BankSepVerificationResponse
