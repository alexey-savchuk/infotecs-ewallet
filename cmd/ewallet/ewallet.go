package main

import (
	"database/sql"
	"os"

	"github.com/alexey-savchuk/infotecs-ewallet/internal/handlers"
	"github.com/alexey-savchuk/infotecs-ewallet/internal/repository/postgres"
	"github.com/alexey-savchuk/infotecs-ewallet/internal/service"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	e := echo.New()
	e.Debug = true

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		e.Logger.Fatal(err)
	}

	wr := postgres.NewWalletRepository(db)
	tr := postgres.NewTransferRepository(db)

	ws := service.NewWalletService(wr)
	ts := service.NewTransferService(tr)

	h := handlers.NewHandlers(ws, ts)

	v1 := e.Group("/api/v1")

	v1.GET("/wallet/:walletId", h.GetWallet)
	v1.GET("/wallet/:walletId/history", h.GetWalletTransfers)

	v1.POST("/wallet", h.CreateWallet)
	v1.POST("/wallet/:walletId/send", h.CreateTransfer)

	err = e.Start(":8080")
	if err != nil {
		e.Logger.Fatal(err)
	}
}
