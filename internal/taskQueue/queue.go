package taskQueue

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"log"
	"sync"
)

// Queue - Очередь задач.
// ReadyToCalc - задачи, готовые для выдачи вычислителям,
// Working - задачи, взятые вычислителем
// NotReady - задачи, ожидающие решения других задач
// WaitingIds - id готовых для вычисления задач
// L - количество элементов в очереди (всего)
type Queue struct {
	AllTasks    map[uint64]*Task `json:"allTasks"`    // мапа всех задач
	ReadyToCalc []uint64         `json:"readytocalc"` // готовые для выдачи вычислителям,
	Working     *Tasks           `json:"working"`     // взятые вычислителем
	NotReady    *Tasks           `json:"notready"`    // ожидающие решения других задач
	L           uint             `json:"l"`           // количество элементов в очереди (всего)
	mu          sync.RWMutex
}

// NewTasks Возвращает указатель на новую очередь задач
func NewQueue() *Queue {
	working := NewTasks()
	notReady := NewTasks()
	return &Queue{
		AllTasks:    make(map[uint64]*Task),
		ReadyToCalc: make([]uint64, 0),
		Working:     working,
		NotReady:    notReady,
		L:           0,
		mu:          sync.RWMutex{},
	}
}
func (q *Queue) Info() string {
	return fmt.Sprintf("В очереди %d задач(и), из них %d - в работе, %d - готовы к вычислению, %d не готовы",
		q.L,
		q.Working.Len(),
		len(q.ReadyToCalc),
		q.NotReady.Len(),
	)
}

// Создает и добавляет таски в очeредь
func (q *Queue) AddExpressionNodes(ctx context.Context, tx *sql.Tx, nodes []*entities.Node) error {
	// Создаем словарь узлов
	nodesMap := make(map[uint64]*entities.Node)
	for _, node := range nodes {
		nodesMap[node.Id] = node
	}
	for _, node := range nodes {
		// Если узел лист==не оператор, то в задачи не попадает
		if node.Sheet {
			continue
		}
		task := NewTask(node)
		// Если оператор X вычислен
		if xNode, ok := nodesMap[node.X]; ok && xNode.Calculated {
			task.X = xNode.Val
			task.XReady = true
		}
		// Если оператор Y вычислен
		if yNode, ok := nodesMap[node.Y]; ok && yNode.Calculated {
			task.Y = yNode.Val
			task.YReady = true
		}
		// Добавляем запись в БД
		err := InsertTask(ctx, tx, task)
		if err != nil {
			return fmt.Errorf("cannot add task to queue: %w", err)
		}
		// Добавляем в очередь
		_ = q.AddTask(task)
	}
	return nil
}

