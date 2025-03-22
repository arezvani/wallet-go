package controllers

import (
	"time"

	"github.com/arezvani/wallet-go/app/models"
	"github.com/arezvani/wallet-go/platform/database"

	"github.com/gofiber/fiber/v2"
)

// PostTransaction handles creating a new transaction
// @Description Create a new transaction for a wallet.
// @Summary create a new transaction
// @Tags Wallet
// @Accept json
// @Produce json
// @Param transaction body models.Transaction true "Transaction details"
// @Success 201 {object} models.Transaction
// @Router /transaction [post]
func PostTransaction(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to connect to database",
		})
	}
	defer db.Close()

	var tx models.Transaction
	if err = c.BodyParser(&tx); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if tx.WalletID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Wallet ID is required",
		})
	}

	if tx.Amount <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Amount must be positive",
		})
	}

	if tx.Type != "credit" && tx.Type != "debit" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Transaction type must be 'credit' or 'debit'",
		})
	}

	tx_db, err := db.Beginx()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to start database transaction",
		})
	}
	defer tx_db.Rollback()

	var wallet models.Wallet
	err = tx_db.Get(&wallet, "SELECT id, balance FROM wallets WHERE id = $1", tx.WalletID)
	if err != nil {
		wallet = models.Wallet{
			ID:      tx.WalletID,
			Balance: 0,
		}
		_, err = tx_db.Exec("INSERT INTO wallets (id, balance) VALUES ($1, $2)", wallet.ID, wallet.Balance)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to create wallet",
			})
		}
	}

	if tx.Type == "credit" {
		wallet.Balance += tx.Amount
	} else if tx.Type == "debit" {
		if wallet.Balance < tx.Amount {
			return c.Status(400).JSON(fiber.Map{
				"error": "Insufficient balance",
			})
		}
		wallet.Balance -= tx.Amount
	}

	_, err = tx_db.Exec("UPDATE wallets SET balance = $1 WHERE id = $2", wallet.Balance, wallet.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update wallet balance",
		})
	}

	tx.Timestamp = time.Now()
	var txID int64
	err = tx_db.QueryRow(
		"INSERT INTO transactions (wallet_id, amount, type, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
		tx.WalletID, tx.Amount, tx.Type, tx.Timestamp,
	).Scan(&txID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to record transaction",
		})
	}
	tx.ID = txID

	if err := tx_db.Commit(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	return c.Status(201).JSON(tx)
}

// GetTransactions handles retrieving all transactions for a wallet
// @Description Get all transactions for a specific wallet.
// @Summary get all wallet transactions
// @Tags Wallet
// @Accept json
// @Produce json
// @Param walletId path string true "Wallet ID"
// @Success 200 {array} models.Transaction
// @Router /transactions/{walletId} [get]
func GetTransactions(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to connect to database",
		})
	}
	defer db.Close()

	walletID := c.Params("walletId")
	if walletID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Wallet ID is required",
		})
	}

	var transactions []models.Transaction
	err = db.Select(&transactions,
		"SELECT id, wallet_id, amount, type, created_at FROM transactions WHERE wallet_id = $1 ORDER BY created_at DESC",
		walletID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve transactions",
		})
	}

	return c.JSON(transactions)
}

// GetBalance handles retrieving the balance for a wallet
// @Description Get the current balance for a specific wallet.
// @Summary get wallet balance
// @Tags Wallet
// @Accept json
// @Produce json
// @Param walletId path string true "Wallet ID"
// @Success 200 {object} models.Wallet "wallet_id and balance"
// @Router /balance/{walletId} [get]
func GetBalance(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to connect to database",
		})
	}
	defer db.Close()

	walletID := c.Params("walletId")
	if walletID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Wallet ID is required",
		})
	}

	var wallet models.Wallet
	err = db.Get(&wallet, "SELECT id, balance FROM wallets WHERE id = $1", walletID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Wallet not found",
		})
	}

	return c.JSON(fiber.Map{
		"wallet_id": wallet.ID,
		"balance":   wallet.Balance,
	})
}
