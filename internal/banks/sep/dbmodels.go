package sep

import "time"

// A terminal/merchant is an entity who is providing services to a customer,
// registered in bank's database.
type BankSepTerminal struct {
	ID       uint64 `gorm:"primarykey"`
	Name     string
	Username string
	Password string
}

type BankSepTransaction struct {
	ID uint64 `gorm:"primarykey"`

	// the merchant/termianl ID
	TerminalId uint64          `gorm:"index:,unique,composite:terminal_resnum_idx"`
	Terminal   BankSepTerminal `gorm:"foreignKey:TerminalId"`

	// amount of payment in IRR
	Amount int64

	// reservation number is a unique number generated in merchant side to
	// prevent double-spending and can be used for inquery
	ResNum string `gorm:"size:50;index:,unique,composite:terminal_resnum_idx"`

	// optional resnum, used for reporting
	ResNum1 *string `gorm:"size:50"`

	// optional resnum, used for reporting
	ResNum2 *string `gorm:"size:50"`

	// optional resnum, used for reporting
	ResNum3 *string `gorm:"size:50"`

	// optional resnum, used for reporting
	ResNum4 *string `gorm:"size:50"`

	// where to redirect the buyer after the transaction finished
	RedirectURL string `gorm:"size:2083"`

	// optional fee of transaction, usually used for business partnership programs
	Wage *int64

	// amount that is reduced from the customer card. this parameter is ignore by the
	// irbankmock service.
	AffectiveAmount *int64

	// optional buyer phone number64 used to store and retrieve card info and auto-fill
	// the payment form
	CellNumber *string

	// validity duration of this token in range 20 to 3600 minutes
	TokenExpiryInMin int `gorm:"check:token_expiry_in_min >= 20 AND token_expiry_in_min <= 3600"`

	// optional md5 hash of the card number for input and sha256 for output.
	// forces user to pick these cards only.
	// separate with one of the |;, characters to send multiple hashes.
	// maximum 10 cards allowed.
	HashedCardNumber *string

	// actuall card after payment was done
	PaidCardNumber *string

	// if provided, you should pass this key to be able to receive the receipt
	TxnRandomSessionKey *int64

	Token string

	TraceNo   *int64
	TraceDate *time.Time

	// reference number used for validation and verification of transaction
	// generated in bank side only after a successful transaction
	RefNum *string

	// Retrieval Reference Number (RRN) is a unique identifier assigned to a specific transaction
	// to facilitate the retrieval of transaction details
	// generated in bank side only after a successful transaction
	Rrn *int64

	CancelledAt *time.Time
	FailedAt    *time.Time
	SubmittedAt *time.Time
	VerifiedAt  *time.Time
	ReversedAt  *time.Time

	// merchant has 30 minutes to verify
	VerifyDeadline *time.Time

	// merchant has 50 minutes to reverse
	ReverseDeadline *time.Time

	Status PaymentReceiptStatus

	CreatedAt time.Time

	// this transaction will be invalidated due this time.
	// calculated using created_at + token_expiry_in_min
	ExpiresAt time.Time

	// receipt will expire an hour past token creation
	ReceiptExpiresAt time.Time
}