// AddTask Добавляет задачу в список задач NotReady и увеличивает счетчик L,
func (q *Queue) AddTask(task *Task) bool {
	id := task.Node.Id
	// Проверка на наличие задачи с таким же id в очереди
	if _, ok := q.AllTasks[id]; ok {
		return false
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.NotReady.Add(id) {
		q.AllTasks[id] = task
		q.L++
		//fmt.Printf("id %d добавлена в очередь\n", id) // TODO delete
	}
	return true
}

// RemoveTask Удаляет задачу из очереди задач
func (q *Queue) RemoveTask(idTask uint64) bool {
	// var msg = "задача удалена"
	q.mu.RLock()
	// проверяем наличие задачи в очереди
	if _, ok := q.AllTasks[idTask]; !ok {
		q.mu.RUnlock()
		//log.Println("задача", idTask, "нет в очереди") // TODO delete?
		return false
	}
	q.mu.RUnlock()

	// удаляем из Working
	if q.Working.Remove(idTask) {
		delete(q.AllTasks, idTask)
		q.L--
		//log.Println(msg, idTask, "из working")
		return true
	}
	// удаляем из NotReady
	if q.NotReady.Remove(idTask) {
		delete(q.AllTasks, idTask)
		q.L--
		// log.Println(msg, idTask, "из NotReady") //TODO delete
		return true
	}
	// удаляем из ReadyToCalc
	q.mu.Lock()
	defer q.mu.Unlock()
	for i := 0; i < len(q.ReadyToCalc); i++ {
		if q.ReadyToCalc[i] == idTask {
			q.ReadyToCalc = append(q.ReadyToCalc[:i], q.ReadyToCalc[i+1:]...) // TODO возможная ошибка
			delete(q.AllTasks, idTask)
			q.L--
			// log.Println(msg, idTask, "из ReadyToCalc") //TODO delete
			return true
		}
	}
	log.Printf("ошибка в очереди (элемент есть в AllTask но нет в других местах) при id=%d", idTask)
	return false
}

// GetTask
// Возвращает свободную задачу для вычислителя,
// переносит эту задачу в мапу работающих задач. При пустой очереди возвращает nil.
// Работа с таймингами и id вычислителя снаружи функции.
// Сначала обновляет очередь: проверяем в NotReady и если задача готова для вычисления
// переносим её в waiting.
func (q *Queue) GetTask() *Task {
	// обновляем очереди
	q.CheckDeadlines()
	q.UpdateReady()

	// берем первый элемент из очереди
	q.mu.Lock()
	defer q.mu.Unlock()
	l := len(q.ReadyToCalc)
	// если очередь пустая возвращаем nil
	if l == 0 {
		return nil
	}
	// Получаем задачу - первый элемент очереди
	resultID := q.ReadyToCalc[0]
	// удаляем из списка ReadyToCalc
	switch l {
	case 1:
		q.ReadyToCalc = q.ReadyToCalc[:0] // при длине 1 опустошаем очередь
	default:
		q.ReadyToCalc = q.ReadyToCalc[1:] // иначе оставляем очередь без первого элемента
	}
	// переносим в список ожидающих решения
	q.Working.Add(resultID)

	return q.AllTasks[resultID]
}

// UpdateReady Обновляет очередь задач, находит среди NotReady готовые к вычислению
// и переносит их в ReadyToCalc
// TODO обновить с учетом выгрузки из БД всех tasks в notReady (проверять поле calcID)
func (q *Queue) UpdateReady() {
	// получаем список ключей
	keys := q.NotReady.GetAllIDs()
	// проходим по ключам и проверяем хранящиеся задачи на готовность к вычислению
	for _, key := range keys {
		q.mu.RLock()
		task := q.AllTasks[key]
		q.mu.RUnlock()
		if task.IsReadyToCalc() {
			q.mu.Lock()
			q.ReadyToCalc = append(q.ReadyToCalc, key)
			q.mu.Unlock()
			q.NotReady.Remove(key)
		}
	}
}

// CheckDeadlines Функция обновления очереди по состоянию таймингов
// Если среди работающих задач есть с простроченным дедлайном,
// то задача переносится в список ожидающих.
// Возвращает количество просроченных задач
func (q *Queue) CheckDeadlines() int {
	// получаем список ключей
	keys := q.Working.GetAllIDs()
	n := 0
	for _, key := range keys {
		q.mu.RLock()
		task := q.AllTasks[key]
		q.mu.RUnlock()
		// если задача с прошедшим дедлайном
		if task.IsTimeout() {
			// увеличиваем счетчик просроченных
			n++
			// устанавливаем дедлайн в далекое будущее
			task.SetDeadline(3600 * 240)
			// и перемещаем задачу в начало очереди ожидающих
			q.mu.Lock()
			q.ReadyToCalc = append([]uint64{key}, q.ReadyToCalc...)
			q.mu.Unlock()
			q.Working.Remove(key)
		}
	}
	return n
}

// Устанавливает ответ полученный от вычислителя
// Возвращает true если вычислен конренвой узел
func (q *Queue) AddAnswer(id uint64, answer float64) (bool, error) {
	// если нет task с таким id выходим с ошибкой
	q.mu.RLock()
	task, ok := q.AllTasks[id]
	q.mu.RUnlock()
	if !ok {
		return false, entities.NoTaskInQueue
	}
	//Если нашли, записываем результат
	task.SetResult(answer)

	// проверяем Working
	if q.Working.Contains(id) {
		// обновляем родительский узел
		rootFlag, err := q.UpdateParent(answer, task)
		// удаляем узел из очереди
		q.RemoveTask(id)
		return rootFlag, err
	}
	// проверяем NotReady - TODO не должно быть таких ситуаций
	if q.NotReady.Contains(id); task != nil {
		// обновляем родительский узел
		rootFlag, err := q.UpdateParent(answer, task)
		// удаляем узел из очереди
		q.RemoveTask(id)
		return rootFlag, err
	}
	// TODO надо замьютить ReadyToCalc через копию возможно
	// проверяем ReadyToCalc - если вычислитель сильно запаздал и задача вернулась в ожидающие
	for _, taskId := range q.ReadyToCalc {
		if taskId == id {
			rootFlag, err := q.UpdateParent(answer, task)
			q.RemoveTask(id)
			return rootFlag, err
		}
	}
	return false, entities.NoTaskInQueue
}

// Устанавливает значения и флаги X/Y у родительского узла
// возвращает true, если у узла нет родителя==он корень дерева
func (q *Queue) UpdateParent(answer float64, task *Task) (bool, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if parent := q.AllTasks[task.Node.Parent]; parent != nil {
		switch task.Node.Id {
		case parent.Node.X:
			parent.X = answer
			parent.XReady = true
		case parent.Node.Y:
			parent.Y = answer
			parent.YReady = true
		default:
			return false, errors.Join(fmt.Errorf("ошибка при добавлении ответа"), entities.QueueError)
		}
		return false, nil
	}
	return true, nil
}
