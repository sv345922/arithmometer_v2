package wSpace

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/dataBase"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"github.com/sv345922/arithmometer_v2/internal/parser"
	"github.com/sv345922/arithmometer_v2/internal/taskQueue"
	"log"
	"sync"
)

type WorkingSpace struct {
	Queue       *taskQueue.Queue                `json:"tasks"`
	Expressions map[uint64]*entities.Expression `json:"expressions"`
	AllNodes    map[uint64]*entities.Node       `json:"allnodes"`
	Timings     *entities.Timings               `json:"timings"`
	DB          *sql.DB
	Users       map[uint64]*entities.User
	Mu          sync.RWMutex
}

//func (ws *WorkingSpace) apply(options *interface{}) {
//	//TODO implement me
//	panic("implement me")
//}

func NewWorkingSpace(db *sql.DB) *WorkingSpace {
	ws := new(WorkingSpace)
	ws.Queue = taskQueue.NewQueue()
	expressions := make(map[uint64]*entities.Expression)
	ws.Expressions = expressions
	allNodes := make(map[uint64]*entities.Node)
	ws.AllNodes = allNodes
	ws.DB = db
	users := make(map[uint64]*entities.User)
	ws.Users = users
	return ws
}

// Сохраняет рабочее пространство
// TODO не будет использоваться

func (ws *WorkingSpace) Save() error {
	ws.Mu.RLock()
	// Создаем и заполняем структуру базы данных для сохранения
	db := dataBase.NewDB()
	// Заполняем структуру БД
	// заполняем тайминги
	db.Timings = ws.Timings
	// заполняем очередь задач
	for _, task := range ws.Queue.AllTasks {
		dTask := entities.Task{
			NodeId:   task.Node.Id,
			X:        task.X,
			XReady:   task.XReady,
			Y:        task.Y,
			YReady:   task.YReady,
			CalcId:   task.CalcId,
			Deadline: task.Deadline,
			Duration: task.Duration,
		}
		db.Tasks = append(db.Tasks, &dTask)
	}

	// заполняем список выражений
	for _, expression := range ws.Expressions {
		db.Expressions = append(db.Expressions, expression)
	}

	// заполняем AllNodes
	for _, node := range ws.AllNodes {
		db.AllNodes = append(db.AllNodes, node)
	}
	// заполняем Users
	for _, user := range ws.Users {
		db.Users = append(db.Users, user)
	}

	ws.Mu.RUnlock()
	err := db.Save()
	if err != nil {
		return err
	}
	log.Println("DB saved")
	return nil
}

// Добавляет новое выражение в структуры,
// обновляет мапу узлов
// обновляет очередь вычислений
func (ws *WorkingSpace) AddToExpressions(ctx context.Context, tx *sql.Tx, expression *entities.Expression) (uint64, error) {
	// Сохраняем выражение в БД с получением его id
	id, err := InsertExpression(ctx, tx, expression)
	if err != nil {
		return 0, fmt.Errorf("cannot add expression: %w", err)
	}
	// добавляем выражение в список выражений
	ws.Mu.Lock()
	expression.Id = id

	// Добавляем выражение в мапу выражений
	ws.Expressions[expression.Id] = expression
	ws.Mu.Unlock()
	return id, nil
}

// Добавляем в allNodes узлы выражения, с сохранением в БД и присвоением id
// возвращаем список узлов выражения (тип entities.Node), ошибку
// tx возвращаем, потомучто база данных блокируется иначе?
func (ws *WorkingSpace) InsertToAllNodes(ctx context.Context,
	tx *sql.Tx,
	parseNodes []*parser.Node) ([]*entities.Node, error, *sql.Tx) {
	// записать в БД и получить id
	tx, err := ws.DB.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}
	for _, parseNode := range parseNodes {
		node := parser.TransformNode(parseNode)
		id, err := InsertNode(ctx, tx, node)
		if err != nil {
			return nil, fmt.Errorf("cannot save node: %w", err), tx
		}
		parseNode.NodeId = id
	}
	// список узлов с установленными id связей
	nodes := make([]*entities.Node, 0, len(parseNodes))
	// записать в ws
	for _, parseNode := range parseNodes {
		node := parser.TransformNode(parseNode)
		ws.AllNodes[node.Id] = node
		nodes = append(nodes, node)
	}
	// обновить ссылочные поля в БД
	for _, node := range nodes {
		err = UpdateNode(ctx, tx, node)
		if err != nil {
			return nil, fmt.Errorf("cannot update node: %w", err), tx
		}
	}

	return nodes, nil, tx
}

// Возвращает корень узла и выражение с этим корнем
func (ws *WorkingSpace) GetRoot(nodeId uint64) (*entities.Node, *entities.Expression, error) {
	ws.Mu.RLock()
	defer ws.Mu.RUnlock()
	expression, ok := ws.Expressions[ws.AllNodes[nodeId].ExpressionId]
	if !ok {
		return nil, nil, fmt.Errorf("cannot find expression")
	}
	root, ok := ws.AllNodes[expression.RootId]
	if !ok {
		return nil, nil, fmt.Errorf("cannot find root")
	}
	return root, expression, nil
}

// Возвращает список id узлов выражения по id корня
func (ws *WorkingSpace) GetExpressionNodesID(nodeId uint64, nodesId *[]uint64) {
	*nodesId = append(*nodesId, nodeId)
	ws.Mu.RLock()
	node, ok := ws.AllNodes[nodeId]
	ws.Mu.RUnlock()
	if !ok || node.IsSheet() {
		return
	}

	ws.GetExpressionNodesID(node.X, nodesId)
	ws.GetExpressionNodesID(node.Y, nodesId)
}

// Загружает структуру из db и возвращает её
func LoadDB(ctx context.Context, db *sql.DB) (*WorkingSpace, error) {
	// Создаем рабочее пространство
	ws := NewWorkingSpace(db)
	dataBase := dataBase.NewDB()
	err := dataBase.Load(ctx, db)
	if err != nil {
		//log.Println(err)
		return ws, fmt.Errorf("load database: %w", err)
	}
	ws.Mu.Lock()
	// заполняем поля
	// Expressions
	for _, expression := range dataBase.Expressions {
		ws.Expressions[expression.Id] = expression
	}
	// AllNodes
	for _, node := range dataBase.AllNodes {
		ws.AllNodes[node.Id] = node
	}
	// Users
	for _, user := range dataBase.Users {
		ws.Users[user.ID] = user
	}
	//Queue
	for _, dtask := range dataBase.Tasks {
		node, ok := ws.AllNodes[dtask.NodeId]
		if !ok {
			return ws, fmt.Errorf("node %d not found", dtask.NodeId)
		}
		t := taskQueue.Task{
			Node:     node,
			X:        dtask.X,
			XReady:   dtask.XReady,
			Y:        dtask.Y,
			YReady:   dtask.YReady,
			CalcId:   dtask.CalcId,
			Deadline: dtask.Deadline,
			Duration: dtask.Duration,
		}

		ok = ws.Queue.AddTask(&t)
		if !ok {
			log.Printf("add task %d fail", dtask.NodeId)
		}

	}
	// timings
	ws.Timings = dataBase.Timings
	ws.Mu.Unlock()
	return ws, nil
}
