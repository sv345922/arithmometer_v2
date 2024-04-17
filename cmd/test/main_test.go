package main

import (
	"context"
	"fmt"
	"github.com/sv345922/arithmometer_v2/client/appClient"
	"github.com/sv345922/arithmometer_v2/internal/app"
	"github.com/sv345922/arithmometer_v2/internal/calculator"
	"github.com/sv345922/arithmometer_v2/internal/configs"
	"github.com/sv345922/arithmometer_v2/internal/dataBase"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

const testDBname = "/testDB.db"

func TestAll(t *testing.T) {

	fmt.Println("main Test running...")
	ctx, stopTest := context.WithCancel(context.Background())
	dbase, err := dataBase.CreateDb(ctx, testDBname)
	if err != nil {
		panic(err)
	}

	ws, err := app.RunTasker(ctx, dbase)
	if err != nil {
		log.Fatalf("main: %v", err)
	}
	app, _ := app.New(ctx, ws)
	//запуск серверов
	go func() {
		err = app.Run()
		if err != nil {
			log.Fatalf("app.Run: %v", err)
		}
	}()
	// Запуск калькуляторов
	var n = 10 // количество калькуляторов
	for i := 0; i < n; i++ {
		time.Sleep(50 * time.Millisecond)
		go func() {
			calc := calculator.NewCalculator()
			err := calc.Do(ctx)
			if err != nil {
				log.Fatalf("calc.Do: %v", err)
			}
		}()
	}
	time.Sleep(5 * time.Second) // время на запуск сервера

	type testCase struct {
		expression string // Выражение
		nOp        int    // Количество операций в выражении
		timing     int    // Тайминг на одну операцию
		result     string // Ответ
		sendOk     bool   // Подтверждение отправки выражения
		resStatus  string
	}

	tests := []testCase{
		{"", 0, 1,
			"id not found", true, "400 Bad Request"},

		{"-5", 1, 1,
			"-5 = -5.000000", true, "200 OK"},

		{"2+2", 1, 5,
			"2+2 = 4.000000", true, "200 OK"},

		{"-1+2-3/(4+5) * 6 -7 * 8 / 0", 9, 1,
			"-1+2-3/(4+5) * 6 -7 * 8 / 0 = zero division error", true, "200 OK"},

		{"1 + 2 * 3", 2, 2,
			"1 + 2 * 3 = 7.000000", true, "200 OK"},

		{"(1 + 2) * 3 * 100 / (6 - 3 * 2)", 6, 1,
			"(1 + 2) * 3 * 100 / (6 - 3 * 2) = zero division error", true, "200 OK"},

		{"(1 + 2) * 3", 2, 1,
			"(1 + 2) * 3 = 9.000000", true, "200 OK"},

		{"(1 + 2 * 3) / 4", 3, 1,
			"(1 + 2 * 3) / 4 = 1.750000", true, "200 OK"},
	}
	maxTime := 0
	for _, t := range tests {
		maxTime += t.nOp
	}
	//maxTime /= 2
	wg := &sync.WaitGroup{}
	wg.Add(len(tests))
	for i, test := range tests {
		go func(test testCase, n int) {
			defer wg.Done()
			id, ok := appClient.SendNewExpression(test.expression, test.timing)
			time.Sleep(2 * configs.TConst)
			if ok != test.sendOk {
				t.Errorf("#%d: send expression error: got %t, expected %t", n, ok, test.sendOk)
			}
			time.Sleep(time.Duration(maxTime) * configs.TConst)
			status, body, err := appClient.GetResult(id)
			if err != nil || status != test.resStatus || body != test.result {
				t.Errorf("#%d: got status '%s' result '%s' err '%v', expected status '%s', result '%s'",
					n, status, body, err, test.resStatus, test.result)
			}
		}(test, i+1)
	}
	wg.Wait()
	stopTest()
	dbase.Close()
	<-ctx.Done()

	time.Sleep(1 * time.Second)
	err = os.Remove("testDB.db")
	if err != nil {
		log.Println(err)
	}

}
