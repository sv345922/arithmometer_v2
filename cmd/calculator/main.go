package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/calculator"
	"github.com/sv345922/arithmometer_v2/internal/configs"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"io"
	"log"
	"net/http"
	"time"
)

var URL = "http://127.0.0.1:" + configs.Port

// запрашивает задачу у оркестратора
func getTask(calcId string) (*entities.MessageTask, error) {
	container := &entities.MessageTask{}
	//container := &calculator.TaskContainer{}
	url := URL + "/gettask?id=" + calcId
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	// Если оркестратор не дал задачу возвращаем nil
	if len(body) == 0 {
		return nil, nil
	}
	// Анмаршалим body в контейнер
	err = json.Unmarshal(body, container)
	if err != nil {
		return nil, err
	}
	return container, nil
}

// Отправляем ответ, если не отправилось, возвращаем ошибку
func SendAnswer(container entities.MessageResult) error {
	url := URL + "/giveanswer"

	data, _ := json.Marshal(container) //ошибку пропускаем
	r := bytes.NewReader(data)

	resp, err := http.Post(url, "application/json", r)
	if err != nil {
		fmt.Printf("ошибка отправки запроса POST", err) //TODO delete
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("ошибка отправки ответа")
}

// Выполняет запросы оркестратору и вычисляет выражение
// TODO периодическое подтверждение работы
func main() {
	calcId := entities.GetDelta(5)
	log.Printf("calculator %d is runing", calcId)
	result := make(chan entities.MessageResult)
	for {
		container, err := getTask(fmt.Sprintf("%d", calcId))
		if err != nil {
			log.Println("ошибка получения задачи", err)
			time.Sleep(5 * time.Second)
			continue
		}
		// Окрестратор не дал задание
		if container == nil {
			log.Println("tik ")
			time.Sleep(5 * time.Second)
			continue
		}
		//log.Println("задача принята")
		// запускаем задачу в горутине
		go func(container *entities.MessageTask) {
			res, err := calculator.Do(container)
			errString := "nil"
			if err != nil {
				errString = err.Error()
			}
			result <- entities.MessageResult{
				Id:     container.Id,
				Result: res,
				Err:    errString,
			}
		}(container)
		answer := <-result
		log.Printf("задача %.3f%s%.3f выполнена, результат %f, ошибка %s\n",
			container.X,
			container.Op,
			container.Y,
			answer.Result,
			answer.Err,
		)
		// отправляем ответ, до тех пор пока он не будет принят
		for {
			err = SendAnswer(answer)
			if err == nil {
				break
			}
		}
		log.Println("отправлен ответ")
		time.Sleep(5 * time.Second)
	}
}
