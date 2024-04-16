package main

import (
	"github.com/sv345922/arithmometer_v2/internal/calculator"
	"log"
)

// Выполняет запросы оркестратору и вычисляет выражение
// TODO периодическое подтверждение работы
func main() {
	calc := calculator.NewCalculator()
	err := calc.Do()
	log.Println(err)
}
