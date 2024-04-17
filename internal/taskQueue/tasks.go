package taskQueue

import (
	"fmt"
)

type Tasks struct {
	AllTask map[uint64]struct{}
}

// Создать новой список узлов
func NewTasks() *Tasks {
	nodes := make(map[uint64]struct{})
	return &Tasks{AllTask: nodes}
}

// Возвращает длину словаря
func (t *Tasks) Len() int {
	return len(t.AllTask)
}

// Если узел есть в списке, возвращает true, иниче false
func (t *Tasks) Contains(key uint64) bool {
	if _, ok := t.AllTask[key]; ok {
		return true
	}
	return false
}

// Добавить узел в список узлов, если узел с таким id существует в списке, возвращает false
func (t *Tasks) Add(id uint64) bool {
	if t.Contains(id) {
		return false
	}
	t.AllTask[id] = struct{}{}
	return true
}

// Удалить узел из списка узлов по его id, если узла с таким id нет, возвращает false
func (t *Tasks) Remove(key uint64) bool {
	if t.Contains(key) {
		delete(t.AllTask, key)
		return true
	}
	return false
}

// Стрингер
func (t *Tasks) String() string {
	result := ""
	for key, _ := range t.AllTask {
		result += fmt.Sprintf("key=%d", key)
	}
	return result
}

// Возвращает id всех узлов
func (t *Tasks) GetAllIDs() []uint64 {
	result := make([]uint64, len(t.AllTask))
	i := 0
	for key := range t.AllTask {
		result[i] = key
		i++
	}
	return result
}
