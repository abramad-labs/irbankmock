package sep

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"slices"

	"github.com/abramad-labs/irbankmock/internal/banks/sep/managementerrors"
	"github.com/abramad-labs/irbankmock/internal/banks/sep/seperrors"
	"github.com/abramad-labs/irbankmock/internal/dbutils"
	"github.com/abramad-labs/irbankmock/internal/usererror"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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

	resp := &BankSepGetTerminalsResponse{
		Terminals: terminals,
	}
	return resp, nil
}

func createTerminal(ctx *fiber.Ctx, req *BankSepCreateTerminalRequest) (*BankSepCreateTerminalResponse, error) {
	db, err := dbutils.GetDb(ctx)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(req.Name) == "" {
		return nil, usererror.NewBadRequest(managementerrors.ErrEmptyName)
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
	return &BankSepCreateTerminalResponse{
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
	_ = db
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
	if req.TokenExpiryInMin <= 20 {
		req.TokenExpiryInMin = 20
	}

	req.TokenExpiryInMin = ClampTokenExpiryMinute(req.TokenExpiryInMin)

	now := time.Now()

	token := uuid.NewString()
	refNum := uuid.NewString()

	trxModel := &BankSepTransaction{
		TerminalId:          req.TerminalId,
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
