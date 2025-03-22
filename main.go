package main

import (
	"log"
	"os"
	"time"

	_ "github.com/arezvani/wallet-go/docs" // load API Docs files (Swagger)
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
	_ "github.com/lib/pq"
)

type Transaction struct {
	ID        int64     `json:"id" db:"id"`
	WalletID  string    `json:"wallet_id" db:"wallet_id"`
	Amount    float64   `json:"amount" db:"amount"`
	Type      string    `json:"type" db:"type"` // "credit" or "debit"
	Timestamp time.Time `json:"timestamp" db:"created_at"`
}

type Wallet struct {
	ID      string  `json:"id" db:"id"`
	Balance float64 `json:"balance" db:"balance"`
}

func main() {
	connStr := "postgres://wallet:123456@external.abriment.com:30473/wallet?sslmode=disable"

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	createTables(db)

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	app.Post("/transaction", func(c *fiber.Ctx) error {
		var tx Transaction
		if err := c.BodyParser(&tx); err != nil {
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

		var wallet Wallet
		err = tx_db.Get(&wallet, "SELECT id, balance FROM wallets WHERE id = $1", tx.WalletID)
		if err != nil {
			wallet = Wallet{
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
	})

	app.Get("/transactions/:walletId", func(c *fiber.Ctx) error {
		walletID := c.Params("walletId")
		if walletID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Wallet ID is required",
			})
		}

		var transactions []Transaction
		err := db.Select(&transactions,
			"SELECT id, wallet_id, amount, type, created_at FROM transactions WHERE wallet_id = $1 ORDER BY created_at DESC",
			walletID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to retrieve transactions",
			})
		}

		return c.JSON(transactions)
	})

	app.Get("/balance/:walletId", func(c *fiber.Ctx) error {
		walletID := c.Params("walletId")
		if walletID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Wallet ID is required",
			})
		}

		var wallet Wallet
		err := db.Get(&wallet, "SELECT id, balance FROM wallets WHERE id = $1", walletID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "Wallet not found",
			})
		}

		return c.JSON(fiber.Map{
			"wallet_id": wallet.ID,
			"balance":   wallet.Balance,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	app.Listen(":" + port)
}

func createTables(db *sqlx.DB) {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS wallets (
            id VARCHAR(50) PRIMARY KEY,
            balance DECIMAL(20, 2) NOT NULL DEFAULT 0
        )
    `)
	if err != nil {
		log.Fatalf("Failed to create wallets table: %v", err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS transactions (
            id SERIAL PRIMARY KEY,
            wallet_id VARCHAR(50) NOT NULL REFERENCES wallets(id),
            amount DECIMAL(20, 2) NOT NULL,
            type VARCHAR(10) NOT NULL CHECK (type IN ('credit', 'debit')),
            created_at TIMESTAMP NOT NULL DEFAULT NOW()
        )
    `)
	if err != nil {
		log.Fatalf("Failed to create transactions table: %v", err)
	}
}
