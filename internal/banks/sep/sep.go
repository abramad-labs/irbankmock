package sep

import (
	"errors"
	"strings"

	"github.com/abramad-labs/irbankmock/internal/dbutils"
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
		return nil, errors.New("name can't be empty")
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
		return nil, errors.New("failed to add terminal")
	}
	return &BankSepCreateTerminalResponse{
		ID:       model.ID,
		Name:     model.Name,
		Username: model.Username,
		Password: model.Password,
	}, nil
}
