package taskQueue

import (
	"time"

	"github.com/sv345922/arithmometer_v2/internal/configs"
	"github.com/sv345922/arithmometer_v2/internal/entities"
)

type Task struct {
	Node     *entities.Node `json:"node"`
	X        float64        `json:"x"`
	XReady   bool           `json:"xready"`
	Y        float64        `json:"y"`
	YReady   bool           `json:"yready"`
	Error    error          `json:"error"`    // ошибка
	CalcId   uint64         `json:"calc_id"`  // id вычислителя задачи
	Deadline time.Time      `json:"deadline"` // дедлайн задачи
	Duration int            `json:"duration"` // длительность операции в условных единицах из конфига
}

func NewTask(node *entities.Node) *Task {
	return &Task{Node: node, Deadline: time.Now().Add(time.Hour * 1000)}
}

// GetID Возвращает id
func (t *Task) GetID() uint64 {
	return t.Node.Id
}

// SetCalc Присваивает id вычислителя
func (t *Task) SetCalc(calcId uint64) {
	t.CalcId = calcId
}

// Проверка на завершение дедлайна задачи, если время вышло, возвращает true
func (t *Task) IsTimeout() bool {
	if t.Deadline.Before(time.Now()) {
		return true
	}
	return false
}

// проверяет готовность задачи для расчетов
func (t *Task) IsReadyToCalc() bool {
	if !t.Node.Calculated {
		if t.XReady && t.YReady {
			return true
		}
	}
	return false
}

// SetDeadline устанавливает дедлайн задаче от текущего момента
func (t *Task) SetDeadline(add int) {
	t.Deadline = time.Now().Add(time.Duration(add) * configs.TConst)
}

func (t *Task) SetDuration(duration int) {
	t.Duration = duration
}

// Заносит результат выражения в узел
func (t *Task) SetResult(result float64) {
	t.Node.Calculated = true
	t.Node.Val = result
}
