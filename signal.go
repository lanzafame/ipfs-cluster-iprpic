package iprpic

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lanzafame/bobblehat/sense/screen"
)

// SignalHandler blocks, waiting for a signal from the system,
// before clearing the screen.
func SignalHandler(ctx context.Context, cancel func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cancel()
	screen.Clear()
}
