package sep

import (
	"fmt"
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
		g.Get("/public/token", GetTokenInfo)
		g.Post("/management/token/submit", SubmitToken)
		g.Post("/management/token/fail", FailToken)
		g.Post("/management/token/cancel", CancelToken)
		g.Post(BankSepPathOnlinePaymentGateway, PaymentGwTransaction)
		g.Post(BankSepPathGetReceipt, GetReceipt)
		g.Post(BankSepPathVerifyTransaction, VerifyTransaction)
		g.Post(BankSepPathReverseTransaction, ReverseTransaction)
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
		routerPrefix := registry.GetRouterPrefix(c)
		target := routerPrefix + BankSepPathOnlinePaymenyTokenRedirect + "?token=" + url.QueryEscape(tokenValue)
		return c.Redirect(target, fiber.StatusTemporaryRedirect)
	}
	txReq := new(BankSepTransactionRequest)
	err := c.BodyParser(txReq)
	if err != nil {
		wErr := fmt.Errorf("%w: %w", seperrors.ErrXInvalidRequest, err)
		return sendJsonFromSamanError(c, wErr, fiber.StatusBadRequest)
	}
	resp, err := processTransactionRequest(c, txReq)
	if err != nil {
		return sendJsonFromSamanError(c, err, fiber.StatusBadRequest)
	}
	return c.JSON(resp)
}

func GetTokenInfo(c *fiber.Ctx) error {
	tokenValue := c.Query("token")
	resp, err := getPublicTokenInfo(c, tokenValue)
	if err != nil {
		return err
	}
	return c.JSON(resp)
}

func SubmitToken(c *fiber.Ctx) error {
	req := new(BankSepSubmitTokenRequest)
	err := c.BodyParser(req)
	if err != nil {
		return err
	}
	resp, err := submitToken(c, req)
	if err != nil {
		return err
	}
	return c.JSON(resp)
}

func CancelToken(c *fiber.Ctx) error {
	req := new(BankSepCancelOrFailTokenRequest)
	err := c.BodyParser(req)
	if err != nil {
		return err
	}

	resp, err := cancelToken(c, req)
	if err != nil {
		return err
	}
	return c.JSON(resp)
}

func FailToken(c *fiber.Ctx) error {
	req := new(BankSepCancelOrFailTokenRequest)
	err := c.BodyParser(req)
	if err != nil {
		return err
	}

	resp, err := failToken(c, req)
	if err != nil {
		return err
	}
	return c.JSON(resp)
}

func GetReceipt(c *fiber.Ctx) error {
	req := new(BankSepGetReceiptRequest)
	err := c.BodyParser(&req)
	if err != nil {
		return err
	}
	resp, err := getReceipt(c, req.TerminalNumber, req.RefNum, req.Token, req.TxnRandomSessionKey, req.Rrn)
	if err != nil {
		return err
	}
	return c.JSON(resp)
}

func VerifyTransaction(c *fiber.Ctx) error {
	req := new(BankSepVerificationRequest)
	err := c.BodyParser(&req)
	if err != nil {
		return err
	}
	resp, err := verifyTransaction(c, req.TerminalNumber, req.RefNum)
	if err != nil {
		return err
	}
	return c.JSON(resp)
}

func ReverseTransaction(c *fiber.Ctx) error {
	req := new(BankSepReverseRequest)
	err := c.BodyParser(&req)
	if err != nil {
		return err
	}
	resp, err := reverseTransaction(c, req.TerminalNumber, req.RefNum)
	if err != nil {
		return err
	}
	return c.JSON(resp)
}
