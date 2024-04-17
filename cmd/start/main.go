package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/sv345922/arithmometer_v2/internal/useCases/autostart"
)

// Запускает все модули сервиса, параметром принимает количество калькуляторов
func main() {
	var nCalc int = 5
	var err error
	if len(os.Args) > 1 {
		nCalc, err = strconv.Atoi(os.Args[1])
		if err != nil {
			nCalc = 5
		}
	}
	ctx, stopCtx := context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)
	// запуск оркестратора
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := autostart.StartOrchestrator(ctx, stopCtx)
		if err != nil {
			log.Fatal("Ошибка запуска оркестратора", err)
		}
	}()
	// запуск калькуляторов
	wg.Add(nCalc)
	for i := 0; i < nCalc; i++ {
		time.Sleep(50 * time.Millisecond)
		go func() {
			defer wg.Done()
			err := autostart.StartCalculator(ctx, stopCtx)
			if err != nil {
				log.Println("Ошибка запуска калькулятора #", i+1, err)
			}
		}()
	}
	wg.Wait()
}
