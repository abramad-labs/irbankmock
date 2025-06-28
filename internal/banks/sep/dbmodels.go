package sep

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
	TerminalId uint64          `gorm:"uniqueIndex:terminal_resnum_idx"`
	Terminal   BankSepTerminal `gorm:"foreignKey:TerminalId"`

	// amount of payment in IRR
	Amount int64

	// a unique number generated in merchant side to prevent double-spending and can be used for inquery
	ResNum string `gorm:"uniqueIndex:terminal_resnum_idx"`

	// where to redirect the buyer after the transaction finished
	RedirectURL string

	// optional fee of transaction, usually used for business partnership programs
	Wage *int64

	// amount that is reduced from the customer card. this parameter is ignore by the
	// irbankmock service.
	AffectiveAmount *int

	// optional buyer phone number64 used to store and retrieve card info and auto-fill
	// the payment form
	CellNumber *string

	// validity duration of this token in range 20 to 3600 minutes
	TokenExpiryInMin int `gorm:"check:token_expiry_in_min >= 20 AND token_expiry_in_min <= 3600"`

	// optional hash of the card number. irbankmock ignores this value.
	HashedCardNumber *string

	// if provided, you should pass this key to be able to receive the receipt
	TxnRandomSessionKey *int64
}
