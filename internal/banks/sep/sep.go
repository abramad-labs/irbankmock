package sep

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
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
	gonanoid "github.com/matoous/go-nanoid/v2"
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
		ReceiptExpiresAt:    now.Add(time.Hour),
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

	if tokenInfo.Status != PaymentReceiptStatusInProgress {
		return nil, usererror.New(managementerrors.ErrTokenNoLongerAvailable)
	}

	return &BankSepPublicTokenInfoResponse{
		TerminalName: tokenInfo.Terminal.Name,
		TerminalId:   tokenInfo.TerminalId,
		Website:      "mock.example.com",
		Amount:       tokenInfo.Amount,
		ExpiresAt:    tokenInfo.ExpiresAt,
	}, nil
}

func cancelToken(c *fiber.Ctx, req *BankSepCancelOrFailTokenRequest) (*BankSepTokenFinalizeResponse, error) {
	db, err := dbutils.GetDb(c)
	if err != nil {
		return nil, err
	}

	var btrx BankSepTransaction
	err = db.Transaction(func(tx *gorm.DB) error {
		txErr := tx.Model(&BankSepTransaction{}).Where("token = ?", req.Token).Take(&btrx).Error
		if txErr != nil {
			return txErr
		}

		if btrx.Status != PaymentReceiptStatusInProgress {
			return usererror.New(managementerrors.ErrTransactionNotFound)
		}

		txErr = tx.Model(btrx).Updates(map[string]any{
			"cancelled_at": time.Now(),
			"status":       PaymentReceiptStateCanceledByUser,
		}).Error
		if txErr != nil {
			return txErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	url, err := url.Parse(btrx.RedirectURL)
	if err != nil {
		return nil, err
	}
	query := url.Query()
	query.Set("Token", req.Token)
	url.RawQuery = query.Encode()

	return &BankSepTokenFinalizeResponse{
		RedirectURL: url.String(),
		CallbackData: &BankSepTokenFinalizeResponseCallbackData{
			MID:        fmt.Sprint(btrx.TerminalId),
			TerminalId: fmt.Sprint(btrx.TerminalId),
			Token:      req.Token,
		},
	}, nil
}

func failToken(c *fiber.Ctx, req *BankSepCancelOrFailTokenRequest) (*BankSepTokenFinalizeResponse, error) {
	db, err := dbutils.GetDb(c)
	if err != nil {
		return nil, err
	}

	var btrx BankSepTransaction
	err = db.Transaction(func(tx *gorm.DB) error {
		txErr := tx.Model(&BankSepTransaction{}).Where("token = ?", req.Token).Take(&btrx).Error
		if txErr != nil {
			return txErr
		}

		if btrx.Status != PaymentReceiptStatusInProgress {
			return usererror.New(managementerrors.ErrTransactionNotFound)
		}

		txErr = tx.Model(btrx).Updates(map[string]any{
			"failed_at": time.Now(),
			"status":    PaymentReceiptStatusFailed,
		}).Error
		if txErr != nil {
			return txErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	url, err := url.Parse(btrx.RedirectURL)
	if err != nil {
		return nil, err
	}
	query := url.Query()
	query.Set("Token", req.Token)
	url.RawQuery = query.Encode()

	return &BankSepTokenFinalizeResponse{
		RedirectURL: url.String(),
		CallbackData: &BankSepTokenFinalizeResponseCallbackData{
			MID:        fmt.Sprint(btrx.TerminalId),
			TerminalId: fmt.Sprint(btrx.TerminalId),
			Token:      req.Token,
		},
	}, nil
}

func submitToken(c *fiber.Ctx, req *BankSepSubmitTokenRequest) (*BankSepTokenFinalizeResponse, error) {
	db, err := dbutils.GetDb(c)
	if err != nil {
		return nil, err
	}

	var btrx BankSepTransaction
	rrn := rand.Int63()
	refNum, err := gonanoid.New()
	if err != nil {
		return nil, err
	}
	cardHashBinary := sha256.Sum256([]byte(req.CardNumber))
	hashedCardNumber := hex.EncodeToString(cardHashBinary[:])

	err = db.Transaction(func(tx *gorm.DB) error {
		txErr := tx.Model(&BankSepTransaction{}).Where("token = ?", req.Token).Take(&btrx).Error
		if txErr != nil {
			return txErr
		}

		if btrx.Status != PaymentReceiptStatusInProgress {
			return usererror.New(managementerrors.ErrTransactionNotFound)
		}

		now := time.Now()

		update := tx.Model(&BankSepTransaction{}).
			Where("id = ?", btrx.ID).
			Updates(map[string]any{
				"status":             PaymentReceiptStatusOK,
				"rrn":                rrn,
				"ref_num":            refNum,
				"submitted_at":       now,
				"verify_deadline":    now.Add(30 * time.Minute),
				"reverse_deadline":   now.Add(50 * time.Minute),
				"paid_card_number":   req.CardNumber,
				"hashed_card_number": hashedCardNumber,
			})
		if update.Error != nil {
			return update.Error
		}
		if update.RowsAffected == 0 {
			return usererror.New(managementerrors.ErrTransactionNotFound)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	url, err := url.Parse(btrx.RedirectURL)
	if err != nil {
		return nil, err
	}
	query := url.Query()
	query.Set("Token", req.Token)
	query.Set("RefNum", refNum)
	url.RawQuery = query.Encode()

	return &BankSepTokenFinalizeResponse{
		RedirectURL: url.String(),
		CallbackData: &BankSepTokenFinalizeResponseCallbackData{
			MID:              fmt.Sprint(btrx.TerminalId),
			TerminalId:       fmt.Sprint(btrx.TerminalId),
			Token:            req.Token,
			RefNum:           refNum,
			Rrn:              fmt.Sprint(rrn),
			State:            string(btrx.Status.GetState()),
			Status:           fmt.Sprint(btrx.Status),
			ResNum:           btrx.ResNum,
			Amount:           fmt.Sprint(btrx.Amount),
			HashedCardNumber: hashedCardNumber,
			SecurePan:        maskThirdQuarter(req.CardNumber),
		},
	}, nil
}
