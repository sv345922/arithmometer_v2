package autostart

import (
	"context"
	"github.com/sv345922/arithmometer_v2/internal/app"
	"github.com/sv345922/arithmometer_v2/internal/calculator"
	"github.com/sv345922/arithmometer_v2/internal/dataBase"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func StartOrchestrator(ctx context.Context, stopCtx context.CancelFunc) error {
	dbase, err := dataBase.CreateEmptyDb(ctx)
	if err != nil {
		return err
	}
	defer dbase.Close()

	ws, err := app.RunTasker(ctx, dbase)
	if err != nil {
		return err
	}

	app, _ := app.New(ctx, ws)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		app.GracefulStop(ctx, sig, stopCtx)
	}()
	err = app.Run()
	if err != nil {
		return err
	}
	<-ctx.Done()

	return nil
}

func StartCalculator(ctx context.Context, stopCtx context.CancelFunc) error {
	calc := calculator.NewCalculator()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		log.Println("Calculator was stopped")
		os.Exit(0)
	}()

	err := calc.Do(ctx)
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}
