package main

import (
	"context"
	"github.com/sv345922/arithmometer_v2/internal/calculator"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Выполняет запросы оркестратору и вычисляет выражение
// TODO периодическое подтверждение работы
func main() {
	ctx := context.Background()
	calc := calculator.NewCalculator()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		log.Println("Calculator was stopped")
		os.Exit(0)
	}()

	err := calc.Do(ctx)
	log.Println(err)
}
