package sep

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"slices"

	"github.com/abramad-labs/irbankmock/internal/banks/registry"
	"github.com/abramad-labs/irbankmock/internal/banks/sep/managementerrors"
	"github.com/abramad-labs/irbankmock/internal/banks/sep/seperrors"
	"github.com/abramad-labs/irbankmock/internal/conf"
	"github.com/abramad-labs/irbankmock/internal/dbutils"
	"github.com/abramad-labs/irbankmock/internal/pointers"
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

	terminalId, err := req.TerminalId.Int64()
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
	traceNo := rand.Int63()
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
				"trace_no":           traceNo,
				"trace_date":         now,
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

func getReceipt(c *fiber.Ctx, terminalId int64, refNum *string, token *string, rndSessionKey *int64, rrn *int64) (*BankSepGetReceiptResponse, error) {
	db, err := dbutils.GetDb(c)
	if err != nil {
		return nil, err
	}

	var terminal BankSepTerminal
	err = db.Model(&BankSepTerminal{}).Where("id = ?", terminalId).Take(&terminal).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &BankSepGetReceiptResponse{
				HasError:     true,
				ErrorCode:    12,
				ErrorMessage: "TerminalNotFound",
			}, nil
		}
		return &BankSepGetReceiptResponse{
			HasError:     true,
			ErrorCode:    -1,
			ErrorMessage: err.Error(),
		}, nil
	}

	var tx BankSepTransaction
	query := db.Model(&BankSepTransaction{})
	if refNum != nil {
		query = query.Where("terminal_id = ? and ref_num = ?", terminalId, refNum)
	} else if token != nil {
		query = query.Where("terminal_id = ? and token = ?", terminalId, token)
	} else {
		return &BankSepGetReceiptResponse{
			HasError:     true,
			ErrorCode:    -1,
			ErrorMessage: "either token or refnum must be provided",
		}, nil
	}
	if rndSessionKey != nil {
		query = query.Where("txn_random_session_key = ?", rndSessionKey)
	}
	if rrn != nil {
		query = query.Where("rrn = ?", rrn)
	}
	err = query.Take(&tx).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &BankSepGetReceiptResponse{
				HasError:     true,
				ErrorCode:    404,
				ErrorMessage: "ResourceNotFound",
			}, nil
		}
		return &BankSepGetReceiptResponse{
			HasError:     true,
			ErrorCode:    -1,
			ErrorMessage: err.Error(),
		}, nil
	}
	if tx.ReceiptExpiresAt.Before(time.Now()) {
		return &BankSepGetReceiptResponse{
			HasError:     true,
			ErrorCode:    404,
			ErrorMessage: "ResourceNotFound",
		}, nil
	}
	return &BankSepGetReceiptResponse{
		Data: BankSepPaymentReceipt{
			State:      tx.Status.GetState(),
			Status:     tx.Status,
			TerminalId: int64(tx.TerminalId),
			Token:      tx.Token,
			RefNum:     pointers.DerefZero(tx.RefNum),
			ResNum:     tx.ResNum,
			TraceNo:    pointers.DerefZero(tx.TraceNo),
			Amount:     int64(tx.Amount),
			AffectiveAmount: func() int64 {
				v := pointers.DerefZero(tx.AffectiveAmount)
				if v == 0 {
					return tx.Amount
				}
				return v
			}(),
			Rrn:              pointers.DerefZero(tx.Rrn),
			HashedCardNumber: pointers.DerefZero(tx.HashedCardNumber),
		},
	}, nil
}

func verifyTransaction(c *fiber.Ctx, terminalId int64, refNum string) (*BankSepVerificationResponse, error) {
	db, err := dbutils.GetDb(c)
	if err != nil {
		return nil, err
	}

	var terminalExists bool
	err = db.Model(&BankSepTerminal{}).
		Select("count(*) > 0").
		Where("id = ?", terminalId).
		Find(&terminalExists).
		Error

	if err != nil {
		return &BankSepVerificationResponse{
			Success:           false,
			ResultCode:        -1,
			ResultDescription: err.Error(),
		}, nil
	}

	if !terminalExists {
		return &BankSepVerificationResponse{
			Success:           false,
			ResultCode:        -105,
			ResultDescription: "ترمینال ارسالی در سیستم موجود نمی باشد.",
		}, nil
	}

	var btx BankSepTransaction
	err = db.Model(&BankSepTransaction{}).Where("terminal_id = ? and ref_num = ?", terminalId, refNum).Take(&btx).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &BankSepVerificationResponse{
				Success:           false,
				ResultCode:        -2,
				ResultDescription: "تراکنش یافت نشد",
			}, nil
		}
		return &BankSepVerificationResponse{
			Success:           false,
			ResultCode:        -1,
			ResultDescription: err.Error(),
		}, nil
	}

	if btx.Status != PaymentReceiptStatusOK {
		return &BankSepVerificationResponse{
			Success:           false,
			ResultCode:        -2,
			ResultDescription: "تراکنش یافت نشد",
		}, nil
	}

	if btx.VerifyDeadline.Before(time.Now()) {
		return &BankSepVerificationResponse{
			Success:           false,
			ResultCode:        -6,
			ResultDescription: "بیش از نیم ساعت از اجرای تراکنش گذشته است.",
		}, nil
	}

	if btx.VerifiedAt != nil {
		return &BankSepVerificationResponse{
			Success:           false,
			ResultCode:        2,
			ResultDescription: "درخواست تکراری می باشد.",
		}, nil
	}

	if btx.ReversedAt != nil {
		return &BankSepVerificationResponse{
			Success:           false,
			ResultCode:        5,
			ResultDescription: "تراکنش برگشت خورده می باشد.",
		}, nil
	}

	now := time.Now()

	update := db.Model(&BankSepTransaction{}).Where("id = ?", btx.ID).Update("verified_at", now)

	if update.Error != nil {
		return &BankSepVerificationResponse{
			Success:           false,
			ResultCode:        -1,
			ResultDescription: update.Error.Error(),
		}, nil
	}
	if update.RowsAffected == 0 {
		return &BankSepVerificationResponse{
			Success:           false,
			ResultCode:        -2,
			ResultDescription: "تراکنش یافت نشد",
		}, nil
	}

	return &BankSepVerificationResponse{
		Success:           true,
		ResultDescription: "عملیات با موفقیت انجام شد.",
		TransactionDetail: &BankSepTransactionDetailResponse{
			RRN:            fmt.Sprint(btx.Rrn),
			RefNum:         *btx.RefNum,
			MaskedPan:      maskThirdQuarter(*btx.PaidCardNumber),
			HashedPan:      *btx.HashedCardNumber,
			TerminalNumber: int32(btx.TerminalId),
			OriginalAmount: int64(btx.Amount),
			AffectiveAmount: func() int64 {
				if btx.AffectiveAmount == nil {
					return btx.Amount
				} else {
					return *btx.AffectiveAmount
				}
			}(),
			StraceDate: *btx.TraceDate,
			StraceNo:   *btx.TraceNo,
		},
	}, nil
}

