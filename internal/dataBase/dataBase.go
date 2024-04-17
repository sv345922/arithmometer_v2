package dataBase

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sv345922/arithmometer_v2/internal/configs"
	"github.com/sv345922/arithmometer_v2/internal/entities"
	"log"
	"os"
)

// структура для взаимождействия с базой данных sql
type DataBase struct {
	// список выражений (с таймингами)
	Tasks       []*entities.Task       `json:"tasks"`
	Expressions []*entities.Expression `json:"expressions"` // []Expression
	AllNodes    []*entities.Node       `json:"allNodes"`
	Timings     *entities.Timings      `json:"timings"`
	Users       []*entities.User       `json:"users"`
}

// Возвращает новый экземпляр структуры
func NewDB() *DataBase {
	db := DataBase{
		Tasks:       make([]*entities.Task, 0),
		Expressions: make([]*entities.Expression, 0),
		AllNodes:    make([]*entities.Node, 0),
		Timings:     &entities.Timings{},
		Users:       make([]*entities.User, 0),
	}
	return &db
}
func (db *DataBase) Save() error {
	jsonBytes, err := json.Marshal(db)
	if err != nil {
		log.Println(err)
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return err
	}
	path := wd + "/db/" + configs.NameDataBase + ".json"
	err = os.WriteFile(path, jsonBytes, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// SafeJSON Сохраняет структуру в базе данных, в папке db
func SafeJSON(name string, dataBase *DataBase) error {
	jsonBytes, err := json.Marshal(dataBase)
	if err != nil {
		log.Println(err)
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return err
	}
	path := wd + "/db/" + configs.NameDataBase + ".json"
	err = os.WriteFile(path, jsonBytes, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Cоздает файл пустой БД
func CreateEmptyDb_v1() error {
	// Создаем файл с пустой БД
	err := SafeJSON(configs.NameDataBase, NewDB())
	if err != nil {
		log.Println("Ошибка создания пустой БД")
	}
	return nil
}

// Проверяет существование файла базы данных, и если он существует удаляет его,
// затем вызывает CreateDb
func CreateEmptyDb(ctx context.Context) (*sql.DB, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	path := wd + configs.DBPath
	ok, err := IsFileExist(path)
	if err != nil {
		return nil, err
	}
	if ok {
		// удаляем файл
		err = os.Remove(path)
	}
	if err != nil {
		return nil, err
	}
	return CreateDb(ctx, configs.DBPath)
}

// Проверяет файл на существование
func IsFileExist(path string) (found bool, err error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
	} else {
		found = true
	}
	return
}

// Создает подключение базы данных, при отсутсвии таблиц сождает новые
func CreateDb(ctx context.Context, nameDB string) (*sql.DB, error) {
	// Создаем файл с пустой БД
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getWD: %w", err)
	}
	path := wd + nameDB
	//fmt.Println("PATH=", path)
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	if err = CreateTables(ctx, db); err != nil {
		return nil, fmt.Errorf("createTables: %w", err)
	}
	return db, nil
}
