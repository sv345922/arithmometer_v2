package appClient

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
)

func SendNewExpression(exprString string, timing int) (string, bool) {
	// Создать запрос
	url := "http://127.0.0.1:" + configs.Port + "/newexpression"
	// Задать тайминги вычислений
	var timings *entities.Timings
	switch timing {
	case 0:
		timings = &entities.Timings{
			Plus:  1,
			Minus: 3,
			Mult:  5,
			Div:   9,
		}
	default:
		timings = &entities.Timings{
			Plus:  timing,
			Minus: timing,
			Mult:  timing,
			Div:   timing,
		}
	}

	var expression = newExpression.NewExpr{
		Expression: exprString,
		Timings:    timings,
	}
	data, _ := json.Marshal(expression) //ошибку пропускаем
	r := bytes.NewReader(data)

	// отправка запроса Post
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		log.Println(err)
		return "", false
	}
	token := ""
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {fmt.Sprintf("Bearer %s", token)},
	}

	resp, err := client.Do(req)
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

	// Отправка запроса Get
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", errTotal
	}
	token := ""
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {fmt.Sprintf("Bearer %s", token)},
	}
	resp, err := client.Do(req)
	//resp, err := http.Get(url)
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
