package calculator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"io"
	"net/http"
)

// Выбирает оператор вычисления, производит вычисления с учетом тайминга
// и возвращает результат с ошибкой

// запрашивает задачу у оркестратора
func (c *Calculator) GetTask() (bool, error) {
	container := &entities.MessageTask{}
	//container := &calculator.TaskContainer{}
	url := fmt.Sprintf("%s/gettask?id=%d", URL, c.Id)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return false, err
	}
	// Если оркестратор не дал задачу возвращаем nil
	if len(body) == 0 {
		return false, nil
	}
	// Анмаршалим body в контейнер
	err = json.Unmarshal(body, container)
	if err != nil {
		return false, err
	}
	c.Task = container
	return true, nil
}

// Отправляем ответ, если не отправилось, возвращаем ошибку
func (c *Calculator) SendAnswer(container entities.MessageResult) error {
	url := URL + "/giveanswer"

	data, _ := json.Marshal(container) //ошибку пропускаем
	r := bytes.NewReader(data)

	resp, err := http.Post(url, "application/json", r)
	if err != nil {
		fmt.Printf("ошибка отправки запроса POST", err) //  POST%!(EXTRA *url.Error=Post "http://127.0.0.1:8000/giveanswer": EOF)
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("ошибка отправки ответа")
}
