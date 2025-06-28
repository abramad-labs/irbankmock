package sep

import (
	"strconv"

	"github.com/abramad-labs/irbankmock/internal/banks/registry"
	"github.com/abramad-labs/irbankmock/internal/banks/sep/seperrors"
	"github.com/abramad-labs/irbankmock/internal/dbutils/migration"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func init() {
	migration.RegisterMigration("samanbank_models", func(m gorm.Migrator) error {
		return m.AutoMigrate(BankSepTerminal{})
	})

	registry.RegisterBank("saman", func(g fiber.Router) {
		g.Post(BankSepPathOnlinePaymentGateway, func(c *fiber.Ctx) error {
			var txReq BankSepTransactionRequest
			err := c.BodyParser(&txReq)
			if err != nil {
				return c.JSON(BankSepTransactionResponse{
					Status:    -1,
					ErrorCode: strconv.Itoa(seperrors.GetBankSepErrorCode(seperrors.ErrTerminalNotFound)),
					ErrorDesc: err.Error(),
				})
			}
			return c.JSON(txReq)
		})
	})
}

const BankSepPathOnlinePaymentGateway = "/OnlinePG/OnlinePG"
const BankSepPathOnlinePaymenyTokenRedirect = "/OnlinePG/SendToken"

const BankSepPathGetReceipt = "/verifyTxnRandomSessionkey/api/v2/ipg/payment/receipt"
const BankSepPathVerifyTransaction = "/verifyTxnRandomSessionkey/ipg/VerifyTransaction"
const BankSepPathReverseTransaction = "/verifyTxnRandomSessionkey/ipg/ReverseTransaction"
