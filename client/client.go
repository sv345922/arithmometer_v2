package main

import (
	"fmt"
	"github.com/sv345922/arithmometer_v2/client/appClient"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// Задача для вычисления
var expressions []string = []string{
	"-1+2-3/(4+5) * 6 -7 * 8 / 0 + 1", // zero division
	"1 + 2 * 3",                       // 7
	"(1 + 2) * 3 * 100 / (6 - 3 * 2)", // zero division
	"(1 + 2) * 3",                     // 9
	"(1 + 2 * 3) / 4",                 // 1.75
}

func main() {
	// отправка выражения
	//flag.Parse()
	if len(os.Args) > 1 {
		expressions = []string{os.Args[1]}
	}
	idExpressions := make([]string, 0)
	for _, expr := range expressions {
		id, _ := appClient.SendNewExpression(expr, 0)
		idExpressions = append(idExpressions, id)
		fmt.Println()
		//fmt.Println(id)
	}
	wg := sync.WaitGroup{}
	wg.Add(len(idExpressions))
	for _, id := range idExpressions {
		go func(id string) {
			defer wg.Done()
			// получение ответа
			_, answer, err := appClient.GetResult(id)
			if err != nil {
				log.Fatal(err)
			}
			for !strings.Contains(answer, "=") {
				time.Sleep(10 * time.Second)
				_, answer, err = appClient.GetResult(id)
			}
		}(id)
	}
	wg.Wait()
}
