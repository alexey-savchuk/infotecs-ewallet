package handlers

import (
	"net/http"

	"github.com/alexey-savchuk/infotecs-ewallet/internal/customerrors"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
	ws WalletService
	ts TransferService
}

func NewHandlers(ws WalletService, ts TransferService) *Handlers {
	return &Handlers{
		ws: ws,
		ts: ts,
	}
}

func (h *Handlers) CreateWallet(c echo.Context) error {
	walletDTO, err := h.ws.Create(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, walletDTO)
}

func (h *Handlers) GetWallet(c echo.Context) error {
	walletID := c.Param("walletId")
	walletDTO, err := h.ws.Get(c.Request().Context(), walletID)
	if err == customerrors.ErrWalletNotExists {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, walletDTO)
}

func (h *Handlers) CreateTransfer(c echo.Context) error {
	transferDTO := &TransferDTO{}
	if err := c.Bind(transferDTO); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	transferDTO, err := h.ts.Create(c.Request().Context(), transferDTO)
	if err == customerrors.ErrFromWalletNotExists {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, transferDTO)
}

func (h *Handlers) GetWalletTransfers(c echo.Context) error {
	walletID := c.Param("walletId")
	transferDTOs, err := h.ts.GetWalletTransfers(c.Request().Context(), walletID)
	if err == customerrors.ErrWalletNotExists {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, transferDTOs)
}
