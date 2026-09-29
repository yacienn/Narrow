// Composition root: the only place that knows about every layer.
package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/yacienn/Game/internal/repository/memory"
	"github.com/yacienn/Game/internal/transport/ws"
	"github.com/yacienn/Game/internal/usecase"
)

func main() {
	players := memory.NewPlayerRepository()
	hub := ws.NewHub()
	game := usecase.NewGameService(players, hub)

	mux := http.NewServeMux()
	mux.Handle("/game", ws.NewHandler(game, hub))

	addr := "0.0.0.0:8080"
	slog.Info("server listening", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}
