package sep

import (
	"net/url"
	"strconv"

	"github.com/abramad-labs/irbankmock/internal/banks/registry"
	"github.com/abramad-labs/irbankmock/internal/banks/sep/seperrors"
	"github.com/abramad-labs/irbankmock/internal/dbutils/migration"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const BankSepPathOnlinePaymentGateway = "/OnlinePG/OnlinePG"
const BankSepPathOnlinePaymenyTokenRedirect = "/OnlinePG/SendToken"

const BankSepPathGetReceipt = "/verifyTxnRandomSessionkey/api/v2/ipg/payment/receipt"
const BankSepPathVerifyTransaction = "/verifyTxnRandomSessionkey/ipg/VerifyTransaction"
const BankSepPathReverseTransaction = "/verifyTxnRandomSessionkey/ipg/ReverseTransaction"

func init() {
	migration.RegisterMigration("samanbank_models", func(m gorm.Migrator) error {
		return m.AutoMigrate(BankSepTerminal{}, BankSepTransaction{})
	})

	registry.RegisterBank("saman", func(g fiber.Router) {
		g.Post("/management/terminal", CreateTerminal)
		g.Get("/management/terminal", GetTerminals)
		g.Post(BankSepPathOnlinePaymentGateway, PaymentGwTransaction)
	})
}

func GetTerminals(c *fiber.Ctx) error {
	resp, err := getTerminals(c)
	if err != nil {
		return err
	}
	return c.JSON(resp)
}

func CreateTerminal(c *fiber.Ctx) error {
	req := new(BankSepCreateTerminalRequest)
	err := c.BodyParser(req)
	if err != nil {
		return err
	}

	resp, err := createTerminal(c, req)
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

func sendJsonFromSamanError(c *fiber.Ctx, err error, status int) error {
	return c.Status(status).JSON(BankSepTransactionResponse{
		Status:    -1,
		ErrorCode: strconv.Itoa(seperrors.GetBankSepErrorCode(err)),
		ErrorDesc: err.Error(),
	})
}

func PaymentGwTransaction(c *fiber.Ctx) error {
	tokenValue := c.FormValue("Token")
	if tokenValue != "" {
		target := "/api/saman" + BankSepPathOnlinePaymenyTokenRedirect + "?token=" + url.QueryEscape(tokenValue)
		return c.Redirect(target, fiber.StatusTemporaryRedirect)
	}
	txReq := new(BankSepTransactionRequest)
	err := c.BodyParser(txReq)
	if err != nil {
		return sendJsonFromSamanError(c, seperrors.ErrXInvalidRequest, fiber.StatusBadRequest)
	}
	resp, err := processTransactionRequest(c, txReq)
	if err != nil {
		return sendJsonFromSamanError(c, err, fiber.StatusBadRequest)
	}
	return c.JSON(resp)
}
