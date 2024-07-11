package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	Logger *slog.Logger
)

func setupLogger() {
	logLevel := &slog.LevelVar{}
	opts := &slog.HandlerOptions{
		Level: logLevel,
		ReplaceAttr: func(group []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				return slog.Attr{Key: "severity", Value: a.Value}
			}
			if a.Key == slog.MessageKey {
				return slog.Attr{Key: "message", Value: a.Value}
			}
			return slog.Attr{Key: a.Key, Value: a.Value}
		},
	}
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

func main() {
	setupLogger()
	e := echo.New()
	e.Use(
		middleware.CORSWithConfig(
			middleware.CORSConfig{
				AllowOrigins: []string{"*"},
			},
		),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Level: 5,
		}),
		middleware.Secure(),
	)
	if os.Getenv("DO_DEBUG") != "" {
		e.Use(middleware.Logger())
	}
	e.IPExtractor = echo.ExtractIPFromXFFHeader(
		echo.TrustLoopback(false),   // e.g. ipv4 start with 127.
		echo.TrustLinkLocal(false),  // e.g. ipv4 start with 169.254
		echo.TrustPrivateNet(false), // e.g. ipv4 start with 10. or 192.168
	)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Setup handlers
	for _, entry := range RoutingMap {
		e.GET(entry.Path, entry.Handler).Name = entry.Name
	}

	// Start server
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal and gracefully shutdown the server after 5 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
