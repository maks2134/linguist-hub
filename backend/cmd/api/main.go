package main

import (
	"log/slog"
	"os"

	"github.com/makskozlov/linguist-hub/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error("failed to run application", "error", err)
		os.Exit(1)
	}
}
