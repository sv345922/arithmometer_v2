package taskQueue

import (
	"fmt"
	"sync"
)

type Tasks struct {
	AllTask map[uint64]struct{}
	mu      sync.RWMutex
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
	// TODO проверить мьютексы
	//t.mu.RLock()
	//defer t.mu.RUnlock()
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
	// TODO проверить мьютексы
	// t.mu.Lock()
	// defer t.mu.Unlock()
	t.AllTask[id] = struct{}{}
	return true
}

// TODO не будет использоваться
//// возвращает узел из списка по его id, если id нет в списке узлов, возвращает nil
//func (t *Tasks) Get(id uint64) *Task {
//	if t.Contains(id) {
//		//t.mu.RLock()
//		// defer t.mu.RUnlock
//		result := t.AllTask[id]
//		return result
//	}
//	return nil
//}

// Удалить узел из списка узлов по его id, если узла с таким id нет, возвращает false
func (t *Tasks) Remove(key uint64) bool {
	if t.Contains(key) {
		t.mu.Lock()
		defer t.mu.Unlock()
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
	// TODO проверить мьютексы
	t.mu.RLock()
	defer t.mu.RUnlock()
	for key := range t.AllTask {
		result[i] = key
		i++
	}
	return result
}