func reverseTransaction(c *fiber.Ctx, terminalId int64, refNum string) (*BankSepReverseResponse, error) {
	db, err := dbutils.GetDb(c)
	if err != nil {
		return nil, err
	}

	var terminalExists bool
	err = db.Model(&BankSepTerminal{}).
		Select("count(*) > 0").
		Where("id = ?", terminalId).
		Find(&terminalExists).
		Error

	if err != nil {
		return &BankSepReverseResponse{
			Success:           false,
			ResultCode:        -1,
			ResultDescription: err.Error(),
		}, nil
	}

	if !terminalExists {
		return &BankSepReverseResponse{
			Success:           false,
			ResultCode:        -105,
			ResultDescription: "ترمینال ارسالی در سیستم موجود نمی باشد.",
		}, nil
	}

	var btx BankSepTransaction
	err = db.Model(&BankSepTransaction{}).Where("terminal_id = ? and ref_num = ?", terminalId, refNum).Take(&btx).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &BankSepReverseResponse{
				Success:           false,
				ResultCode:        -2,
				ResultDescription: "تراکنش یافت نشد",
			}, nil
		}
		return &BankSepReverseResponse{
			Success:           false,
			ResultCode:        -1,
			ResultDescription: err.Error(),
		}, nil
	}

	if btx.Status != PaymentReceiptStatusOK {
		return &BankSepReverseResponse{
			Success:           false,
			ResultCode:        -2,
			ResultDescription: "تراکنش یافت نشد",
		}, nil
	}

	if btx.ReverseDeadline.Before(time.Now()) {
		return &BankSepReverseResponse{
			Success:           false,
			ResultCode:        -6,
			ResultDescription: "بیش از 50 دقیقه از اجرای تراکنش گذشته است.",
		}, nil
	}

	if btx.VerifiedAt == nil {
		return &BankSepReverseResponse{
			Success:           false,
			ResultCode:        2,
			ResultDescription: "درخواست تایید نشده می باشد.",
		}, nil
	}

	if btx.ReversedAt != nil {
		return &BankSepReverseResponse{
			Success:           false,
			ResultCode:        5,
			ResultDescription: "تراکنش برگشت خورده می باشد.",
		}, nil
	}

	now := time.Now()

	update := db.Model(&BankSepTransaction{}).Where("id = ?", btx.ID).Update("reversed_at", now)

	if update.Error != nil {
		return &BankSepReverseResponse{
			Success:           false,
			ResultCode:        -1,
			ResultDescription: update.Error.Error(),
		}, nil
	}
	if update.RowsAffected == 0 {
		return &BankSepReverseResponse{
			Success:           false,
			ResultCode:        -2,
			ResultDescription: "تراکنش یافت نشد",
		}, nil
	}

	return &BankSepReverseResponse{
		Success:           true,
		ResultDescription: "عملیات با موفقیت انجام شد.",
		TransactionDetail: &BankSepTransactionDetailResponse{
			RRN:            fmt.Sprint(btx.Rrn),
			RefNum:         *btx.RefNum,
			MaskedPan:      maskThirdQuarter(*btx.PaidCardNumber),
			HashedPan:      *btx.HashedCardNumber,
			TerminalNumber: int32(btx.TerminalId),
			OriginalAmount: int64(btx.Amount),
			AffectiveAmount: func() int64 {
				if btx.AffectiveAmount == nil {
					return btx.Amount
				} else {
					return *btx.AffectiveAmount
				}
			}(),
			StraceDate: *btx.TraceDate,
			StraceNo:   *btx.TraceNo,
		},
	}, nil
}
