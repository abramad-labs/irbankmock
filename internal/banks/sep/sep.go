package sep

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"slices"

	"github.com/abramad-labs/irbankmock/internal/banks/registry"
	"github.com/abramad-labs/irbankmock/internal/banks/sep/managementerrors"
	"github.com/abramad-labs/irbankmock/internal/banks/sep/seperrors"
	"github.com/abramad-labs/irbankmock/internal/conf"
	"github.com/abramad-labs/irbankmock/internal/dbutils"
	"github.com/abramad-labs/irbankmock/internal/security"
	"github.com/abramad-labs/irbankmock/internal/usererror"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getTerminalEndpoints(ctx *fiber.Ctx) *BankSepGetTerminalsResponseEndpoints {
	routerPrefix := registry.GetRouterPrefix(ctx)
	fullPrefix := conf.GetPublicHostname() + routerPrefix
	return &BankSepGetTerminalsResponseEndpoints{
		PaymentGateway:     fullPrefix + BankSepPathOnlinePaymentGateway,
		PaymentToken:       fullPrefix + BankSepPathOnlinePaymenyTokenRedirect,
		Receipt:            fullPrefix + BankSepPathGetReceipt,
		VerifyTransaction:  fullPrefix + BankSepPathVerifyTransaction,
		ReverseTransaction: fullPrefix + BankSepPathReverseTransaction,
	}
}

func getTerminals(ctx *fiber.Ctx) (*BankSepGetTerminalsResponse, error) {
	db, err := dbutils.GetDb(ctx)
	if err != nil {
		return nil, err
	}

	var terminals []BankSepTerminal

	err = db.Find(&terminals).Error
	if err != nil {
		return nil, errors.New("failed to fetch terminals")
	}

	terminalResponse := make([]*BankSepTerminalResponse, len(terminals))
	for i, t := range terminals {
		terminalResponse[i] = &BankSepTerminalResponse{
			ID:       t.ID,
			Name:     t.Name,
			Username: t.Username,
			Password: t.Password,
		}
	}

	resp := &BankSepGetTerminalsResponse{
		Terminals: terminalResponse,
		Endpoints: getTerminalEndpoints(ctx),
	}
	return resp, nil
}

func createTerminal(ctx *fiber.Ctx, req *BankSepCreateTerminalRequest) (*BankSepTerminalResponse, error) {
	db, err := dbutils.GetDb(ctx)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(req.Name) == "" {
		return nil, usererror.NewBadRequest(managementerrors.ErrEmptyName)
	}

	if security.StringHasInsecureCharacters(req.Name) {
		return nil, usererror.NewBadRequest(managementerrors.ErrInvalidName)
	}

	username := uuid.NewString()
	password := uuid.NewString()

	model := &BankSepTerminal{
		Name:     req.Name,
		Username: username,
		Password: password,
	}

	err = db.Create(&model).Error
	if err != nil {
		return nil, fmt.Errorf("failed creating terminal: %w", err)
	}
	return &BankSepTerminalResponse{
		ID:       model.ID,
		Name:     model.Name,
		Username: model.Username,
		Password: model.Password,
	}, nil
}

func processTransactionRequest(ctx *fiber.Ctx, req *BankSepTransactionRequest) (*BankSepTransactionResponse, error) {
	db, err := dbutils.GetDb(ctx)
	if err != nil {
		return nil, err
	}
	if req.Action != "token" {
		return nil, seperrors.ErrXInvalidAction
	}
	if req.Amount <= 0 {
		return nil, seperrors.ErrXInvalidAmount
	}
	if req.CellNumber != nil && !IsValidPhoneNumber(*req.CellNumber) {
		return nil, seperrors.ErrXInvalidPhoneNumber
	}
	if req.HashedCardNumber != nil {
		cardHashes := SplitByDelimiters(*req.HashedCardNumber)
		if len(cardHashes) > 10 {
			return nil, seperrors.ErrXInvalidNumberOfCards
		}
		if slices.Contains(cardHashes, "") {
			return nil, seperrors.ErrXInvalidCardHash
		}
	}
	err = ValidateURL(req.RedirectURL)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.ResNum) == "" {
		return nil, seperrors.ErrXEmptyResNum
	}
	req.TokenExpiryInMin = ClampTokenExpiryMinute(req.TokenExpiryInMin)

	now := time.Now()

	token := uuid.NewString()
	refNum := uuid.NewString()

	terminalId, err := strconv.ParseUint(req.TerminalId, 10, 64)
	if err != nil {
		return nil, seperrors.ErrTerminalNotFound
	}

	// TODO: do it in a transaction
	var exists bool
	err = db.Model(&BankSepTerminal{}).
		Select("count(*) > 0").
		Where("id = ?", terminalId).
		Find(&exists).
		Error

	if err != nil || !exists {
		return nil, seperrors.ErrTerminalNotFound
	}

	trxModel := &BankSepTransaction{
		TerminalId:          terminalId,
		Amount:              req.Amount,
		ResNum:              req.ResNum,
		ResNum1:             req.ResNum1,
		ResNum2:             req.ResNum2,
		ResNum3:             req.ResNum3,
		ResNum4:             req.ResNum4,
		RedirectURL:         req.RedirectURL,
		Wage:                req.Wage,
		AffectiveAmount:     req.AffectiveAmount,
		CellNumber:          req.CellNumber,
		TokenExpiryInMin:    req.TokenExpiryInMin,
		HashedCardNumber:    req.HashedCardNumber,
		TxnRandomSessionKey: req.TxnRandomSessionKey,
		CreatedAt:           now,
		ExpiresAt:           now.Add(time.Duration(req.TokenExpiryInMin) * time.Minute),
		Token:               token,
		Verified:            false,
		ReceiptExpiresAt:    now.Add(time.Hour),
		RefNum:              refNum,
	}

	err = db.Create(&trxModel).Error
	if err != nil {
		return nil, err
	}

	return &BankSepTransactionResponse{
		Status: 1,
		Token:  token,
	}, nil
}

func getPublicTokenInfo(c *fiber.Ctx, token string) (*BankSepPublicTokenInfoResponse, error) {
	db, err := dbutils.GetDb(c)
	if err != nil {
		return nil, err
	}
	var tokenInfo BankSepTransaction
	err = db.Model(&BankSepTransaction{}).Preload("Terminal").
		Where("token = ?", token).Take(&tokenInfo).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, usererror.New(managementerrors.ErrTokenNotFound)
		}
		return nil, err
	}

	if tokenInfo.ExpiresAt.Before(time.Now()) {
		return nil, usererror.New(managementerrors.ErrTokenExpired)
	}

	return &BankSepPublicTokenInfoResponse{
		TerminalName: tokenInfo.Terminal.Name,
		TerminalId:   tokenInfo.TerminalId,
		Website:      "mock.example.com",
		Amount:       tokenInfo.Amount,
		ExpiresAt:    tokenInfo.ExpiresAt,
	}, nil
}
