package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/configs"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"github.com/sv345922/arithmometer_v2/internal/useCases/newExpression"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Задача для вычисления
var expressions []string = []string{
	"-1+2-3/(4+5) * 6 -7 * 8 / 0",     // zero division
	"1 + 2 * 3",                       // 7
	"(1 + 2) * 3 * 100 / (6 - 3 * 2)", // zero division
	"(1 + 2) * 3",                     // 9
	"(1 + 2 * 3) / 4",                 //
}

// var expr = "2 / 0"

func SendNewExpression(exprString string) (string, bool) {
	// Создать запрос
	url := "http://127.0.0.1:" + configs.Port + "/newexpression"
	// Задать тайминги вычислений
	timing := &entities.Timings{
		Plus:  1,
		Minus: 1,
		Mult:  2,
		Div:   2,
	}
	//timing = nil
	var expression = newExpression.NewExpr{
		Expression: exprString,
		Timings:    timing,
	}
	data, _ := json.Marshal(expression) //ошибку пропускаем
	r := bytes.NewReader(data)
	resp, err := http.Post(url, "application/json", r)
	if err != nil {
		return "", false
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false
	}
	fmt.Printf("Постановка задачи\nStatus: %s\tBody:\t%s\n", resp.Status, string(body))
	id := string(body)
	//fmt.Println("Задача отправлена")
	return id, true
}

// Получает результат вычислений
func GetResult(id string) (string, string, error) {
	errTotal := errors.New("ошибка получения результата")
	// Создать запрос
	url := "http://127.0.0.1:" + configs.Port + "/getresult" + "?id=" + id
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", "", errTotal
	}
	fmt.Printf("Получение результата\nStatus: %s\tBody:\t%s\n", resp.Status, string(body))
	return resp.Status, string(body), nil
}
func main() {
	// отправка выражения
	//flag.Parse()
	if len(os.Args) > 1 {
		expressions = []string{os.Args[1]}
	}
	idExpressions := make([]string, 0)
	for _, expr := range expressions {
		id, _ := SendNewExpression(expr)
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
			_, answer, err := GetResult(id)
			if err != nil {
				log.Fatal(err)
			}
			for !strings.Contains(answer, "=") {
				time.Sleep(10 * time.Second)
				_, answer, err = GetResult(id)
			}
		}(id)
	}
	wg.Wait()
}
