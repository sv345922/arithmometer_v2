package main

import (
	"context"
	"github.com/sv345922/arithmometer_v2/internal/calculator"
	"log"
)

// Выполняет запросы оркестратору и вычисляет выражение
// TODO периодическое подтверждение работы
func main() {
	ctx := context.Background()
	calc := calculator.NewCalculator()
	err := calc.Do(ctx)
	log.Println(err)
}
